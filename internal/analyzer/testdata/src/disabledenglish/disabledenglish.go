package disabledenglish

import "log/slog"

func check() {
	slog.Info("запуск сервера")
	slog.Info("Starting server") // want "log message must start with a lowercase letter"
	slog.Info("request failed!") // want "log message must not contain special symbols or emoji"
}
