package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/micnncim/stardust/pkg/app"
	"github.com/micnncim/stardust/pkg/berglas"
	"github.com/micnncim/stardust/pkg/config"
	"github.com/micnncim/stardust/pkg/github"
	"github.com/micnncim/stardust/pkg/logger"
	"github.com/micnncim/stardust/pkg/reporter"
	"github.com/micnncim/stardust/pkg/reporter/slack"
)

var local = flag.Bool("local", false, "run app locally")

func main() {
	flag.Parse()

	if !*local {
		if err := resolveBerglasEnvs(); err != nil {
			fmt.Fprint(os.Stderr, err)
			os.Exit(1)
		}
	}

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

	log, err := logger.New(cfg.LogLevel)
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

func resolveBerglasEnvs() error {
	l, err := logger.New("info")
	if err != nil {
		return err
	}
	ctx := context.Background()
	b, err := berglas.NewClient(ctx, l)
	if err != nil {
		return err
	}
	b.Resolve(ctx)
	return nil
}
