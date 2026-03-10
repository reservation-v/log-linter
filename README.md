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
- Go 1.24+

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

Typical output looks like this:

```text
examples/manualcheck/main.go:15:12: log message must start with a lowercase letter
examples/manualcheck/main.go:16:12: log message must contain English text only
examples/manualcheck/main.go:17:13: log message must not contain special symbols or emoji
examples/manualcheck/main.go:18:13: log message must not contain potentially sensitive data
```
