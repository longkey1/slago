package config

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	Token   string
	Author  string
	Mention []string
}

func Load() (*Config, error) {
	cfg := &Config{
		Token:  os.Getenv("SLACK_API_TOKEN"),
		Author: os.Getenv("SLACK_AUTHOR"),
	}

	if mention := os.Getenv("SLACK_MENTION"); mention != "" {
		cfg.Mention = strings.Split(mention, ",")
	}

	return cfg, nil
}

func (c *Config) Validate() error {
	if c.Token == "" {
		return fmt.Errorf("slack API token is required (set SLACK_API_TOKEN or use --token flag)")
	}
	return nil
}
