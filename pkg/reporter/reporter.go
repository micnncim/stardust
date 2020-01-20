package reporter

import (
	"context"

	"github.com/micnncim/stardust/pkg/github"
)

// Reporter reports repositories information.
type Reporter interface {
	// Report shows repositories information on the report platform.
	Report(ctx context.Context, repos []*github.Repo) error
}
