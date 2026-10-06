package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yiminlab/quick-workflow/internal/ui"
	"github.com/yiminlab/quick-workflow/pkg/config"
)

var aiCmd = &cobra.Command{
	Use:   "ai",
	Short: "Manage AI provider settings",
	Long:  `View and switch between AI providers (Cerebras, DeepSeek, OpenAI).`,
	Run:   runAIStatus,
}

var aiSwitchCmd = &cobra.Command{
	Use:   "switch [provider]",
	Short: "Switch AI provider",
	Long: `Switch to a specific AI provider.
Available providers: auto, cerebras, deepseek, openai

Examples:
  qkflow ai switch cerebras   # Switch to Cerebras
  qkflow ai switch deepseek   # Switch to DeepSeek
  qkflow ai switch openai     # Switch to OpenAI
  qkflow ai switch auto       # Auto-select best available`,
	Args: cobra.ExactArgs(1),
	Run:  runAISwitch,
}

var aiSetCmd = &cobra.Command{
	Use:   "set [key] [value]",
	Short: "Set AI configuration",
	Long: `Set AI provider configuration values.

Examples:
  qkflow ai set cerebras-key YOUR_API_KEY
  qkflow ai set cerebras-url https://api.cerebras.ai/v1
  qkflow ai set deepseek-key YOUR_API_KEY
  qkflow ai set openai-key YOUR_API_KEY`,
	Args: cobra.ExactArgs(2),
	Run:  runAISet,
}

func init() {
	aiCmd.AddCommand(aiSwitchCmd)
	aiCmd.AddCommand(aiSetCmd)
}

func runAIStatus(cmd *cobra.Command, args []string) {
	cfg := config.Get()
	if cfg == nil {
		ui.Error("No configuration found. Run 'qkflow init' first.")
		return
	}

	fmt.Println("🤖 AI Provider Status:")
	fmt.Println()

	// Current provider
	provider := cfg.AIProvider
	if provider == "" {
		provider = "auto"
	}
	fmt.Printf("  Current Mode: %s\n", provider)
	fmt.Println()

	// Available providers
	fmt.Println("  Available Providers:")

	// Cerebras
	if cfg.CerebrasKey != "" {
		if provider == "cerebras" || (provider == "auto" && cfg.CerebrasKey != "") {
			fmt.Printf("    ✅ Cerebras: %s (Active)\n", maskToken(cfg.CerebrasKey))
		} else {
			fmt.Printf("    ⚪ Cerebras: %s\n", maskToken(cfg.CerebrasKey))
		}
		if cfg.CerebrasURL != "" {
			fmt.Printf("       URL: %s\n", cfg.CerebrasURL)
		}
	} else {
		fmt.Printf("    ❌ Cerebras: not configured\n")
	}

	// DeepSeek
	if cfg.DeepSeekKey != "" {
		isActive := provider == "deepseek" || (provider == "auto" && cfg.CerebrasKey == "" && cfg.DeepSeekKey != "")
		if isActive {
			fmt.Printf("    ✅ DeepSeek: %s (Active)\n", maskToken(cfg.DeepSeekKey))
		} else {
			fmt.Printf("    ⚪ DeepSeek: %s\n", maskToken(cfg.DeepSeekKey))
		}
	} else {
		fmt.Printf("    ❌ DeepSeek: not configured\n")
	}

	// OpenAI
	if cfg.OpenAIKey != "" {
		isActive := provider == "openai" || (provider == "auto" && cfg.CerebrasKey == "" && cfg.DeepSeekKey == "" && cfg.OpenAIKey != "")
		if isActive {
			fmt.Printf("    ✅ OpenAI: %s (Active)\n", maskToken(cfg.OpenAIKey))
		} else {
			fmt.Printf("    ⚪ OpenAI: %s\n", maskToken(cfg.OpenAIKey))
		}
		if cfg.OpenAIProxyURL != "" {
			fmt.Printf("       Proxy: %s\n", cfg.OpenAIProxyURL)
		}
	} else {
		fmt.Printf("    ❌ OpenAI: not configured\n")
	}

	fmt.Println()
	fmt.Println("  Commands:")
	fmt.Println("    qkflow ai switch <provider>  - Switch AI provider (auto/cerebras/deepseek/openai)")
	fmt.Println("    qkflow ai set <key> <value>  - Set API key or URL")
}

func runAISwitch(cmd *cobra.Command, args []string) {
	provider := args[0]

	// Validate provider
	validProviders := map[string]bool{
		"auto":     true,
		"cerebras": true,
		"deepseek": true,
		"openai":   true,
	}

	if !validProviders[provider] {
		ui.Error(fmt.Sprintf("Invalid provider: %s. Valid options: auto, cerebras, deepseek, openai", provider))
		return
	}

	cfg := config.Get()
	if cfg == nil {
		ui.Error("No configuration found. Run 'qkflow init' first.")
		return
	}

	// Check if the selected provider has a key configured
	switch provider {
	case "cerebras":
		if cfg.CerebrasKey == "" {
			ui.Warning("Cerebras API key not configured. Set it with: qkflow ai set cerebras-key YOUR_KEY")
		}
	case "deepseek":
		if cfg.DeepSeekKey == "" {
			ui.Warning("DeepSeek API key not configured. Set it with: qkflow ai set deepseek-key YOUR_KEY")
		}
	case "openai":
		if cfg.OpenAIKey == "" {
			ui.Warning("OpenAI API key not configured. Set it with: qkflow ai set openai-key YOUR_KEY")
		}
	}

	cfg.AIProvider = provider

	if err := config.Save(cfg); err != nil {
		ui.Error(fmt.Sprintf("Failed to save configuration: %v", err))
		return
	}

	ui.Success(fmt.Sprintf("Switched AI provider to: %s", provider))

	// Show which provider will be active
	if provider == "auto" {
		if cfg.CerebrasKey != "" {
			ui.Info("Auto mode will use: Cerebras (highest priority)")
		} else if cfg.DeepSeekKey != "" {
			ui.Info("Auto mode will use: DeepSeek")
		} else if cfg.OpenAIKey != "" {
			ui.Info("Auto mode will use: OpenAI")
		} else {
			ui.Warning("No AI keys configured. Run 'qkflow ai set' to add one.")
		}
	}
}

func runAISet(cmd *cobra.Command, args []string) {
	key := args[0]
	value := args[1]

	cfg := config.Get()
	if cfg == nil {
		ui.Error("No configuration found. Run 'qkflow init' first.")
		return
	}

	switch key {
	case "cerebras-key":
		cfg.CerebrasKey = value
		ui.Success("Cerebras API key set successfully")
	case "cerebras-url":
		cfg.CerebrasURL = value
		ui.Success(fmt.Sprintf("Cerebras URL set to: %s", value))
	case "deepseek-key":
		cfg.DeepSeekKey = value
		ui.Success("DeepSeek API key set successfully")
	case "openai-key":
		cfg.OpenAIKey = value
		ui.Success("OpenAI API key set successfully")
	case "openai-proxy-url":
		cfg.OpenAIProxyURL = value
		ui.Success(fmt.Sprintf("OpenAI Proxy URL set to: %s", value))
	default:
		ui.Error(fmt.Sprintf("Unknown key: %s", key))
		fmt.Println("Valid keys: cerebras-key, cerebras-url, deepseek-key, openai-key, openai-proxy-url")
		return
	}

	if err := config.Save(cfg); err != nil {
		ui.Error(fmt.Sprintf("Failed to save configuration: %v", err))
		return
	}
}
