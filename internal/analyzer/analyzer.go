package analyzer

import (
	"go/ast"

	"github.com/reservation-v/log-linter/internal/matchers"
	"github.com/reservation-v/log-linter/internal/rules"
	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "loglinter",
	Doc:  "checks log message rules in supported logger calls",
	URL:  "https://github.com/reservation-v/log-linter",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(node ast.Node) bool {
			call, ok := node.(*ast.CallExpr)
			if !ok {
				return true
			}

			logCall, ok := matchers.ExtractLogCall(pass, call)
			if !ok {
				return true
			}

			message := stringConstant(pass, logCall.Message)

			for _, diagnostic := range rules.Check(logCall, message) {
				pass.Report(*diagnostic)
			}

			return true
		})
	}

	return nil, nil
}
