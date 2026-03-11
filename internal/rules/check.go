package rules

import (
	"go/ast"
	"go/token"

	"github.com/reservation-v/log-linter/internal/config"
	"github.com/reservation-v/log-linter/internal/matchers"
	"golang.org/x/tools/go/analysis"
)

func Check(cfg config.Config, logCall matchers.LogCall, message string) []*analysis.Diagnostic {
	var diagnostics []*analysis.Diagnostic

	diagnostics = append(diagnostics, CheckMessage(cfg, logCall.Message, message)...)

	diagnostic := CheckSensitive(logCall, cfg.ExtraSensitiveKeywords)
	if cfg.Sensitive && diagnostic != nil {
		diagnostics = append(diagnostics, diagnostic)
	}

	return diagnostics
}

func CheckMessage(cfg config.Config, expr ast.Expr, message string) []*analysis.Diagnostic {
	var diagnostics []*analysis.Diagnostic
	pos := diagnosticPos(expr)
	suggestedFixes := messageSuggestedFixes(expr, message)

	if cfg.Lowercase && !CheckLowercase(message) {
		diagnostics = append(diagnostics, &analysis.Diagnostic{
			Pos:            pos,
			Message:        LowercaseMessageDiagnostic,
			SuggestedFixes: suggestedFixes,
		})
	}

	if cfg.English && !CheckEnglish(message) {
		diagnostics = append(diagnostics, &analysis.Diagnostic{
			Pos:     pos,
			Message: EnglishMessageDiagnostic,
		})
	}

	if cfg.Symbols && !CheckSymbols(message) {
		diagnostics = append(diagnostics, &analysis.Diagnostic{
			Pos:            pos,
			Message:        SymbolsMessageDiagnostic,
			SuggestedFixes: suggestedFixes,
		})
	}

	return diagnostics
}

func diagnosticPos(expr ast.Expr) token.Pos {
	if expr == nil {
		return token.NoPos
	}

	return expr.Pos()
}
