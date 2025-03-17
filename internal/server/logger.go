package server

import (
	"errors"
	"log/slog"
	"os"
)

func SetupLogger(env string) (*slog.Logger, error) {
	var logger *slog.Logger
	switch env {
	case EnvLocal:
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case EnvDev:
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case EnvProd:
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	default:
		return nil, errors.New("invalid env variable")
	}
	return logger, nil
}
