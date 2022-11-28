
default: build

.PHONY: clean
clean:
	rm -rf build/bin

.PHONY: wails
wails:
	which wails || go install github.com/wailsapp/wails/cmd/wails@latest

.PHONY: build
build: clean wails
	wails build

.PHONY: test
test:
	go test ./... -race

.PHONY: run
run: clean wails
	wails dev
