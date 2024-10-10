
all: lint test


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
