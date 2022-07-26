package integration

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
	"testing"

	"github.com/newrelic-forks/newrelic-prometheus/test/integration/mocks"

	"github.com/stretchr/testify/require"
)

const (
	prometheusVersion = "2.36.2"
)

func TestMain(m *testing.M) {
	if err := fetchPrometheusBinary(prometheusVersion); err != nil {
		log.Fatalf("fail to fetch prometheus binary: %s", err)
	}

	os.Exit(m.Run())
}

func Test_ServerReady(t *testing.T) {
	t.Parallel()

	ps := newPrometheusServer(t)

	asserter := newAsserter(ps.port)

	inputConfig := `
data_source_name: "data-source"
newrelic_remote_write:
  license_key: nrLicenseKey
`

	outputConfigPath := runConfigurator(t, inputConfig)

	ps.start(t, outputConfigPath)

	asserter.prometheusServerReady(t)
}

func Test_SelfMetricsAreScrapedCorrectly(t *testing.T) {
	t.Parallel()

	ps := newPrometheusServer(t)

	asserter := newAsserter(ps.port)

	rw := mocks.StartRemoteWriteEndpoint(t, asserter.appendable)

	inputConfig := fmt.Sprintf(`
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

	outputConfigPath := runConfigurator(t, inputConfig)

	ps.start(t, outputConfigPath)

	asserter.metricName(t, "prometheus_build_info")
}

func Test_ExtraScapeConfig(t *testing.T) {
	t.Parallel()

	ps := newPrometheusServer(t)

	asserter := newAsserter(ps.port)

	rw := mocks.StartRemoteWriteEndpoint(t, asserter.appendable)
	ex := mocks.StartExporter(t)

	mockExporterTarget := strings.Replace(ex.URL, "http://", "", 1)
	inputConfig := fmt.Sprintf(`
static_targets:
  jobs:
    - job_name: metrics-a
      scrape_interval: 1s
      targets:
        - "%s"
      metrics_path: /metrics-a
      extra_relabel_config:
        - replacement: my-value
          target_label: custom_label
          action: replace

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

	outputConfigPath := runConfigurator(t, inputConfig)

	ps.start(t, outputConfigPath)

	asserter.metricLabels(t, map[string]string{"custom_label": "my-value", "job": "metrics-a"}, "custom_metric_a")
	asserter.metricLabels(t, map[string]string{"custom_label": "", "job": "metrics-b"}, "custom_metric_b")
}

func Test_ExternalLabelsAreAddedToEachSample(t *testing.T) {
	t.Parallel()

	ps := newPrometheusServer(t)
	asserter := newAsserter(ps.port)
	rw := mocks.StartRemoteWriteEndpoint(t, asserter.appendable)

	inputConfig := fmt.Sprintf(`
static_targets:
  jobs:
    - job_name: self-metrics
      scrape_interval: 1s
      targets:
        - "localhost:%s"
data_source_name: "data-source"
newrelic_remote_write:
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

	outputConfigPath := runConfigurator(t, inputConfig)

	ps.start(t, outputConfigPath)

	asserter.metricName(t, "prometheus_build_info")
	asserter.metricLabels(t, map[string]string{"cluster_name": "test", "one": "two", "three": "four"}, "prometheus_build_info")
}

func runConfigurator(t *testing.T, inputConfig string) string {
	t.Helper()

	tempDir := t.TempDir()
	inputConfigPath := path.Join(tempDir, "input.yml")
	outputConfigPath := path.Join(tempDir, "output.yml")

	err := ioutil.WriteFile(inputConfigPath, []byte(inputConfig), 0o444)
	require.NoError(t, err)

	// nolint:gosec
	configurator := exec.Command(
		"go",
		"run",
		"../../cmd/configurator",
		fmt.Sprintf("--input=%s", inputConfigPath),
		fmt.Sprintf("--output=%s", outputConfigPath),
	)

	out, err := configurator.CombinedOutput()
	require.NoError(t, err, string(out))

	return outputConfigPath
}
