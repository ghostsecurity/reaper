
default: build

.PHONY: clean
clean:
	rm -rf build/bin || true
	rm -rf frontend/dist/* || true

.PHONY: test
test: test-go test-js

.PHONY: test-go
test-go:
	go clean -testcache
	go test ./... -race

.PHONY: test-js
test-js:
	cd frontend && npm install && npm test

.PHONY: lint
lint: lint-go lint-js

.PHONY: lint-js
lint-js:
	cd frontend && npm install && npm run lint

.PHONY: lint-go
lint-go:
	which golangci-lint || go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.50.1
	golangci-lint run --timeout 3m --verbose

.PHONY: build-js
build-js: clean
	cd frontend && npm install && npm run build

.PHONY: build-bindings
build-bindings: clean
	go run ./cmd/gen-bindings

.PHONY: build
build: clean build-bindings build-js

.PHONY: install
install: build
	go install ./cmd/reaper

.PHONY: run
run: build
	go run ./cmd/reaper

.PHONY: fix
fix:
	cd frontend && npm install && npm run fix

.PHONY: docs
docs:
	cd docs && bundle install && bundle exec jekyll serve --livereload

docker:
	docker build -t "reaper" .
	docker run -v $$HOME/.reaper:/.reaper -p 8080:8080 -p 31337:31337 reaper
