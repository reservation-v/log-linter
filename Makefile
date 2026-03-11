.PHONY: build run run-example build-custom-gcl lint-example lint test-analyzer

build:
	go build -o ./logLinter ./cmd/loglinter/

run-example: build
	./logLinter ./examples/manualcheck/

build-custom-gcl:
	go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2 custom

lint-example: build-custom-gcl
	./custom-gcl run ./examples/manualcheck/

test-analyzer:
	go test ./internal/analyzer -run TestAnalyzer

test:
	go test ./...
