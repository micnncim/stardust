package slack

import (
	"bytes"
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/nlopes/slack"
	"go.uber.org/zap"

	"github.com/micnncim/stardust/pkg/github"
)

type fakeSlackClient struct {
	buf *bytes.Buffer

	postAttachmentsMessageCalls int // count of PostAttachmentsMessage method calls
}

func (s *fakeSlackClient) PostAttachmentsMessage(_ context.Context, _ string, attachments ...slack.Attachment) error {
	s.postAttachmentsMessageCalls++
	for _, a := range attachments {
		for _, f := range a.Fields {
			s.buf.WriteString(fmt.Sprintf("%s\n%s\n\n", f.Title, f.Value))
		}
	}
	return nil
}

func TestClient_Report(t *testing.T) {
	fakeStarredAt := time.Date(2020, 4, 1, 0, 0, 0, 0, time.UTC)

	type args struct {
		repos []*github.Repo
	}
	tests := []struct {
		name                            string
		args                            args
		maxAttachments                  int
		want                            string
		wantPostAttachmentsMessageCalls int
		wantErr                         bool
	}{
		{
			name: "report repos less than max",
			args: args{
				repos: []*github.Repo{
					{Name: "repo1", Owner: "owner1", URL: "https://github.com/owner1/repo1", Description: "desc1", StarredAt: fakeStarredAt},
					{Name: "repo2", Owner: "owner2", URL: "https://github.com/owner2/repo2", Description: "desc2", StarredAt: fakeStarredAt},
				},
			},
			maxAttachments: maxAttachments,
			want: `
owner1/repo1
https://github.com/owner1/repo1
desc1
Starred at 2020/04/01 00:00

owner2/repo2
https://github.com/owner2/repo2
desc2
Starred at 2020/04/01 00:00

`[1:],
			wantPostAttachmentsMessageCalls: 1,
			wantErr:                         false,
		},
		{
			name: "report repos more than max",
			args: args{
				repos: []*github.Repo{
					{Name: "repo1", Owner: "owner1", URL: "https://github.com/owner1/repo1", Description: "desc1", StarredAt: fakeStarredAt},
					{Name: "repo2", Owner: "owner2", URL: "https://github.com/owner2/repo2", Description: "desc2", StarredAt: fakeStarredAt},
				},
			},
			maxAttachments: 1,
			want: `
owner1/repo1
https://github.com/owner1/repo1
desc1
Starred at 2020/04/01 00:00

owner2/repo2
https://github.com/owner2/repo2
desc2
Starred at 2020/04/01 00:00

`[1:],
			wantPostAttachmentsMessageCalls: 2,
			wantErr:                         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var orgMaxAttachments int
			orgMaxAttachments, maxAttachments = maxAttachments, tt.maxAttachments
			buf := &bytes.Buffer{}
			fake := &fakeSlackClient{buf: buf}
			c := &Client{
				client: fake,
				logger: zap.NewNop(),
			}
			if err := c.Report(context.Background(), tt.args.repos); (err != nil) != tt.wantErr {
				t.Errorf("err: %v\nwantErr: %v", err, tt.wantErr)
			}
			if buf.String() != tt.want {
				t.Errorf("got: %s\nwant: %s", buf.String(), tt.want)
			}
			if fake.postAttachmentsMessageCalls != tt.wantPostAttachmentsMessageCalls {
				t.Errorf("got: %d\nwant: %d", fake.postAttachmentsMessageCalls, tt.wantPostAttachmentsMessageCalls)
			}
			maxAttachments = orgMaxAttachments
		})
	}
}
