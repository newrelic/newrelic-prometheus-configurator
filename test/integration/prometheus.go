//go:build integration_test

package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

const (
	prometheusBinaryPath = "./prometheus"
	chartDefinitionFile  = "../../charts/newrelic-prometheus/Chart.yaml"
)

type prometheusServer struct {
	port string
}

func prometheusChartVersion() (string, error) {
	f, err := ioutil.ReadFile(chartDefinitionFile)
	if err != nil {
		return "", fmt.Errorf("reading Chart.yaml: %w", err)
	}

	nrConfig := &struct {
		AppVersion string `yaml:"appVersion"`
	}{}

	if err = yaml.Unmarshal(f, nrConfig); err != nil {
		return "", fmt.Errorf("unmarshalling Chart.yaml: %w", err)
	}

	return strings.TrimPrefix(nrConfig.AppVersion, "v"), nil
}

func newPrometheusServer(t *testing.T) *prometheusServer {
	t.Helper()

	ps := &prometheusServer{}

	ps.port = freePort(t)

	return ps
}

// Use when expected to pass a valid config and want to have the server running.
func (ps *prometheusServer) start(t *testing.T, configPath string) {
	t.Helper()

	//nolint: gosec
	prom := exec.Command(
		prometheusBinaryPath,
		"--enable-feature=agent",
		fmt.Sprintf("--config.file=%s", configPath),
		fmt.Sprintf("--web.listen-address=0.0.0.0:%s", ps.port),
		fmt.Sprintf("--storage.agent.path=%s", t.TempDir()),
		"--log.level=debug",
	)

	stderr := &bytes.Buffer{}
	prom.Stderr = stderr

	err := prom.Start()
	require.NoError(t, err, stderr)

	go func() {
		// This is used to print the prometheus errors when it fails.
		err := prom.Wait()
		assert.NoError(t, err, stderr)
	}()

	t.Cleanup(func() {
		err := prom.Process.Signal(os.Interrupt)
		assert.NoError(t, err, stderr)
	})
}

func (ps *prometheusServer) healthy(t *testing.T) bool {
	t.Helper()

	resp, err := http.Get(fmt.Sprintf("http://localhost:%s/-/healthy", ps.port))
	if err != nil {
		t.Logf("Fail to Get healthy API: %s", err)
		return false
	}

	return resp.StatusCode == http.StatusOK
}

func (ps *prometheusServer) targets(t *testing.T) (*targetDiscovery, bool) {
	t.Helper()

	targetsURL := fmt.Sprintf("http://localhost:%s/api/v1/targets", ps.port)
	resp, err := http.Get(targetsURL)
	if err != nil {
		t.Logf("Fail to Get targets API: %s", err)
		return nil, false
	}
	defer resp.Body.Close()

	decodedResponse := &response{}
	err = json.NewDecoder(resp.Body).Decode(decodedResponse)
	require.NoError(t, err)

	t.Logf("Targets API response: %s", decodedResponse)

	targets := &targetDiscovery{}
	err = json.Unmarshal(decodedResponse.Data, targets)
	require.NoError(t, err)

	return targets, true
}

// freePort returns an available TCP port. Basically returns the port provided by the
// kernel when trying to bind to port 0 in a similar way as httptest.NewServer does.
func freePort(t *testing.T) string {
	t.Helper()

	add, err := net.ResolveTCPAddr("tcp", "localhost:0")
	require.NoError(t, err)

	l, err := net.ListenTCP("tcp", add)
	require.NoError(t, err)

	defer l.Close()

	return fmt.Sprintf("%d", l.Addr().(*net.TCPAddr).Port) //nolint: forcetypeassert
}

// fetchPrometheusBinary check that the binary on prometheusBinaryPath is correct or try to fetch it from Prometheus repo.
func fetchPrometheusBinary(version string) error {
	binaryTarget := fmt.Sprintf("prometheus-%s.%s-%s", version, runtime.GOOS, runtime.GOARCH)
	tarName := fmt.Sprintf("%s.tar.gz", binaryTarget)
	tarURL := fmt.Sprintf("https://github.com/prometheus/prometheus/releases/download/v%s/%s", version, tarName)
	tarPath := path.Join(os.TempDir(), tarName)

	// binary already exists and has the correct version.
	if ok, _ := checkVersion(prometheusBinaryPath, version); ok {
		return nil
	}

	fetchTar := exec.Command(
		"curl", "-v",
		"--retry", "5",
		"--retry-delay", "3",
		"-L", tarURL,
		"--output", tarPath)
	if out, err := fetchTar.CombinedOutput(); err != nil {
		return fmt.Errorf("downloading the prometheus binary: command %s: prometheusConfig %s: %w", fetchTar.String(), out, err)
	}

	extract := exec.Command(
		"tar", "-x",
		"-f", tarPath,
		"--strip-components", "1", // remove the parent directory when extracting.
		"-C", ".", // change directory.
		path.Join(binaryTarget, "prometheus")) // selects only the 'prometheus' file to be extracted.
	if out, err := extract.CombinedOutput(); err != nil {
		return fmt.Errorf("un-compressing the prometheus binary: command %s: prometheusConfig %s: %w", extract.String(), out, err)
	}

	if ok, err := checkVersion(prometheusBinaryPath, version); !ok {
		return fmt.Errorf("prometheus binary version different than %s: %w", version, err)
	}

	return nil
}

func checkVersion(path string, version string) (bool, error) {
	if _, err := os.Stat(prometheusBinaryPath); os.IsNotExist(err) {
		return false, nil
	}

	checkVersion := exec.Command(path, "--version")

	out, err := checkVersion.CombinedOutput()
	if err != nil {
		return false, fmt.Errorf("executing the prometheus binary: command %s: prometheusConfig %s: %w", checkVersion.String(), out, err)
	}

	return strings.Contains(string(out), version), nil
}

type response struct {
	Status    string          `json:"status"`
	Data      json.RawMessage `json:"data,omitempty"`
	ErrorType string          `json:"errorType,omitempty"`
	Error     string          `json:"error,omitempty"`
	Warnings  []string        `json:"warnings,omitempty"`
}

// target has the information for one target.
type target struct {
	// Labels before any processing.
	DiscoveredLabels map[string]string `json:"discoveredLabels"`
	// Any labels that are added to this target and its metrics.
	Labels map[string]string `json:"labels"`

	ScrapePool string `json:"scrapePool"`
	ScrapeURL  string `json:"scrapeUrl"`
	GlobalURL  string `json:"globalUrl"`

	LastError          string    `json:"lastError"`
	LastScrape         time.Time `json:"lastScrape"`
	LastScrapeDuration float64   `json:"lastScrapeDuration"`
	Health             string    `json:"health"`

	ScrapeInterval string `json:"scrapeInterval"`
	ScrapeTimeout  string `json:"scrapeTimeout"`
}

// droppedTarget has the information for one target that was dropped during relabelling.
type droppedTarget struct {
	// Labels before any processing.
	DiscoveredLabels map[string]string `json:"discoveredLabels"`
}

// targetDiscovery has all the active targets.
type targetDiscovery struct {
	ActiveTargets  []*target        `json:"activeTargets"`
	DroppedTargets []*droppedTarget `json:"droppedTargets"`
}
