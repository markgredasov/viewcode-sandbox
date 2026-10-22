package logger

import (
	"context"
	"log/slog"
	"os"
	"time"
)

type key string

var loggerKey key = "log"

type Logger struct {
	*slog.Logger
}

func FromContext(ctx context.Context) *Logger {
	log, ok := ctx.Value(loggerKey).(*Logger)
	if !ok {
		panic("no logger in context")
	}

	return log
}

func NewLogger(cfg Config) (*Logger, error) {
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: cfg.Level,
	}))

	return &Logger{
		Logger: log,
	}, nil
}

func (l *Logger) With(fields ...any) *Logger {
	return &Logger{
		Logger: l.Logger.With(fields...),
	}
}

func (l *Logger) InContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, loggerKey, l)
}

func Err(value error) slog.Attr {
	return slog.String("error", value.Error())
}

func String(key, value string) slog.Attr {
	return slog.String(key, value)
}

func Int(key string, value int) slog.Attr {
	return slog.Int(key, value)
}

func Duration(key string, value time.Duration) slog.Attr {
	return slog.String(key, value.String())
}
