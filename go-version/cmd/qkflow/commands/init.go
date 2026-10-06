package commands

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"github.com/yiminlab/quick-workflow/internal/ui"
	"github.com/yiminlab/quick-workflow/internal/utils"
	"github.com/yiminlab/quick-workflow/pkg/config"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize configuration for quick-workflow",
	Long: `Initialize configuration by prompting for required credentials and settings.
This will create a configuration file at ~/.config/quick-workflow/config.yaml`,
	Run: runInit,
}

func runInit(cmd *cobra.Command, args []string) {
	ui.Info("Welcome to Quick Workflow Setup!")
	fmt.Println()

	cfg := &config.Config{}

	// Email
	email, err := ui.PromptInput("Enter your email address:", true)
	if err != nil {
		ui.Error(fmt.Sprintf("Failed to get email: %v", err))
		return
	}
	cfg.Email = email

	// GitHub Token
	ui.Info("Getting GitHub token from gh CLI...")
	ghToken, err := getGitHubToken()
	if err != nil {
		ui.Warning("Failed to get GitHub token from gh CLI")
		ghToken, err = ui.PromptPassword("Enter your GitHub personal access token:")
		if err != nil {
			ui.Error(fmt.Sprintf("Failed to get GitHub token: %v", err))
			return
		}
	} else {
		ui.Success("GitHub token obtained from gh CLI")
	}
	cfg.GitHubToken = ghToken

	// GitHub Owner
	githubOwner, err := ui.PromptInput("Enter your GitHub username or organization:", true)
	if err != nil {
		ui.Error(fmt.Sprintf("Failed to get GitHub owner: %v", err))
		return
	}
	cfg.GitHubOwner = githubOwner

	// GitHub Repo
	githubRepo, err := ui.PromptInput("Enter your GitHub repository name:", true)
	if err != nil {
		ui.Error(fmt.Sprintf("Failed to get GitHub repo: %v", err))
		return
	}
	cfg.GitHubRepo = githubRepo

	// Jira Service Address
	jiraAddr, err := ui.PromptInput("Enter your Jira service address (e.g., https://your-domain.atlassian.net):", true)
	if err != nil {
		ui.Error(fmt.Sprintf("Failed to get Jira address: %v", err))
		return
	}
	cfg.JiraServiceAddress = strings.TrimRight(jiraAddr, "/")

	// Jira API Token
	ui.Info("Get your Jira API token from: https://id.atlassian.com/manage-profile/security/api-tokens")
	jiraToken, err := ui.PromptPassword("Enter your Jira API token:")
	if err != nil {
		ui.Error(fmt.Sprintf("Failed to get Jira token: %v", err))
		return
	}
	cfg.JiraAPIToken = jiraToken

	// Branch Prefix (optional)
	branchPrefix, err := ui.PromptInput("Enter branch prefix (optional, e.g., 'feature' or 'username'):", false)
	if err != nil {
		ui.Warning("Skipping branch prefix")
	} else {
		cfg.BranchPrefix = strings.TrimRight(branchPrefix, "/")
	}

	// OpenAI Key (optional)
	useAI, err := ui.PromptConfirm("Do you want to configure AI features (OpenAI/DeepSeek)?", false)
	if err == nil && useAI {
		aiKey, err := ui.PromptPassword("Enter OpenAI or DeepSeek API key:")
		if err == nil {
			cfg.OpenAIKey = aiKey
		}
	}

	// Auto Update (default: true)
	autoUpdate, err := ui.PromptConfirm("Enable automatic updates? (recommended)", true)
	if err == nil {
		cfg.AutoUpdate = autoUpdate
	} else {
		cfg.AutoUpdate = true // Default to true
	}
	if cfg.AutoUpdate {
		ui.Success("✅ Auto-update enabled - qkflow will keep itself up to date")
	} else {
		ui.Info("ℹ️  Auto-update disabled - run 'qkflow update-cli' to update manually")
	}

	// Save configuration
	if err := config.Save(cfg); err != nil {
		ui.Error(fmt.Sprintf("Failed to save configuration: %v", err))
		return
	}

	ui.Success("Configuration saved successfully!")
	fmt.Println()
	
	// Show storage location
	location := utils.GetConfigLocation()
	configDir, _ := utils.GetQuickWorkflowConfigDir()
	ui.Info(fmt.Sprintf("Storage location: %s", location))
	if configDir != "" {
		fmt.Printf("  📁 Config: %s/config.yaml\n", configDir)
	}
	fmt.Println()
	
	ui.Info("You can now use the following commands:")
	fmt.Println("  qkflow pr create   - Create a PR and update Jira")
	fmt.Println("  qkflow pr merge    - Merge a PR and update Jira")
	fmt.Println("  qkflow update      - Quick commit and push with PR title")
	fmt.Println("  qkflow update-cli  - Update qkflow to the latest version")
	fmt.Println("  qkflow jira list   - List Jira status mappings")
	fmt.Println()
}

func getGitHubToken() (string, error) {
	cmd := exec.Command("gh", "auth", "token")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

