// Copyright 2022 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/newrelic-forks/newrelic-prometheus/configurator"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

const (
	inputErrCode = iota + 1
	outputErrCode
	parseErrCode
)

func main() {
	logger := log.StandardLogger()

	inputFlag := flag.String("input", "", "Input file to load the configuration from, defaults to stdin.")
	outputFlag := flag.String("output", "", "Output file to use as output, defaults to stdout.")
	flag.Parse()

	input, err := readInput(*inputFlag)
	if err != nil {
		logger.Errorf("Error loading the input: %s", err)
		os.Exit(inputErrCode)
	}

	output, err := configurator.BuildOutput(input)
	if err != nil {
		logger.Errorf("Error parsing the configuration: %s", err)
		os.Exit(parseErrCode)
	}

	if err := writeOutput(*outputFlag, output); err != nil {
		logger.Errorf("Error writing the output configuration: %s", err)
		os.Exit(outputErrCode)
	}
}

func readInput(inputPath string) (*configurator.Input, error) {
	input := &configurator.Input{}

	if inputPath == "" {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			return nil, fmt.Errorf("could not read from stdin: %w", err)
		}

		err = yaml.Unmarshal(data, input)
		if err != nil {
			return nil, fmt.Errorf("yaml input could not be loaded: %w", err)
		}

		return input, nil
	}

	fileReader, err := os.Open(inputPath)
	if err != nil {
		return nil, fmt.Errorf("the input file could not be opened: %w", err)
	}

	data, err := io.ReadAll(fileReader)
	if err != nil {
		return nil, fmt.Errorf("could not read from the input file: %w", err)
	}

	if err = fileReader.Close(); err != nil {
		return nil, fmt.Errorf("could not close the input file: %w", err)
	}

	err = yaml.Unmarshal(data, input)
	if err != nil {
		return nil, fmt.Errorf("yaml input could not be loaded: %w", err)
	}

	return input, nil
}

func writeOutput(outputPath string, output *configurator.Output) error {
	data, err := yaml.Marshal(output)
	if err != nil {
		return fmt.Errorf("marshaling output: %w", err)
	}

	if outputPath == "" {
		if _, err = os.Stdout.Write(data); err != nil {
			return fmt.Errorf("could not to stdout: %w", err)
		}
		return nil
	}

	fileWriter, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("the output file cannot be created: %w", err)
	}

	if _, err := fileWriter.Write(data); err != nil {
		return fmt.Errorf("could not write the output: %w", err)
	}

	if err := fileWriter.Close(); err != nil {
		return fmt.Errorf("could not close the output file: %w", err)
	}

	return nil
}
