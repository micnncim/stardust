package slack

import (
	"context"
	"errors"
	"fmt"

	"github.com/nlopes/slack"
	"go.uber.org/zap"

	"github.com/micnncim/stardust/pkg/github"
	"github.com/micnncim/stardust/pkg/reporter"
)

const (
	timeFormat = "2006/01/02 15:04"

	githubColor = "#24292E"

	attachmentTitle = "GitHub Star Report"
)

// Slack API limits the number of attachments.
// https://api.slack.com/methods/chat.postMessage
var maxAttachments = 100

type Client struct {
	client  slackClientInterface
	channel string
	logger  *zap.Logger
}

var _ reporter.Reporter = (*Client)(nil)

type Option func(*Client)

func WithLogger(l *zap.Logger) Option {
	return func(c *Client) {
		c.logger = l.Named("slack")
	}
}

type slackClientInterface interface {
	PostAttachmentsMessage(ctx context.Context, channel string, attachments ...slack.Attachment) error
}

type slackClient struct {
	client *slack.Client
}

var _ slackClientInterface = (*slackClient)(nil)

func NewClient(token, channel string, opts ...Option) (reporter.Reporter, error) {
	if token == "" {
		return nil, errors.New("missing slack access token")
	}
	if channel == "" {
		return nil, errors.New("missing slack channel id")
	}

	c := &Client{
		client:  &slackClient{client: slack.New(token)},
		channel: channel,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c, nil
}

func (c *Client) Report(ctx context.Context, repos []*github.Repo) error {
	c.logger.Info("started sending message")

	repoSum := len(repos)
	var count int // the number of messages already sent
	for {
		if len(repos) == 0 {
			break
		}

		var rest []*github.Repo
		if len(repos) > maxAttachments {
			repos, rest = repos[0:maxAttachments], repos[maxAttachments:]
		}

		attachments := make([]slack.Attachment, 0, len(repos))
		for i, repo := range repos {
			sent := maxAttachments * count
			title := fmt.Sprintf("%s (%d/%d)", attachmentTitle, i+1+sent, repoSum)
			attachments = append(attachments, makeSlackAttachment(title, repo))
		}

		if err := c.client.PostAttachmentsMessage(ctx, c.channel, attachments...); err != nil {
			c.logger.Error("failed to send message", zap.Error(err))
			return err
		}
		count++

		repos = rest
	}

	c.logger.Info("successfully sent message")
	return nil
}

func (s *slackClient) PostAttachmentsMessage(ctx context.Context, channel string, attachments ...slack.Attachment) error {
	if _, _, err := s.client.PostMessageContext(ctx, channel, slack.MsgOptionAttachments(attachments...)); err != nil {
		return err
	}
	return nil
}

func makeSlackAttachmentField(repo *github.Repo) slack.AttachmentField {
	var v string // v should contains URL, description and starred_at.
	if repo.Description == "" {
		v = repo.URL
	} else {
		v = fmt.Sprintf("%s\n%s", repo.URL, repo.Description)
	}
	v = fmt.Sprintf("%s\nStarred at %s", v, repo.StarredAt.Format(timeFormat))
	return slack.AttachmentField{
		Title: repo.String(),
		Value: v,
	}
}

func makeSlackAttachment(title string, repo *github.Repo) slack.Attachment {
	return slack.Attachment{
		Title:    title,
		Fields:   []slack.AttachmentField{makeSlackAttachmentField(repo)},
		Color:    githubColor,
		ThumbURL: repo.OwnerImageURL,
	}
}
