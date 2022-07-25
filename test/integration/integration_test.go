package integration

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path"
	"strings"
	"testing"

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

func Test_SelfMetrics(t *testing.T) {
	t.Parallel()

	ps := newPrometheusServer(t)

	asserter := newAsserter(ps.port)

	rw := startRemoteWriteEndpoint(t, asserter.appendable)

	inputConfig := fmt.Sprintf(`
static_targets:
  jobs:
    - job_name: self-metrics
      scrape_interval: 1s
      targets:
        - "localhost:%s"

extra_remote_write:
  - url: %s

newrelic_remote_write:
  license_key: nrLicenseKey
`, ps.port, rw.URL)

	outputConfigPath := runConfigurator(t, inputConfig)

	ps.start(t, outputConfigPath)

	asserter.metricName(t, "prometheus_build_info")
}

func Test_ExtraScapeConfig(t *testing.T) {
	t.Parallel()

	ps := newPrometheusServer(t)

	asserter := newAsserter(ps.port)

	rw := startRemoteWriteEndpoint(t, asserter.appendable)

	ex := startMockExporter(t)

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

extra_remote_write:
  - url: %s

newrelic_remote_write:
  license_key: nrLicenseKey
`, mockExporterTarget, mockExporterTarget, rw.URL)

	outputConfigPath := runConfigurator(t, inputConfig)

	ps.start(t, outputConfigPath)

	asserter.metricWithLabels(t, "custom_metric_a", []string{"custom_label", "instance", "job"})
	asserter.metricWithLabels(t, "custom_metric_b", []string{"instance", "job"})
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

func startMockExporter(t *testing.T) *httptest.Server {
	t.Helper()

	mux := http.NewServeMux()

	mux.HandleFunc("/metrics-a/", func(w http.ResponseWriter, r *http.Request) {
		response := "custom_metric_a 46"
		_, _ = fmt.Fprintln(w, response)
	})

	mux.HandleFunc("/metrics-b/", func(w http.ResponseWriter, r *http.Request) {
		response := "custom_metric_b 88"
		_, _ = fmt.Fprintln(w, response)
	})

	mockExporterServer := httptest.NewServer(mux)

	t.Cleanup(func() {
		mockExporterServer.Close()
	})

	return mockExporterServer
}
