package rules

import (
	"go/token"

	"github.com/reservation-v/log-linter/internal/matchers"
	"golang.org/x/tools/go/analysis"
)

func Check(logCall matchers.LogCall, message string) []*analysis.Diagnostic {
	var diagnostics []*analysis.Diagnostic

	diagnostics = append(diagnostics, CheckMessage(logCall.Message.Pos(), message)...)

	diagnostic := CheckSensitive(logCall)
	if diagnostic != nil {
		diagnostics = append(diagnostics, diagnostic)
	}

	return diagnostics
}

func CheckMessage(pos token.Pos, message string) []*analysis.Diagnostic {
	var diagnostics []*analysis.Diagnostic

	if !CheckLowercase(message) {
		diagnostics = append(diagnostics, &analysis.Diagnostic{
			Pos:     pos,
			Message: LowercaseMessageDiagnostic,
		})
	}

	if !CheckEnglish(message) {
		diagnostics = append(diagnostics, &analysis.Diagnostic{
			Pos:     pos,
			Message: EnglishMessageDiagnostic,
		})
	}

	if !CheckSymbols(message) {
		diagnostics = append(diagnostics, &analysis.Diagnostic{
			Pos:     pos,
			Message: SymbolsMessageDiagnostic,
		})
	}

	return diagnostics
}
