BIN_DIR = ./bin
TOOLS_DIR := $(BIN_DIR)/dev-tools
BINARY_NAME = prometheus-configurator

GOOS ?=
GOARCH ?=
CGO_ENABLED ?= 0

BUILD_DATE := $(shell date)
COMMIT := $(shell git rev-parse HEAD)
TAG ?= dev
COMMIT ?= $(shell git rev-parse HEAD || echo "unknown")

LD_FLAGS ?= -ldflags="-X 'main.integrationVersion=$(TAG)' -X 'main.gitCommit=$(COMMIT)' -X 'main.buildDate=$(BUILD_DATE)' "

.PHONY: all
all: clean build-multiarch

.PHONY: clean
clean:
	rm -rf $(BIN_DIR)

.PHONY: test
test:
	go test ./... -count=1 -race

.PHONY: build
build: BINARY_NAME := $(if $(GOOS),$(BINARY_NAME)-$(GOOS),$(BINARY_NAME))
build: BINARY_NAME := $(if $(GOARCH),$(BINARY_NAME)-$(GOARCH),$(BINARY_NAME))
build:
	CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(LD_FLAGS) -o $(BIN_DIR)/$(BINARY_NAME) ./cmd/configurator/configurator.go

.PHONY: build-multiarch
build-multiarch: clean
	$(MAKE) build GOOS=linux GOARCH=amd64
	$(MAKE) build GOOS=linux GOARCH=arm64
	$(MAKE) build GOOS=linux GOARCH=arm

.PHONY: start-local-cluster
start-local-cluster:
	minikube start

.PHONY: helm-deps
helm-deps:
	helm repo add newrelic https://helm-charts.newrelic.com
	helm dependency update ./charts/newrelic-prometheus

.PHONY: tilt-up
tilt-up:
	$(MAKE) helm-deps
	eval $$(minikube docker-env) && tilt up ; tilt down

.PHONY: tilt-ci
tilt-ci:
	$(MAKE) helm-deps
	eval $$(minikube docker-env) && tilt ci

.PHONY: integration-test
integration-test:
	KUBECONFIG='./.kubeconfig-dev' minikube update-context
	go test ./... -tags=integration_test -count=1 -race

.PHONY: chart-unit-test
chart-unit-test:
	ct --config .github/ct.yaml lint --debug
	helm dependency update ./charts/newrelic-prometheus
	helm unittest ./charts/newrelic-prometheus -3

.PHONY: e2e-test
e2e-test:
	newrelic-integration-e2e \
		--commit_sha=test-string \
		--retry_attempts=5 \
		--retry_seconds=60 \
        --account_id=${ACCOUNT_ID} \
		--api_key=${API_REST_KEY} \
		--license_key=${LICENSE_KEY} \
        --spec_path=./test/e2e/test-specs.yml \
		--verbose_mode=true \
		--agent_enabled="false"
