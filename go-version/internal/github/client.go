package github

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/go-github/v57/github"
	"github.com/yiminlab/quick-workflow/pkg/config"
	"golang.org/x/oauth2"
)

// Client wraps the GitHub API client
type Client struct {
	client *github.Client
	ctx    context.Context
}

// NewClient creates a new GitHub client
func NewClient() (*Client, error) {
	cfg := config.Get()
	if cfg == nil || cfg.GitHubToken == "" {
		return nil, fmt.Errorf("GitHub token not configured")
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: cfg.GitHubToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	return &Client{
		client: github.NewClient(tc),
		ctx:    ctx,
	}, nil
}

// PullRequest represents a pull request
type PullRequest struct {
	Number   int
	Title    string
	Body     string
	HTMLURL  string
	Head     string
	Base     string
	State    string
	MergedAt string
	MergedBy string
}

// CreatePullRequestInput contains the input for creating a PR
type CreatePullRequestInput struct {
	Owner string
	Repo  string
	Title string
	Body  string
	Head  string // branch name
	Base  string // target branch, usually "main" or "master"
}

// CreatePullRequest creates a new pull request
func (c *Client) CreatePullRequest(input CreatePullRequestInput) (*PullRequest, error) {
	newPR := &github.NewPullRequest{
		Title: github.String(input.Title),
		Body:  github.String(input.Body),
		Head:  github.String(input.Head),
		Base:  github.String(input.Base),
	}

	pr, _, err := c.client.PullRequests.Create(c.ctx, input.Owner, input.Repo, newPR)
	if err != nil {
		return nil, fmt.Errorf("failed to create pull request: %w", err)
	}

	return &PullRequest{
		Number:  pr.GetNumber(),
		Title:   pr.GetTitle(),
		Body:    pr.GetBody(),
		HTMLURL: pr.GetHTMLURL(),
		Head:    pr.GetHead().GetRef(),
		Base:    pr.GetBase().GetRef(),
		State:   pr.GetState(),
	}, nil
}

// ListPullRequests lists pull requests with the given state and optional author filter
func (c *Client) ListPullRequests(owner, repo, state, author string) ([]PullRequest, error) {
	opts := &github.PullRequestListOptions{
		State: state,
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	prs, _, err := c.client.PullRequests.List(c.ctx, owner, repo, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to list pull requests: %w", err)
	}

	result := make([]PullRequest, 0)
	for _, pr := range prs {
		// Filter by author if specified
		if author != "" && pr.User != nil && pr.User.GetLogin() != author {
			continue
		}

		mergedAt := ""
		if pr.MergedAt != nil {
			mergedAt = pr.MergedAt.Format("2006-01-02T15:04:05Z")
		}

		mergedBy := ""
		if pr.MergedBy != nil {
			mergedBy = pr.MergedBy.GetLogin()
		}

		result = append(result, PullRequest{
			Number:   pr.GetNumber(),
			Title:    pr.GetTitle(),
			Body:     pr.GetBody(),
			HTMLURL:  pr.GetHTMLURL(),
			Head:     pr.GetHead().GetRef(),
			Base:     pr.GetBase().GetRef(),
			State:    pr.GetState(),
			MergedAt: mergedAt,
			MergedBy: mergedBy,
		})
	}

	return result, nil
}

// GetPullRequest gets a pull request by number
func (c *Client) GetPullRequest(owner, repo string, number int) (*PullRequest, error) {
	pr, _, err := c.client.PullRequests.Get(c.ctx, owner, repo, number)
	if err != nil {
		return nil, fmt.Errorf("failed to get pull request: %w", err)
	}

	return &PullRequest{
		Number:  pr.GetNumber(),
		Title:   pr.GetTitle(),
		Body:    pr.GetBody(),
		HTMLURL: pr.GetHTMLURL(),
		Head:    pr.GetHead().GetRef(),
		Base:    pr.GetBase().GetRef(),
		State:   pr.GetState(),
	}, nil
}

// MergePullRequest merges a pull request using squash merge by default
func (c *Client) MergePullRequest(owner, repo string, number int, commitMessage string) error {
	options := &github.PullRequestOptions{
		CommitTitle: commitMessage,
		MergeMethod: "squash", // Use squash merge by default
	}

	_, _, err := c.client.PullRequests.Merge(c.ctx, owner, repo, number, commitMessage, options)
	if err != nil {
		return fmt.Errorf("failed to merge pull request: %w", err)
	}

	return nil
}

// GetCurrentRepository gets the owner and repo from git remote
func GetCurrentRepository() (owner, repo string, err error) {
	// 这里可以通过执行 git remote get-url origin 来获取
	// 简化实现：从环境变量或配置读取
	// 实际应该解析 git remote 输出

	// TODO: 实现通过 git 命令获取
	return "", "", fmt.Errorf("not implemented: should parse git remote")
}

// ParseRepositoryFromURL parses owner and repo from GitHub URL
func ParseRepositoryFromURL(url string) (owner, repo string, err error) {
	// 支持多种格式：
	// https://github.com/owner/repo
	// git@github.com:owner/repo.git
	// owner/repo

	url = strings.TrimSpace(url)

	// 移除协议部分
	url = strings.TrimPrefix(url, "https://")
	url = strings.TrimPrefix(url, "http://")
	url = strings.TrimPrefix(url, "git@")

	// 移除 github.com
	url = strings.TrimPrefix(url, "github.com/")
	url = strings.TrimPrefix(url, "github.com:")

	// 移除 .git 后缀
	url = strings.TrimSuffix(url, ".git")

	// 分割 owner/repo
	parts := strings.Split(url, "/")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid repository format: %s", url)
	}

	return parts[0], parts[1], nil
}

// ParsePRFromURL parses owner, repo and PR number from GitHub PR URL
// Supports formats like:
// - https://github.com/owner/repo/pull/123
// - https://github.com/owner/repo/pull/123/files
// - https://github.com/owner/repo/pull/123/commits
// - https://github.com/owner/repo/pull/123/checks
// - http://github.com/owner/repo/pull/123
// - github.com/owner/repo/pull/123
func ParsePRFromURL(url string) (owner, repo string, prNumber int, err error) {
	url = strings.TrimSpace(url)

	// 移除协议部分
	url = strings.TrimPrefix(url, "https://")
	url = strings.TrimPrefix(url, "http://")

	// 移除 github.com
	url = strings.TrimPrefix(url, "github.com/")

	// 移除可能的 query params 和 fragments
	if idx := strings.IndexAny(url, "?#"); idx != -1 {
		url = url[:idx]
	}

	// 分割路径: owner/repo/pull/123 或 owner/repo/pull/123/files
	parts := strings.Split(url, "/")
	if len(parts) < 4 {
		return "", "", 0, fmt.Errorf("invalid PR URL format: expected github.com/owner/repo/pull/number")
	}

	if parts[2] != "pull" {
		return "", "", 0, fmt.Errorf("invalid PR URL format: missing 'pull' segment")
	}

	owner = parts[0]
	repo = parts[1]

	// 解析 PR 号（parts[3] 可能后面还有 /files, /commits 等）
	prStr := parts[3]

	prNumber, err = strconv.Atoi(prStr)
	if err != nil {
		return "", "", 0, fmt.Errorf("invalid PR number in URL: %s", prStr)
	}

	return owner, repo, prNumber, nil
}

// IsPRURL checks if a string looks like a GitHub PR URL
func IsPRURL(s string) bool {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "https://")
	s = strings.TrimPrefix(s, "http://")
	return strings.HasPrefix(s, "github.com/") && strings.Contains(s, "/pull/")
}

// GetPRByBranch gets a pull request by branch name
func (c *Client) GetPRByBranch(owner, repo, branch string) (*PullRequest, error) {
	// 构建查询条件：head 应该是 owner:branch
	head := fmt.Sprintf("%s:%s", owner, branch)

	opts := &github.PullRequestListOptions{
		State: "open",
		Head:  head,
		ListOptions: github.ListOptions{
			PerPage: 10,
		},
	}

	prs, _, err := c.client.PullRequests.List(c.ctx, owner, repo, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to list pull requests: %w", err)
	}

	if len(prs) == 0 {
		return nil, fmt.Errorf("no open pull request found for branch %s", branch)
	}

	// 返回第一个匹配的 PR
	pr := prs[0]
	return &PullRequest{
		Number:  pr.GetNumber(),
		Title:   pr.GetTitle(),
		Body:    pr.GetBody(),
		HTMLURL: pr.GetHTMLURL(),
		Head:    pr.GetHead().GetRef(),
		Base:    pr.GetBase().GetRef(),
		State:   pr.GetState(),
	}, nil
}

// AddPRComment adds a comment to a pull request
func (c *Client) AddPRComment(owner, repo string, prNumber int, body string) error {
	comment := &github.IssueComment{
		Body: github.String(body),
	}

	_, _, err := c.client.Issues.CreateComment(c.ctx, owner, repo, prNumber, comment)
	if err != nil {
		return fmt.Errorf("failed to add comment to PR #%d: %w", prNumber, err)
	}

	return nil
}

// ApprovePullRequest approves a pull request
func (c *Client) ApprovePullRequest(owner, repo string, number int, body string) error {
	event := "APPROVE"
	review := &github.PullRequestReviewRequest{
		Event: &event,
	}

	if body != "" {
		review.Body = &body
	}

	_, _, err := c.client.PullRequests.CreateReview(c.ctx, owner, repo, number, review)
	if err != nil {
		return fmt.Errorf("failed to approve pull request: %w", err)
	}

	return nil
}

// IsPRMergeable checks if a PR is mergeable
func (c *Client) IsPRMergeable(owner, repo string, number int) (bool, error) {
	pr, _, err := c.client.PullRequests.Get(c.ctx, owner, repo, number)
	if err != nil {
		return false, fmt.Errorf("failed to get pull request: %w", err)
	}

	// Check if PR is open
	if pr.GetState() != "open" {
		return false, fmt.Errorf("PR is not open (state: %s)", pr.GetState())
	}

	// Check if mergeable
	if pr.Mergeable != nil && !*pr.Mergeable {
		return false, fmt.Errorf("PR has conflicts and cannot be merged")
	}

	return true, nil
}
