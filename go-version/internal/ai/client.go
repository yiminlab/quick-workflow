package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/yiminlab/quick-workflow/pkg/config"
)

// Client represents an AI client for translation and optimization
type Client struct {
	apiKey   string
	apiURL   string
	model    string
	provider string // "openai", "deepseek", or "cerebras"
}

// TranslationResult represents the result of AI translation
type TranslationResult struct {
	OriginalTitle  string
	NeedTranslate  bool
	TranslatedTitle string
}

// NewClient creates a new AI client
func NewClient() (*Client, error) {
	cfg := config.Get()
	if cfg == nil {
		return nil, fmt.Errorf("config not loaded")
	}

	// 根据 ai_provider 配置选择 AI 服务
	provider := cfg.AIProvider
	if provider == "" {
		provider = "auto"
	}

	switch provider {
	case "cerebras":
		if cfg.CerebrasKey != "" {
			apiURL := cfg.CerebrasURL
			if apiURL == "" {
				apiURL = "https://cerebras-proxy.brain.loocaa.com:1443/v1/chat/completions"
			} else if !strings.HasSuffix(apiURL, "/chat/completions") {
				apiURL = strings.TrimSuffix(apiURL, "/") + "/chat/completions"
			}
			return &Client{
				apiKey:   cfg.CerebrasKey,
				apiURL:   apiURL,
				model:    "llama-3.3-70b",
				provider: "cerebras",
			}, nil
		}
		return nil, fmt.Errorf("cerebras selected but CEREBRAS_API_KEY not configured")

	case "deepseek":
		if cfg.DeepSeekKey != "" {
			return &Client{
				apiKey:   cfg.DeepSeekKey,
				apiURL:   "https://api.deepseek.com/v1/chat/completions",
				model:    "deepseek-chat",
				provider: "deepseek",
			}, nil
		}
		return nil, fmt.Errorf("deepseek selected but DEEPSEEK_KEY not configured")

	case "openai":
		if cfg.OpenAIKey != "" {
			apiURL := "https://api.openai.com/v1/chat/completions"
			if cfg.OpenAIProxyURL != "" {
				apiURL = cfg.OpenAIProxyURL
				if !strings.HasSuffix(apiURL, "/chat/completions") {
					apiURL = strings.TrimSuffix(apiURL, "/") + "/chat/completions"
				}
			}
			return &Client{
				apiKey:   cfg.OpenAIKey,
				apiURL:   apiURL,
				model:    "gpt-3.5-turbo",
				provider: "openai",
			}, nil
		}
		return nil, fmt.Errorf("openai selected but OPENAI_KEY not configured")

	default: // "auto" - 自动选择可用的 AI 服务
		// 优先级: Cerebras > DeepSeek > OpenAI
		if cfg.CerebrasKey != "" {
			apiURL := cfg.CerebrasURL
			if apiURL == "" {
				apiURL = "https://cerebras-proxy.brain.loocaa.com:1443/v1/chat/completions"
			} else if !strings.HasSuffix(apiURL, "/chat/completions") {
				apiURL = strings.TrimSuffix(apiURL, "/") + "/chat/completions"
			}
			return &Client{
				apiKey:   cfg.CerebrasKey,
				apiURL:   apiURL,
				model:    "llama-3.3-70b",
				provider: "cerebras",
			}, nil
		}

		if cfg.DeepSeekKey != "" {
			return &Client{
				apiKey:   cfg.DeepSeekKey,
				apiURL:   "https://api.deepseek.com/v1/chat/completions",
				model:    "deepseek-chat",
				provider: "deepseek",
			}, nil
		}

		if cfg.OpenAIKey != "" {
			apiURL := "https://api.openai.com/v1/chat/completions"
			if cfg.OpenAIProxyURL != "" {
				apiURL = cfg.OpenAIProxyURL
				if !strings.HasSuffix(apiURL, "/chat/completions") {
					apiURL = strings.TrimSuffix(apiURL, "/") + "/chat/completions"
				}
			}
			return &Client{
				apiKey:   cfg.OpenAIKey,
				apiURL:   apiURL,
				model:    "gpt-3.5-turbo",
				provider: "openai",
			}, nil
		}
	}

	return nil, fmt.Errorf("no AI API key configured (CEREBRAS_API_KEY, DEEPSEEK_KEY, or OPENAI_KEY)")
}

// ChatCompletionRequest represents the request to chat completion API
type ChatCompletionRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

// Message represents a chat message
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatCompletionResponse represents the response from chat completion API
type ChatCompletionResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

// GeneratePRTitle generates a concise PR title based on Jira summary and PR type
func (c *Client) GeneratePRTitle(jiraSummary, prType, shortDesc string) (string, error) {
	context := jiraSummary
	if shortDesc != "" {
		context = fmt.Sprintf("%s. Additional context: %s", jiraSummary, shortDesc)
	}
	
	prompt := fmt.Sprintf(`You are a professional software engineer writing a GitHub PR title.
Based on the following information, generate a concise PR title (max 60 characters).

PR Type: %s
Jira Summary: %s

Requirements:
1. Start with the type prefix (e.g., "fix:", "feat:", "docs:")
2. Be concise and clear (max 60 characters)
3. Focus on WHAT was changed, not WHY
4. Use English
5. Use lowercase after the colon
6. Only return the PR title, nothing else

Example good titles:
- fix: empty text in Dynamic Island during download
- feat: add user authentication
- docs: update installation guide

Generate the PR title now:`, prType, context)

	title, err := c.callChatAPI(prompt)
	if err != nil {
		return "", err
	}
	
	// 确保标题不超过 60 个字符
	if len(title) > 60 {
		title = title[:57] + "..."
	}
	
	return title, nil
}

// TranslateAndOptimize translates and optimizes a Jira issue title for PR
func (c *Client) TranslateAndOptimize(title string) (*TranslationResult, error) {
	// 检测是否需要翻译（包含中文字符）
	needTranslate := containsChinese(title)

	if !needTranslate {
		// 如果不需要翻译，仍然可以优化标题
		optimized, err := c.optimizeTitle(title)
		if err != nil {
			return &TranslationResult{
				OriginalTitle:   title,
				NeedTranslate:   false,
				TranslatedTitle: title,
			}, nil
		}
		return &TranslationResult{
			OriginalTitle:   title,
			NeedTranslate:   false,
			TranslatedTitle: optimized,
		}, nil
	}

	// 需要翻译
	translated, err := c.translateTitle(title)
	if err != nil {
		return &TranslationResult{
			OriginalTitle:   title,
			NeedTranslate:   true,
			TranslatedTitle: title, // 翻译失败，返回原标题
		}, err
	}

	return &TranslationResult{
		OriginalTitle:   title,
		NeedTranslate:   true,
		TranslatedTitle: translated,
	}, nil
}

// translateTitle translates a Chinese title to English
func (c *Client) translateTitle(title string) (string, error) {
	prompt := fmt.Sprintf(`You are a professional translator for software development. 
Translate the following Jira issue title from Chinese to English. 
Make it concise, clear, and suitable for a GitHub PR title.
Keep technical terms in English if they already exist in the original text.
Only return the translated title, nothing else.

Title: %s`, title)

	return c.callChatAPI(prompt)
}

// optimizeTitle optimizes an English title
func (c *Client) optimizeTitle(title string) (string, error) {
	prompt := fmt.Sprintf(`You are a professional software engineer. 
Optimize the following PR title to be more concise and clear.
Follow conventional commit format if applicable (feat:, fix:, docs:, etc.).
Only return the optimized title, nothing else.

Title: %s`, title)

	return c.callChatAPI(prompt)
}

// callChatAPI makes a request to the chat completion API
func (c *Client) callChatAPI(prompt string) (string, error) {
	reqBody := ChatCompletionRequest{
		Model: c.model,
		Messages: []Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", c.apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to call API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	var response ChatCompletionResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}

	result := strings.TrimSpace(response.Choices[0].Message.Content)
	return result, nil
}

// containsChinese checks if a string contains Chinese characters
func containsChinese(s string) bool {
	for _, r := range s {
		if r >= 0x4E00 && r <= 0x9FFF {
			return true
		}
	}
	return false
}

