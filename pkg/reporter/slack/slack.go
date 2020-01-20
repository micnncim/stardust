package slack

import (
	"context"
	"errors"
	"fmt"

	"github.com/micnncim/stardust/pkg/github"
	"github.com/micnncim/stardust/pkg/reporter"

	"github.com/nlopes/slack"
	"go.uber.org/zap"
)

const (
	timeFormat = "2006/01/02 15:04"

	githubColor = "#24292E"

	attachmentTitle = "GitHub Star Report"

	// Slack API limits the number of attachments.
	// https://api.slack.com/methods/chat.postMessage
	maxAttachments = 100
)

type Client struct {
	client    slackClientInterface
	channelID string
	logger    *zap.Logger
}

var _ reporter.Reporter = (*Client)(nil)

type slackClientInterface interface {
	PostAttachmentsMessage(ctx context.Context, channelID string, attachments ...slack.Attachment) error
}

type slackClient struct {
	client *slack.Client
}

var _ slackClientInterface = (*slackClient)(nil)

func NewClient(token, channelID string, logger *zap.Logger) (reporter.Reporter, error) {
	if token == "" {
		return nil, errors.New("missing slack access token")
	}
	if channelID == "" {
		return nil, errors.New("missing slack channel id")
	}

	return &Client{
		client:    &slackClient{client: slack.New(token)},
		channelID: channelID,
		logger:    logger.Named("slack"),
	}, nil
}

func (c *Client) Report(ctx context.Context, repos []*github.Repo) error {
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

		if err := c.client.PostAttachmentsMessage(ctx, c.channelID, attachments...); err != nil {
			c.logger.Error("failed to send message", zap.Error(err))
			return err
		}
		count++

		repos = rest
	}

	c.logger.Info("successfully sent message")
	return nil
}

func (s *slackClient) PostAttachmentsMessage(ctx context.Context, channelID string, attachments ...slack.Attachment) error {
	if _, _, err := s.client.PostMessageContext(ctx, channelID, slack.MsgOptionAttachments(attachments...)); err != nil {
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
