package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/micnncim/stardust/pkg/app"
	"github.com/micnncim/stardust/pkg/config"
	"github.com/micnncim/stardust/pkg/github"
	"github.com/micnncim/stardust/pkg/logging"
	"github.com/micnncim/stardust/pkg/reporter"
	"github.com/micnncim/stardust/pkg/reporter/slack"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	if err := run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "failed to run stardust: %v", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	cfg, err := config.New()
	if err != nil {
		return err
	}

	logger, err := logging.NewLogger(os.Stdout, logging.LevelInfo, logging.FormatColorConsole)
	if err != nil {
		return err
	}

	gitHubClient, err := github.NewClient(cfg.GitHubToken, github.WithLogger(logger))
	if err != nil {
		return err
	}

	var reporters []reporter.Reporter
	if cfg.SlackEnabled {
		logger.Info("slack enabled")
		slackClient, err := slack.NewClient(cfg.SlackToken, cfg.SlackChannel, slack.WithLogger(logger))
		if err != nil {
			return err
		}
		reporters = append(reporters, slackClient)
	}

	return app.New(cfg, gitHubClient, reporters, logger).Run(ctx)
}
