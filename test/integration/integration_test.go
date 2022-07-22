package integration

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
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
    - name: self-metrics
      scrape_interval: 1s
      urls:
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
