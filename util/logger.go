package util

import (
	"context"
	"log/slog"
)

type NilHandler struct{}

// Enabled returns always false for [NilHandler].
func (h *NilHandler) Enabled(_ context.Context, level slog.Level) bool {
	return false
}

// WithAttrs returns a new [NilHandler] whose attributes consists
// of h's attributes followed by attrs.
// [NilHandler] returns a new empty handler.
func (h *NilHandler) WithAttrs(_ []slog.Attr) slog.Handler {
	return &NilHandler{}
}

func (h *NilHandler) WithGroup(_ string) slog.Handler {
	return &NilHandler{}
}

// Handle does nothing for [NilHandler].
func (h *NilHandler) Handle(_ context.Context, r slog.Record) error {
	return nil
}
