package github

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/google/go-github/v27/github"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

type Client struct {
	client *github.Client
	logger *zap.Logger
}

func NewClient(token string, logger *zap.Logger) (*Client, error) {
	if token == "" {
		return nil, errors.New("missing github access token")
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(context.Background(), ts)

	return &Client{
		client: github.NewClient(tc),
		logger: logger.Named("github"),
	}, nil
}

type Repo struct {
	Name          string
	Owner         string
	URL           string
	Description   string
	OwnerImageURL string
	StarredAt     time.Time
}

func (r *Repo) String() string {
	return fmt.Sprintf("%s/%s", r.Owner, r.Name)
}

func (c *Client) ListStarredRepos(ctx context.Context, username string, from time.Time, interval time.Duration) ([]*Repo, error) {
	c.logger.Info("started fetching starred repositories", zap.String("username", username))

	var srepos []*github.StarredRepository
	opt := github.ListOptions{PerPage: 50}

	for {
		repos, resp, err := c.client.Activity.ListStarred(ctx, username, &github.ActivityListStarredOptions{
			Sort:        "created",
			Direction:   "desc",
			ListOptions: opt,
		})
		if err != nil {
			c.logger.Error("failed list starred repositories", zap.Error(err))
			return nil, err
		}

		srepos = append(srepos, repos...)
		if repos[len(repos)-1].StarredAt.Before(from.Add(-interval)) {
			break
		}
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	repos := make([]*Repo, 0, len(srepos))
	for _, srepo := range srepos {
		if srepo.StarredAt.Before(from.Add(-interval)) {
			break
		}
		repo := &Repo{
			Name:          srepo.Repository.GetName(),
			Owner:         srepo.Repository.Owner.GetLogin(),
			URL:           srepo.Repository.GetHTMLURL(),
			Description:   srepo.Repository.GetDescription(),
			OwnerImageURL: srepo.Repository.Owner.GetAvatarURL(),
			StarredAt:     srepo.StarredAt.Time,
		}
		repos = append(repos, repo)
		c.logger.Info("fetched repository",
			zap.String("owner", repo.Owner),
			zap.String("repo", repo.Name),
		)
	}

	sort.Slice(repos, func(i, j int) bool {
		return repos[i].StarredAt.Before(repos[j].StarredAt)
	})

	return repos, nil
}
