package updater

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/yiminlab/quick-workflow/internal/ui"
	"github.com/yiminlab/quick-workflow/pkg/config"
)

const (
	githubAPIURL     = "https://api.github.com/repos/yiminlab/quick-workflow/releases/latest"
	updateCheckFile  = ".last_update_check"
	checkInterval    = 24 * time.Hour // Check once per day
)

// GitHubRelease represents a GitHub release
type GitHubRelease struct {
	TagName string `json:"tag_name"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

// CheckAndUpdate checks for updates and updates if auto_update is enabled
func CheckAndUpdate(currentVersion string, autoUpdate bool) error {
	// Check if we should check for updates (once per day)
	if !shouldCheckForUpdates() {
		return nil
	}

	latestVersion, downloadURL, err := getLatestVersion()
	if err != nil {
		// Silently fail if we can't check for updates
		return nil
	}

	// Update last check time
	updateLastCheckTime()

	if !isNewerVersion(currentVersion, latestVersion) {
		return nil
	}

	// New version available
	ui.Info(fmt.Sprintf("🎉 New version available: %s (current: %s)", latestVersion, currentVersion))

	if !autoUpdate {
		ui.Info(fmt.Sprintf("Run 'qkflow update-cli' to update, or visit: https://github.com/yiminlab/quick-workflow/releases"))
		return nil
	}

	// Auto update
	ui.Info("⬇️  Downloading update...")
	if err := downloadAndInstall(downloadURL); err != nil {
		ui.Error(fmt.Sprintf("Failed to auto-update: %v", err))
		ui.Info("Please update manually: https://github.com/yiminlab/quick-workflow/releases")
		return nil
	}

	ui.Success(fmt.Sprintf("✅ Successfully updated to version %s! Please restart qkflow.", latestVersion))
	os.Exit(0)
	return nil
}

// shouldCheckForUpdates checks if enough time has passed since last check
func shouldCheckForUpdates() bool {
	configDir, err := config.GetConfigDir()
	if err != nil {
		return true
	}

	checkFile := filepath.Join(configDir, updateCheckFile)
	info, err := os.Stat(checkFile)
	if err != nil {
		return true // File doesn't exist, check for updates
	}

	return time.Since(info.ModTime()) > checkInterval
}

// updateLastCheckTime updates the timestamp of last update check
func updateLastCheckTime() {
	configDir, err := config.GetConfigDir()
	if err != nil {
		return
	}

	checkFile := filepath.Join(configDir, updateCheckFile)
	f, err := os.Create(checkFile)
	if err != nil {
		return
	}
	defer f.Close()
	f.WriteString(time.Now().Format(time.RFC3339))
}

// getLatestVersion fetches the latest version from GitHub
func getLatestVersion() (string, string, error) {
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(githubAPIURL)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("failed to fetch latest version: %s", resp.Status)
	}

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", "", err
	}

	// Find the appropriate asset for current platform
	assetName := getBinaryName()
	var downloadURL string
	for _, asset := range release.Assets {
		if asset.Name == assetName {
			downloadURL = asset.BrowserDownloadURL
			break
		}
	}

	if downloadURL == "" {
		return "", "", fmt.Errorf("no binary found for current platform")
	}

	return release.TagName, downloadURL, nil
}

// getBinaryName returns the binary name for current platform
func getBinaryName() string {
	os := runtime.GOOS
	arch := runtime.GOARCH

	switch os {
	case "darwin":
		return fmt.Sprintf("qkflow-darwin-%s", arch)
	case "linux":
		return fmt.Sprintf("qkflow-linux-%s", arch)
	case "windows":
		return fmt.Sprintf("qkflow-windows-%s.exe", arch)
	default:
		return ""
	}
}

// isNewerVersion checks if the remote version is newer than current
func isNewerVersion(current, latest string) bool {
	// Remove 'v' prefix if present
	current = strings.TrimPrefix(current, "v")
	latest = strings.TrimPrefix(latest, "v")

	// Simple version comparison (works for semantic versioning)
	return latest > current
}

// downloadAndInstall downloads and installs the new binary
func downloadAndInstall(url string) error {
	// Get current executable path
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}

	// Follow symlinks to get the real path
	execPath, err = filepath.EvalSymlinks(execPath)
	if err != nil {
		return fmt.Errorf("failed to resolve executable path: %w", err)
	}

	// Download to temporary file
	tmpFile, err := os.CreateTemp("", "qkflow-*")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()
	defer os.Remove(tmpPath)

	client := &http.Client{Timeout: 5 * time.Minute}
	resp, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download update: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download update: %s", resp.Status)
	}

	if _, err := io.Copy(tmpFile, resp.Body); err != nil {
		tmpFile.Close()
		return fmt.Errorf("failed to save update: %w", err)
	}
	tmpFile.Close()

	// Make it executable
	if err := os.Chmod(tmpPath, 0755); err != nil {
		return fmt.Errorf("failed to make binary executable: %w", err)
	}

	// Backup current binary
	backupPath := execPath + ".backup"
	if err := os.Rename(execPath, backupPath); err != nil {
		return fmt.Errorf("failed to backup current binary: %w", err)
	}

	// Move new binary to current location
	if err := os.Rename(tmpPath, execPath); err != nil {
		// Restore backup on failure
		os.Rename(backupPath, execPath)
		return fmt.Errorf("failed to install update: %w", err)
	}

	// Remove backup
	os.Remove(backupPath)

	return nil
}

// ManualUpdate performs a manual update check and install
func ManualUpdate(currentVersion string) error {
	latestVersion, downloadURL, err := getLatestVersion()
	if err != nil {
		return fmt.Errorf("failed to check for updates: %w", err)
	}

	if !isNewerVersion(currentVersion, latestVersion) {
		ui.Success(fmt.Sprintf("✅ You are already running the latest version (%s)", currentVersion))
		return nil
	}

	ui.Info(fmt.Sprintf("🎉 New version available: %s (current: %s)", latestVersion, currentVersion))
	ui.Info("⬇️  Downloading update...")

	if err := downloadAndInstall(downloadURL); err != nil {
		return fmt.Errorf("failed to update: %w", err)
	}

	ui.Success(fmt.Sprintf("✅ Successfully updated to version %s! Please restart qkflow.", latestVersion))
	return nil
}

