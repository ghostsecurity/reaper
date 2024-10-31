.DEFAULT_GOAL := help

ENV_FILE := $(shell touch .env) # create .env if it doesn't exist

# export all .env vars
include .env
$(eval export $(shell sed -ne 's/ *#.*$$//; /./ s/=.*$$// p' .env))

.PHONY: help
help: ## Show this help
	@echo "Usage: make [target]\n"
	@cat ${MAKEFILE_LIST} | grep "[#]# " | grep -v grep | sort | column -t -s '##' | sed -e 's/^/ /'
	@echo ""

.PHONY: run
run: ## Run the local server
	go run ./cmd/reaper
    
.PHONY: dev
dev: ## Run the local server with air watcher
	air
    
.PHONY: build
build: ## Build the local server
	go build -ldflags="-s -w" -o reaper ./cmd/reaper

.PHONY: lint
lint: ## Run linters
	go vet ./...
	which golangci-lint || go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.52.2
	golangci-lint run --timeout 3m --verbose

.PHONY: test
test: ## Run tests
	go test -v ./...
