package main

import (
	"context"
	"log/slog"
)

func main() {
	logger := slog.Default()
	ctx := context.Background()
	password := "secret"
	token := "token-value"
	msg := "Starting from variable"

	slog.Info("Starting server")
	slog.Warn("запуск warning path")
	slog.Error("server failed...")
	slog.Debug("user password: " + password)

	slog.InfoContext(ctx, "Started request")
	slog.WarnContext(ctx, "ошибка request path")
	slog.ErrorContext(ctx, "request failed!")
	slog.DebugContext(ctx, "request done", "token", token)

	logger.Info("Started logger")
	logger.Warn("logger warning...")
	logger.Error("token: " + token)
	logger.DebugContext(ctx, "request done", "api_key", password)

	slog.Info("request done")
	slog.Info(msg)
	logger.InfoContext(ctx, "request completed")
}
