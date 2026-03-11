# log-linter

Custom Go linter for validating log messages in application code.

The project is built on top of `go/analysis` and currently checks log calls from:
- `log/slog`
- `go.uber.org/zap`

## Implemented Rules

The linter reports these violations:

1. log message must start with a lowercase letter
2. log message must contain English text only
3. log message must not contain special symbols or emoji
4. log message must not contain potentially sensitive data

Examples:

```go
slog.Info("Starting server")             // lowercase violation
slog.Info("запуск сервера")              // English-only violation
slog.Info("server started!")             // special symbols violation
slog.Info("token: " + token)             // sensitive data violation
zapLogger.Info("request done", zap.String("token", token)) // sensitive data violation
```

## Supported Calls

### `log/slog`

Supported package-level and receiver methods:
- `Debug`
- `Info`
- `Warn`
- `Error`
- `DebugContext`
- `InfoContext`
- `WarnContext`
- `ErrorContext`

### `go.uber.org/zap`

Supported receiver methods:
- `Debug`
- `Info`
- `Warn`
- `Error`

## Quick Start

Requirements:
- Go 1.25+

Build the linter:

```bash
make build
```

Run it against the included example package:

```bash
make run-example
```

This runs the built binary on [`examples/manualcheck/main.go`](/home/phobos/golangProjects/logLinter/examples/manualcheck/main.go) and prints diagnostics for all implemented rules.

## Manual Run

Build the binary:

```bash
go build -o ./logLinter ./cmd/loglinter
```

Run the linter on a package:

```bash
./logLinter ./examples/manualcheck
```

Run it on the whole module:

```bash
./logLinter ./...
```

## Suggested Fixes

The analyzer provides `SuggestedFixes` for these rules only:
- log message must start with a lowercase letter
- log message must not contain special symbols or emoji

The fixes are generated only for direct string literal messages. `english` and `sensitive` diagnostics are still report-only.

Preview suggested fixes as a unified diff with the standalone binary:

```bash
make diff-fixes-example
```

Or run the binary directly (in project root):

```bash
./logLinter -fix -diff /path/to/your/file_or_directory
```

## `golangci-lint` Integration

This repository exposes `loglinter` as a private module plugin for `golangci-lint`.

Build a custom `golangci-lint` binary with the plugin compiled in:

```bash
make build-custom-gcl
```

This uses the local [`.custom-gcl.yml`](/home/phobos/golangProjects/logLinter/.custom-gcl.yml) configuration and produces `/tmp/loglinter-custom-gcl`.

Run it against the included example package:

```bash
make lint-example
```

Run it on your own package or module with the built binary:

```bash
/tmp/loglinter-custom-gcl run --config /path/to/your/.golangci.yml /path/to/your/package
```

The included [`.golangci.yml`](/home/phobos/golangProjects/logLinter/.golangci.yml) shows how to enable `loglinter` as a module-based custom linter.

The integration entrypoint lives in [`pkg/golangci/plugin.go`](/home/phobos/golangProjects/logLinter/pkg/golangci/plugin.go) and delegates directly to the existing analyzer in [`internal/analyzer/analyzer.go`](/home/phobos/golangProjects/logLinter/internal/analyzer/analyzer.go), keeping plugin code isolated from analyzer logic.

### **Config**

You can disable specific checks in `golangci-lint` configuration:

```yaml
version: "2"

linters:
  settings:
    custom:
      loglinter:
        type: module
        settings:
          disable:
            - english
            - sensitive
          extra_sensitive_keywords:
            - session_id
            - client_secret
```

Supported values inside `disable`:
- `lowercase`
- `english`
- `symbols`
- `sensitive`

`extra_sensitive_keywords` extends the built-in sensitive keyword list. Values are normalized before matching, so entries like `client_secret`, `client-secret`, and `ClientSecret` behave the same.

Important: changing [`.golangci.yml`](/home/phobos/golangProjects/logLinter/.golangci.yml) only requires rerunning the built binary.

## Testing

Run analyzer tests based on `analysistest`:

```bash
make test-analyzer
```

Run all tests:

```bash
make test
```

The analyzer testdata is located under:
- [`internal/analyzer/testdata/src/a`](/home/phobos/golangProjects/logLinter/internal/analyzer/testdata/src/a)

## Project Layout

```text
cmd/loglinter/         singlechecker entrypoint
internal/analyzer/     analyzer wiring and analysistest cases
internal/matchers/     logger call extraction
internal/rules/        individual rule implementations
examples/manualcheck/  package for manual binary checks
```

## Example Output

Typical `golangci-lint` output looks like this:

```text
examples/manualcheck/main.go:15:12: log message must start with a lowercase letter (loglinter)
	slog.Info("Starting server")
	          ^
examples/manualcheck/main.go:16:12: log message must contain English text only (loglinter)
	slog.Warn("запуск warning path")
	          ^
examples/manualcheck/main.go:17:13: log message must not contain special symbols or emoji (loglinter)
	slog.Error("server failed...")
	           ^
examples/manualcheck/main.go:18:13: log message must not contain potentially sensitive data (loglinter)
	slog.Debug("user password: " + password)
	           ^
11 issues:
* loglinter: 11
```

Standalone binary output looks like this:

```text
/home/phobos/golangProjects/logLinter/examples/manualcheck/main.go:15:12: log message must start with a lowercase letter
/home/phobos/golangProjects/logLinter/examples/manualcheck/main.go:16:12: log message must contain English text only
/home/phobos/golangProjects/logLinter/examples/manualcheck/main.go:17:13: log message must not contain special symbols or emoji
/home/phobos/golangProjects/logLinter/examples/manualcheck/main.go:18:13: log message must not contain potentially sensitive data
```

Suggested fixes preview from `./logLinter -fix -diff ./examples/manualcheck` looks like this:

```diff
--- /home/phobos/golangProjects/logLinter/examples/manualcheck/main.go (old)
+++ /home/phobos/golangProjects/logLinter/examples/manualcheck/main.go (new)
@@ -12,18 +12,18 @@
 	token := "token-value"
 	msg := "Starting from variable"
 
-	slog.Info("Starting server")
+	slog.Info("starting server")
 	slog.Warn("запуск warning path")
-	slog.Error("server failed...")
+	slog.Error("server failed")
 	slog.Debug("user password: " + password)
 ```
