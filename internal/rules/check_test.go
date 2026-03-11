package rules

import (
	"go/parser"
	"testing"

	"github.com/reservation-v/log-linter/internal/config"
	"github.com/reservation-v/log-linter/internal/matchers"
)

func TestCheckMessageDisableEnglish(t *testing.T) {
	cfg := config.Default()
	cfg.English = false

	expr, err := parser.ParseExpr(`"запуск сервера"`)
	if err != nil {
		t.Fatalf("ParseExpr() error = %v", err)
	}

	diagnostics := CheckMessage(cfg, expr, "запуск сервера")
	if len(diagnostics) != 0 {
		t.Fatalf("CheckMessage returned %d diagnostics, want 0", len(diagnostics))
	}
}

func TestCheckDisableSensitive(t *testing.T) {
	cfg := config.Default()
	cfg.Sensitive = false

	call := parseCallExpr(t, `slog.Info("request done", "token", password)`)

	logCall := matchers.LogCall{
		Kind:         matchers.LoggerSlog,
		Call:         call,
		Message:      call.Args[0],
		MessageIndex: 0,
	}

	diagnostics := Check(cfg, logCall, "request done")
	if len(diagnostics) != 0 {
		t.Fatalf("Check returned %d diagnostics, want 0", len(diagnostics))
	}
}
