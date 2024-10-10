
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif


LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

ENVTEST_K8S_VERSION = 1.28.0
ENVTEST ?= $(LOCALBIN)/setup-envtest


all: lint test


.PHONY: lint
lint:
	@echo "🪮 Linting..."
	@golangci-lint run --timeout 5m


.PHONY: init-env
init-env: $(LOCALBIN)
	test -s $(LOCALBIN)/setup-envtest || GOBIN=$(LOCALBIN) go install sigs.k8s.io/controller-runtime/tools/setup-envtest@latest


.PHONY: test
test:
	@echo "🧪 Running tests with coverage..."
	@KUBEBUILDER_ASSETS="$(shell $(ENVTEST) use $(ENVTEST_K8S_VERSION) --bin-dir $(LOCALBIN) -p path)" go test ./... -coverprofile cover.out -covermode atomic -coverpkg ./...


.PHONY: calc
calc:
	@echo "🧪 Checking test coverage threshold..."
	@go-test-coverage --config=./.coverage.yaml


.PHONY: coverage
coverage: test calc


.PHONY: cov
cov: coverage
