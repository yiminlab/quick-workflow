package jira

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/yiminlab/quick-workflow/pkg/config"
)

// Exporter handles exporting Jira issues to files
type Exporter struct {
	client    *Client
	formatter *Formatter
}

// NewExporter creates a new exporter
func NewExporter(client *Client) *Exporter {
	return &Exporter{
		client:    client,
		formatter: NewFormatter(client),
	}
}

// ExportOptions contains options for exporting
type ExportOptions struct {
	OutputDir  string
	WithImages bool
	IssueKey   string
}

// ExportResult contains the result of an export operation
type ExportResult struct {
	ExportPath  string
	ContentFile string
	ImageFiles  []string
	Success     bool
	Error       error
}

// Export exports a Jira issue
func (e *Exporter) Export(opts ExportOptions) (*ExportResult, error) {
	// Get issue with full details
	issue, err := e.client.GetIssueDetailed(opts.IssueKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get issue: %w", err)
	}

	// Determine export directory
	exportDir := opts.OutputDir
	if exportDir == "" {
		exportDir = filepath.Join("/tmp", "qkflow", "jira", opts.IssueKey)
	} else {
		exportDir = filepath.Join(exportDir, opts.IssueKey)
	}

	// Create directory
	if err := os.MkdirAll(exportDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create export directory: %w", err)
	}

	result := &ExportResult{
		ExportPath: exportDir,
		Success:    true,
	}

	// Generate Markdown content
	markdown := e.formatter.FormatMarkdown(issue, exportDir, opts.WithImages)
	contentFile := filepath.Join(exportDir, "content.md")
	if err := os.WriteFile(contentFile, []byte(markdown), 0644); err != nil {
		return nil, fmt.Errorf("failed to write content file: %w", err)
	}
	result.ContentFile = contentFile

	// Download images if requested
	if opts.WithImages && len(issue.Attachments) > 0 {
		attachmentsDir := filepath.Join(exportDir, "attachments")
		if err := os.MkdirAll(attachmentsDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create attachments directory: %w", err)
		}

		for _, att := range issue.Attachments {
			if isImageFile(att.Filename) || isDocumentFile(att.Filename) {
				localPath := filepath.Join(attachmentsDir, att.Filename)
				if err := e.downloadAttachment(att, localPath); err != nil {
					// Log error but continue with other attachments
					fmt.Fprintf(os.Stderr, "Warning: failed to download %s: %v\n", att.Filename, err)
					continue
				}
				result.ImageFiles = append(result.ImageFiles, localPath)
			}
		}
	}

	// Create README
	readme := e.generateREADME(opts.IssueKey, exportDir, opts.WithImages)
	readmeFile := filepath.Join(exportDir, "README.md")
	if err := os.WriteFile(readmeFile, []byte(readme), 0644); err != nil {
		// Non-fatal, just log
		fmt.Fprintf(os.Stderr, "Warning: failed to write README: %v\n", err)
	}

	return result, nil
}

// downloadAttachment downloads a Jira attachment
func (e *Exporter) downloadAttachment(att Attachment, localPath string) error {
	cfg := config.Get()
	if cfg == nil {
		return fmt.Errorf("config not loaded")
	}

	// Create HTTP client with basic auth
	client := &http.Client{}
	req, err := http.NewRequest("GET", att.Content, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Add basic auth
	req.SetBasicAuth(cfg.Email, cfg.JiraAPIToken)

	// Execute request
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status: %d", resp.StatusCode)
	}

	// Create local file
	out, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("failed to create local file: %w", err)
	}
	defer out.Close()

	// Copy content
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// generateREADME generates a README file for the export
func (e *Exporter) generateREADME(issueKey, exportPath string, withImages bool) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("# How to Use This Jira Export in Cursor\n\n"))
	sb.WriteString(fmt.Sprintf("This directory contains exported content from **%s**.\n\n", issueKey))

	sb.WriteString("## 📁 Files\n\n")
	sb.WriteString("- `content.md` - Main content (start here)\n")
	sb.WriteString("- `README.md` - This file\n")
	if withImages {
		sb.WriteString("- `attachments/` - Downloaded images and files\n")
	}
	sb.WriteString("\n")

	sb.WriteString("## 💡 Usage in Cursor\n\n")
	sb.WriteString("### Quick Start\n\n")
	sb.WriteString("Simply tell Cursor:\n")
	sb.WriteString("```\n")
	sb.WriteString(fmt.Sprintf("Read the exported Jira ticket: %s/content.md\n", exportPath))
	sb.WriteString("```\n\n")

	if withImages {
		sb.WriteString("### With Images\n\n")
		sb.WriteString("To let Cursor see images:\n")
		sb.WriteString("```\n")
		sb.WriteString(fmt.Sprintf("Read %s/content.md and all images in ./attachments/\n", exportPath))
		sb.WriteString("```\n\n")

		sb.WriteString("Or attach them manually:\n")
		sb.WriteString(fmt.Sprintf("1. `@%s/content.md`\n", exportPath))
		sb.WriteString(fmt.Sprintf("2. `@%s/attachments/` (attach specific images)\n", exportPath))
		sb.WriteString("\n")
	}

	sb.WriteString("### Full Context\n\n")
	sb.WriteString("```\n")
	sb.WriteString(fmt.Sprintf("Read all content from %s/\n", exportPath))
	sb.WriteString("```\n\n")

	sb.WriteString("## 🗑️ Clean Up\n\n")
	sb.WriteString("When done:\n")
	sb.WriteString("```bash\n")
	sb.WriteString(fmt.Sprintf("qkflow jira clean %s\n", issueKey))
	sb.WriteString("```\n\n")

	sb.WriteString("## 🔄 Re-export\n\n")
	sb.WriteString("To get latest changes:\n")
	sb.WriteString("```bash\n")
	sb.WriteString(fmt.Sprintf("qkflow jira export %s --with-images\n", issueKey))
	sb.WriteString("```\n\n")

	sb.WriteString("---\n\n")
	sb.WriteString(fmt.Sprintf("Exported: %s\n", formatDate("")))

	return sb.String()
}

// HasImages checks if an issue has images/attachments
func (e *Exporter) HasImages(issueKey string) (bool, error) {
	issue, err := e.client.GetIssueDetailed(issueKey)
	if err != nil {
		return false, err
	}
	return len(issue.Attachments) > 0, nil
}

// Helper functions

func isDocumentFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	docExts := []string{".pdf", ".doc", ".docx", ".txt", ".md"}
	for _, docExt := range docExts {
		if ext == docExt {
			return true
		}
	}
	return false
}

