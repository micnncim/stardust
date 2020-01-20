package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/micnncim/stardust/pkg/app"
	"github.com/micnncim/stardust/pkg/config"
	"github.com/micnncim/stardust/pkg/github"
	"github.com/micnncim/stardust/pkg/logger"
	"github.com/micnncim/stardust/pkg/reporter"
	"github.com/micnncim/stardust/pkg/reporter/slack"
)

const defaultLogLevel = "info"

func main() {
	http.HandleFunc("/", handler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
}

func handler(w http.ResponseWriter, _ *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	if err := run(ctx); err != nil {
		w.Write([]byte(err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func run(ctx context.Context) error {
	cfg, err := config.New()
	if err != nil {
		return err
	}

	lvl := defaultLogLevel
	if l := cfg.LogLevel; l != "" {
		lvl = l
	}
	log, err := logger.New(lvl)
	if err != nil {
		return err
	}

	gitHubClient, err := github.NewClient(cfg.GitHubToken, log)
	if err != nil {
		return err
	}

	var reporters []reporter.Reporter
	if cfg.EnableSlack {
		slackClient, err := slack.NewClient(cfg.SlackToken, cfg.SlackChannelID, log)
		if err != nil {
			return err
		}
		reporters = append(reporters, slackClient)
	}

	a := app.New(cfg, gitHubClient, reporters, log)
	return a.Run(ctx)
}
