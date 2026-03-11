package rules

import (
	"go/parser"
	"testing"

	"github.com/reservation-v/log-linter/internal/config"
)

func TestNormalizeMessageLiteral(t *testing.T) {
	tests := []struct {
		name    string
		message string
		want    string
	}{
		{
			name:    "lowercase only",
			message: "Starting server",
			want:    "starting server",
		},
		{
			name:    "symbols and emoji",
			message: "server started! 🚀",
			want:    "server started",
		},
		{
			name:    "combined lowercase and symbols",
			message: "Starting_request!",
			want:    "starting request",
		},
		{
			name:    "digit prefix stays unchanged",
			message: "8080 started",
			want:    "8080 started",
		},
		{
			name:    "symbols collapse to single spaces",
			message: "warning: something... odd",
			want:    "warning something odd",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := normalizeMessageLiteral(tt.message)
			if got != tt.want {
				t.Fatalf("normalizeMessageLiteral(%q) = %q, want %q", tt.message, got, tt.want)
			}
		})
	}
}

func TestMessageSuggestedFixes(t *testing.T) {
	tests := []struct {
		name    string
		expr    string
		message string
		want    string
	}{
		{
			name:    "basic literal rewrite",
			expr:    `"Starting_request!"`,
			message: "Starting_request!",
			want:    `"starting request"`,
		},
		{
			name:    "dynamic expression has no fix",
			expr:    `"token: " + token`,
			message: "token: secret",
			want:    "",
		},
		{
			name:    "unchanged message has no fix",
			expr:    `"request done"`,
			message: "request done",
			want:    "",
		},
		{
			name:    "unfixable empty result has no fix",
			expr:    `"!!!"`,
			message: "!!!",
			want:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := parser.ParseExpr(tt.expr)
			if err != nil {
				t.Fatalf("ParseExpr(%q) error = %v", tt.expr, err)
			}

			fixes := messageSuggestedFixes(expr, tt.message)
			if tt.want == "" {
				if len(fixes) != 0 {
					t.Fatalf("messageSuggestedFixes(%q) returned %d fixes, want 0", tt.expr, len(fixes))
				}
				return
			}

			if len(fixes) != 1 {
				t.Fatalf("messageSuggestedFixes(%q) returned %d fixes, want 1", tt.expr, len(fixes))
			}

			fix := fixes[0]
			if fix.Message != normalizeMessageFix {
				t.Fatalf("fix.Message = %q, want %q", fix.Message, normalizeMessageFix)
			}
			if len(fix.TextEdits) != 1 {
				t.Fatalf("len(fix.TextEdits) = %d, want 1", len(fix.TextEdits))
			}
			if got := string(fix.TextEdits[0].NewText); got != tt.want {
				t.Fatalf("fix.TextEdits[0].NewText = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestCheckMessageSuggestedFixesSharedAcrossDiagnostics(t *testing.T) {
	cfg := config.Default()
	expr, err := parser.ParseExpr(`"Starting_request!"`)
	if err != nil {
		t.Fatalf("ParseExpr() error = %v", err)
	}

	diagnostics := CheckMessage(cfg, expr, "Starting_request!")
	if len(diagnostics) != 2 {
		t.Fatalf("CheckMessage returned %d diagnostics, want 2", len(diagnostics))
	}

	first := diagnostics[0].SuggestedFixes
	second := diagnostics[1].SuggestedFixes
	if len(first) != 1 || len(second) != 1 {
		t.Fatalf("suggested fixes lengths = %d and %d, want 1 and 1", len(first), len(second))
	}
	if string(first[0].TextEdits[0].NewText) != string(second[0].TextEdits[0].NewText) {
		t.Fatalf("suggested fixes differ: %q vs %q", first[0].TextEdits[0].NewText, second[0].TextEdits[0].NewText)
	}
}
