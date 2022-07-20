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
	@echo "[clean] Removing integration binaries"
	@rm -rf $(BIN_DIR)

.PHONY: test
test:
	@echo "[test] Running tests"
	@go test ./... -count=1

.PHONY: compile
compile:
	@echo "[compile] Building $(BINARY_NAME)"
	CGO_ENABLED=$(CGO_ENABLED) go build $(LDFLAGS) -o $(BIN_DIR)/$(BINARY_NAME) ./cmd/configurator/configurator.go

.PHONY: compile-multiarch
compile-multiarch:
	$(MAKE) compile GOOS=linux GOARCH=amd64
	$(MAKE) compile GOOS=linux GOARCH=arm64
	$(MAKE) compile GOOS=linux GOARCH=arm

.PHONY: local-env-start
local-env-start:
	minikube start
	helm repo add newrelic https://helm-charts.newrelic.com
	helm dependency update ./charts/newrelic-prometheus
	$(MAKE) tilt-up

.PHONY: tilt-up
tilt-up:
	eval $$(minikube docker-env); tilt up ; tilt down

.PHONY: tilt-ci
tilt-ci:
	helm repo add newrelic https://helm-charts.newrelic.com
	helm dependency update ./charts/newrelic-prometheus
# tilt ci has a non configurable timeout of 30m that is why using 'timeout'.
	eval $$(minikube docker-env); timeout 5m tilt ci
