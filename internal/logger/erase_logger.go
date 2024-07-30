package logger

import (
	"context"
	"log/slog"
)

func NewEraseLogger() *slog.Logger {
	return slog.New(NewDiscardHandler())
}

type EraseHandler struct{}

func NewDiscardHandler() *EraseHandler {
	return &EraseHandler{}
}

func (h *EraseHandler) Handle(_ context.Context, _ slog.Record) error {
	return nil
}

func (h *EraseHandler) WithAttrs(_ []slog.Attr) slog.Handler {
	return h
}

func (h *EraseHandler) WithGroup(_ string) slog.Handler {
	return h
}

func (h *EraseHandler) Enabled(_ context.Context, _ slog.Level) bool {
	return false
}
