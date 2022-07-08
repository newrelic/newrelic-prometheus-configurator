package configurator

// Output holds all configuration information in prometheus format which can be directly marshaled to a valid yaml configuration.
type Output struct {
	RemoteWrite []interface{} `yaml:"remote_write"`
}

// BuildOutput builds the prometheus config output from the provided input, it holds "first level" transformations required
// to obtain a valid prometheus configuration.
func BuildOutput(input *Input) (Output, error) {
	output := Output{
		RemoteWrite: []interface{}{BuildRemoteWriteOutput(input)},
	}
	// Include extra remote-write configs
	for _, extraRemoteWriteConfig := range input.ExtraRemoteWrite {
		output.RemoteWrite = append(output.RemoteWrite, extraRemoteWriteConfig)
	}
	return output, nil
}
