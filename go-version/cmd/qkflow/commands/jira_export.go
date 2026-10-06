package commands

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/yiminlab/quick-workflow/internal/jira"
	"github.com/yiminlab/quick-workflow/internal/ui"
)

var (
	exportWithImages bool
	exportOutputDir  string
)

var jiraExportCmd = &cobra.Command{
	Use:   "export <ISSUE-KEY>",
	Short: "Export a Jira issue to files",
	Long: `Export a Jira issue to local files in Markdown format.

By default, exports to /tmp/qkflow/jira/<ISSUE-KEY>/
Use --with-images to download all attachments and images.

The export creates:
  - content.md: Main content in Markdown format
  - README.md: How to use in Cursor
  - attachments/: Downloaded files (if --with-images)

Examples:
  qkflow jira export NA-9245                    # Export text only
  qkflow jira export NA-9245 --with-images      # Export with all images
  qkflow jira export NA-9245 -o ~/exports/      # Custom output directory`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		issueKey := args[0]

		// Validate issue key
		if !jira.ValidateIssueKey(issueKey) {
			ui.Error(fmt.Sprintf("Invalid issue key: %s", issueKey))
			return
		}

		// Create Jira client
		client, err := jira.NewClient()
		if err != nil {
			ui.Error(fmt.Sprintf("Failed to create Jira client: %v", err))
			return
		}

		// Create exporter
		exporter := jira.NewExporter(client)

		// Prepare options
		opts := jira.ExportOptions{
			IssueKey:   issueKey,
			WithImages: exportWithImages,
			OutputDir:  exportOutputDir,
		}

		// Show progress
		fmt.Printf("🔍 Fetching %s from Jira...\n", issueKey)

		// Export
		result, err := exporter.Export(opts)
		if err != nil {
			ui.Error(fmt.Sprintf("Failed to export: %v", err))
			return
		}

		// Display results
		fmt.Println()
		ui.Success(fmt.Sprintf("Jira issue %s exported successfully!", issueKey))
		fmt.Println()
		fmt.Printf("📁 Location: %s/\n", result.ExportPath)
		fmt.Printf("📄 Main content: %s\n", result.ContentFile)

		if exportWithImages && len(result.ImageFiles) > 0 {
			fmt.Printf("🖼️  Images: %d files downloaded\n", len(result.ImageFiles))
			fmt.Println()
			fmt.Println("Downloaded files:")
			for _, imgPath := range result.ImageFiles {
				fmt.Printf("  - %s\n", filepath.Base(imgPath))
			}
		}

		// Display Cursor instructions
		fmt.Println()
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Println("💡 How to use in Cursor:")
		fmt.Println()
		fmt.Println("1. In Cursor, attach the exported content:")
		fmt.Printf("   @%s\n", result.ContentFile)
		fmt.Println()

		if exportWithImages && len(result.ImageFiles) > 0 {
			fmt.Println("2. To include images, tell Cursor:")
			fmt.Printf("   \"Please read the images in %s/attachments/\"\n", result.ExportPath)
			fmt.Println()
			fmt.Println("3. Or simply tell Cursor:")
			fmt.Printf("   \"Read the exported Jira ticket at %s/\"\n", result.ExportPath)
			fmt.Println()
		} else {
			fmt.Println("2. Or simply tell Cursor:")
			fmt.Printf("   \"Read %s\"\n", result.ContentFile)
			fmt.Println()
		}

		fmt.Println("4. When done, clean up:")
		fmt.Printf("   qkflow jira clean %s\n", issueKey)
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Println()

		// Show hint if images available but not downloaded
		if !exportWithImages {
			hasImages, _ := exporter.HasImages(issueKey)
			if hasImages {
				fmt.Println("💡 Tip: This issue has attachments. Use --with-images to download them.")
			}
		}
	},
}

func init() {
	jiraExportCmd.Flags().BoolVarP(&exportWithImages, "with-images", "i", false, "Download all attachments and images")
	jiraExportCmd.Flags().StringVarP(&exportOutputDir, "output-dir", "o", "", "Output directory (default: /tmp/qkflow/jira/)")
}

