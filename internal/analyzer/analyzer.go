package analyzer

import (
	"go/ast"

	"github.com/reservation-v/log-linter/internal/matchers"
	"golang.org/x/tools/go/analysis"
)

var Analyzer = &analysis.Analyzer{
	Name: "loglinter",
	Doc:  "checks log message rules in supported logger calls",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	var calls []matchers.LogCall

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

			calls = append(calls, logCall)
			return true
		})
	}

	return calls, nil
}
