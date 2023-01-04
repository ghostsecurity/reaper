
default: build

.PHONY: clean
clean:
	rm -rf build/bin || true
	rm -rf frontend/dist || true

.PHONY: wails
wails:
	(which wails >/dev/null && wails update) || go install github.com/wailsapp/wails/v2/cmd/wails@latest

.PHONY: build
build: clean wails
	wails build

.PHONY: test
test: test-go test-js

.PHONY: test-go
test-go:
	go test ./... -race

.PHONY: test-js
test-js:
	cd frontend && npm install && npm test

.PHONY: lint
lint: lint-go

.PHONY: lint-go
lint-go:
	which golangci-lint || go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.50.1
	go env
	echo $PATH
	golangci-lint run --timeout 3m --verbose

.PHONY: run
run: clean wails
	wails dev
