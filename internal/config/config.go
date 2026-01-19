package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	Token   string   `mapstructure:"token"`
	Author  string   `mapstructure:"author"`
	Mention []string `mapstructure:"mention"`
}

func Load(configPath string) (*Config, error) {
	v := viper.New()

	// Set defaults from environment variables
	v.SetDefault("token", os.Getenv("SLACK_API_TOKEN"))
	v.SetDefault("author", os.Getenv("SLACK_AUTHOR"))
	if mention := os.Getenv("SLACK_MENTION"); mention != "" {
		v.SetDefault("mention", []string{mention})
	}

	// Bind environment variables
	v.SetEnvPrefix("")
	v.BindEnv("token", "SLACK_API_TOKEN")
	v.BindEnv("author", "SLACK_AUTHOR")
	v.BindEnv("mention", "SLACK_MENTION")

	// Load config file if specified or exists at default location
	if configPath != "" {
		v.SetConfigFile(configPath)
	} else {
		home, err := os.UserHomeDir()
		if err == nil {
			v.SetConfigFile(filepath.Join(home, ".slago.yaml"))
		}
	}

	v.SetConfigType("yaml")

	// Read config file (ignore error if file doesn't exist)
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// Only return error if it's not a "file not found" error
			if _, err := os.Stat(v.ConfigFileUsed()); err == nil {
				return nil, fmt.Errorf("failed to read config file: %w", err)
			}
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}

func (c *Config) Validate() error {
	if c.Token == "" {
		return fmt.Errorf("slack API token is required (set SLACK_API_TOKEN or use --token flag)")
	}
	return nil
}
