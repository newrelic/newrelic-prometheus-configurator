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
	prometheusConfigErrCode
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
	prometheusConfigFlag := flag.String("output", "", "Output file to use as prometheus config, defaults to stdout.")
	verboseLog := flag.Bool("verbose", false, "Sets log level to debug.")
	flag.Parse()

	if *verboseLog {
		logger.SetLevel(log.DebugLevel)
	}

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

	prometheusConfig, err := configurator.BuildPromConfig(nrConfig)
	if err != nil {
		logger.Errorf("Error parsing the configuration: %s", err)
		os.Exit(parseErrCode)
	}

	if err := writePromConfig(*prometheusConfigFlag, prometheusConfig); err != nil {
		logger.Errorf("Error writing the prometheusConfig configuration: %s", err)
		os.Exit(prometheusConfigErrCode)
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

func writePromConfig(prometheusConfigPath string, prometheusConfig *configurator.PromConfig) error {
	data, err := yaml.Marshal(prometheusConfig)
	if err != nil {
		return fmt.Errorf("marshaling prometheusConfig: %w", err)
	}

	if prometheusConfigPath == "" {
		if _, err = os.Stdout.Write(data); err != nil {
			return fmt.Errorf("could not to stdout: %w", err)
		}
		return nil
	}

	fileWriter, err := os.Create(prometheusConfigPath)
	if err != nil {
		return fmt.Errorf("the prometheusConfig file cannot be created: %w", err)
	}

	if _, err := fileWriter.Write(data); err != nil {
		return fmt.Errorf("could not write the prometheusConfig: %w", err)
	}

	if err := fileWriter.Close(); err != nil {
		return fmt.Errorf("could not close the prometheusConfig file: %w", err)
	}

	return nil
}
