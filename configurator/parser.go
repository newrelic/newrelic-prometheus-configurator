package configurator

import (
	"bytes"

	"gopkg.in/yaml.v3"
)

const (
	yamlEncoderIndent = 2
)

// Parse loads a yaml input and returns the corresponding prometheus-agent yaml.
func Parse(in []byte) ([]byte, error) {
	// load the yaml input
	input, err := loadYaml(in)
	if err != nil {
		return nil, err
	}
	// builds the corresponding output
	output, err := BuildOutput(input)
	if err != nil {
		return nil, err
	}
	// parse it to yml
	return toYaml(&output)
}

func loadYaml(in []byte) (*Input, error) {
	input := &Input{}
	err := yaml.Unmarshal(in, input)
	return input, err
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
