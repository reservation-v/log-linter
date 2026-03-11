package rules

import (
	"go/ast"
	"go/token"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/tools/go/analysis"
)

const normalizeMessageFix = "Normalize log message"

func messageSuggestedFixes(expr ast.Expr, message string) []analysis.SuggestedFix {
	literal := directStringLiteral(expr)
	if literal == nil {
		return nil
	}

	normalized := normalizeMessageLiteral(message)
	if normalized == message || strings.TrimSpace(normalized) == "" {
		return nil
	}

	return []analysis.SuggestedFix{{
		Message: normalizeMessageFix,
		TextEdits: []analysis.TextEdit{{
			Pos:     literal.Pos(),
			End:     literal.End(),
			NewText: []byte(strconv.Quote(normalized)),
		}},
	}}
}

func directStringLiteral(expr ast.Expr) *ast.BasicLit {
	switch value := expr.(type) {
	case *ast.BasicLit:
		if value.Kind == token.STRING {
			return value
		}
	case *ast.ParenExpr:
		return directStringLiteral(value.X)
	}

	return nil
}

func normalizeMessageLiteral(message string) string {
	normalized := lowercaseMessageLiteral(message)
	normalized = stripSpecialSymbols(normalized)
	return normalized
}

func lowercaseMessageLiteral(message string) string {
	if message == "" {
		return message
	}

	r, size := utf8.DecodeRuneInString(message)
	if !unicode.IsUpper(r) {
		return message
	}

	lowered := unicode.ToLower(r)
	if lowered == r {
		return message
	}

	return string(lowered) + message[size:]
}

func stripSpecialSymbols(message string) string {
	var builder strings.Builder
	changed := false

	for _, r := range message {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsSpace(r) {
			builder.WriteRune(r)
			continue
		}

		builder.WriteByte(' ')
		changed = true
	}

	if !changed {
		return message
	}

	return strings.Join(strings.Fields(builder.String()), " ")
}
