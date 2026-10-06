package watcher

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yiminlab/quick-workflow/internal/github"
	"github.com/yiminlab/quick-workflow/internal/jira"
	"github.com/yiminlab/quick-workflow/pkg/config"
)

// Daemon represents the watch daemon
type Daemon struct {
	config       *config.Config
	scheduleConf *ScheduleConfig
	ghClient     *github.Client
	jiraClient   *jira.Client
	statusCache  *jira.StatusCache
	checker      *Checker
	processor    *Processor
	scheduler    *Scheduler
	state        *State
	logger       *Logger
	watchingList *WatchingList
	ctx          context.Context
	cancel       context.CancelFunc
}

// NewDaemon creates a new Daemon instance
func NewDaemon(cfg *config.Config, scheduleConf *ScheduleConfig) (*Daemon, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config is required")
	}

	// Create GitHub client
	ghClient, err := github.NewClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create GitHub client: %w", err)
	}

	// Create Jira client
	jiraClient, err := jira.NewClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create Jira client: %w", err)
	}

	// Load status cache
	statusCache, err := jira.NewStatusCache()
	if err != nil {
		return nil, fmt.Errorf("failed to load status cache: %w", err)
	}

	// Create logger
	logger, err := NewLogger()
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %w", err)
	}

	// Load state
	state, err := NewState()
	if err != nil {
		return nil, fmt.Errorf("failed to load state: %w", err)
	}

	// Load watching list
	watchingList, err := NewWatchingList()
	if err != nil {
		return nil, fmt.Errorf("failed to load watching list: %w", err)
	}

	// Create checker and processor
	checker := NewChecker(ghClient, logger)
	processor := NewProcessor(jiraClient, statusCache, logger)

	// Create scheduler
	if scheduleConf == nil {
		scheduleConf = DefaultScheduleConfig()
	}
	scheduler := NewScheduler(scheduleConf)

	ctx, cancel := context.WithCancel(context.Background())

	return &Daemon{
		config:       cfg,
		scheduleConf: scheduleConf,
		ghClient:     ghClient,
		jiraClient:   jiraClient,
		statusCache:  statusCache,
		checker:      checker,
		processor:    processor,
		scheduler:    scheduler,
		state:        state,
		logger:       logger,
		watchingList: watchingList,
		ctx:          ctx,
		cancel:       cancel,
	}, nil
}

// Start starts the daemon
func (d *Daemon) Start() error {
	d.logger.Info("Watch daemon starting...")
	d.logger.Infof("Watching %d PRs", d.watchingList.Count())
	d.logger.Infof("Schedule: every %d minutes (daytime), checks at %v (night)", 
		d.scheduleConf.DaytimeInterval, d.scheduleConf.NightChecks)

	// Save daemon PID
	if err := d.state.SetDaemonInfo(os.Getpid()); err != nil {
		d.logger.Warningf("Failed to save daemon PID: %v", err)
	}

	// Setup signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)

	// Main loop
	go d.mainLoop()

	// Wait for signal
	for {
		select {
		case sig := <-sigChan:
			switch sig {
			case syscall.SIGTERM, syscall.SIGINT:
				d.logger.Info("Received shutdown signal, stopping daemon...")
				return d.Stop()
			case syscall.SIGHUP:
				d.logger.Info("Received SIGHUP, reloading configuration...")
				if err := d.reloadConfig(); err != nil {
					d.logger.Errorf("Failed to reload config: %v", err)
				}
			}
		case <-d.ctx.Done():
			return nil
		}
	}
}

// Stop stops the daemon gracefully
func (d *Daemon) Stop() error {
	d.logger.Info("Stopping watch daemon...")
	d.cancel()
	
	// Clear daemon PID
	if err := d.state.ClearDaemonInfo(); err != nil {
		d.logger.Warningf("Failed to clear daemon PID: %v", err)
	}

	// Close logger
	if err := d.logger.Close(); err != nil {
		return fmt.Errorf("failed to close logger: %w", err)
	}

	return nil
}

// mainLoop is the main daemon loop
func (d *Daemon) mainLoop() {
	for {
		select {
		case <-d.ctx.Done():
			return
		default:
			// Perform check
			d.performCheck()

			// Calculate next check time
			now := time.Now()
			nextCheck := d.scheduler.CalculateNextCheckTime(now)
			d.logger.Infof("Next check at: %s (%s)", 
				nextCheck.Format("2006-01-02 15:04:05"), 
				d.scheduler.FormatNextCheckTime(nextCheck))

			// Sleep until next check
			select {
			case <-d.scheduler.SleepUntilNextCheck(now):
				// Time to check again
			case <-d.ctx.Done():
				return
			}
		}
	}
}

// performCheck performs a PR check
func (d *Daemon) performCheck() {
	d.logger.Infof("Checking %d watching PRs", d.watchingList.Count())

	// Check for merged PRs from watching list
	mergedPRs, err := d.checker.CheckMergedPRs(d.watchingList, d.state)
	if err != nil {
		d.logger.Errorf("Failed to check PRs: %v", err)
		d.state.UpdateLastCheckTime()
		return
	}

	if len(mergedPRs) == 0 {
		d.logger.Info("No newly merged PRs")
		d.state.UpdateLastCheckTime()
		return
	}

	d.logger.Infof("Found %d newly merged PR(s)", len(mergedPRs))

	// Process each PR
	if err := d.processor.ProcessBatch(mergedPRs, d.state, d.watchingList); err != nil {
		d.logger.Errorf("Failed to process PRs: %v", err)
	}

	// Update last check time
	d.state.UpdateLastCheckTime()

	// Clean old records
	retentionDays := 7 // Could be from config
	if err := d.state.CleanOldRecords(retentionDays); err != nil {
		d.logger.Warningf("Failed to clean old records: %v", err)
	}

	if err := d.logger.CleanOldLogs(retentionDays); err != nil {
		d.logger.Warningf("Failed to clean old logs: %v", err)
	}

	// Clean old watching PRs (older than 30 days)
	if err := d.watchingList.Clean(30); err != nil {
		d.logger.Warningf("Failed to clean old watching PRs: %v", err)
	}
}

// reloadConfig reloads the configuration
func (d *Daemon) reloadConfig() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to reload config: %w", err)
	}

	d.config = cfg
	d.logger.Info("Configuration reloaded")
	return nil
}

// IsRunning checks if a daemon is already running
func IsRunning() (bool, int, error) {
	state, err := NewState()
	if err != nil {
		return false, 0, err
	}

	if state.DaemonPID == 0 {
		return false, 0, nil
	}

	// Check if process exists
	process, err := os.FindProcess(state.DaemonPID)
	if err != nil {
		return false, 0, nil
	}

	// Send signal 0 to check if process is alive (doesn't actually send a signal)
	err = process.Signal(syscall.Signal(0))
	if err != nil {
		// Process doesn't exist
		return false, 0, nil
	}

	return true, state.DaemonPID, nil
}

