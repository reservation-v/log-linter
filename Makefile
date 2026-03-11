CUSTOM_GCL_BIN := /tmp/loglinter-custom-gcl

.PHONY: build run run-example diff-fixes-example build-custom-gcl lint-example test-analyzer

build:
	go build -o ./logLinter ./cmd/loglinter/

run-example: build
	./logLinter ./examples/manualcheck/

diff-fixes-example: build
	./logLinter -fix -diff ./examples/manualcheck/

build-custom-gcl:
	go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2 custom

lint-example: build-custom-gcl
	$(CUSTOM_GCL_BIN) run ./examples/manualcheck/

test-analyzer:
	go test ./internal/analyzer -run TestAnalyzer

test:
	go test ./...
