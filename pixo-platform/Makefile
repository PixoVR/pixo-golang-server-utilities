
GO_BIN ?= $$(go env GOPATH)/bin

all: lint test

update:
	@echo "🔄 Updating Pixo utilities..."
	@go get github.com/PixoVR/pixo-golang-clients/pixo-platform@latest
	@go mod tidy

lint:
	@echo "🪮 Linting..."
	@golangci-lint run --timeout 5m

test:
	@echo "🧪 Running tests with coverage..."
	@IS_LOCAL=true go test ./... -coverprofile=./cover.out -covermode=atomic -coverpkg=./...

coverage: test
	@echo "🧪 Checking test coverage threshold..."
	@${GO_BIN}/go-test-coverage --config=./.coverage.yaml

