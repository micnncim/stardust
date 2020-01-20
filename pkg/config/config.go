package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	GitHubToken    string `envconfig:"GITHUB_TOKEN"`
	GitHubUsername string `envconfig:"GITHUB_USERNAME"`

	EnableSlack    bool   `envconfig:"ENABLE_SLACK"`
	SlackToken     string `envconfig:"SLACK_TOKEN"`
	SlackChannelID string `envconfig:"SLACK_CHANNEL_ID"`

	LogLevel string `envconfig:"LOG_LEVEL"`

	Interval time.Duration `envconfig:"INTERVAL"`
}

func New() (*Config, error) {
	c := &Config{}
	err := envconfig.Process("", c)
	if err != nil {
		return nil, err
	}
	return c, nil
}
