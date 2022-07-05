package configurator

import (
	"io/ioutil"
	"testing"

	prometheusConfig "github.com/prometheus/prometheus/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestParser(t *testing.T) {
	// it relies on testdata/<placeholder>.yaml and testdata/<placeholder>.expected.yaml
	testCases := []string{
		"remote-write-test",
	}
	for _, c := range testCases {
		t.Run(c, func(t *testing.T) {
			inputFile := "testdata/" + c + ".yaml"
			expectedFile := "testdata/" + c + ".expected.yaml"
			input, err := ioutil.ReadFile(inputFile)
			require.NoError(t, err)
			expected, err := ioutil.ReadFile(expectedFile)
			require.NoError(t, err)
			output, err := Parse(input)
			require.NoError(t, err)
			assertYamlOutputsAreEqual(t, expected, output)
			assertIsPrometheusConfig(t, output)
		})
	}
}

func TestParserInvalidInputYamlError(t *testing.T) {
	input := []byte(`}invalid-yml`)
	_, err := Parse(input)
	assert.Error(t, err)
}

func assertYamlOutputsAreEqual(t *testing.T, y1, y2 []byte) {
	var o1, o2 Output
	require.NoError(t, yaml.Unmarshal(y1, &o1))
	require.NoError(t, yaml.Unmarshal(y2, &o2))
	assert.EqualValues(t, o1, o2)
}

func assertIsPrometheusConfig(t *testing.T, y []byte) {
	cfg := prometheusConfig.Config{}
	require.NoError(t, yaml.Unmarshal(y, &cfg), "At least config should be marshaled to prometheus config")
}
