package config

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/viper"
	"github.com/yiminlab/quick-workflow/internal/utils"
)

// Config holds all configuration for the application
type Config struct {
	Email              string `mapstructure:"email"`
	JiraAPIToken       string `mapstructure:"jira_api_token"`
	JiraServiceAddress string `mapstructure:"jira_service_address"`
	GitHubToken        string `mapstructure:"github_token"`
	GitHubOwner        string `mapstructure:"github_owner"`
	GitHubRepo         string `mapstructure:"github_repo"`
	BranchPrefix       string `mapstructure:"branch_prefix"`
	OpenAIKey          string `mapstructure:"openai_key"`
	DeepSeekKey        string `mapstructure:"deepseek_key"`
	OpenAIProxyURL     string `mapstructure:"openai_proxy_url"`
	OpenAIProxyKey     string `mapstructure:"openai_proxy_key"`
	CerebrasKey        string `mapstructure:"cerebras_key"`
	CerebrasURL        string `mapstructure:"cerebras_url"`
	AIProvider         string `mapstructure:"ai_provider"` // "auto", "deepseek", "openai", "cerebras"
	AutoUpdate         bool   `mapstructure:"auto_update"`
}

var globalConfig *Config

// Load loads configuration from file and environment variables
func Load() (*Config, error) {
	if globalConfig != nil {
		return globalConfig, nil
	}

	// 设置配置文件路径 (优先使用 iCloud Drive on macOS)
	configDir, err := utils.GetQuickWorkflowConfigDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get config directory: %w", err)
	}

	configFile := filepath.Join(configDir, "config.yaml")

	viper.SetConfigFile(configFile)
	viper.SetConfigType("yaml")

	// 环境变量前缀和绑定
	viper.SetEnvPrefix("QK")
	viper.AutomaticEnv()
	
	// 绑定环境变量（支持常用的环境变量名）
	viper.BindEnv("github_token", "GITHUB_TOKEN", "GH_TOKEN")
	viper.BindEnv("github_owner", "GITHUB_OWNER")
	viper.BindEnv("github_repo", "GITHUB_REPO")
	viper.BindEnv("jira_api_token", "JIRA_API_TOKEN")
	viper.BindEnv("jira_service_address", "JIRA_SERVICE_ADDRESS")
	viper.BindEnv("branch_prefix", "GH_BRANCH_PREFIX")
	viper.BindEnv("openai_key", "OPENAI_KEY")
	viper.BindEnv("deepseek_key", "DEEPSEEK_KEY")
	viper.BindEnv("openai_proxy_url", "OPENAI_PROXY_URL")
	viper.BindEnv("openai_proxy_key", "OPENAI_PROXY_KEY")
	viper.BindEnv("cerebras_key", "CEREBRAS_API_KEY")
	viper.BindEnv("cerebras_url", "CEREBRAS_BASE_URL")
	viper.BindEnv("ai_provider", "AI_PROVIDER")
	viper.BindEnv("email", "EMAIL")
	viper.BindEnv("auto_update", "AUTO_UPDATE")

	// 尝试读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
		// 配置文件不存在是正常的，将创建默认配置
	}

	// 设置默认值
	setDefaults()

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	globalConfig = &cfg
	return globalConfig, nil
}

// Get returns the global configuration
func Get() *Config {
	if globalConfig == nil {
		// 如果还没加载，尝试加载
		cfg, _ := Load()
		return cfg
	}
	return globalConfig
}

// Save saves the current configuration to file
func Save(cfg *Config) error {
	// 获取配置目录 (优先使用 iCloud Drive on macOS)
	configDir, err := utils.GetQuickWorkflowConfigDir()
	if err != nil {
		return fmt.Errorf("failed to get config directory: %w", err)
	}

	configFile := filepath.Join(configDir, "config.yaml")

	// 设置值
	viper.Set("email", cfg.Email)
	viper.Set("jira_api_token", cfg.JiraAPIToken)
	viper.Set("jira_service_address", cfg.JiraServiceAddress)
	viper.Set("github_token", cfg.GitHubToken)
	viper.Set("github_owner", cfg.GitHubOwner)
	viper.Set("github_repo", cfg.GitHubRepo)
	viper.Set("branch_prefix", cfg.BranchPrefix)
	viper.Set("openai_key", cfg.OpenAIKey)
	viper.Set("deepseek_key", cfg.DeepSeekKey)
	viper.Set("openai_proxy_url", cfg.OpenAIProxyURL)
	viper.Set("openai_proxy_key", cfg.OpenAIProxyKey)
	viper.Set("cerebras_key", cfg.CerebrasKey)
	viper.Set("cerebras_url", cfg.CerebrasURL)
	viper.Set("ai_provider", cfg.AIProvider)
	viper.Set("auto_update", cfg.AutoUpdate)

	// 写入文件
	if err := viper.WriteConfigAs(configFile); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	globalConfig = cfg
	return nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.Email == "" {
		return fmt.Errorf("email is required")
	}
	if c.JiraAPIToken == "" {
		return fmt.Errorf("jira_api_token is required")
	}
	if c.JiraServiceAddress == "" {
		return fmt.Errorf("jira_service_address is required")
	}
	if c.GitHubToken == "" {
		return fmt.Errorf("github_token is required (run: gh auth token)")
	}
	return nil
}

// IsConfigured checks if the application is configured
func IsConfigured() bool {
	cfg := Get()
	if cfg == nil {
		return false
	}
	return cfg.Validate() == nil
}

func setDefaults() {
	viper.SetDefault("branch_prefix", "")
	viper.SetDefault("auto_update", true)                                                  // 默认启用自动更新
	viper.SetDefault("ai_provider", "auto")                                                // 默认自动选择 AI provider
	viper.SetDefault("cerebras_url", "https://cerebras-proxy.brain.loocaa.com:1443/v1")    // 默认 Cerebras URL
}

// GetConfigDir returns the config directory path
func GetConfigDir() (string, error) {
	return utils.GetQuickWorkflowConfigDir()
}

