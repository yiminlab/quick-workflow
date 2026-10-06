package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yiminlab/quick-workflow/internal/git"
	"github.com/yiminlab/quick-workflow/internal/github"
	"github.com/yiminlab/quick-workflow/internal/ui"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Quick update - commit and push with PR title as commit message",
	Long: `Automatically get the current PR title and use it as the commit message.
This command will:
1. Get the PR title for the current branch
2. Stage all changes (git add --all)
3. Commit with the PR title as the message
4. Push to origin

If no PR is found, it will use "update" as the default commit message.`,
	Run: runUpdate,
}

func runUpdate(cmd *cobra.Command, args []string) {
	// 检查是否是 git 仓库
	if !git.IsGitRepository() {
		ui.Error("Not a git repository")
		return
	}

	// 检查是否有未提交的更改
	hasChanges, err := git.HasUncommittedChanges()
	if err != nil {
		ui.Error(fmt.Sprintf("Failed to check git status: %v", err))
		return
	}

	if !hasChanges {
		ui.Warning("No changes to commit")
		return
	}

	// 获取当前分支
	branch, err := git.GetCurrentBranch()
	if err != nil {
		ui.Error(fmt.Sprintf("Failed to get current branch: %v", err))
		return
	}

	ui.Info(fmt.Sprintf("Current branch: %s", branch))

	// 获取 PR 标题作为 commit message
	commitMessage := "update" // 默认 commit message
	
	// 尝试从 GitHub 获取 PR 标题
	remoteURL, err := git.GetRemoteURL()
	if err == nil {
		owner, repo, err := github.ParseRepositoryFromURL(remoteURL)
		if err == nil {
			// 创建 GitHub 客户端
			ghClient, err := github.NewClient()
			if err == nil {
				// 尝试获取当前分支的 PR
				pr, err := ghClient.GetPRByBranch(owner, repo, branch)
				if err == nil && pr != nil {
					commitMessage = pr.Title
					ui.Success(fmt.Sprintf("Got PR title: %s", commitMessage))
				} else {
					ui.Warning(fmt.Sprintf("No open PR found for branch %s, using default message 'update'", branch))
				}
			} else {
				ui.Warning(fmt.Sprintf("Failed to create GitHub client: %v, using default message", err))
			}
		}
	}

	// Stage all changes
	ui.Info("Staging all changes...")
	if err := git.AddAll(); err != nil {
		ui.Error(fmt.Sprintf("Failed to stage changes: %v", err))
		return
	}

	// Commit
	ui.Info(fmt.Sprintf("Committing with message: '%s'", commitMessage))
	if err := git.Commit(commitMessage); err != nil {
		ui.Error(fmt.Sprintf("Failed to commit: %v", err))
		return
	}

	// Push
	ui.Info("Pushing to origin...")
	if err := git.Push(branch); err != nil {
		ui.Error(fmt.Sprintf("Failed to push: %v", err))
		return
	}

	ui.Success("✅ Successfully committed and pushed changes!")
}

