
all: lint test


.PHONY: update
update:
	@echo "🔄 Updating Pixo utilities..."
	@go get github.com/PixoVR/pixo-golang-server-utilities/pixo-platform@0.2.34
	@go mod tidy


.PHONY: lint
lint:
	@echo "🪮 Linting..."
	@golangci-lint run --timeout 5m


.PHONY: test
test:
	@echo "🧪 Running tests with coverage..."
	@go test ./... -coverprofile cover.out -covermode atomic -coverpkg ./...


.PHONY: calc
calc:
	@echo "🧪 Checking test coverage threshold..."
	@go-test-coverage --config=./.coverage.yaml


.PHONY: coverage
coverage: test calc


.PHONY: cov
cov: coverage
