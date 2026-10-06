package watcher

import (
	"regexp"
	"strings"

	"github.com/yiminlab/quick-workflow/internal/github"
)

// MergedPR represents a merged pull request with Jira tickets
type MergedPR struct {
	Number      int
	Title       string
	URL         string
	Branch      string
	MergedAt    string
	MergedBy    string
	JiraTickets []string
}

// Checker handles PR checking logic
type Checker struct {
	client *github.Client
	logger *Logger
}

// NewChecker creates a new Checker instance
func NewChecker(client *github.Client, logger *Logger) *Checker {
	return &Checker{
		client: client,
		logger: logger,
	}
}

// CheckMergedPRs checks for newly merged PRs from the watching list
func (c *Checker) CheckMergedPRs(watchingList *WatchingList, state *State) ([]MergedPR, error) {
	watchingPRs := watchingList.GetAll()
	c.logger.Infof("Checking %d watching PRs", len(watchingPRs))

	if len(watchingPRs) == 0 {
		c.logger.Info("No PRs in watching list")
		return []MergedPR{}, nil
	}

	// Check each watching PR
	mergedPRs := make([]MergedPR, 0)
	
	for _, watchingPR := range watchingPRs {
		// Skip if already processed
		if state.IsPRProcessed(watchingPR.PRNumber) {
			c.logger.Infof("PR #%d already processed, skipping", watchingPR.PRNumber)
			continue
		}

		// Get PR details from GitHub
		pr, err := c.client.GetPullRequest(watchingPR.Owner, watchingPR.Repo, watchingPR.PRNumber)
		if err != nil {
			c.logger.Warningf("Failed to get PR #%d from %s/%s: %v", watchingPR.PRNumber, watchingPR.Owner, watchingPR.Repo, err)
			continue
		}

		// Check if merged
		if pr.MergedAt == "" {
			c.logger.Infof("PR #%d not merged yet", watchingPR.PRNumber)
			continue
		}

		// Use Jira tickets from watching list (already extracted when PR was created)
		jiraTickets := watchingPR.JiraTickets
		if len(jiraTickets) == 0 {
			c.logger.Warningf("PR #%d: No Jira tickets found", watchingPR.PRNumber)
			continue
		}

		mergedPR := MergedPR{
			Number:      watchingPR.PRNumber,
			Title:       pr.Title,
			URL:         pr.HTMLURL,
			Branch:      pr.Head,
			MergedAt:    pr.MergedAt,
			MergedBy:    pr.MergedBy,
			JiraTickets: jiraTickets,
		}

		mergedPRs = append(mergedPRs, mergedPR)
		c.logger.Infof("✅ Found newly merged PR #%d: %s (Jira: %v)", pr.Number, pr.Title, jiraTickets)
	}

	if len(mergedPRs) == 0 {
		c.logger.Info("No newly merged PRs from watching list")
	}

	return mergedPRs, nil
}

// extractJiraTickets extracts Jira ticket IDs from branch name and PR title
func (c *Checker) extractJiraTickets(branch, title string) []string {
	tickets := make(map[string]bool) // Use map to deduplicate

	// Pattern: PROJECT-123 format
	pattern := regexp.MustCompile(`([A-Z][A-Z0-9]+-\d+)`)

	// Extract from branch name
	branchMatches := pattern.FindAllString(branch, -1)
	for _, match := range branchMatches {
		tickets[match] = true
	}

	// Extract from title
	titleMatches := pattern.FindAllString(title, -1)
	for _, match := range titleMatches {
		tickets[match] = true
	}

	// Convert map to slice
	result := make([]string, 0, len(tickets))
	for ticket := range tickets {
		result = append(result, ticket)
	}

	return result
}

// GetProjectFromTicket extracts project key from ticket ID
func GetProjectFromTicket(ticket string) string {
	parts := strings.Split(ticket, "-")
	if len(parts) > 0 {
		return parts[0]
	}
	return ""
}

