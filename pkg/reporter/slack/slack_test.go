package slack

import (
	"bytes"
	"context"
	"fmt"
	"testing"

	"github.com/nlopes/slack"
	"go.uber.org/zap"

	"github.com/micnncim/stardust/pkg/github"
)

type fakeSlackClient struct {
	buf *bytes.Buffer
}

func (s *fakeSlackClient) PostAttachmentsMessage(_ context.Context, _ string, attachments ...slack.Attachment) error {
	for _, a := range attachments {
		for _, f := range a.Fields {
			s.buf.WriteString(fmt.Sprintf("%s\n", f.Title))
		}
	}
	return nil
}

func TestClient_Report(t *testing.T) {
	type args struct {
		repos []*github.Repo
	}
	tests := []struct {
		name           string
		args           args
		maxAttachments int
		want           string
		wantErr        bool
	}{
		{
			name: "report repos less than max",
			args: args{
				repos: []*github.Repo{
					{Name: "repo1", Owner: "owner1"},
					{Name: "repo2", Owner: "owner2"},
				},
			},
			maxAttachments: maxAttachments,
			want: `
owner1/repo1
owner2/repo2
`[1:],
			wantErr: false,
		},
		{
			name: "report repos more than max",
			args: args{
				repos: []*github.Repo{
					{Name: "repo1", Owner: "owner1"},
					{Name: "repo2", Owner: "owner2"},
				},
			},
			maxAttachments: 1,
			want: `
owner1/repo1
owner2/repo2
`[1:],
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			c := &Client{
				client: &fakeSlackClient{buf: buf},
				logger: zap.NewNop(),
			}
			if err := c.Report(context.Background(), tt.args.repos); (err != nil) != tt.wantErr {
				t.Errorf("err: %v\nwantErr: %v", err, tt.wantErr)
			}
			if buf.String() != tt.want {
				t.Errorf("got: %s\nwant: %s", buf.String(), tt.want)
			}
		})
	}
}
