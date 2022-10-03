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

HELM_VALUES_FILE ?= "./tilt-chart-values.yaml"

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
	helm dependency update ./charts/newrelic-prometheus-agent

.PHONY: tilt-up
tilt-up:
	$(MAKE) helm-deps
	eval $$(minikube docker-env) && tilt up -- --helm_values_file=$(HELM_VALUES_FILE) ; tilt down -- --helm_values_file=$(HELM_VALUES_FILE)

.PHONY: tilt-ci
tilt-ci:
	$(MAKE) helm-deps
	eval $$(minikube docker-env) && tilt ci -- --helm_values_file=$(HELM_VALUES_FILE)

.PHONY: integration-test
integration-test:
	KUBECONFIG='./.kubeconfig-dev' minikube update-context
	go test ./... -tags=integration_test -count=1 -race

.PHONY: chart-unit-test
chart-unit-test:
	ct --config .github/ct.yaml lint --debug
	helm dependency update ./charts/newrelic-prometheus-agent
	helm unittest ./charts/newrelic-prometheus-agent -3

PROMETHEUS_VERSION_CHART := $(shell grep "appVersion" charts/newrelic-prometheus-agent/Chart.yaml | grep -o -E "(\.\d+.\d+)")
PROMETHEUS_VERSION_GO := $(shell grep "github.com/prometheus/prometheus" go.mod| grep -o -E "(\.\d+.\d+)")
.PHONY: check-prometheus-version
check-prometheus-version:
ifneq ($(PROMETHEUS_VERSION_CHART), $(PROMETHEUS_VERSION_GO))
	@echo "Prometheus server version defined in chart does not match with the one in Go dependency"
	exit 1
endif


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

.PHONY: build-license-notice
build-license-notice:
	@go list -mod=mod -m -json all | go-licence-detector -noticeOut=NOTICE.txt -rules ./assets/licence/rules.json  -noticeTemplate ./assets/licence/THIRD_PARTY_NOTICES.md.tmpl -noticeOut THIRD_PARTY_NOTICES.md -overrides ./assets/licence/overrides -includeIndirect

.PHONY: generate-chart-docs
build-chart-docs:
	helm-docs -c charts/newrelic-prometheus-agent
