package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yiminlab/quick-workflow/internal/jira"
	"github.com/yiminlab/quick-workflow/internal/ui"
)

var (
	readAuto bool
)

var jiraReadCmd = &cobra.Command{
	Use:   "read <ISSUE-KEY>",
	Short: "Intelligently read a Jira issue (recommended for Cursor)",
	Long: `Intelligently read a Jira issue - optimized for Cursor AI.

This command automatically decides the best way to present the issue:
- If the issue has images/attachments: exports to files for Cursor to read
- If the issue has only text: displays directly in terminal

This is the recommended command for use with Cursor AI.

Examples:
  qkflow jira read NA-9245           # Auto mode (smart decision)
  qkflow jira read NA-9245 --auto    # Explicit auto mode

Usage in Cursor:
  Simply tell Cursor: "通过 qkflow 工具读取 NA-9245 所有内容并总结"
  Or: "qkflow 读 NA-9245"`,
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

		// Check if issue has images
		fmt.Printf("🔍 Analyzing %s...\n", issueKey)
		hasImages, err := exporter.HasImages(issueKey)
		if err != nil {
			ui.Error(fmt.Sprintf("Failed to analyze issue: %v", err))
			return
		}

		if hasImages {
			// Has images: export to files
			fmt.Printf("📸 Found attachments, exporting with images...\n\n")

			// Use export command logic
			opts := jira.ExportOptions{
				IssueKey:   issueKey,
				WithImages: true,
				OutputDir:  "",
			}

			result, err := exporter.Export(opts)
			if err != nil {
				ui.Error(fmt.Sprintf("Failed to export: %v", err))
				return
			}

			// Display results (Cursor-optimized output)
			fmt.Println()
			ui.Success(fmt.Sprintf("Exported to: %s/", result.ExportPath))
			fmt.Println()
			fmt.Printf("Main file: %s\n", result.ContentFile)
			if len(result.ImageFiles) > 0 {
				fmt.Printf("Images: %s/attachments/ (%d files)\n", result.ExportPath, len(result.ImageFiles))
			}
			fmt.Println()

			// Cursor-specific hints
			fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
			fmt.Println("💡 CURSOR: Please read the following files:")
			fmt.Printf("1. %s\n", result.ContentFile)
			if len(result.ImageFiles) > 0 {
				fmt.Printf("2. All images in %s/attachments/\n", result.ExportPath)
			}
			fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
			fmt.Println()

		} else {
			// No images: display in terminal
			fmt.Printf("ℹ️  No attachments found, showing content directly...\n\n")

			// Get detailed issue
			issue, err := client.GetIssueDetailed(issueKey)
			if err != nil {
				ui.Error(fmt.Sprintf("Failed to get issue: %v", err))
				return
			}

			// Create formatter and display
			formatter := jira.NewFormatter(client)
			output := formatter.FormatIssueFull(issue)
			fmt.Println(output)
		}
	},
}

func init() {
	jiraReadCmd.Flags().BoolVarP(&readAuto, "auto", "a", true, "Automatically decide between terminal output and file export (default: true)")
}

