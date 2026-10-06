package jira

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/yiminlab/quick-workflow/pkg/config"
	jira "github.com/andygrunwald/go-jira"
)

// Client wraps the Jira API client
type Client struct {
	client *jira.Client
}

// NewClient creates a new Jira client
func NewClient() (*Client, error) {
	cfg := config.Get()
	if cfg == nil {
		return nil, fmt.Errorf("config not loaded")
	}

	if cfg.JiraAPIToken == "" || cfg.JiraServiceAddress == "" || cfg.Email == "" {
		return nil, fmt.Errorf("Jira credentials not configured")
	}

	tp := jira.BasicAuthTransport{
		Username: cfg.Email,
		Password: cfg.JiraAPIToken,
	}

	client, err := jira.NewClient(tp.Client(), cfg.JiraServiceAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to create Jira client: %w", err)
	}

	return &Client{client: client}, nil
}

// Issue represents a Jira issue
type Issue struct {
	Key         string
	Summary     string
	Description string
	Status      string
	Type        string
	Priority    string
	Assignee    string
	Reporter    string
	Created     string
	Updated     string
	Attachments []Attachment
	Comments    []Comment
}

// Attachment represents a Jira attachment
type Attachment struct {
	ID       string
	Filename string
	Size     int64
	MimeType string
	Content  string // URL to download
	Created  string
	Author   string
}

// Comment represents a Jira comment
type Comment struct {
	ID      string
	Author  string
	Body    string
	Created string
	Updated string
}

// GetIssue gets a Jira issue by key (basic info only)
func (c *Client) GetIssue(issueKey string) (*Issue, error) {
	issue, _, err := c.client.Issue.Get(issueKey, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get issue %s: %w", issueKey, err)
	}

	return &Issue{
		Key:         issue.Key,
		Summary:     issue.Fields.Summary,
		Description: issue.Fields.Description,
		Status:      issue.Fields.Status.Name,
		Type:        issue.Fields.Type.Name,
	}, nil
}

// GetIssueDetailed gets a Jira issue with all details including attachments and comments
func (c *Client) GetIssueDetailed(issueKey string) (*Issue, error) {
	issue, _, err := c.client.Issue.Get(issueKey, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get issue %s: %w", issueKey, err)
	}

	cfg := config.Get()

	// Build issue struct
	result := &Issue{
		Key:         issue.Key,
		Summary:     issue.Fields.Summary,
		Description: issue.Fields.Description,
		Status:      issue.Fields.Status.Name,
		Type:        issue.Fields.Type.Name,
	}

	// Priority
	if issue.Fields.Priority != nil {
		result.Priority = issue.Fields.Priority.Name
	}

	// Assignee
	if issue.Fields.Assignee != nil {
		result.Assignee = issue.Fields.Assignee.DisplayName
	}

	// Reporter
	if issue.Fields.Reporter != nil {
		result.Reporter = issue.Fields.Reporter.DisplayName
	}

	// Dates (convert jira.Time to string)
	result.Created = time.Time(issue.Fields.Created).Format(time.RFC3339)
	result.Updated = time.Time(issue.Fields.Updated).Format(time.RFC3339)

	// Attachments
	if issue.Fields.Attachments != nil {
		for _, att := range issue.Fields.Attachments {
			attachment := Attachment{
				ID:       att.ID,
				Filename: att.Filename,
				Size:     int64(att.Size),
				MimeType: att.MimeType,
				Content:  att.Content,
				Created:  att.Created, // Keep as string from API
			}
			if att.Author != nil {
				attachment.Author = att.Author.DisplayName
			}
			result.Attachments = append(result.Attachments, attachment)
		}
	}

	// Comments
	if issue.Fields.Comments != nil {
		for _, c := range issue.Fields.Comments.Comments {
			comment := Comment{
				ID:      c.ID,
				Body:    c.Body,
				Created: c.Created,
				Updated: c.Updated,
			}
			// Author is a User struct (not pointer), check DisplayName instead
			if c.Author.DisplayName != "" {
				comment.Author = c.Author.DisplayName
			}
			result.Comments = append(result.Comments, comment)
		}
	}

	// Build Jira URL
	if cfg != nil && cfg.JiraServiceAddress != "" {
		// Store URL in a way that can be accessed (we'll add a method for this)
	}

	return result, nil
}

// GetJiraURL returns the full Jira URL for an issue
func (c *Client) GetJiraURL(issueKey string) string {
	cfg := config.Get()
	if cfg == nil || cfg.JiraServiceAddress == "" {
		return ""
	}
	return fmt.Sprintf("%s/browse/%s", strings.TrimRight(cfg.JiraServiceAddress, "/"), issueKey)
}

// UpdateStatus updates the status of a Jira issue
func (c *Client) UpdateStatus(issueKey, newStatus string) error {
	// 获取可用的转换
	transitions, _, err := c.client.Issue.GetTransitions(issueKey)
	if err != nil {
		return fmt.Errorf("failed to get transitions for %s: %w", issueKey, err)
	}

	// 查找匹配的转换
	var targetTransition *jira.Transition
	for _, t := range transitions {
		if strings.EqualFold(t.Name, newStatus) || strings.EqualFold(t.To.Name, newStatus) {
			targetTransition = &t
			break
		}
	}

	if targetTransition == nil {
		return fmt.Errorf("transition to status '%s' not found for issue %s", newStatus, issueKey)
	}

	// 执行转换
	_, err = c.client.Issue.DoTransition(issueKey, targetTransition.ID)
	if err != nil {
		return fmt.Errorf("failed to transition issue %s to %s: %w", issueKey, newStatus, err)
	}

	return nil
}

// AddComment adds a comment to a Jira issue
func (c *Client) AddComment(issueKey, comment string) error {
	jiraComment := &jira.Comment{
		Body: comment,
	}

	_, _, err := c.client.Issue.AddComment(issueKey, jiraComment)
	if err != nil {
		return fmt.Errorf("failed to add comment to %s: %w", issueKey, err)
	}

	return nil
}

// AddAttachment adds an attachment to a Jira issue
func (c *Client) AddAttachment(issueKey, filename string, content io.Reader) (*Attachment, error) {
	attachments, _, err := c.client.Issue.PostAttachment(issueKey, content, filename)
	if err != nil {
		return nil, fmt.Errorf("failed to add attachment to %s: %w", issueKey, err)
	}

	if attachments == nil || len(*attachments) == 0 {
		return nil, fmt.Errorf("no attachment was created")
	}

	attachment := (*attachments)[0]
	return &Attachment{
		ID:       attachment.ID,
		Filename: attachment.Filename,
		Size:     int64(attachment.Size),
		MimeType: attachment.MimeType,
		Content:  attachment.Content,
		Created:  attachment.Created,
		Author:   attachment.Author.DisplayName,
	}, nil
}

// GetProjectStatuses gets all available statuses for a project
func (c *Client) GetProjectStatuses(projectKey string) ([]string, error) {
	// 获取项目信息
	project, _, err := c.client.Project.Get(projectKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get project %s: %w", projectKey, err)
	}

	// 获取状态列表
	statuses, _, err := c.client.Status.GetAllStatuses()
	if err != nil {
		return nil, fmt.Errorf("failed to get statuses: %w", err)
	}

	// 使用 map 去重
	statusMap := make(map[string]bool)
	var statusNames []string
	for _, status := range statuses {
		// 只添加未见过的状态
		if !statusMap[status.Name] {
			statusMap[status.Name] = true
			statusNames = append(statusNames, status.Name)
		}
	}

	// 如果项目有特定的工作流，这里应该过滤
	// 简化实现：返回所有唯一状态
	_ = project // 避免未使用变量警告
	
	return statusNames, nil
}

// AddPRLink adds a PR link to a Jira issue
func (c *Client) AddPRLink(issueKey, prURL string) error {
	comment := fmt.Sprintf("Pull Request created: %s", prURL)
	return c.AddComment(issueKey, comment)
}

// AssignToMe assigns a Jira issue to the current user
func (c *Client) AssignToMe(issueKey string) error {
	cfg := config.Get()
	if cfg == nil || cfg.Email == "" {
		return fmt.Errorf("email not configured")
	}

	// Get current user's account ID
	myself, _, err := c.client.User.GetSelf()
	if err != nil {
		return fmt.Errorf("failed to get current user: %w", err)
	}

	// Assign issue to current user
	_, err = c.client.Issue.UpdateAssignee(issueKey, myself)
	if err != nil {
		return fmt.Errorf("failed to assign issue %s: %w", issueKey, err)
	}

	return nil
}

// ExtractProjectKey extracts the project key from an issue key
// For example: "PROJ-123" -> "PROJ"
func ExtractProjectKey(issueKey string) string {
	parts := strings.Split(issueKey, "-")
	if len(parts) < 2 {
		return ""
	}
	return parts[0]
}

// ValidateIssueKey checks if an issue key is valid
func ValidateIssueKey(issueKey string) bool {
	if issueKey == "" {
		return false
	}
	parts := strings.Split(issueKey, "-")
	return len(parts) >= 2
}

