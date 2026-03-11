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

var defaultSensitiveKeywords = []string{
	"password",
	"passwd",
	"pwd",
	"token",
	"secret",
	"apikey",
}

func CheckSensitive(logCall matchers.LogCall, extraKeywords []string) *analysis.Diagnostic {
	keywords := sensitiveKeywords(extraKeywords)

	if diagnostic := CheckSensitiveMessage(logCall.Message, keywords); diagnostic != nil {
		return diagnostic
	}

	return CheckSensitiveFields(logCall, keywords)
}

func CheckSensitiveMessage(expr ast.Expr, keywords []string) *analysis.Diagnostic {
	if !exprContainsSensitiveData(expr, keywords) {
		return nil
	}

	return &analysis.Diagnostic{
		Pos:     expr.Pos(),
		Message: SensitiveMessageDiagnostic,
	}
}

func CheckSensitiveFields(logCall matchers.LogCall, keywords []string) *analysis.Diagnostic {
	args := logCall.Call.Args
	if len(args) <= logCall.MessageIndex+1 {
		return nil
	}

	switch logCall.Kind {
	case matchers.LoggerSlog:
		return checkSensitiveSlogFields(args[logCall.MessageIndex+1:], keywords)
	case matchers.LoggerZap:
		return checkSensitiveZapFields(args[logCall.MessageIndex+1:], keywords)
	default:
		return nil
	}
}

func checkSensitiveSlogFields(args []ast.Expr, keywords []string) *analysis.Diagnostic {
	for _, arg := range args {
		if !exprContainsSensitiveData(arg, keywords) {
			continue
		}

		return &analysis.Diagnostic{
			Pos:     arg.Pos(),
			Message: SensitiveMessageDiagnostic,
		}
	}

	return nil
}

func checkSensitiveZapFields(args []ast.Expr, keywords []string) *analysis.Diagnostic {
	for _, arg := range args {
		fieldCall, ok := arg.(*ast.CallExpr)
		if !ok {
			continue
		}

		for _, fieldArg := range fieldCall.Args {
			if !exprContainsSensitiveData(fieldArg, keywords) {
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

func exprContainsSensitiveData(expr ast.Expr, keywords []string) bool {
	switch e := expr.(type) {
	case *ast.BasicLit:
		return basicLitContainsSensitiveData(e, keywords)
	case *ast.BinaryExpr:
		if e.Op != token.ADD {
			return false
		}
		return exprContainsSensitiveData(e.X, keywords) || exprContainsSensitiveData(e.Y, keywords)
	case *ast.Ident:
		return containsSensitiveKeyword(e.Name, keywords)
	case *ast.SelectorExpr:
		return containsSensitiveKeyword(e.Sel.Name, keywords)
	case *ast.ParenExpr:
		return exprContainsSensitiveData(e.X, keywords)
	default:
		return false
	}
}

func basicLitContainsSensitiveData(lit *ast.BasicLit, keywords []string) bool {
	if lit.Kind != token.STRING {
		return false
	}

	value, err := strconv.Unquote(lit.Value)
	if err != nil {
		return containsSensitiveKeyword(lit.Value, keywords)
	}

	return containsSensitiveKeyword(value, keywords)
}

func containsSensitiveKeyword(value string, keywords []string) bool {
	normalized := normalizeSensitiveText(value)
	if normalized == "" {
		return false
	}

	for _, keyword := range keywords {
		if strings.Contains(normalized, keyword) {
			return true
		}
	}

	return false
}

func sensitiveKeywords(extraKeywords []string) []string {
	keywords := make([]string, 0, len(defaultSensitiveKeywords)+len(extraKeywords))
	seen := make(map[string]struct{}, len(defaultSensitiveKeywords)+len(extraKeywords))

	appendNormalizedKeywords(&keywords, seen, defaultSensitiveKeywords)
	appendNormalizedKeywords(&keywords, seen, extraKeywords)

	return keywords
}

func appendNormalizedKeywords(target *[]string, seen map[string]struct{}, rawKeywords []string) {
	for _, raw := range rawKeywords {
		keyword := normalizeSensitiveText(raw)
		if keyword == "" {
			continue
		}
		if _, ok := seen[keyword]; ok {
			continue
		}

		seen[keyword] = struct{}{}
		*target = append(*target, keyword)
	}
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
