package rules

import (
	"go/parser"
	"testing"
)

func TestCheckSensitiveMessage(t *testing.T) {
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

			got := CheckSensitiveMessage(expr) != nil
			if got != tt.want {
				t.Fatalf("CheckSensitiveMessage(%q) = %v, want %v", tt.expr, got, tt.want)
			}
		})
	}
}
