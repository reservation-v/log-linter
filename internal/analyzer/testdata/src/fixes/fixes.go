package fixes

import (
	"log/slog"

	"go.uber.org/zap"
)

func check(logger *zap.Logger) {
	slog.Info("Starting server")     // want "log message must start with a lowercase letter"
	slog.Info("server started! 🚀")   // want "log message must not contain special symbols or emoji"
	slog.Info("Starting_request!")   // want "log message must start with a lowercase letter" "log message must not contain special symbols or emoji"
	slog.Info("8080 started")        // want "log message must start with a lowercase letter"
	slog.Info("запуск сервера")      // want "log message must contain English text only"
	slog.Info("token: " + password)  // want "log message must not contain potentially sensitive data"
	logger.Info("Failed request!?")  // want "log message must start with a lowercase letter" "log message must not contain special symbols or emoji"
	logger.Info("request_id failed") // want "log message must not contain special symbols or emoji"
}

var password = "secret"
