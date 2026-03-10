package analyzer

import (
	"go/ast"
	"go/constant"

	"golang.org/x/tools/go/analysis"
)

func stringConstant(pass *analysis.Pass, expr ast.Expr) string {
	tv, ok := pass.TypesInfo.Types[expr]
	if !ok || tv.Value == nil || tv.Value.Kind() != constant.String {
		return ""
	}

	return constant.StringVal(tv.Value)
}
