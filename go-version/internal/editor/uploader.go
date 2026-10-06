package editor

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/yiminlab/quick-workflow/internal/github"
	"github.com/yiminlab/quick-workflow/internal/jira"
)

// UploadResult contains the uploaded file URLs
type UploadResult struct {
	LocalPath  string
	Filename   string
	GithubURL  string
	JiraURL    string
	IsImage    bool
	IsVideo    bool
}

// UploadFiles uploads files to GitHub and Jira, returns updated markdown and results
func UploadFiles(files []string, ghClient *github.Client, jiraClient *jira.Client, prNumber int, owner, repo, jiraTicket string) ([]UploadResult, error) {
	results := make([]UploadResult, 0, len(files))

	for _, filePath := range files {
		result, err := uploadSingleFile(filePath, ghClient, jiraClient, prNumber, owner, repo, jiraTicket)
		if err != nil {
			return nil, fmt.Errorf("failed to upload %s: %w", filepath.Base(filePath), err)
		}
		results = append(results, *result)
	}

	return results, nil
}

func uploadSingleFile(filePath string, ghClient *github.Client, jiraClient *jira.Client, prNumber int, owner, repo, jiraTicket string) (*UploadResult, error) {
	// Read file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	filename := filepath.Base(filePath)
	ext := strings.ToLower(filepath.Ext(filename))

	result := &UploadResult{
		LocalPath: filePath,
		Filename:  filename,
		IsImage:   isImageFile(ext),
		IsVideo:   isVideoFile(ext),
	}

	// Upload to GitHub (create an issue comment with the file as attachment)
	// GitHub doesn't have a direct file upload API, so we'll use the issue comment API
	// which accepts image URLs. We'll need to upload to GitHub's asset CDN first.
	if ghClient != nil {
		githubURL, err := uploadToGitHub(ghClient, owner, repo, prNumber, filename, data)
		if err != nil {
			return nil, fmt.Errorf("failed to upload to GitHub: %w", err)
		}
		result.GithubURL = githubURL
	}

	// Upload to Jira
	if jiraClient != nil && jiraTicket != "" {
		jiraURL, err := uploadToJira(jiraClient, jiraTicket, filename, data)
		if err != nil {
			return nil, fmt.Errorf("failed to upload to Jira: %w", err)
		}
		result.JiraURL = jiraURL
	}

	return result, nil
}

// uploadToGitHub uploads a file as a base64-encoded data URL
// GitHub PR comments support inline images via data URLs or external URLs
func uploadToGitHub(client *github.Client, owner, repo string, prNumber int, filename string, data []byte) (string, error) {
	// For images, we can use base64 data URLs in markdown
	// For videos, we need to upload them somewhere and link them
	ext := strings.ToLower(filepath.Ext(filename))
	
	if isImageFile(ext) {
		// Create a data URL
		mimeType := getMimeType(ext)
		encoded := base64.StdEncoding.EncodeToString(data)
		dataURL := fmt.Sprintf("data:%s;base64,%s", mimeType, encoded)
		return dataURL, nil
	}

	// For videos, we'll return a note that it needs to be uploaded separately
	// In a real implementation, you might want to upload to a cloud storage service
	return "", fmt.Errorf("video upload not yet implemented - please upload manually")
}

// uploadToJira uploads a file as an attachment to a Jira issue
func uploadToJira(client *jira.Client, issueKey, filename string, data []byte) (string, error) {
	// Create a reader from the data
	reader := bytes.NewReader(data)

	// Upload attachment
	attachment, err := client.AddAttachment(issueKey, filename, reader)
	if err != nil {
		return "", err
	}

	return attachment.Content, nil
}

// ReplaceLocalPathsWithURLs replaces local file paths in markdown with actual URLs
func ReplaceLocalPathsWithURLs(content string, results []UploadResult) string {
	for _, result := range results {
		// Replace image references
		if result.IsImage && result.GithubURL != "" {
			// Replace markdown image syntax: ![...](./ filename)
			localRef := fmt.Sprintf("(./%s)", result.Filename)
			content = strings.ReplaceAll(content, localRef, fmt.Sprintf("(%s)", result.GithubURL))
			
			// Also replace without ./
			localRef2 := fmt.Sprintf("(%s)", result.Filename)
			content = strings.ReplaceAll(content, localRef2, fmt.Sprintf("(%s)", result.GithubURL))
		}
	}
	return content
}

func isImageFile(ext string) bool {
	imageExts := map[string]bool{
		".png":  true,
		".jpg":  true,
		".jpeg": true,
		".gif":  true,
		".webp": true,
		".svg":  true,
	}
	return imageExts[ext]
}

func isVideoFile(ext string) bool {
	videoExts := map[string]bool{
		".mp4":  true,
		".mov":  true,
		".webm": true,
		".avi":  true,
	}
	return videoExts[ext]
}

func getMimeType(ext string) string {
	mimeTypes := map[string]string{
		".png":  "image/png",
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".gif":  "image/gif",
		".webp": "image/webp",
		".svg":  "image/svg+xml",
	}
	if mime, ok := mimeTypes[ext]; ok {
		return mime
	}
	return "application/octet-stream"
}

