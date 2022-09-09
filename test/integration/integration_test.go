//go:build integration_test

package integration

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/newrelic/newrelic-prometheus-configurator/test/integration/mocks"
)

func TestMain(m *testing.M) {
	prometheusVersion, err := prometheusChartVersion()
	if err != nil {
		log.Fatalf("fail to fetch prometheus version: %s", err)
	}

	if err := fetchPrometheusBinary(prometheusVersion); err != nil {
		log.Fatalf("fail to fetch prometheus binary: %s", err)
	}

	os.Exit(m.Run())
}

func Test_ServerReady(t *testing.T) {
	t.Parallel()

	ps := newPrometheusServer(t)

	asserter := newAsserter(ps)

	nrConfigConfig := `
newrelic_remote_write:
  data_source_name: "data-source"
  license_key: nrLicenseKey
`

	prometheusConfigConfigPath := runConfigurator(t, nrConfigConfig)

	ps.start(t, prometheusConfigConfigPath)

	asserter.prometheusServerReady(t)
}

func Test_SelfMetricsAreScrapedCorrectly(t *testing.T) {
	t.Parallel()

	ps := newPrometheusServer(t)

	asserter := newAsserter(ps)

	rw := mocks.StartRemoteWriteEndpoint(t, asserter.appendable)

	nrConfigConfig := fmt.Sprintf(`
static_targets:
  jobs:
    - job_name: self-metrics
      scrape_interval: 1s
      targets:
        - "localhost:%s"

newrelic_remote_write:
  license_key: nrLicenseKey
  proxy_url: %s
  tls_config:
    insecure_skip_verify: true
common:
  scrape_interval: 1s
`, ps.port, rw.URL)

	prometheusConfigConfigPath := runConfigurator(t, nrConfigConfig)

	ps.start(t, prometheusConfigConfigPath)

	asserter.metricName(t, "prometheus_build_info")
}

func Test_ExtraScapeConfig(t *testing.T) {
	t.Parallel()

	ps := newPrometheusServer(t)

	asserter := newAsserter(ps)

	rw := mocks.StartRemoteWriteEndpoint(t, asserter.appendable)
	ex := mocks.StartExporter(t)

	mockExporterTarget := ex.Listener.Addr().String()
	nrConfigConfig := fmt.Sprintf(`
static_targets:
  jobs:
    - job_name: metrics-a
      scrape_interval: 1s
      targets:
        - "%s"
      metrics_path: /metrics-a
      labels:
        custom_label: foo

extra_scrape_configs:
  - job_name: metrics-b
    static_configs:
      - targets:
        - "%s"
    metrics_path: /metrics-b
    honor_timestamps: false
    scrape_interval: 1s

newrelic_remote_write:
  license_key: nrLicenseKey
  proxy_url: %s
  tls_config:
    insecure_skip_verify: true
`, mockExporterTarget, mockExporterTarget, rw.URL)

	prometheusConfigConfigPath := runConfigurator(t, nrConfigConfig)

	ps.start(t, prometheusConfigConfigPath)

	asserter.metricLabels(t, map[string]string{"custom_label": "foo", "instance": mockExporterTarget, "job": "metrics-a"}, "custom_metric_a")
	asserter.metricLabels(t, map[string]string{"instance": mockExporterTarget, "job": "metrics-b"}, "custom_metric_b")
}

func Test_ExternalLabelsAreAddedToEachSample(t *testing.T) {
	t.Parallel()

	ps := newPrometheusServer(t)
	asserter := newAsserter(ps)
	rw := mocks.StartRemoteWriteEndpoint(t, asserter.appendable)

	nrConfigConfig := fmt.Sprintf(`
static_targets:
  jobs:
    - job_name: self-metrics
      scrape_interval: 1s
      targets:
        - "localhost:%s"
newrelic_remote_write:
  data_source_name: "data-source"
  license_key: nrLicenseKey
  proxy_url: %s
  tls_config:
    insecure_skip_verify: true
common:
  scrape_interval: 1s
  external_labels:
    cluster_name: test
    one: two
    three: four
`, ps.port, rw.URL)

	prometheusConfigConfigPath := runConfigurator(t, nrConfigConfig)

	ps.start(t, prometheusConfigConfigPath)

	asserter.metricName(t, "prometheus_build_info")
	asserter.metricLabels(t, map[string]string{"cluster_name": "test", "one": "two", "three": "four"}, "prometheus_build_info")
}
