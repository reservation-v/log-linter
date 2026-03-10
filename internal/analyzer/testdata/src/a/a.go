package a

import (
	"context"
	"log/slog"

	"go.uber.org/zap"
)

// remove any "want" and run "go test ./internal/analyzer/ -run TestAnalyzer" in project root
func checkSlog(logger *slog.Logger) {
	slog.Info("Starting server") // want "log message must start with a lowercase letter"
	slog.Info("запуск сервера")  // want "log message must contain English text only"
	slog.Info("server started!") // want "log message must not contain special symbols or emoji"
	slog.Warn("disk almost full")
	slog.Error("request failed")
	slog.Debug("debug path")

	password := "secret"
	slog.Info("user password: " + password)      // want "log message must not contain potentially sensitive data"
	slog.Info("request done", "token", password) // want "log message must not contain potentially sensitive data"

	ctx := context.Background()
	slog.InfoContext(ctx, "Started request")                    // want "log message must start with a lowercase letter"
	slog.WarnContext(ctx, "запуск warning path")                // want "log message must contain English text only"
	slog.ErrorContext(ctx, "request failed...")                 // want "log message must not contain special symbols or emoji"
	slog.DebugContext(ctx, "secret: "+password)                 // want "log message must not contain potentially sensitive data"
	logger.InfoContext(ctx, "Started logger ctx")               // want "log message must start with a lowercase letter"
	logger.WarnContext(ctx, "ошибка logger ctx")                // want "log message must contain English text only"
	logger.ErrorContext(ctx, "logger ctx failed!")              // want "log message must not contain special symbols or emoji"
	logger.DebugContext(ctx, "request done", "token", password) // want "log message must not contain potentially sensitive data"

	msg := "Starting server"
	slog.Info(msg)

	logger.Info("request done")
}

func checkZap(logger *zap.Logger) {
	logger.Info("Failed request")    // want "log message must start with a lowercase letter"
	logger.Info("ошибка запроса")    // want "log message must contain English text only"
	logger.Info("request failed...") // want "log message must not contain special symbols or emoji"
	logger.Debug("debug request")
	logger.Warn("warn request")
	logger.Error("error request")

	token := "secret"
	logger.Info("token: " + token)                          // want "log message must not contain potentially sensitive data"
	logger.Info("request done", zap.String("token", token)) // want "log message must not contain potentially sensitive data"
	logger.Debug("Password leaked")                         // want "log message must start with a lowercase letter" "log message must not contain potentially sensitive data"
	logger.Warn("request done", zap.Any("password", token)) // want "log message must not contain potentially sensitive data"

	msg := "Failed request"
	logger.Info(msg)
	logger.Info("request done", zap.String("user", "123"))
}
