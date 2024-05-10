package client

import (
	"log/slog"

	"go.tradeforge.dev/oxr/internal/client"
)

type Option func(*client.Client)

func WithLogger(logger *slog.Logger) Option {
	return func(c *client.Client) {
		c.SetLogger(logger)
	}
}
