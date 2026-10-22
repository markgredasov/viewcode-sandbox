package logger

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
)

var defaultLevel = "info"

type Config struct {
	Level slog.Level
}

func NewConfig() (Config, error) {
	level, err := parseLevel(getEnv("LOGGER_LEVEL", defaultLevel))
	if err != nil {
		return Config{}, nil
	}

	return Config{
		Level: level,
	}, nil
}

func NewConfigMust() Config {
	config, err := NewConfig()
	if err != nil {
		err = fmt.Errorf("get logger config: %w", err)
		panic(err)
	}

	return config
}

func getEnv(key, fallback string) string {
	value, ok := os.LookupEnv(key)
	if !ok || strings.TrimSpace(value) == "" {
		return fallback
	}

	return value
}

func parseLevel(level string) (slog.Level, error) {
	level = strings.ToLower(strings.TrimSpace(level))

	switch level {
	case "debug":
		return slog.LevelDebug, nil
	case "info":
		return slog.LevelInfo, nil
	case "warn":
		return slog.LevelWarn, nil
	case "error":
		return slog.LevelError, nil
	default:
		return 0, fmt.Errorf("unsupported log level: %s", level)
	}
}
