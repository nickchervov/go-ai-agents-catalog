package logger

import (
	"log/slog"
	"os"
)

func NewLogger(format string) *slog.Logger {
	switch format {
	case "json":
		return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{})).With("service", "ai-agents")
	default:
		return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{})).With("service", "ai-agents")
	}
}
