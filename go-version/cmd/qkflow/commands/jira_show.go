package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yiminlab/quick-workflow/internal/jira"
	"github.com/yiminlab/quick-workflow/internal/ui"
)

var (
	showFull bool
)

var jiraShowCmd = &cobra.Command{
	Use:   "show <ISSUE-KEY>",
	Short: "Show a Jira issue in the terminal",
	Long: `Display a Jira issue in the terminal.

Use --full to show complete details including description and comments.
Without --full, shows only basic metadata.

Examples:
  qkflow jira show NA-9245           # Quick view
  qkflow jira show NA-9245 --full    # Full details`,
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

		// Create formatter
		formatter := jira.NewFormatter(client)

		if showFull {
			// Get detailed issue
			issue, err := client.GetIssueDetailed(issueKey)
			if err != nil {
				ui.Error(fmt.Sprintf("Failed to get issue: %v", err))
				return
			}

			// Format and display
			output := formatter.FormatIssueFull(issue)
			fmt.Println(output)
		} else {
			// Get basic issue
			issue, err := client.GetIssue(issueKey)
			if err != nil {
				ui.Error(fmt.Sprintf("Failed to get issue: %v", err))
				return
			}

			// Format and display
			output := formatter.FormatIssueSimple(issue)
			fmt.Println(output)
		}
	},
}

func init() {
	jiraShowCmd.Flags().BoolVarP(&showFull, "full", "f", false, "Show full details including description and comments")
}

