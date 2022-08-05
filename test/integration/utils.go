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

func runConfigurator(t *testing.T, inputConfig string) string {
	t.Helper()

	tempDir := t.TempDir()
	inputConfigPath := path.Join(tempDir, "input.yml")
	outputConfigPath := path.Join(tempDir, "output.yml")

	readOnly := 0o444
	err := ioutil.WriteFile(inputConfigPath, []byte(inputConfig), fs.FileMode(readOnly))
	require.NoError(t, err)

	//nolint:gosec
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
