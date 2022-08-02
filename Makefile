BIN_DIR = ./bin
TOOLS_DIR := $(BIN_DIR)/dev-tools
BINARY_NAME = configurator

GOOS ?=
GOARCH ?=
CGO_ENABLED ?= 0

BUILD_DATE := $(shell date)
COMMIT := $(shell git rev-parse HEAD)
TAG ?= dev
COMMIT ?= $(shell git rev-parse HEAD || echo "unknown")

LDFLAGS ?= -ldflags="-X 'main.integrationVersion=$(TAG)' -X 'main.gitCommit=$(COMMIT)' -X 'main.buildDate=$(BUILD_DATE)' "

ifneq ($(strip $(GOOS)), )
BINARY_NAME := $(BINARY_NAME)-$(GOOS)
endif

ifneq ($(strip $(GOARCH)), )
BINARY_NAME := $(BINARY_NAME)-$(GOARCH)
endif

.PHONY: all
all: build

.PHONY: build
build: clean test compile-multiarch

.PHONY: clean
clean:
	rm -rf $(BIN_DIR)

.PHONY: test
test:
	go test ./... -count=1 -race

.PHONY: compile
compile:
	CGO_ENABLED=$(CGO_ENABLED) go build $(LDFLAGS) -o $(BIN_DIR)/$(BINARY_NAME) ./cmd/configurator/configurator.go

.PHONY: compile-multiarch
compile-multiarch:
	$(MAKE) compile GOOS=linux GOARCH=amd64
	$(MAKE) compile GOOS=linux GOARCH=arm64
	$(MAKE) compile GOOS=linux GOARCH=arm

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
	helm dependency update ./charts/newrelic-prometheus
	helm unittest ./charts/newrelic-prometheus -3
