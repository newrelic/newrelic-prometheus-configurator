package configurator

import (
	"bytes"
	"fmt"

	"gopkg.in/yaml.v3"
)

const (
	yamlEncoderIndent = 2
)

// Parse loads a yaml input and returns the corresponding prometheus-agent yaml.
func Parse(in []byte) ([]byte, error) {
	// load the yaml input
	input := &Input{}
	err := yaml.Unmarshal(in, input)
	if err != nil {
		return nil, fmt.Errorf("yaml input could not be loaded: %w", err)
	}
	// builds the corresponding output
	output, err := BuildOutput(input)
	if err != nil {
		return nil, fmt.Errorf("output could not be built: %w", err)
	}
	// parse it to yml
	parsed, err := toYaml(&output)
	if err != nil {
		return nil, fmt.Errorf("output could not be encoded to yaml: %w", err)
	}
	return parsed, nil
}

func toYaml(output *Output) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := yaml.NewEncoder(&buffer)
	encoder.SetIndent(yamlEncoderIndent)
	if err := encoder.Encode(output); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}
