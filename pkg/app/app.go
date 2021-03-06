package app

import (
	"context"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/micnncim/stardust/pkg/config"
	"github.com/micnncim/stardust/pkg/github"
	"github.com/micnncim/stardust/pkg/reporter"
)

type App struct {
	config       *config.Config
	githubClient *github.Client
	reporters    []reporter.Reporter
	logger       *zap.Logger
}

func New(config *config.Config, githubClient *github.Client, reporters []reporter.Reporter, logger *zap.Logger) *App {
	return &App{
		config:       config,
		githubClient: githubClient,
		reporters:    reporters,
		logger:       logger.Named("app"),
	}
}

func (a *App) Run(ctx context.Context) error {
	repos, err := a.githubClient.ListStarredRepos(ctx, a.config.GitHubUsername, time.Now(), a.config.Interval)
	if err != nil {
		return err
	}
	if len(repos) == 0 {
		a.logger.Info("no starred repositories")
		return nil
	}

	wg := &sync.WaitGroup{}

	for _, r := range a.reporters {
		r := r
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := r.Report(ctx, repos); err != nil {
				a.logger.Error("report failed", zap.Error(err))
			}
		}()
	}

	wg.Wait()

	return nil
}
