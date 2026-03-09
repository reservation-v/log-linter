package analyzer

import "golang.org/x/tools/go/analysis"

var Analyzer = &analysis.Analyzer{
	Name: "loglinter",
	Doc:  "checks log message rules in supported logger calls",
	Run:  run,
}

func run(pass *analysis.Pass) (any, error) {
	return nil, nil
}
