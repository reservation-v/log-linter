package rules

import (
	"go/ast"
	"go/parser"
	"testing"

	"github.com/reservation-v/log-linter/internal/config"
	"github.com/reservation-v/log-linter/internal/matchers"
)

func TestCheckSensitiveMessage(t *testing.T) {
	cfg := config.Default()
	keywords := sensitiveKeywords(cfg.ExtraSensitiveKeywords)

	tests := []struct {
		name string
		expr string
		want bool
	}{
		{
			name: "plain safe string",
			expr: `"starting server"`,
			want: false,
		},
		{
			name: "string literal with password",
			expr: `"user password updated"`,
			want: true,
		},
		{
			name: "binary concat with token ident",
			expr: `"token: " + token`,
			want: true,
		},
		{
			name: "identifier password",
			expr: `password`,
			want: true,
		},
		{
			name: "selector token",
			expr: `req.Token`,
			want: true,
		},
		{
			name: "parenthesized concat",
			expr: `("api_key=" + apiKey)`,
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := parser.ParseExpr(tt.expr)
			if err != nil {
				t.Fatalf("ParseExpr(%q) error = %v", tt.expr, err)
			}

			got := CheckSensitiveMessage(expr, keywords) != nil
			if got != tt.want {
				t.Fatalf("CheckSensitiveMessage(%q) = %v, want %v", tt.expr, got, tt.want)
			}
		})
	}
}

func TestCheckSensitiveMessageCustomKeyword(t *testing.T) {
	cfg := config.Default()
	cfg.ExtraSensitiveKeywords = []string{"session_id"}

	expr, err := parser.ParseExpr(`"session id leaked"`)
	if err != nil {
		t.Fatalf("ParseExpr() error = %v", err)
	}

	if CheckSensitiveMessage(expr, sensitiveKeywords(cfg.ExtraSensitiveKeywords)) == nil {
		t.Fatal("CheckSensitiveMessage() = nil, want diagnostic for custom keyword")
	}
}

func TestCheckSensitiveFields(t *testing.T) {
	cfg := config.Default()
	keywords := sensitiveKeywords(cfg.ExtraSensitiveKeywords)

	tests := []struct {
		name string
		kind matchers.LoggerKind
		expr string
		want bool
	}{
		{
			name: "slog sensitive key",
			kind: matchers.LoggerSlog,
			expr: `slog.Info("request done", "token", token)`,
			want: true,
		},
		{
			name: "slog sensitive value ident",
			kind: matchers.LoggerSlog,
			expr: `slog.Info("request done", "user", apiKey)`,
			want: true,
		},
		{
			name: "slog safe fields",
			kind: matchers.LoggerSlog,
			expr: `slog.Info("request done", "user", userID)`,
			want: false,
		},
		{
			name: "zap sensitive key",
			kind: matchers.LoggerZap,
			expr: `logger.Info("request done", zap.String("token", value))`,
			want: true,
		},
		{
			name: "zap sensitive value",
			kind: matchers.LoggerZap,
			expr: `logger.Info("request done", zap.Any("user", password))`,
			want: true,
		},
		{
			name: "zap safe field",
			kind: matchers.LoggerZap,
			expr: `logger.Info("request done", zap.String("user", userID))`,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			call := parseCallExpr(t, tt.expr)
			logCall := matchers.LogCall{
				Kind:         tt.kind,
				Call:         call,
				Message:      call.Args[0],
				MessageIndex: 0,
			}

			got := CheckSensitiveFields(logCall, keywords) != nil
			if got != tt.want {
				t.Fatalf("CheckSensitiveFields(%q) = %v, want %v", tt.expr, got, tt.want)
			}
		})
	}
}

func TestCheckSensitiveFieldsCustomKeyword(t *testing.T) {
	cfg := config.Default()
	cfg.ExtraSensitiveKeywords = []string{"client_secret_id"}

	call := parseCallExpr(t, `slog.Info("request done", "client-secret-id", value)`)
	logCall := matchers.LogCall{
		Kind:         matchers.LoggerSlog,
		Call:         call,
		Message:      call.Args[0],
		MessageIndex: 0,
	}

	if CheckSensitiveFields(logCall, sensitiveKeywords(cfg.ExtraSensitiveKeywords)) == nil {
		t.Fatal("CheckSensitiveFields() = nil, want diagnostic for custom keyword")
	}
}

func parseCallExpr(t *testing.T, expr string) *ast.CallExpr {
	t.Helper()

	parsed, err := parser.ParseExpr(expr)
	if err != nil {
		t.Fatalf("ParseExpr(%q) error = %v", expr, err)
	}

	call, ok := parsed.(*ast.CallExpr)
	if !ok {
		t.Fatalf("ParseExpr(%q) returned %T, want *ast.CallExpr", expr, parsed)
	}

	return call
}
