package rules

import (
	"go/token"

	"github.com/reservation-v/log-linter/internal/matchers"
	"golang.org/x/tools/go/analysis"
)

func Check(logCall matchers.LogCall, message string) []analysis.Diagnostic {
	var diagnostics []analysis.Diagnostic

	diagnostics = append(diagnostics, CheckMessage(logCall.Message.Pos(), message)...)

	return diagnostics
}

func CheckMessage(pos token.Pos, message string) []analysis.Diagnostic {
	var diagnostics []analysis.Diagnostic

	if !CheckLowercase(message) {
		diagnostics = append(diagnostics, analysis.Diagnostic{
			Pos:     pos,
			Message: LowercaseMessageDiagnostic,
		})
	}

	return diagnostics
}
