// Copyright 2019 The Berglas Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// This file contains the modified content of https://github.com/GoogleCloudPlatform/berglas/blob/master/pkg/auto/importer.go.

package berglas

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/GoogleCloudPlatform/berglas/pkg/berglas"
	"go.uber.org/zap"
)

type Client struct {
	berglasClient *berglas.Client
	logger        *zap.Logger
}

func NewClient(ctx context.Context, logger *zap.Logger) (*Client, error) {
	client, err := berglas.New(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize berglas client: %s", err)
	}
	return &Client{
		berglasClient: client,
		logger:        logger,
	}, nil
}

func (c *Client) Resolve(ctx context.Context) {
	for _, e := range os.Environ() {
		p := strings.SplitN(e, "=", 2)
		if len(p) < 2 {
			continue
		}

		k, v := p[0], p[1]
		if !berglas.IsReference(v) {
			continue
		}

		s, err := c.berglasClient.Resolve(ctx, v)
		if err != nil {
			c.logger.Warn("failed to parse env var", zap.String("key", k))
			continue
		}

		if err := os.Setenv(k, string(s)); err != nil {
			c.logger.Warn("failed to set env var", zap.String("key", k))
			continue
		}
	}
}
