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
	go test ./... -count=1 -race -coverprofile=coverage.out -covermode=atomic

.PHONY: build
build: BINARY_NAME := $(if $(GOOS),$(BINARY_NAME)-$(GOOS),$(BINARY_NAME))
build: BINARY_NAME := $(if $(GOARCH),$(BINARY_NAME)-$(GOARCH),$(BINARY_NAME))
build:
	CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(LD_FLAGS) -o $(BIN_DIR)/$(BINARY_NAME) ./cmd/configurator/configurator.go
compile: build

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
	go test ./... -tags=integration_test -count=1 -race -coverprofile=integration-coverage.out -covermode=atomic

.PHONY: chart-unit-test
chart-unit-test:
	ct --config .github/ct.yaml lint --debug
	helm dependency update ./charts/newrelic-prometheus-agent
	helm unittest ./charts/newrelic-prometheus-agent -3

PROMETHEUS_VERSION_CHART := $(shell grep "appVersion" charts/newrelic-prometheus-agent/Chart.yaml | grep -o -E "(\.[0-9]+\.[0-9]+)")
PROMETHEUS_VERSION_GO := $(shell grep "github.com/prometheus/prometheus" go.mod| grep -o -E "(\.[0-9]+\.[0-9]+)")
.PHONY: check-prometheus-version
check-prometheus-version:
	@echo PROMETHEUS_VERSION_CHART=$(PROMETHEUS_VERSION_CHART), PROMETHEUS_VERSION_GO=$(PROMETHEUS_VERSION_GO)
ifneq ($(PROMETHEUS_VERSION_CHART), $(PROMETHEUS_VERSION_GO))
	@echo "Prometheus server version defined in chart does not match with the one in Go dependency"
	exit 1
endif

NEWRELIC_E2E ?= go run github.com/newrelic/newrelic-integration-e2e-action@latest
.PHONY: e2e-test
e2e-test:
	$(NEWRELIC_E2E) \
		--commit_sha=test-string \
		--retry_attempts=5 \
		--retry_seconds=60 \
        --account_id=${ACCOUNT_ID} \
		--api_key=${API_REST_KEY} \
		--license_key=${LICENSE_KEY} \
        --spec_path=./test/e2e/test-specs.yml \
		--verbose_mode=true \
		--agent_enabled="false"

LICENSE_DETECTOR ?= go run go.elastic.co/go-licence-detector@latest
.PHONY: build-license-notice
build-license-notice:
	@go list -mod=mod -m -json all | $(LICENSE_DETECTOR) \
	  -noticeOut=NOTICE.txt \
	  -rules ./assets/licence/rules.json \
	  -noticeTemplate ./assets/licence/THIRD_PARTY_NOTICES.md.tmpl \
	  -noticeOut THIRD_PARTY_NOTICES.md \
	  -overrides ./assets/licence/overrides \
	  -includeIndirect

HELM_DOCS ?= go run github.com/norwoodj/helm-docs/cmd/helm-docs@latest
.PHONY: generate-chart-docs
build-chart-docs:
	$(HELM_DOCS) -c charts/newrelic-prometheus-agent

## Release Toolkit targets
RT_BIN ?= go run github.com/newrelic/release-toolkit@latest
DICTIONARY_DIRECTORY = .github/rt-dictionary.yaml

.PHONY: release-notes
release-notes: _release-changes
	@$(RT_BIN) render-changelog
	@echo "RELEASE NOTES:\n"
	@cat CHANGELOG.partial.md

### Upgrades the CHANGELOG.md as if a Release is being triggered.
.PHONY: release-changelog
release-changelog: _release-changes
	$(RT_BIN) update-markdown --markdown CHANGELOG.md --version $$($(RT_BIN) next-version --tag-prefix "v")
	@git --no-pager diff CHANGELOG.md

### Prints out the Release Notes as if a Release is being triggered.
.PHONY: _release-changes
_release-changes:
	@$(RT_BIN) validate-markdown
	@$(RT_BIN) generate-yaml --excluded-dirs "charts,.github" --tag-prefix "v"
	@$(RT_BIN) link-dependencies --dictionary ${DICTIONARY_DIRECTORY}


CHART_DIRECTORY = "charts/newrelic-prometheus-agent"
MARKDOWN_FILE = ${CHART_DIRECTORY}/CHANGELOG.md
CHART_PREFIX = newrelic-prometheus-agent-

.PHONY: release-notes-chart
release-notes-chart: _release-changes-chart
	@$(RT_BIN) render-changelog --markdown ${CHART_DIRECTORY}/release-notes.md
	@echo "RELEASE NOTES:\n"
	@cat ${CHART_DIRECTORY}/release-notes.md

### Upgrades the CHANGELOG.md as if a Release is being triggered.
.PHONY: release-changelog-chart
release-changelog-chart: _release-changes-chart
	@$(RT_BIN) update-markdown --markdown ${MARKDOWN_FILE} --version $$($(RT_BIN) next-version --tag-prefix ${CHART_PREFIX})
	@git --no-pager diff ${CHART_DIRECTORY}/CHANGELOG.md

### Prints out the Release Notes as if a Release is being triggered.
.PHONY: _release-changes-chart
_release-changes-chart:
	@$(RT_BIN) validate-markdown --markdown ${MARKDOWN_FILE}
	@$(RT_BIN) generate-yaml --markdown ${MARKDOWN_FILE} --included-dirs ${CHART_DIRECTORY}  --tag-prefix ${CHART_PREFIX}
	@$(RT_BIN) link-dependencies --dictionary ${DICTIONARY_DIRECTORY}

# rt-update-changelog runs the release-toolkit run.sh script by piping it into bash to update the CHANGELOG.md.
# It also passes down to the script all the flags added to the make target. To check all the accepted flags,
# see: https://github.com/newrelic/release-toolkit/blob/main/contrib/ohi-release-notes/run.sh
#  e.g. `make rt-update-changelog -- -v`
.PHONY: rt-update-changelog
rt-update-changelog:
	curl "https://raw.githubusercontent.com/newrelic/release-toolkit/v1/contrib/ohi-release-notes/run.sh" | bash -s -- $(filter-out $@,$(MAKECMDGOALS))
