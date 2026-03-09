// this package is made for extracting package names,
// finding needed packages - slog/zap and giving 
// analyzer full inf about called function
package matchers

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

type LoggerKind string

const (
	LoggerSlog LoggerKind = "slog"
	LoggerZap  LoggerKind = "zap"
)

// main log structure
type LogCall struct {
	Kind    LoggerKind
	Method  string
	Call    *ast.CallExpr
	Message ast.Expr
}

// extrats func package and in switch checks
// chooses between slog and zap logger packages
func ExtractLogCall(pass *analysis.Pass, call *ast.CallExpr) (LogCall, bool) {
	fn := resolveFunc(pass.TypesInfo, call)
	if fn == nil || fn.Pkg() == nil {
		return LogCall{}, false
	}

	switch fn.Pkg().Path() {
	case "log/slog":
		idx, ok := slogMessageIndex(fn)
		if !ok || len(call.Args) <= idx {
			return LogCall{}, false
		}
		return LogCall{
			Kind:    LoggerSlog,
			Method:  fn.Name(),
			Call:    call,
			Message: call.Args[idx],
		}, true
	case "go.uber.org/zap":
		idx, ok := zapMessageIndex(fn)
		if !ok || len(call.Args) <= idx {
			return LogCall{}, false
		}
		return LogCall{
			Kind:    LoggerZap,
			Method:  fn.Name(),
			Call:    call,
			Message: call.Args[idx],
		}, true
	default:
		return LogCall{}, false
	}
}

// resolves func package name (protection from aliases)
func resolveFunc(info *types.Info, call *ast.CallExpr) *types.Func {
	switch fun := call.Fun.(type) {
	case *ast.Ident:
		fn, _ := info.Uses[fun].(*types.Func)
		return fn
	case *ast.SelectorExpr:
		if selection := info.Selections[fun]; selection != nil {
			if fn, ok := selection.Obj().(*types.Func); ok {
				return fn
			}
		}
		fn, _ := info.Uses[fun.Sel].(*types.Func)
		return fn
	default:
		return nil
	}
}

// cover basic slog cases
func slogMessageIndex(fn *types.Func) (int, bool) {
	switch fn.Name() {
	case "Debug", "Info", "Warn", "Error":
		if !isPackageOrReceiver(fn, "log/slog", "Logger") {
			return 0, false
		}
		return 0, hasStringParam(fn, 0)
	case "DebugContext", "InfoContext", "WarnContext", "ErrorContext":
		if !isPackageOrReceiver(fn, "log/slog", "Logger") {
			return 0, false
		}
		return 1, hasStringParam(fn, 1)
	default:
		return 0, false
	}
}

// cover basic slog cases
func zapMessageIndex(fn *types.Func) (int, bool) {
	switch fn.Name() {
	case "Debug", "Info", "Warn", "Error":
		if !isReceiver(fn, "go.uber.org/zap", "Logger") {
			return 0, false
		}
		return 0, hasStringParam(fn, 0)
	default:
		return 0, false
	}
}

// checks if idx parameter of func is string
func hasStringParam(fn *types.Func, idx int) bool {
	sig, ok := fn.Type().(*types.Signature)
	if !ok || sig.Params() == nil || sig.Params().Len() <= idx {
		return false
	}

	basic, ok := sig.Params().At(idx).Type().(*types.Basic)
	return ok && basic.Kind() == types.String
}

// slog has 2 types of output: 
// - package level like slog.Info()
// - methods like logger.Info (logger is an object of type *slog.Logger)
func isPackageOrReceiver(fn *types.Func, pkgPath, receiver string) bool {
	sig, ok := fn.Type().(*types.Signature)
	if !ok {
		return false
	}

	if sig.Recv() == nil {
		return fn.Pkg() != nil && fn.Pkg().Path() == pkgPath
	}

	return isReceiver(fn, pkgPath, receiver)
}

// zap has only methods calls and this func checks
// is it receiver or not
func isReceiver(fn *types.Func, pkgPath, receiver string) bool {
	sig, ok := fn.Type().(*types.Signature)
	if !ok || sig.Recv() == nil {
		return false
	}

	recvType := sig.Recv().Type()
	if ptr, ok := recvType.(*types.Pointer); ok {
		recvType = ptr.Elem()
	}

	named, ok := recvType.(*types.Named)
	if !ok || named.Obj() == nil || named.Obj().Pkg() == nil {
		return false
	}

	return named.Obj().Pkg().Path() == pkgPath && named.Obj().Name() == receiver
}
