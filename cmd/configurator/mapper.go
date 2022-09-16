package main

import (
	"fmt"
	"net/url"
	"time"

	"github.com/newrelic/newrelic-prometheus-configurator/internal/configurator"
	"github.com/newrelic/newrelic-prometheus-configurator/internal/guesser"
	"github.com/newrelic/newrelic-prometheus-configurator/internal/promapi"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type mapper struct {
	prometheusReload bool
	prometheusURL    url.URL
	interval         time.Duration
	relabelsFilePath string

	nrConfig       *configurator.NrConfig
	promConfigPath string

	logger *log.Logger
}

func (m mapper) run() error {
	pc := promapi.New(m.prometheusURL)

	mrg := guesser.New(pc)

	ticker := time.NewTicker(m.interval)
	for range ticker.C {
		relabelRules, err := mrg.UpdatedMetricTypeRules()
		if err != nil {
			return fmt.Errorf("building metric type rules: %w", err)
		}

		// TODO it could happen that after some rules has.
		if len(relabelRules) == 0 {
			m.logger.Debug("metric type rules has not changed from last iteration ")
			continue
		}

		relabelRulesYaml, err := yaml.Marshal(relabelRules)
		if err != nil {
			m.logger.Errorf("Error marshaling relabel rules mappings: %s", err)
		}

		if err := writeConfig(m.relabelsFilePath, relabelRulesYaml); err != nil {
			m.logger.Errorf("Error writing relabel rules mappings file: %s", err)
		}

		m.logger.Infof("metric rules mapping updated on file: %s", m.relabelsFilePath)

		if m.prometheusReload {
			newNrConfig := m.nrConfig

			newNrConfig.RemoteWrite.ExtraWriteRelabelConfigs = append(newNrConfig.RemoteWrite.ExtraWriteRelabelConfigs, relabelRules...)

			prometheusConfig, err := configurator.BuildPromConfig(newNrConfig)
			if err != nil {
				return fmt.Errorf("parsing the configuration: %w", err)
			}

			if err := writePromConfig(m.promConfigPath, prometheusConfig); err != nil {
				m.logger.Errorf("Error writing the prometheusConfig configuration: %s", err)
			}

			m.logger.Infof("Sending request to reload Prometheus")
			if err := pc.Reload(); err != nil {
				m.logger.Errorf("Error reloading Prometheus Server: %s", err)
			}
		}
	}

	return nil
}
