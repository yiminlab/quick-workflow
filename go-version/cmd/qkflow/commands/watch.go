package commands

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/yiminlab/quick-workflow/internal/github"
	"github.com/yiminlab/quick-workflow/internal/jira"
	"github.com/yiminlab/quick-workflow/internal/ui"
	"github.com/yiminlab/quick-workflow/internal/utils"
	"github.com/yiminlab/quick-workflow/internal/watcher"
	"github.com/yiminlab/quick-workflow/pkg/config"
)

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Watch daemon for auto-updating Jira",
	Long:  `Watch daemon monitors your PRs and automatically updates Jira when merged.`,
}

var watchCheckCmd = &cobra.Command{
	Use:   "check",
	Short: "Manually check PRs once",
	Long:  `Manually check for merged PRs and update Jira status.`,
	Run:   runWatchCheck,
}

var watchStartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start watch daemon",
	Long:  `Start the watch daemon to monitor PRs in the background.`,
	Run:   runWatchStart,
}

var watchStopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop watch daemon",
	Long:  `Stop the running watch daemon.`,
	Run:   runWatchStop,
}

var watchRestartCmd = &cobra.Command{
	Use:   "restart",
	Short: "Restart watch daemon",
	Long:  `Restart the watch daemon.`,
	Run:   runWatchRestart,
}

var watchStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show watch daemon status",
	Long:  `Show the current status of the watch daemon.`,
	Run:   runWatchStatus,
}

var watchDaemonCmd = &cobra.Command{
	Use:    "daemon",
	Short:  "Run as daemon (internal use)",
	Long:   `Run as daemon process. This command is for internal use only.`,
	Hidden: true,
	Run:    runWatchDaemon,
}

var watchInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install and start watch daemon with auto-start",
	Long:  `Install watch daemon as a system service (launchd on macOS) with auto-start on login.`,
	Run:   runWatchInstall,
}

var watchUninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall watch daemon and remove auto-start",
	Long:  `Uninstall the watch daemon service and remove auto-start configuration.`,
	Run:   runWatchUninstall,
}

var watchLogCmd = &cobra.Command{
	Use:   "log",
	Short: "View watch daemon logs",
	Long:  `View the watch daemon log file.`,
	Run:   runWatchLog,
}

var watchHistoryCmd = &cobra.Command{
	Use:   "history",
	Short: "View PR processing history",
	Long:  `View the history of processed PRs.`,
	Run:   runWatchHistory,
}

var watchConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Show watch configuration",
	Long:  `Show the current watch daemon configuration.`,
	Run:   runWatchConfig,
}

var dryRun bool
var followLog bool
var logLines int
var historyDays int

func init() {
	watchCheckCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Simulate run without actually updating Jira")
	watchLogCmd.Flags().BoolVar(&followLog, "follow", false, "Follow log output (like tail -f)")
	watchLogCmd.Flags().IntVar(&logLines, "last", 50, "Show last N lines")
	watchHistoryCmd.Flags().IntVar(&historyDays, "days", 7, "Show history for last N days")
	
	watchCmd.AddCommand(watchCheckCmd)
	watchCmd.AddCommand(watchStartCmd)
	watchCmd.AddCommand(watchStopCmd)
	watchCmd.AddCommand(watchRestartCmd)
	watchCmd.AddCommand(watchStatusCmd)
	watchCmd.AddCommand(watchInstallCmd)
	watchCmd.AddCommand(watchUninstallCmd)
	watchCmd.AddCommand(watchLogCmd)
	watchCmd.AddCommand(watchHistoryCmd)
	watchCmd.AddCommand(watchConfigCmd)
	watchCmd.AddCommand(watchDaemonCmd)
}

func runWatchCheck(cmd *cobra.Command, args []string) {
	cfg := config.Get()
	if cfg == nil {
		ui.Error("Configuration not found. Please run 'qkflow init' first")
		return
	}

	// Validate configuration
	if cfg.GitHubToken == "" {
		ui.Error("GitHub token not configured")
		return
	}

	if cfg.JiraServiceAddress == "" || cfg.Email == "" || cfg.JiraAPIToken == "" {
		ui.Error("Jira not configured. Please run 'qkflow init' first")
		return
	}

	// Get GitHub client
	ghClient, err := github.NewClient()
	if err != nil {
		ui.Error(fmt.Sprintf("Failed to create GitHub client: %v", err))
		return
	}

	if dryRun {
		ui.Info("🔍 Dry-run mode: No Jira updates will be made")
		fmt.Println()
	}

	// Initialize components
	logger, err := watcher.NewLogger()
	if err != nil {
		ui.Error(fmt.Sprintf("Failed to create logger: %v", err))
		return
	}
	defer logger.Close()

	state, err := watcher.NewState()
	if err != nil {
		ui.Error(fmt.Sprintf("Failed to load state: %v", err))
		return
	}

	watchingList, err := watcher.NewWatchingList()
	if err != nil {
		ui.Error(fmt.Sprintf("Failed to load watching list: %v", err))
		return
	}

	ui.Info(fmt.Sprintf("Checking %d watching PRs...", watchingList.Count()))
	fmt.Println()

	if watchingList.Count() == 0 {
		ui.Info("No PRs in watching list")
		ui.Info("PRs will be added automatically when you create them with 'qkflow pr create'")
		return
	}

	jiraClient, err := jira.NewClient()
	if err != nil {
		ui.Error(fmt.Sprintf("Failed to create Jira client: %v", err))
		return
	}

	statusCache, err := jira.NewStatusCache()
	if err != nil {
		ui.Error(fmt.Sprintf("Failed to load Jira status cache: %v", err))
		return
	}

	// Create checker and processor
	checker := watcher.NewChecker(ghClient, logger)
	processor := watcher.NewProcessor(jiraClient, statusCache, logger)

	// Check for merged PRs from watching list
	mergedPRs, err := checker.CheckMergedPRs(watchingList, state)
	if err != nil {
		ui.Error(fmt.Sprintf("Failed to check PRs: %v", err))
		logger.Errorf("Failed to check PRs: %v", err)
		return
	}

	if len(mergedPRs) == 0 {
		ui.Success("✅ No newly merged PRs found")
		logger.Info("No newly merged PRs found")
		return
	}

	ui.Info(fmt.Sprintf("Found %d newly merged PR(s) with Jira tickets", len(mergedPRs)))
	fmt.Println()

	// Process each PR
	for _, pr := range mergedPRs {
		ui.Info(fmt.Sprintf("📋 PR #%d: %s", pr.Number, pr.Title))
		ui.Info(fmt.Sprintf("   Branch: %s", pr.Branch))
		ui.Info(fmt.Sprintf("   Jira: %v", pr.JiraTickets))
		ui.Info(fmt.Sprintf("   Merged: %s by %s", pr.MergedAt, pr.MergedBy))
		fmt.Println()

		if dryRun {
			// Dry-run: just log what would happen
			for _, ticket := range pr.JiraTickets {
				projectKey := watcher.GetProjectFromTicket(ticket)
				mapping, err := statusCache.GetProjectStatus(projectKey)
				if err != nil || mapping == nil {
					ui.Warning(fmt.Sprintf("   ⚠️  %s: No status mapping configured for project %s", ticket, projectKey))
					continue
				}

				ui.Info(fmt.Sprintf("   Would update %s → %s", ticket, mapping.PRMergedStatus))
			}
			fmt.Println()
			continue
		}

		// Actually process the PR
		processedPR := processor.ProcessMergedPR(pr)

		// Display results
		for _, update := range processedPR.JiraUpdates {
			if update.Success {
				ui.Success(fmt.Sprintf("   ✅ %s: %s → %s", update.Ticket, update.OldStatus, update.NewStatus))
			} else {
				ui.Error(fmt.Sprintf("   ❌ %s: %s", update.Ticket, update.Error))
			}
		}

		fmt.Println()

		// Save to state
		if err := state.AddProcessedPR(processedPR); err != nil {
			ui.Warning(fmt.Sprintf("Failed to save processed PR to state: %v", err))
		}

		// Remove from watching list
		for _, watchingPR := range watchingList.GetAll() {
			if watchingPR.PRNumber == pr.Number {
				if err := watchingList.Remove(watchingPR.Owner, watchingPR.Repo, pr.Number); err != nil {
					logger.Warningf("Failed to remove PR #%d from watching list: %v", pr.Number, err)
				} else {
					logger.Infof("Removed PR #%d from watching list", pr.Number)
				}
				break
			}
		}
	}

	// Update last check time
	if err := state.UpdateLastCheckTime(); err != nil {
		logger.Warningf("Failed to update last check time: %v", err)
	}

	// Clean old records
	retentionDays := 7 // Default from config
	if err := state.CleanOldRecords(retentionDays); err != nil {
		logger.Warningf("Failed to clean old records: %v", err)
	}

	if err := logger.CleanOldLogs(retentionDays); err != nil {
		logger.Warningf("Failed to clean old logs: %v", err)
	}

	if dryRun {
		ui.Info("🔍 Dry-run completed. No changes were made.")
	} else {
		ui.Success(fmt.Sprintf("✅ Processed %d PR(s)", len(mergedPRs)))
	}

	ui.Info(fmt.Sprintf("\n📝 Logs: %s", logger.GetFilePath()))
}

func runWatchDaemon(cmd *cobra.Command, args []string) {
	cfg := config.Get()
	if cfg == nil {
		fmt.Fprintf(os.Stderr, "Configuration not found. Please run 'qkflow init' first\n")
		os.Exit(1)
	}

	// Create daemon with default schedule config
	daemon, err := watcher.NewDaemon(cfg, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create daemon: %v\n", err)
		os.Exit(1)
	}

	// Start daemon (blocks until signal)
	if err := daemon.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Daemon error: %v\n", err)
		os.Exit(1)
	}
}

func runWatchStart(cmd *cobra.Command, args []string) {
	// Check if already running
	running, pid, err := watcher.IsRunning()
	if err != nil {
		ui.Error(fmt.Sprintf("Failed to check daemon status: %v", err))
		return
	}

	if running {
		ui.Warning(fmt.Sprintf("Watch daemon is already running (PID: %d)", pid))
		return
	}

	ui.Info("🚀 Starting watch daemon...")

	// Fork daemon process
	execPath, err := os.Executable()
	if err != nil {
		ui.Error(fmt.Sprintf("Failed to get executable path: %v", err))
		return
	}

	// Start daemon in background
	procAttr := &os.ProcAttr{
		Files: []*os.File{nil, nil, nil}, // Detach from terminal
	}

	process, err := os.StartProcess(execPath, []string{execPath, "watch", "daemon"}, procAttr)
	if err != nil {
		ui.Error(fmt.Sprintf("Failed to start daemon: %v", err))
		return
	}

	// Release the process
	process.Release()

	// Give it a moment to start
	time.Sleep(time.Second)

	// Verify it's running
	running, pid, _ = watcher.IsRunning()
	if running {
		ui.Success(fmt.Sprintf("✅ Watch daemon started successfully (PID: %d)", pid))
	} else {
		ui.Warning("Daemon may have failed to start. Check logs for details.")
	}
}

func runWatchStop(cmd *cobra.Command, args []string) {
	running, pid, err := watcher.IsRunning()
	if err != nil {
		ui.Error(fmt.Sprintf("Failed to check daemon status: %v", err))
		return
	}

	if !running {
		ui.Info("Watch daemon is not running")
		return
	}

	ui.Info(fmt.Sprintf("🛑 Stopping watch daemon (PID: %d)...", pid))

	// Find process
	process, err := os.FindProcess(pid)
	if err != nil {
		ui.Error(fmt.Sprintf("Failed to find process: %v", err))
		return
	}

	// Send SIGTERM
	if err := process.Signal(syscall.SIGTERM); err != nil {
		ui.Error(fmt.Sprintf("Failed to send stop signal: %v", err))
		return
	}

	// Wait for process to stop (with timeout)
	for i := 0; i < 10; i++ {
		time.Sleep(500 * time.Millisecond)
		running, _, _ := watcher.IsRunning()
		if !running {
			ui.Success("✅ Watch daemon stopped")
			return
		}
	}

	ui.Warning("Daemon did not stop within timeout. It may still be shutting down.")
}

func runWatchRestart(cmd *cobra.Command, args []string) {
	ui.Info("🔄 Restarting watch daemon...")

	// Stop if running
	running, _, _ := watcher.IsRunning()
	if running {
		runWatchStop(cmd, args)
		time.Sleep(time.Second)
	}

	// Start
	runWatchStart(cmd, args)
}

func runWatchStatus(cmd *cobra.Command, args []string) {
	state, err := watcher.NewState()
	if err != nil {
		ui.Error(fmt.Sprintf("Failed to load state: %v", err))
		return
	}

	logger, err := watcher.NewLogger()
	if err != nil {
		ui.Error(fmt.Sprintf("Failed to create logger: %v", err))
		return
	}
	defer logger.Close()

	cfg := config.Get()

	fmt.Println()
	fmt.Println("🔍 Watch Daemon Status")
	fmt.Println()

	// Check if running
	running, pid, _ := watcher.IsRunning()
	
	if running {
		uptime := time.Since(state.DaemonStartTime)
		fmt.Printf("Status: Running ✅\n")
		fmt.Printf("PID: %d\n", pid)
		fmt.Printf("Started: %s (uptime: %s)\n", 
			state.DaemonStartTime.Format("2006-01-02 15:04:05"),
			formatDuration(uptime))
		
		if !state.LastCheckTime.IsZero() {
			lastCheck := time.Since(state.LastCheckTime)
			fmt.Printf("Last Check: %s (%s ago)\n", 
				state.LastCheckTime.Format("2006-01-02 15:04:05"),
				formatDuration(lastCheck))
			
			// Calculate next check
			scheduler := watcher.NewScheduler(nil)
			nextCheck := scheduler.CalculateNextCheckTime(time.Now())
			fmt.Printf("Next Check: %s (%s)\n",
				nextCheck.Format("2006-01-02 15:04:05"),
				scheduler.FormatNextCheckTime(nextCheck))
		}
	} else {
		fmt.Println("Status: Stopped ⏸️")
		fmt.Println()
		ui.Warning("Watch daemon is not running")
	}

	fmt.Println()
	fmt.Println("📊 Statistics (last 7 days):")
	recentPRs := state.GetRecentPRs(7)
	fmt.Printf("  PRs Processed: %d\n", len(recentPRs))
	
	successCount := 0
	for _, pr := range recentPRs {
		for _, update := range pr.JiraUpdates {
			if update.Success {
				successCount++
			}
		}
	}
	fmt.Printf("  Jira Updated: %d\n", successCount)
	fmt.Printf("  Errors: %d\n", state.Stats.TotalErrors)

	if cfg != nil {
		fmt.Println()
		fmt.Println("📋 Configuration:")
		fmt.Printf("  Repository: %s/%s\n", cfg.GitHubOwner, cfg.GitHubRepo)
		fmt.Printf("  Author: %s (only your PRs)\n", cfg.GitHubOwner)
		fmt.Printf("  Check Interval: 15 minutes (daytime)\n")
		fmt.Printf("  Night Checks: 02:00, 06:00\n")
		fmt.Printf("  Log Retention: 7 days\n")
	}

	fmt.Println()
	fmt.Println("📝 Files:")
	fmt.Printf("  Log: %s\n", logger.GetFilePath())
	
	if cfg != nil {
		configDir, _ := utils.GetConfigDir()
		if configDir != "" {
			fmt.Printf("  State: %s/watch-state.json\n", configDir)
			fmt.Printf("  Config: %s/config.yaml\n", configDir)
		}
	}

	fmt.Println()
}

func formatDuration(d time.Duration) string {
	if d < time.Minute {
		return "less than a minute"
	}
	if d < time.Hour {
		return fmt.Sprintf("%d minute(s)", int(d.Minutes()))
	}
	if d < 24*time.Hour {
		hours := int(d.Hours())
		minutes := int(d.Minutes()) % 60
		if minutes == 0 {
			return fmt.Sprintf("%d hour(s)", hours)
		}
		return fmt.Sprintf("%d hour(s) %d minute(s)", hours, minutes)
	}
	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	if hours == 0 {
		return fmt.Sprintf("%d day(s)", days)
	}
	return fmt.Sprintf("%d day(s) %d hour(s)", days, hours)
}

func runWatchInstall(cmd *cobra.Command, args []string) {
	cfg := config.Get()
	if cfg == nil {
		ui.Error("Configuration not found. Please run 'qkflow init' first")
		return
	}

	ui.Info("📦 Installing watch daemon...")
	fmt.Println()

	// Check prerequisites
	ui.Info("Checking prerequisites...")
	if cfg.GitHubToken == "" {
		ui.Error("✗ GitHub token not configured")
		return
	}
	ui.Info("  ✓ GitHub token configured")

	if cfg.JiraServiceAddress == "" || cfg.Email == "" || cfg.JiraAPIToken == "" {
		ui.Error("✗ Jira not configured")
		return
	}
	ui.Info("  ✓ Jira configured")

	// Check Jira status mappings
	statusCache, err := jira.NewStatusCache()
	if err != nil {
		ui.Error(fmt.Sprintf("✗ Failed to load Jira status cache: %v", err))
		return
	}

	mappings, err := statusCache.ListAllMappings()
	if err != nil || len(mappings) == 0 {
		ui.Warning("✗ No Jira status mappings found")
		ui.Info("  Please run 'qkflow jira setup' to configure status mappings")
		return
	}
	ui.Info(fmt.Sprintf("  ✓ Jira status mappings found (%d project(s))", len(mappings)))

	fmt.Println()

	// Check if already installed
	installed, err := watcher.IsLaunchAgentInstalled()
	if err != nil {
		ui.Error(fmt.Sprintf("Failed to check installation status: %v", err))
		return
	}

	if installed {
		ui.Warning("Watch daemon is already installed")
		ui.Info("To reinstall, run 'qkflow watch uninstall' first")
		return
	}

	// Get executable path
	execPath, err := os.Executable()
	if err != nil {
		ui.Error(fmt.Sprintf("Failed to get executable path: %v", err))
		return
	}

	// Install launch agent
	ui.Info("Creating launch agent...")
	if err := watcher.InstallLaunchAgent(execPath); err != nil {
		ui.Error(fmt.Sprintf("✗ Failed to install launch agent: %v", err))
		return
	}
	ui.Info("  ✓ Generated plist file")
	ui.Info("  ✓ Installed to ~/Library/LaunchAgents/")

	// Give it a moment to start
	time.Sleep(2 * time.Second)

	// Check if running
	running, pid, _ := watcher.IsRunning()
	if running {
		ui.Info(fmt.Sprintf("  ✓ Watch daemon started (PID: %d)", pid))
	} else {
		ui.Warning("  ⚠️  Daemon may take a moment to start")
	}

	fmt.Println()
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()
	ui.Success("🎉 Installation complete!")
	fmt.Println()
	fmt.Println("The watch daemon is now running and will:")
	fmt.Println("  • Check your PRs every 15 minutes (8:30-24:00)")
	fmt.Println("  • Check twice during night (2:00, 6:00)")
	fmt.Println("  • Update Jira when PRs merge")
	fmt.Println("  • Start automatically on login")
	fmt.Println()
	fmt.Println("📊 Useful commands:")
	fmt.Println("  • Check status: qkflow watch status")
	fmt.Println("  • View logs: qkflow watch log --follow")
	fmt.Println("  • View history: qkflow watch history")
	fmt.Println("  • Stop daemon: qkflow watch stop")
	fmt.Println("  • Uninstall: qkflow watch uninstall")
	fmt.Println()

	logger, _ := watcher.NewLogger()
	if logger != nil {
		fmt.Printf("📝 Logs: %s\n", logger.GetFilePath())
		logger.Close()
	}
}

func runWatchUninstall(cmd *cobra.Command, args []string) {
	ui.Info("🗑️  Uninstalling watch daemon...")
	fmt.Println()

	// Check if installed
	installed, err := watcher.IsLaunchAgentInstalled()
	if err != nil {
		ui.Error(fmt.Sprintf("Failed to check installation status: %v", err))
		return
	}

	if !installed {
		ui.Info("Watch daemon is not installed")
		return
	}

	// Stop daemon if running
	running, pid, _ := watcher.IsRunning()
	if running {
		ui.Info(fmt.Sprintf("Stopping daemon (PID: %d)...", pid))
		
		process, err := os.FindProcess(pid)
		if err == nil {
			process.Signal(syscall.SIGTERM)
			time.Sleep(time.Second)
		}
		
		ui.Info("  ✓ Daemon stopped")
	}

	// Uninstall launch agent
	ui.Info("Removing auto-start...")
	if err := watcher.UninstallLaunchAgent(); err != nil {
		ui.Error(fmt.Sprintf("✗ Failed to uninstall launch agent: %v", err))
		return
	}
	
	plistPath, _ := watcher.GetLaunchAgentPath()
	ui.Info("  ✓ Removed launch agent")
	if plistPath != "" {
		ui.Info(fmt.Sprintf("  ✓ Deleted %s", plistPath))
	}

	fmt.Println()
	ui.Success("✅ Watch daemon completely uninstalled")
	fmt.Println()
	fmt.Println("The daemon will NOT start automatically anymore.")
	
	configDir, _ := utils.GetConfigDir()
	if configDir != "" {
		fmt.Println("Logs and history are preserved at ~/.qkflow/")
	}
	
	fmt.Println()
	fmt.Println("💡 To re-enable later, use 'qkflow watch install'")
}

func runWatchLog(cmd *cobra.Command, args []string) {
	logger, err := watcher.NewLogger()
	if err != nil {
		ui.Error(fmt.Sprintf("Failed to create logger: %v", err))
		return
	}
	defer logger.Close()

	logPath := logger.GetFilePath()

	if followLog {
		// Follow log (tail -f style)
		ui.Info(fmt.Sprintf("Following logs: %s", logPath))
		ui.Info("Press Ctrl+C to stop")
		fmt.Println()

		cmd := exec.Command("tail", "-f", logPath)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
		return
	}

	// Show last N lines
	lines, err := watcher.ReadLastLines(logPath, logLines)
	if err != nil {
		ui.Error(fmt.Sprintf("Failed to read log file: %v", err))
		return
	}

	fmt.Println()
	fmt.Printf("📝 Last %d lines from %s:\n", logLines, logPath)
	fmt.Println()
	for _, line := range lines {
		fmt.Println(line)
	}
	fmt.Println()
}

func runWatchHistory(cmd *cobra.Command, args []string) {
	state, err := watcher.NewState()
	if err != nil {
		ui.Error(fmt.Sprintf("Failed to load state: %v", err))
		return
	}

	prs := state.GetRecentPRs(historyDays)

	if len(prs) == 0 {
		ui.Info(fmt.Sprintf("No PRs processed in the last %d days", historyDays))
		return
	}

	fmt.Println()
	fmt.Printf("📋 PR Processing History (Last %d days)\n", historyDays)
	fmt.Println()

	for i, pr := range prs {
		if i > 0 {
			fmt.Println()
		}

		// Determine overall success
		allSuccess := true
		successCount := 0
		for _, update := range pr.JiraUpdates {
			if update.Success {
				successCount++
			} else {
				allSuccess = false
			}
		}

		status := "✅"
		if !allSuccess {
			if successCount == 0 {
				status = "❌"
			} else {
				status = "⚠️ "
			}
		}

		fmt.Printf("%s PR #%d: %s\n", status, pr.PRNumber, pr.PRTitle)
		fmt.Printf("   Branch: %s\n", pr.Branch)
		fmt.Printf("   Merged: %s by %s\n", pr.MergedAt.Format("2006-01-02 15:04:05"), pr.MergedBy)

		if len(pr.JiraTickets) > 0 {
			fmt.Printf("   Jira: ")
			for j, ticket := range pr.JiraTickets {
				if j > 0 {
					fmt.Printf(", ")
				}
				fmt.Printf("%s", ticket)
			}
			fmt.Println()
		}

		if len(pr.JiraUpdates) > 0 {
			for _, update := range pr.JiraUpdates {
				if update.Success {
					fmt.Printf("   ✓ %s: %s → %s\n", update.Ticket, update.OldStatus, update.NewStatus)
				} else {
					fmt.Printf("   ✗ %s: %s\n", update.Ticket, update.Error)
				}
			}
		}

		fmt.Printf("   Processed: %s\n", pr.ProcessedAt.Format("2006-01-02 15:04:05"))
	}

	fmt.Println()
	fmt.Printf("Total: %d PRs processed, %d Jira updates\n", len(prs), state.Stats.TotalJiraUpdated)
	fmt.Println()
}

func runWatchConfig(cmd *cobra.Command, args []string) {
	cfg := config.Get()
	if cfg == nil {
		ui.Error("Configuration not found")
		return
	}

	fmt.Println()
	fmt.Println("⚙️  Watch Daemon Configuration")
	fmt.Println()

	// Show watching list status
	watchingList, err := watcher.NewWatchingList()
	if err == nil {
		fmt.Println("📋 Watching List:")
		fmt.Printf("  Total PRs: %d\n", watchingList.Count())
		if watchingList.Count() > 0 {
			for _, pr := range watchingList.GetAll() {
				fmt.Printf("  • PR #%d: %s/%s\n", pr.PRNumber, pr.Owner, pr.Repo)
			}
		}
	}

	fmt.Println()
	fmt.Println("⏰ Schedule:")
	fmt.Println("  Daytime (8:30-24:00): Every 15 minutes")
	fmt.Println("  Night (0:00-8:30): 2:00, 6:00")
	fmt.Println("  Log Retention: 7 days")

	fmt.Println()
	fmt.Println("🔔 Notifications:")
	fmt.Println("  Desktop Notify: Enabled (macOS)")

	// Check installation status
	installed, _ := watcher.IsLaunchAgentInstalled()
	running, pid, _ := watcher.IsRunning()

	fmt.Println()
	fmt.Println("📊 Status:")
	if installed {
		fmt.Println("  Auto-start: Enabled ✅")
	} else {
		fmt.Println("  Auto-start: Not installed ❌")
	}

	if running {
		fmt.Printf("  Daemon: Running (PID: %d) ✅\n", pid)
	} else {
		fmt.Println("  Daemon: Stopped ⏸️")
	}

	// Show Jira mappings
	statusCache, err := jira.NewStatusCache()
	if err == nil {
		mappings, err := statusCache.ListAllMappings()
		if err == nil && len(mappings) > 0 {
			fmt.Println()
			fmt.Println("🎫 Jira Status Mappings:")
			for _, mapping := range mappings {
				fmt.Printf("  %s:\n", mapping.ProjectKey)
				fmt.Printf("    PR Created → %s\n", mapping.PRCreatedStatus)
				fmt.Printf("    PR Merged  → %s\n", mapping.PRMergedStatus)
			}
		}
	}

	fmt.Println()
}

