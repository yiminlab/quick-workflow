package watcher

import (
	"fmt"
	"time"

	"github.com/yiminlab/quick-workflow/internal/jira"
)

// Processor handles PR processing and Jira updates
type Processor struct {
	jiraClient  *jira.Client
	statusCache *jira.StatusCache
	logger      *Logger
}

// NewProcessor creates a new Processor instance
func NewProcessor(jiraClient *jira.Client, statusCache *jira.StatusCache, logger *Logger) *Processor {
	return &Processor{
		jiraClient:  jiraClient,
		statusCache: statusCache,
		logger:      logger,
	}
}

// ProcessMergedPR processes a merged PR and updates Jira
func (p *Processor) ProcessMergedPR(pr MergedPR) ProcessedPR {
	p.logger.Infof("Processing PR #%d: %s", pr.Number, pr.Title)

	jiraUpdates := make([]JiraUpdateResult, 0)

	for _, ticket := range pr.JiraTickets {
		updateResult := p.updateJiraTicket(ticket)
		jiraUpdates = append(jiraUpdates, updateResult)

		if updateResult.Success {
			p.logger.Successf("✅ PR #%d merged → Updated %s: %s → %s", 
				pr.Number, ticket, updateResult.OldStatus, updateResult.NewStatus)
		} else {
			p.logger.Errorf("❌ PR #%d: Failed to update %s: %s", 
				pr.Number, ticket, updateResult.Error)
		}
	}

	mergedAt, _ := time.Parse(time.RFC3339, pr.MergedAt)

	return ProcessedPR{
		PRNumber:    pr.Number,
		PRTitle:     pr.Title,
		PRURL:       pr.URL,
		Branch:      pr.Branch,
		JiraTickets: pr.JiraTickets,
		MergedAt:    mergedAt,
		MergedBy:    pr.MergedBy,
		JiraUpdates: jiraUpdates,
	}
}

// updateJiraTicket updates a single Jira ticket status
func (p *Processor) updateJiraTicket(ticket string) JiraUpdateResult {
	result := JiraUpdateResult{
		Ticket:  ticket,
		Success: false,
	}

	// Get project key from ticket
	projectKey := GetProjectFromTicket(ticket)
	if projectKey == "" {
		result.Error = "invalid ticket format"
		return result
	}

	// Get status mapping for the project
	mapping, err := p.statusCache.GetProjectStatus(projectKey)
	if err != nil {
		result.Error = fmt.Sprintf("failed to get status mapping: %v", err)
		return result
	}

	if mapping == nil {
		result.Error = fmt.Sprintf("no status mapping configured for project %s", projectKey)
		return result
	}

	targetStatus := mapping.PRMergedStatus
	if targetStatus == "" {
		result.Error = "PR merged status not configured"
		return result
	}

	// Get current issue status
	issue, err := p.jiraClient.GetIssue(ticket)
	if err != nil {
		result.Error = fmt.Sprintf("failed to get issue: %v", err)
		return result
	}

	result.OldStatus = issue.Status

	// Check if already in target status
	if result.OldStatus == targetStatus {
		result.NewStatus = targetStatus
		result.Success = true
		p.logger.Infof("Ticket %s already in status %s", ticket, targetStatus)
		return result
	}

	// Update issue status
	if err := p.jiraClient.UpdateStatus(ticket, targetStatus); err != nil {
		result.Error = fmt.Sprintf("failed to update status: %v", err)
		return result
	}

	result.NewStatus = targetStatus
	result.Success = true
	return result
}

// ProcessBatch processes multiple merged PRs and removes them from watching list
func (p *Processor) ProcessBatch(prs []MergedPR, state *State, watchingList *WatchingList) error {
	for _, pr := range prs {
		processedPR := p.ProcessMergedPR(pr)
		
		// Add to state
		if err := state.AddProcessedPR(processedPR); err != nil {
			p.logger.Errorf("Failed to save processed PR #%d: %v", pr.Number, err)
			continue
		}

		// Find owner and repo from watching list to remove it
		for _, watchingPR := range watchingList.GetAll() {
			if watchingPR.PRNumber == pr.Number {
				if err := watchingList.Remove(watchingPR.Owner, watchingPR.Repo, pr.Number); err != nil {
					p.logger.Warningf("Failed to remove PR #%d from watching list: %v", pr.Number, err)
				} else {
					p.logger.Infof("Removed PR #%d from watching list", pr.Number)
				}
				break
			}
		}
	}

	return nil
}

