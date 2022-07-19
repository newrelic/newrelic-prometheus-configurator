package integration

import (
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
	"testing"
)

const (
	prometheusBinaryPath = "./prometheus"
)

type prometheusServer struct {
	port string
}

func newPrometheusServer(t *testing.T, options ...func(*prometheusServer)) *prometheusServer {
	t.Helper()

	ps := &prometheusServer{}

	ps.port = freePort(t)

	for _, op := range options {
		op(ps)
	}

	return ps
}

// Use when expected to pass a valid config and want to have the server running.
func (ps *prometheusServer) start(t *testing.T, configPath string) {
	t.Helper()

	// nolint: gosec
	prom := exec.Command(
		prometheusBinaryPath,
		"--enable-feature=agent",
		fmt.Sprintf("--config.file=%s", configPath),
		fmt.Sprintf("--web.listen-address=0.0.0.0:%s", ps.port),
		fmt.Sprintf("--storage.agent.path=%s", t.TempDir()),
	)

	// Log stderr in case of failure.
	stderr, err := prom.StderrPipe()
	if err != nil {
		t.Error(err)
	}

	if err := prom.Start(); err != nil {
		t.Error(err)
	}

	go func() {
		if err := prom.Wait(); err != nil {
			t.Errorf("unexpected server error: %s", err)
		}
	}()

	t.Cleanup(func() {
		if err := prom.Process.Signal(os.Interrupt); err != nil {
			t.Error(err)
		}

		// logs in case the any test fails.
		slurp, _ := io.ReadAll(stderr)
		t.Log(string(slurp))
	})
}

// freePort returns an available TCP port. Basically returns the port provided by the
// kernel when trying to bind to port 0 in a similar way as httptest.NewServer does.
func freePort(t *testing.T) string {
	t.Helper()

	add, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		t.Error(err)
	}

	l, err := net.ListenTCP("tcp", add)
	if err != nil {
		t.Error(err)
	}

	defer l.Close()

	return fmt.Sprintf("%d", l.Addr().(*net.TCPAddr).Port) // nolint: forcetypeassert
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
		"-L", tarURL,
		"--output", tarPath)
	if out, err := fetchTar.CombinedOutput(); err != nil {
		return fmt.Errorf("downloading the prometheus binary: command %s: output %s: %w", fetchTar.String(), out, err)
	}

	extract := exec.Command(
		"tar", "-x",
		"-f", tarPath,
		"--strip-components", "1", // remove the parent directory when extracting.
		"-C", ".", // change directory.
		path.Join(binaryTarget, "prometheus")) // selects only the 'prometheus' file to be extracted.
	if out, err := extract.CombinedOutput(); err != nil {
		return fmt.Errorf("un-compressing the prometheus binary: command %s: output %s: %w", extract.String(), out, err)
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
		return false, fmt.Errorf("executing the prometheus binary: command %s: output %s: %w", checkVersion.String(), out, err)
	}

	return strings.Contains(string(out), version), nil
}
