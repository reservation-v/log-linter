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
	return CheckSensitiveMessage(logCall.Message)
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
