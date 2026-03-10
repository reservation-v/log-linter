package rules

import (
	"go/ast"
	"go/token"
	"strconv"
	"strings"

	"github.com/reservation-v/log-linter/internal/matchers"
	"golang.org/x/tools/go/analysis"
)

const SensitiveMessageDiagnostic = "log message must not contain potentially sensitive data"

var sensitiveKeywords = []string{
	"password",
	"passwd",
	"pwd",
	"token",
	"secret",
	"apikey",
}

func CheckSensitive(logCall matchers.LogCall) *analysis.Diagnostic {
	if diagnostic := CheckSensitiveMessage(logCall.Message); diagnostic != nil {
		return diagnostic
	}

	return CheckSensitiveFields(logCall)
}

func CheckSensitiveMessage(expr ast.Expr) *analysis.Diagnostic {
	if !exprContainsSensitiveData(expr) {
		return nil
	}

	return &analysis.Diagnostic{
		Pos:     expr.Pos(),
		Message: SensitiveMessageDiagnostic,
	}
}

func CheckSensitiveFields(logCall matchers.LogCall) *analysis.Diagnostic {
	args := logCall.Call.Args
	if len(args) <= logCall.MessageIndex+1 {
		return nil
	}

	switch logCall.Kind {
	case matchers.LoggerSlog:
		return checkSensitiveSlogFields(args[logCall.MessageIndex+1:])
	case matchers.LoggerZap:
		return checkSensitiveZapFields(args[logCall.MessageIndex+1:])
	default:
		return nil
	}
}

func checkSensitiveSlogFields(args []ast.Expr) *analysis.Diagnostic {
	for _, arg := range args {
		if !exprContainsSensitiveData(arg) {
			continue
		}

		return &analysis.Diagnostic{
			Pos:     arg.Pos(),
			Message: SensitiveMessageDiagnostic,
		}
	}

	return nil
}

func checkSensitiveZapFields(args []ast.Expr) *analysis.Diagnostic {
	for _, arg := range args {
		fieldCall, ok := arg.(*ast.CallExpr)
		if !ok {
			continue
		}

		for _, fieldArg := range fieldCall.Args {
			if !exprContainsSensitiveData(fieldArg) {
				continue
			}

			return &analysis.Diagnostic{
				Pos:     fieldArg.Pos(),
				Message: SensitiveMessageDiagnostic,
			}
		}
	}

	return nil
}

func exprContainsSensitiveData(expr ast.Expr) bool {
	switch e := expr.(type) {
	case *ast.BasicLit:
		return basicLitContainsSensitiveData(e)
	// checks 2 parts of binary expr
	// f.e. "token: " + token
	// it checks parts separated by +
	case *ast.BinaryExpr:
		if e.Op != token.ADD {
			return false
		}
		return exprContainsSensitiveData(e.X) || exprContainsSensitiveData(e.Y)
	case *ast.Ident:
		return containsSensitiveKeyword(e.Name)
	case *ast.SelectorExpr:
		return containsSensitiveKeyword(e.Sel.Name)
	case *ast.ParenExpr:
		return exprContainsSensitiveData(e.X)
	default:
		return false
	}
}

func basicLitContainsSensitiveData(lit *ast.BasicLit) bool {
	if lit.Kind != token.STRING {
		return false
	}

	value, err := strconv.Unquote(lit.Value)
	if err != nil {
		return containsSensitiveKeyword(lit.Value)
	}

	return containsSensitiveKeyword(value)
}

func containsSensitiveKeyword(value string) bool {
	normalized := normalizeSensitiveText(value)
	if normalized == "" {
		return false
	}

	for _, keyword := range sensitiveKeywords {
		if strings.Contains(normalized, keyword) {
			return true
		}
	}

	return false
}

// makes string lowercase with nums
func normalizeSensitiveText(value string) string {
	var builder strings.Builder

	for _, r := range value {
		switch {
		case 'a' <= r && r <= 'z':
			builder.WriteRune(r)
		case 'A' <= r && r <= 'Z':
			builder.WriteRune(r + ('a' - 'A'))
		case '0' <= r && r <= '9':
			builder.WriteRune(r)
		}
	}

	return builder.String()
}
