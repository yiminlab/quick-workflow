package main

import (
	"fmt"
	"os"

	"github.com/yiminlab/quick-workflow/cmd/qkflow/commands"
)

var (
	// Version is set during build
	Version = "dev"
	// BuildTime is set during build
	BuildTime = "unknown"
)

func main() {
	// Set version info in commands package
	commands.Version = Version
	commands.BuildTime = BuildTime
	
	// Update root command version for --version flag
	commands.SetVersion(fmt.Sprintf("%s (built: %s)", Version, BuildTime))

	if err := commands.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

