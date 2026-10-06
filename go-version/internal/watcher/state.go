package watcher

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/yiminlab/quick-workflow/internal/utils"
)

// ProcessedPR represents a processed pull request
type ProcessedPR struct {
	PRNumber    int               `json:"pr_number"`
	PRTitle     string            `json:"pr_title"`
	PRURL       string            `json:"pr_url"`
	Branch      string            `json:"branch"`
	JiraTickets []string          `json:"jira_tickets"`
	MergedAt    time.Time         `json:"merged_at"`
	MergedBy    string            `json:"merged_by"`
	ProcessedAt time.Time         `json:"processed_at"`
	JiraUpdates []JiraUpdateResult `json:"jira_updates"`
}

// JiraUpdateResult represents the result of updating a Jira ticket
type JiraUpdateResult struct {
	Ticket    string `json:"ticket"`
	OldStatus string `json:"old_status"`
	NewStatus string `json:"new_status"`
	Success   bool   `json:"success"`
	Error     string `json:"error,omitempty"`
}

// Statistics holds daemon statistics
type Statistics struct {
	TotalPRsProcessed int         `json:"total_prs_processed"`
	TotalJiraUpdated  int         `json:"total_jira_updated"`
	TotalErrors       int         `json:"total_errors"`
	LastError         *ErrorInfo  `json:"last_error,omitempty"`
}

// ErrorInfo represents error information
type ErrorInfo struct {
	Time    time.Time `json:"time"`
	Message string    `json:"message"`
}

// State represents the watch daemon state
type State struct {
	LastCheckTime    time.Time      `json:"last_check_time"`
	DaemonPID        int            `json:"daemon_pid"`
	DaemonStartTime  time.Time      `json:"daemon_start_time"`
	ProcessedPRs     []ProcessedPR  `json:"processed_prs"`
	Stats            Statistics     `json:"stats"`
	filePath         string         `json:"-"`
}

// NewState creates a new State instance
func NewState() (*State, error) {
	configDir, err := utils.GetConfigDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get config directory: %w", err)
	}

	filePath := filepath.Join(configDir, "watch-state.json")

	// Load existing state if file exists
	if _, err := os.Stat(filePath); err == nil {
		data, err := os.ReadFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read state file: %w", err)
		}

		var state State
		if err := json.Unmarshal(data, &state); err != nil {
			return nil, fmt.Errorf("failed to parse state file: %w", err)
		}

		state.filePath = filePath
		return &state, nil
	}

	// Create new state
	state := &State{
		ProcessedPRs: make([]ProcessedPR, 0),
		Stats:        Statistics{},
		filePath:     filePath,
	}

	// Save initial state
	if err := state.Save(); err != nil {
		return nil, err
	}

	return state, nil
}

// Save saves the state to file
func (s *State) Save() error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal state: %w", err)
	}

	if err := os.WriteFile(s.filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write state file: %w", err)
	}

	return nil
}

// AddProcessedPR adds a processed PR to the state
func (s *State) AddProcessedPR(pr ProcessedPR) error {
	pr.ProcessedAt = time.Now()
	s.ProcessedPRs = append(s.ProcessedPRs, pr)

	// Update statistics
	s.Stats.TotalPRsProcessed++
	for _, update := range pr.JiraUpdates {
		if update.Success {
			s.Stats.TotalJiraUpdated++
		} else {
			s.Stats.TotalErrors++
			s.Stats.LastError = &ErrorInfo{
				Time:    time.Now(),
				Message: update.Error,
			}
		}
	}

	return s.Save()
}

// IsPRProcessed checks if a PR has been processed
func (s *State) IsPRProcessed(prNumber int) bool {
	for _, pr := range s.ProcessedPRs {
		if pr.PRNumber == prNumber {
			return true
		}
	}
	return false
}

// CleanOldRecords removes records older than the specified days
func (s *State) CleanOldRecords(retentionDays int) error {
	if retentionDays <= 0 {
		return nil
	}

	cutoffTime := time.Now().AddDate(0, 0, -retentionDays)
	newProcessedPRs := make([]ProcessedPR, 0)

	for _, pr := range s.ProcessedPRs {
		if pr.ProcessedAt.After(cutoffTime) {
			newProcessedPRs = append(newProcessedPRs, pr)
		}
	}

	s.ProcessedPRs = newProcessedPRs
	return s.Save()
}

// UpdateLastCheckTime updates the last check time
func (s *State) UpdateLastCheckTime() error {
	s.LastCheckTime = time.Now()
	return s.Save()
}

// SetDaemonInfo sets daemon process information
func (s *State) SetDaemonInfo(pid int) error {
	s.DaemonPID = pid
	s.DaemonStartTime = time.Now()
	return s.Save()
}

// ClearDaemonInfo clears daemon process information
func (s *State) ClearDaemonInfo() error {
	s.DaemonPID = 0
	s.DaemonStartTime = time.Time{}
	return s.Save()
}

// GetRecentPRs returns PRs processed within the last N days
func (s *State) GetRecentPRs(days int) []ProcessedPR {
	if days <= 0 {
		return s.ProcessedPRs
	}

	cutoffTime := time.Now().AddDate(0, 0, -days)
	recent := make([]ProcessedPR, 0)

	for _, pr := range s.ProcessedPRs {
		if pr.ProcessedAt.After(cutoffTime) {
			recent = append(recent, pr)
		}
	}

	return recent
}

