// Copyright 2022 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"flag"
	"fmt"
	"io"
	"newrelic-prometheus/configurator"
	"os"

	"github.com/scaleway/scaleway-sdk-go/logger"
	log "github.com/sirupsen/logrus"
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

	output, err := configurator.Parse(input)
	if err != nil {
		logger.Errorf("Error parsing the configuration: %s", err)
		os.Exit(parseErrCode)
	}

	err = writeOutput(*outputFlag, output)
	if err != nil {
		logger.Errorf("Error writing the output configuration: %s", err)
		os.Exit(outputErrCode)
	}
}

func readInput(inputPath string) ([]byte, error) {
	if inputPath == "" {
		input, err := io.ReadAll(os.Stdin)
		if err != nil {
			return nil, fmt.Errorf("could not read from stdin: %w", err)
		}

		return input, nil
	}

	fileReader, err := os.Open(inputPath)
	if err != nil {
		return nil, fmt.Errorf("the input file could not be opened: %w", err)
	}

	defer closeLoggingErr(fileReader)

	input, err := io.ReadAll(fileReader)
	if err != nil {
		return nil, fmt.Errorf("could not read from the input file: %w", err)
	}

	return input, nil
}

func writeOutput(outputPath string, output []byte) error {
	if outputPath == "" {
		_, err := os.Stdout.Write(output)

		return fmt.Errorf("could not to stdout: %w", err)
	}

	fileWriter, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("the output file cannot be created: %w", err)
	}

	defer closeLoggingErr(fileWriter)

	_, err = fileWriter.Write(output)
	if err != nil {
		return fmt.Errorf("could not write the output: %w", err)
	}

	return nil
}

func closeLoggingErr(f *os.File) {
	if err := f.Close(); err != nil {
		logger.Errorf("Fail closing file: %s", err)
	}
}
