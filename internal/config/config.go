package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Editor represents a single IDE configuration
type Editor struct {
	Enabled     bool   `json:"enabled"`
	Command     string `json:"command"`
	DisplayName string `json:"displayName"`
}

// Config represents the full application configuration
type Config struct {
	Editors map[string]Editor `json:"editors"`
}

// DefaultConfig returns the default configuration with common editors
func DefaultConfig() *Config {
	return &Config{
		Editors: map[string]Editor{
			"cursor": {
				Enabled:     true,
				Command:     "cursor",
				DisplayName: "Cursor",
			},
			"vscode": {
				Enabled:     true,
				Command:     "code",
				DisplayName: "VS Code",
			},
			"zed": {
				Enabled:     false,
				Command:     "zed",
				DisplayName: "Zed",
			},
			"sublime": {
				Enabled:     false,
				Command:     "subl",
				DisplayName: "Sublime Text",
			},
			"neovim": {
				Enabled:     false,
				Command:     "nvim",
				DisplayName: "Neovim",
			},
		},
	}
}

// GetEnabledEditors returns only the editors that are enabled
func (c *Config) GetEnabledEditors() []Editor {
	var enabled []Editor
	for _, editor := range c.Editors {
		if editor.Enabled {
			enabled = append(enabled, editor)
		}
	}
	return enabled
}

// GetConfigDir returns the config directory path
func GetConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".config", "worktree"), nil
}

// GetConfigPath returns the full path to config.json
func GetConfigPath() (string, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "config.json"), nil
}

// LoadConfig reads the configuration from disk or returns default
func LoadConfig() (*Config, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return DefaultConfig(), nil
	}

	data, err := os.ReadFile(configPath)
	if os.IsNotExist(err) {
		// First run - create default config
		defaultCfg := DefaultConfig()
		_ = SaveConfig(defaultCfg)
		return defaultCfg, nil
	}
	if err != nil {
		return DefaultConfig(), err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return DefaultConfig(), err
	}

	return &cfg, nil
}

// SaveConfig writes the configuration to disk
func SaveConfig(cfg *Config) error {
	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}

	configDir, err := GetConfigDir()
	if err != nil {
		return err
	}

	// Create config directory if it doesn't exist
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}
