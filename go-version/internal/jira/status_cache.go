package jira

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/yiminlab/quick-workflow/internal/utils"
)

// StatusCache holds Jira status mappings for different projects
type StatusCache struct {
	filePath string
}

// StatusMapping represents status for a project
type StatusMapping struct {
	ProjectKey      string `json:"project_key"`
	PRCreatedStatus string `json:"pr_created_status"`
	PRMergedStatus  string `json:"pr_merged_status"`
}

// CacheData represents the entire cache file structure
type CacheData struct {
	Mappings map[string]StatusMapping `json:"mappings"`
}

// NewStatusCache creates a new status cache instance
func NewStatusCache() (*StatusCache, error) {
	// 获取配置目录 (优先使用 iCloud Drive on macOS)
	configDir, err := utils.GetConfigDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get config directory: %w", err)
	}

	filePath := filepath.Join(configDir, "jira-status.json")
	
	// Create file with empty data if it doesn't exist
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		emptyData := CacheData{
			Mappings: make(map[string]StatusMapping),
		}
		data, _ := json.MarshalIndent(emptyData, "", "  ")
		if err := os.WriteFile(filePath, data, 0644); err != nil {
			return nil, fmt.Errorf("failed to create status cache file: %w", err)
		}
	}

	return &StatusCache{
		filePath: filePath,
	}, nil
}

// readCache reads the entire cache from file
func (sc *StatusCache) readCache() (*CacheData, error) {
	data, err := os.ReadFile(sc.filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read cache file: %w", err)
	}

	var cache CacheData
	if err := json.Unmarshal(data, &cache); err != nil {
		return nil, fmt.Errorf("failed to parse cache file: %w", err)
	}

	if cache.Mappings == nil {
		cache.Mappings = make(map[string]StatusMapping)
	}

	return &cache, nil
}

// writeCache writes the entire cache to file
func (sc *StatusCache) writeCache(cache *CacheData) error {
	data, err := json.MarshalIndent(cache, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal cache data: %w", err)
	}

	if err := os.WriteFile(sc.filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write cache file: %w", err)
	}

	return nil
}

// GetProjectStatus retrieves cached status for a project
func (sc *StatusCache) GetProjectStatus(projectKey string) (*StatusMapping, error) {
	cache, err := sc.readCache()
	if err != nil {
		return nil, err
	}

	if mapping, exists := cache.Mappings[projectKey]; exists {
		return &mapping, nil
	}

	return nil, nil
}

// SaveProjectStatus saves status mapping for a project
func (sc *StatusCache) SaveProjectStatus(mapping *StatusMapping) error {
	cache, err := sc.readCache()
	if err != nil {
		return err
	}

	cache.Mappings[mapping.ProjectKey] = *mapping

	return sc.writeCache(cache)
}

// DeleteProjectStatus removes status mapping for a project
func (sc *StatusCache) DeleteProjectStatus(projectKey string) error {
	cache, err := sc.readCache()
	if err != nil {
		return err
	}

	delete(cache.Mappings, projectKey)

	return sc.writeCache(cache)
}

// ListAllMappings returns all status mappings
func (sc *StatusCache) ListAllMappings() ([]StatusMapping, error) {
	cache, err := sc.readCache()
	if err != nil {
		return nil, err
	}

	result := make([]StatusMapping, 0, len(cache.Mappings))
	for _, mapping := range cache.Mappings {
		result = append(result, mapping)
	}

	return result, nil
}


