package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	GitHubToken    string `envconfig:"GITHUB_TOKEN"`
	GitHubUsername string `envconfig:"GITHUB_USERNAME"`

	SlackEnabled bool   `envconfig:"SLACK_ENABLED"`
	SlackToken   string `envconfig:"SLACK_TOKEN"`
	SlackChannel string `envconfig:"SLACK_CHANNEL"`

	LogLevel string `envconfig:"LOG_LEVEL" default:"info"`

	Interval time.Duration `envconfig:"INTERVAL" default:"168h"`
}

func New() (*Config, error) {
	c := &Config{}
	err := envconfig.Process("stardust", c)
	if err != nil {
		return nil, err
	}

	return c, nil
}
