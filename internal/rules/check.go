package rules

import (
	"go/token"

	"github.com/reservation-v/log-linter/internal/config"
	"github.com/reservation-v/log-linter/internal/matchers"
	"golang.org/x/tools/go/analysis"
)

func Check(cfg config.Config, logCall matchers.LogCall, message string) []*analysis.Diagnostic {
	var diagnostics []*analysis.Diagnostic

	diagnostics = append(diagnostics, CheckMessage(cfg, logCall.Message.Pos(), message)...)

	diagnostic := CheckSensitive(logCall)
	if cfg.Sensitive && diagnostic != nil {
		diagnostics = append(diagnostics, diagnostic)
	}

	return diagnostics
}

func CheckMessage(cfg config.Config, pos token.Pos, message string) []*analysis.Diagnostic {
	var diagnostics []*analysis.Diagnostic

	if cfg.Lowercase && !CheckLowercase(message) {
		diagnostics = append(diagnostics, &analysis.Diagnostic{
			Pos:     pos,
			Message: LowercaseMessageDiagnostic,
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
			Pos:     pos,
			Message: SymbolsMessageDiagnostic,
		})
	}

	return diagnostics
}
