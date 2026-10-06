package editor

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/yiminlab/quick-workflow/internal/ui"
)

// EditorResult contains the markdown content and uploaded files
type EditorResult struct {
	Content string   `json:"content"`
	Files   []string `json:"files"` // Paths to uploaded files
}

// StartEditor starts a web-based editor and returns the user's input
func StartEditor() (*EditorResult, error) {
	// Create temporary directory for uploads
	tempDir, err := os.MkdirTemp("", "qkflow-editor-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}

	server := &editorServer{
		tempDir: tempDir,
		done:    make(chan *EditorResult, 1),
		errCh:   make(chan error, 1),
	}

	// Start HTTP server
	mux := http.NewServeMux()
	mux.HandleFunc("/", server.handleIndex)
	mux.HandleFunc("/upload", server.handleUpload)
	mux.HandleFunc("/save", server.handleSave)
	mux.HandleFunc("/cancel", server.handleCancel)

	httpServer := &http.Server{
		Addr:    "127.0.0.1:0", // Random available port
		Handler: mux,
	}

	// Start server in background
	listener, err := getAvailableListener()
	if err != nil {
		return nil, fmt.Errorf("failed to get available port: %w", err)
	}
	server.port = listener.Addr().(*net.TCPAddr).Port

	go func() {
		if err := httpServer.Serve(listener); err != nil && err != http.ErrServerClosed {
			server.errCh <- err
		}
	}()

	// Open browser
	url := fmt.Sprintf("http://127.0.0.1:%d", server.port)
	ui.Info(fmt.Sprintf("🌐 Opening editor in your browser: %s", url))
	
	if err := openBrowser(url); err != nil {
		ui.Warning("Could not open browser automatically. Please open the URL manually.")
		fmt.Println(url)
	}

	ui.Info("📝 Please edit your content in the browser and click 'Save and Continue'")

	// Wait for result or error
	var result *EditorResult
	select {
	case result = <-server.done:
		// Success
	case err := <-server.errCh:
		httpServer.Shutdown(context.Background())
		os.RemoveAll(tempDir)
		return nil, err
	case <-time.After(30 * time.Minute):
		httpServer.Shutdown(context.Background())
		os.RemoveAll(tempDir)
		return nil, fmt.Errorf("editor timeout after 30 minutes")
	}

	// Shutdown server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	httpServer.Shutdown(ctx)

	return result, nil
}

type editorServer struct {
	tempDir string
	port    int
	done    chan *EditorResult
	errCh   chan error
}

func (s *editorServer) handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(editorHTML))
}

func (s *editorServer) handleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse multipart form (max 100MB)
	if err := r.ParseMultipartForm(100 << 20); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Validate file type
	if !isAllowedFileType(header.Filename) {
		http.Error(w, "File type not allowed", http.StatusBadRequest)
		return
	}

	// Save file to temp directory
	filename := filepath.Base(header.Filename)
	filePath := filepath.Join(s.tempDir, filename)
	
	dst, err := os.Create(filePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return file info
	response := map[string]interface{}{
		"filename": filename,
		"path":     filePath,
		"size":     header.Size,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *editorServer) handleSave(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var data struct {
		Content string   `json:"content"`
		Files   []string `json:"files"`
	}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result := &EditorResult{
		Content: data.Content,
		Files:   data.Files,
	}

	// Send result
	select {
	case s.done <- result:
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	default:
		http.Error(w, "Already processed", http.StatusConflict)
	}
}

func (s *editorServer) handleCancel(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Send empty result (user cancelled)
	select {
	case s.done <- &EditorResult{Content: "", Files: []string{}}:
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	default:
		http.Error(w, "Already processed", http.StatusConflict)
	}
}

func getAvailableListener() (net.Listener, error) {
	return net.Listen("tcp", "127.0.0.1:0")
}

func openBrowser(url string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	default:
		return fmt.Errorf("unsupported platform")
	}

	return cmd.Start()
}

func isAllowedFileType(filename string) bool {
	ext := filepath.Ext(filename)
	allowed := map[string]bool{
		".png":  true,
		".jpg":  true,
		".jpeg": true,
		".gif":  true,
		".webp": true,
		".svg":  true,
		".mp4":  true,
		".mov":  true,
		".webm": true,
		".avi":  true,
	}
	return allowed[ext]
}

