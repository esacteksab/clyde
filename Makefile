MAKEFLAGS += --warn-undefined-variables
SHELL := bash
.SHELLFLAGS := -eu -o pipefail -c
.DEFAULT_GOAL := all
.DELETE_ON_ERROR:
.SUFFIXES:

.PHONY: audit
audit:
	go vet ./...
	go tool -modfile=go.tool.mod staticcheck ./...
	go tool -modfile=go.tool.mod govulncheck ./...
	golangci-lint run -v

.PHONY: clean
clean:

.PHONY: build
build:

.PHONY: format
format:
	gofumpt -l -w -extra .

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: test
test:
	go test ./...
