.PHONY: build run run-example test-analyzer

build:
	go build -o ./logLinter ./cmd/loglinter/

run-example: build
	./logLinter ./examples/manualcheck/

test-analyzer:
	go test ./internal/analyzer -run TestAnalyzer

test:
	go test ./...
