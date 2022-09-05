package configurator

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

const (
	DataSourceNameEnvKey = "NR_PROM_DATA_SOURCE_NAME"
	LicenseKeyEnvKey     = "NR_PROM_LICENSE_KEY"
)

var (
	ErrNoLicenseKeyFound = fmt.Errorf(
		"licenseKey was not set neither in yaml config or %s environment variable", LicenseKeyEnvKey,
	)
	ErrInvalidShardingKind = errors.New("the only supported kind of sharding is hash")
)

// BuildPromConfig builds the prometheus config prometheusConfig from the provided nrConfig, it holds "first level" transformations
// required to obtain a valid prometheus configuration.
func BuildPromConfig(nrConfig *NrConfig) (*PromConfig, error) {
	expand(nrConfig)

	if err := validate(nrConfig); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	prometheusConfig := &PromConfig{
		RemoteWrite:  []RawPromConfig{nrConfig.RemoteWrite.Build(nrConfig.DataSourceName)},
		GlobalConfig: nrConfig.Common,
	}

	prometheusConfig.RemoteWrite = append(prometheusConfig.RemoteWrite, nrConfig.ExtraRemoteWrite...)

	for _, staticTargets := range nrConfig.StaticTargets.Build() {
		job := nrConfig.Sharding.IncludeShardingRules(staticTargets)
		prometheusConfig.ScrapeConfigs = append(prometheusConfig.ScrapeConfigs, job)
	}

	k8sJobs, err := nrConfig.Kubernetes.Build()
	if err != nil {
		return prometheusConfig, fmt.Errorf("building k8s config: %w", err)
	}

	for _, K8sJob := range k8sJobs {
		j := nrConfig.Sharding.IncludeShardingRules(K8sJob)
		prometheusConfig.ScrapeConfigs = append(prometheusConfig.ScrapeConfigs, j)
	}

	prometheusConfig.ScrapeConfigs = append(prometheusConfig.ScrapeConfigs, nrConfig.ExtraScrapeConfigs...)

	return prometheusConfig, nil
}

// expand replace some specifics configs that can be defined by env variables.
func expand(config *NrConfig) {
	if licenseKey := os.Getenv(LicenseKeyEnvKey); licenseKey != "" {
		config.RemoteWrite.LicenseKey = licenseKey
	}

	dataSourceName := os.Getenv(DataSourceNameEnvKey)

	if dataSourceName != "" {
		config.DataSourceName = dataSourceName
	}

	if dataSourceName != "" && config.Sharding.TotalShardsCount > 1 {
		shardIndex := getIndexFromDataSourceName(dataSourceName)
		config.Sharding.ShardIndex = shardIndex
	}
}

func validate(config *NrConfig) error {
	if config.RemoteWrite.LicenseKey == "" {
		return ErrNoLicenseKeyFound
	}

	// Defaults to kind hash in case it's empty.
	if config.Sharding.Kind == "" {
		config.Sharding.Kind = "hash"
	}

	if config.Sharding.Kind != "hash" {
		return ErrInvalidShardingKind
	}

	return nil
}

// getIndexFromDataSourceName returns the corresponding shard index from the DataSourceNameEnvKey env var.
// This function assumes the name follows the k8s name convention being `-` separated.
// E.g. by running two shards the pod names will be the following:
//
//	1: newrelic-prometheus-0
//	2: newrelic-prometheus-1
func getIndexFromDataSourceName(dataSourceName string) string {
	parts := strings.Split(dataSourceName, "-")
	if len(parts) <= 1 {
		return ""
	}

	return parts[len(parts)-1]
}
