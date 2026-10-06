package watcher

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/yiminlab/quick-workflow/internal/utils"
)

// WatchingPR represents a PR that is being watched
type WatchingPR struct {
	PRNumber    int      `json:"pr_number"`
	Owner       string   `json:"owner"`
	Repo        string   `json:"repo"`
	Branch      string   `json:"branch"`
	Title       string   `json:"title"`
	PRURL       string   `json:"pr_url"`
	JiraTickets []string `json:"jira_tickets"`
	CreatedAt   string   `json:"created_at"`
}

// WatchingList holds the list of PRs being watched
type WatchingList struct {
	PRs      []WatchingPR `json:"prs"`
	filePath string       `json:"-"`
}

// NewWatchingList creates or loads the watching list
func NewWatchingList() (*WatchingList, error) {
	configDir, err := utils.GetConfigDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get config directory: %w", err)
	}

	filePath := filepath.Join(configDir, "watching-prs.json")

	// Load existing file if exists
	if _, err := os.Stat(filePath); err == nil {
		data, err := os.ReadFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read watching list: %w", err)
		}

		var list WatchingList
		if err := json.Unmarshal(data, &list); err != nil {
			return nil, fmt.Errorf("failed to parse watching list: %w", err)
		}

		list.filePath = filePath
		return &list, nil
	}

	// Create new list
	list := &WatchingList{
		PRs:      make([]WatchingPR, 0),
		filePath: filePath,
	}

	if err := list.Save(); err != nil {
		return nil, err
	}

	return list, nil
}

// Save saves the watching list to file
func (w *WatchingList) Save() error {
	data, err := json.MarshalIndent(w, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal watching list: %w", err)
	}

	if err := os.WriteFile(w.filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write watching list: %w", err)
	}

	return nil
}

// Add adds a PR to the watching list
func (w *WatchingList) Add(pr WatchingPR) error {
	// Check if already exists
	for _, existing := range w.PRs {
		if existing.Owner == pr.Owner && existing.Repo == pr.Repo && existing.PRNumber == pr.PRNumber {
			// Already exists, update it
			return w.Update(pr)
		}
	}

	// Add new PR
	if pr.CreatedAt == "" {
		pr.CreatedAt = time.Now().Format(time.RFC3339)
	}

	w.PRs = append(w.PRs, pr)
	return w.Save()
}

// Update updates an existing PR in the watching list
func (w *WatchingList) Update(pr WatchingPR) error {
	for i, existing := range w.PRs {
		if existing.Owner == pr.Owner && existing.Repo == pr.Repo && existing.PRNumber == pr.PRNumber {
			// Keep original created_at if not provided
			if pr.CreatedAt == "" {
				pr.CreatedAt = existing.CreatedAt
			}
			w.PRs[i] = pr
			return w.Save()
		}
	}

	// Not found, add as new
	return w.Add(pr)
}

// Remove removes a PR from the watching list
func (w *WatchingList) Remove(owner, repo string, prNumber int) error {
	newPRs := make([]WatchingPR, 0)
	for _, pr := range w.PRs {
		if pr.Owner == owner && pr.Repo == repo && pr.PRNumber == prNumber {
			continue // Skip this one
		}
		newPRs = append(newPRs, pr)
	}

	w.PRs = newPRs
	return w.Save()
}

// GetAll returns all watching PRs
func (w *WatchingList) GetAll() []WatchingPR {
	return w.PRs
}

// Count returns the number of watching PRs
func (w *WatchingList) Count() int {
	return len(w.PRs)
}

// Exists checks if a PR is in the watching list
func (w *WatchingList) Exists(owner, repo string, prNumber int) bool {
	for _, pr := range w.PRs {
		if pr.Owner == owner && pr.Repo == repo && pr.PRNumber == prNumber {
			return true
		}
	}
	return false
}

// Clean removes PRs older than the specified days
func (w *WatchingList) Clean(retentionDays int) error {
	if retentionDays <= 0 {
		return nil
	}

	cutoffTime := time.Now().AddDate(0, 0, -retentionDays)
	newPRs := make([]WatchingPR, 0)

	for _, pr := range w.PRs {
		createdAt, err := time.Parse(time.RFC3339, pr.CreatedAt)
		if err != nil {
			// Keep PRs with invalid timestamps
			newPRs = append(newPRs, pr)
			continue
		}

		if createdAt.After(cutoffTime) {
			newPRs = append(newPRs, pr)
		}
	}

	w.PRs = newPRs
	return w.Save()
}

