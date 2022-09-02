// Copyright 2022 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"runtime"

	"github.com/newrelic/newrelic-prometheus-configurator/internal/configurator"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

const (
	integrationName = "Prometheus-Configurator"

	nrConfigErrCode = iota + 1
	promConfigErrCode
	parseErrCode
)

var (
	//nolint:gochecknoglobals
	integrationVersion = "0.0.0"
	//nolint:gochecknoglobals
	gitCommit = ""
	//nolint:gochecknoglobals
	buildDate = ""
)

func main() {
	logger := log.StandardLogger()

	nrConfigFlag := flag.String("input", "", "Input file to load the configuration from, defaults to stdin.")
	promConfigFlag := flag.String("output", "", "Output file to use as prometheus config, defaults to stdout.")
	flag.Parse()

	logger.Infof(
		"New Relic %s integration Version: %s, Platform: %s, GoVersion: %s, GitCommit: %s, BuildDate: %s\n",
		integrationName,
		integrationVersion,
		fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
		runtime.Version(),
		gitCommit,
		buildDate)

	nrConfig, err := readNrConfig(*nrConfigFlag)
	if err != nil {
		logger.Errorf("Error loading the nrConfig: %s", err)
		os.Exit(nrConfigErrCode)
	}

	promConfig, err := configurator.BuildPromConfig(nrConfig)
	if err != nil {
		logger.Errorf("Error parsing the configuration: %s", err)
		os.Exit(parseErrCode)
	}

	if err := writePromConfig(*promConfigFlag, promConfig); err != nil {
		logger.Errorf("Error writing the promConfig configuration: %s", err)
		os.Exit(promConfigErrCode)
	}
}

func readNrConfig(nrConfigPath string) (*configurator.NrConfig, error) {
	nrConfig := &configurator.NrConfig{}

	if nrConfigPath == "" {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			return nil, fmt.Errorf("could not read from stdin: %w", err)
		}

		err = yaml.Unmarshal(data, nrConfig)
		if err != nil {
			return nil, fmt.Errorf("yaml nrConfig could not be loaded: %w", err)
		}

		return nrConfig, nil
	}

	fileReader, err := os.Open(nrConfigPath)
	if err != nil {
		return nil, fmt.Errorf("the nrConfig file could not be opened: %w", err)
	}

	data, err := io.ReadAll(fileReader)
	if err != nil {
		return nil, fmt.Errorf("could not read from the nrConfig file: %w", err)
	}

	if err = fileReader.Close(); err != nil {
		return nil, fmt.Errorf("could not close the nrConfig file: %w", err)
	}

	err = yaml.Unmarshal(data, nrConfig)
	if err != nil {
		return nil, fmt.Errorf("yaml nrConfig could not be loaded: %w", err)
	}

	return nrConfig, nil
}

func writePromConfig(promConfigPath string, promConfig *configurator.PromConfig) error {
	data, err := yaml.Marshal(promConfig)
	if err != nil {
		return fmt.Errorf("marshaling promConfig: %w", err)
	}

	if promConfigPath == "" {
		if _, err = os.Stdout.Write(data); err != nil {
			return fmt.Errorf("could not to stdout: %w", err)
		}
		return nil
	}

	fileWriter, err := os.Create(promConfigPath)
	if err != nil {
		return fmt.Errorf("the promConfig file cannot be created: %w", err)
	}

	if _, err := fileWriter.Write(data); err != nil {
		return fmt.Errorf("could not write the promConfig: %w", err)
	}

	if err := fileWriter.Close(); err != nil {
		return fmt.Errorf("could not close the promConfig file: %w", err)
	}

	return nil
}
