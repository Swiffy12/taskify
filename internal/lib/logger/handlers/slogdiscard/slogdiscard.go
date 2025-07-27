package slogdiscard

import (
	"context"
	"log/slog"
)

func NewDiscardLogger() *slog.Logger {
	return slog.New(NewDiscardHandler())
}

type DiscradHandler struct{}

func NewDiscardHandler() *DiscradHandler {
	return &DiscradHandler{}
}

func (h *DiscradHandler) Enabled(_ context.Context, _ slog.Level) bool {
	return false
}

func (h *DiscradHandler) Handle(_ context.Context, _ slog.Record) error {
	return nil
}

func (h *DiscradHandler) WithAttrs(_ []slog.Attr) slog.Handler {
	return h
}

func (h *DiscradHandler) WithGroup(_ string) slog.Handler {
	return h
}
