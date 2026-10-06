package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yiminlab/quick-workflow/internal/jira"
	"github.com/yiminlab/quick-workflow/internal/ui"
)

var jiraCmd = &cobra.Command{
	Use:   "jira",
	Short: "Manage Jira issues and status mappings",
	Long: `Work with Jira issues and manage status mappings.

Available commands:
  show    - Display a Jira issue in the terminal
  export  - Export a Jira issue to files (with optional images)
  read    - Intelligently read an issue (recommended for Cursor)
  clean   - Clean up exported files
  list    - List status mappings
  setup   - Setup status mappings
  delete  - Delete status mappings`,
}

var jiraListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all Jira status mappings",
	Run: func(cmd *cobra.Command, args []string) {
		statusCache, err := jira.NewStatusCache()
		if err != nil {
			ui.Error(fmt.Sprintf("Failed to create status cache: %v", err))
			return
		}

		mappings, err := statusCache.ListAllMappings()
		if err != nil {
			ui.Error(fmt.Sprintf("Failed to list mappings: %v", err))
			return
		}

		if len(mappings) == 0 {
			ui.Info("No Jira status mappings configured yet")
			return
		}

		fmt.Println("\n📋 Jira Status Mappings:")
		fmt.Println()
		for _, mapping := range mappings {
			fmt.Printf("Project: %s\n", mapping.ProjectKey)
			fmt.Printf("  PR Created → %s\n", mapping.PRCreatedStatus)
			fmt.Printf("  PR Merged  → %s\n", mapping.PRMergedStatus)
			fmt.Println()
		}
	},
}

var jiraSetupCmd = &cobra.Command{
	Use:   "setup [project-key]",
	Short: "Setup or update status mapping for a project",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectKey := args[0]

		jiraClient, err := jira.NewClient()
		if err != nil {
			ui.Error(fmt.Sprintf("Failed to create Jira client: %v", err))
			return
		}

		mapping, err := setupProjectStatusMapping(jiraClient, projectKey)
		if err != nil {
			ui.Error(fmt.Sprintf("Failed to setup mapping: %v", err))
			return
		}

		if mapping == nil {
			ui.Info("Setup cancelled")
			return
		}

		statusCache, err := jira.NewStatusCache()
		if err != nil {
			ui.Error(fmt.Sprintf("Failed to create status cache: %v", err))
			return
		}

		if err := statusCache.SaveProjectStatus(mapping); err != nil {
			ui.Error(fmt.Sprintf("Failed to save mapping: %v", err))
			return
		}

		ui.Success(fmt.Sprintf("Status mapping saved for project %s!", projectKey))
	},
}

var jiraDeleteCmd = &cobra.Command{
	Use:   "delete [project-key]",
	Short: "Delete status mapping for a project",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectKey := args[0]

		confirm, err := ui.PromptConfirm(fmt.Sprintf("Delete status mapping for %s?", projectKey), false)
		if err != nil || !confirm {
			ui.Info("Deletion cancelled")
			return
		}

		statusCache, err := jira.NewStatusCache()
		if err != nil {
			ui.Error(fmt.Sprintf("Failed to create status cache: %v", err))
			return
		}

		if err := statusCache.DeleteProjectStatus(projectKey); err != nil {
			ui.Error(fmt.Sprintf("Failed to delete mapping: %v", err))
			return
		}

		ui.Success(fmt.Sprintf("Status mapping deleted for project %s", projectKey))
	},
}

func init() {
	// Issue commands (new)
	jiraCmd.AddCommand(jiraShowCmd)
	jiraCmd.AddCommand(jiraExportCmd)
	jiraCmd.AddCommand(jiraReadCmd)
	jiraCmd.AddCommand(jiraCleanCmd)
	
	// Status mapping commands (existing)
	jiraCmd.AddCommand(jiraListCmd)
	jiraCmd.AddCommand(jiraSetupCmd)
	jiraCmd.AddCommand(jiraDeleteCmd)
}

