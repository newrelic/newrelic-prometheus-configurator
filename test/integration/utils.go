//go:build integration_test

package integration

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"os/exec"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

func runConfigurator(t *testing.T, nrConfigConfig string) string {
	t.Helper()

	tempDir := t.TempDir()
	nrConfigConfigPath := path.Join(tempDir, "nrConfig.yml")
	promConfigConfigPath := path.Join(tempDir, "promConfig.yml")

	readOnly := 0o444
	err := ioutil.WriteFile(nrConfigConfigPath, []byte(nrConfigConfig), fs.FileMode(readOnly))
	require.NoError(t, err)

	//nolint:gosec
	configurator := exec.Command(
		"go",
		"run",
		"../../cmd/configurator",
		fmt.Sprintf("--input=%s", nrConfigConfigPath),
		fmt.Sprintf("--output=%s", promConfigConfigPath),
	)

	out, err := configurator.CombinedOutput()
	require.NoError(t, err, string(out))

	return promConfigConfigPath
}
