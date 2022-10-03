// Copyright 2022 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"flag"
	"fmt"
	"io"
	"net/url"
	"os"
	"runtime"
	"time"

	"github.com/newrelic/newrelic-prometheus-configurator/internal/configurator"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

const (
	integrationName = "Prometheus-Configurator"

	nrConfigErrCode = iota + 1
	prometheusConfigErrCode
	parseErrCode
	mapperErrCode
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

	mapperEnable := flag.Bool("mapper_enable", false, "Activates the Metric Type Mapper.")
	mapperPrometheusURL := flag.String("mapper_prometheus_url", "localhost:9090", "Prometheus server 'host:port'.")
	mapperPrometheusReload := flag.Bool("mapper_reload", false, "[Experimental] Reloads Prometheus after any change in metrics type mappings.")
	mapperInterval := flag.String("mapper_interval", "", "Metric Type Mapper executions interval.")
	mapperRelabelsFilePath := flag.String("mapper_file_path", "", "File path of the write_relabel_config snippet")

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

	prometheusConfigData, err := yaml.Marshal(prometheusConfig)
	if err != nil {
		logger.Errorf("Error encoding the prometheusConfig configuration: %s", err)
		os.Exit(prometheusConfigErrCode)
	}

	if err := writeConfig(*prometheusConfigFlag, prometheusConfigData); err != nil {
		logger.Errorf("Error writing the prometheusConfig configuration: %s", err)
		os.Exit(prometheusConfigErrCode)
	}

	if *mapperEnable {
		relabelInterval, err := time.ParseDuration(*mapperInterval)
		if err != nil {
			logger.Fatalf("Error mapper_interval '%s' invalid time expression: %s", *mapperInterval, err)
		}

		prometheusURL, err := url.ParseRequestURI(*mapperPrometheusURL)
		if err != nil {
			logger.Fatalf("Error parsing mapper_prometheus_url: %s", err)
		}

		m := mapper{
			prometheusReload: *mapperPrometheusReload,
			prometheusURL:    *prometheusURL,
			interval:         relabelInterval,
			relabelsFilePath: *mapperRelabelsFilePath,
			nrConfig:         nrConfig,
			promConfigPath:   *prometheusConfigFlag,
			logger:           logger,
		}

		if err := m.run(); err != nil {
			logger.Errorf("Error running the mapper: %s", err)
			os.Exit(mapperErrCode)
		}
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

	if err := writeConfig(prometheusConfigPath, data); err != nil {
		return fmt.Errorf("writing config: %w", err)
	}

	return nil
}

func writeConfig(path string, data []byte) error {
	if path == "" {
		if _, err := os.Stdout.Write(data); err != nil {
			return fmt.Errorf("could not to stdout: %w", err)
		}
		return nil
	}

	fileWriter, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("creating file: %w", err)
	}

	if _, err := fileWriter.Write(data); err != nil {
		return fmt.Errorf("writing file: %w", err)
	}

	if err := fileWriter.Close(); err != nil {
		return fmt.Errorf("closing file: %w", err)
	}

	return nil
}
