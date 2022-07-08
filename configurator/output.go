package configurator

// Output holds all configuration information in prometheus format which can be directly exported valid prometheus config yaml.
type Output struct {
	RemoteWrite []interface{} `yaml:"remote_write"`
}

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
