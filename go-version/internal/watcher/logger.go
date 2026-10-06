package watcher

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/yiminlab/quick-workflow/internal/utils"
)

// Logger handles watch daemon logging
type Logger struct {
	filePath string
	file     *os.File
}

// NewLogger creates a new Logger instance
func NewLogger() (*Logger, error) {
	configDir, err := utils.GetConfigDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get config directory: %w", err)
	}

	filePath := filepath.Join(configDir, "watch.log")

	// Open file in append mode, create if not exists
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	return &Logger{
		filePath: filePath,
		file:     file,
	}, nil
}

// Close closes the log file
func (l *Logger) Close() error {
	if l.file != nil {
		return l.file.Close()
	}
	return nil
}

// log writes a log entry with the given level
func (l *Logger) log(level, message string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logLine := fmt.Sprintf("[%s %s] %s\n", timestamp, level, message)
	
	if l.file != nil {
		l.file.WriteString(logLine)
	}
}

// Info logs an info message
func (l *Logger) Info(message string) {
	l.log("INFO", message)
}

// Infof logs a formatted info message
func (l *Logger) Infof(format string, args ...interface{}) {
	l.Info(fmt.Sprintf(format, args...))
}

// Success logs a success message
func (l *Logger) Success(message string) {
	l.log("SUCCESS", message)
}

// Successf logs a formatted success message
func (l *Logger) Successf(format string, args ...interface{}) {
	l.Success(fmt.Sprintf(format, args...))
}

// Warning logs a warning message
func (l *Logger) Warning(message string) {
	l.log("WARNING", message)
}

// Warningf logs a formatted warning message
func (l *Logger) Warningf(format string, args ...interface{}) {
	l.Warning(fmt.Sprintf(format, args...))
}

// Error logs an error message
func (l *Logger) Error(message string) {
	l.log("ERROR", message)
}

// Errorf logs a formatted error message
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.Error(fmt.Sprintf(format, args...))
}

// GetFilePath returns the log file path
func (l *Logger) GetFilePath() string {
	return l.filePath
}

// CleanOldLogs removes log entries older than the specified days
func (l *Logger) CleanOldLogs(retentionDays int) error {
	if retentionDays <= 0 {
		return nil
	}

	// Read all logs
	data, err := os.ReadFile(l.filePath)
	if err != nil {
		return fmt.Errorf("failed to read log file: %w", err)
	}

	// Parse and filter logs
	cutoffTime := time.Now().AddDate(0, 0, -retentionDays)
	lines := string(data)
	
	// Create temporary file
	tmpFile := l.filePath + ".tmp"
	tmpF, err := os.OpenFile(tmpFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}

	// Write filtered logs to temp file
	var currentLine string
	for i := 0; i < len(lines); i++ {
		if lines[i] == '\n' {
			if currentLine != "" {
				// Try to parse timestamp
				if len(currentLine) > 20 {
					timestampStr := currentLine[1:20] // Extract timestamp part
					if t, err := time.Parse("2006-01-02 15:04:05", timestampStr); err == nil {
						if t.After(cutoffTime) {
							tmpF.WriteString(currentLine + "\n")
						}
					} else {
						// If parsing fails, keep the line (might be a continuation line)
						tmpF.WriteString(currentLine + "\n")
					}
				}
			}
			currentLine = ""
		} else {
			currentLine += string(lines[i])
		}
	}

	tmpF.Close()

	// Replace original file with temp file
	if err := os.Rename(tmpFile, l.filePath); err != nil {
		return fmt.Errorf("failed to replace log file: %w", err)
	}

	// Reopen the file
	l.file.Close()
	l.file, err = os.OpenFile(l.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to reopen log file: %w", err)
	}

	return nil
}

// ReadLastLines reads the last N lines from the log file
func ReadLastLines(filePath string, n int) ([]string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read log file: %w", err)
	}

	lines := make([]string, 0)
	currentLine := ""
	
	for i := 0; i < len(data); i++ {
		if data[i] == '\n' {
			if currentLine != "" {
				lines = append(lines, currentLine)
			}
			currentLine = ""
		} else {
			currentLine += string(data[i])
		}
	}
	
	if currentLine != "" {
		lines = append(lines, currentLine)
	}

	// Return last n lines
	if len(lines) <= n {
		return lines, nil
	}
	
	return lines[len(lines)-n:], nil
}

