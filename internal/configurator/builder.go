package configurator

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

const (
	ChartVersionEnvKey   = "NR_PROM_CHART_VERSION"
	DataSourceNameEnvKey = "NR_PROM_DATA_SOURCE_NAME"
	LicenseKeyEnvKey     = "NR_PROM_LICENSE_KEY"
	ProxyUrlEnvKey       = "NR_PROM_PROXY_URL"
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

	remoteWrite, err := nrConfig.RemoteWrite.Build()
	if err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	prometheusConfig := &PromConfig{
		RemoteWrite:  []RawPromConfig{remoteWrite},
		GlobalConfig: nrConfig.Common,
	}

	prometheusConfig.RemoteWrite = append(prometheusConfig.RemoteWrite, nrConfig.ExtraRemoteWrite...)

	for _, job := range nrConfig.StaticTargets.Build(nrConfig.Sharding) {
		prometheusConfig.ScrapeConfigs = append(prometheusConfig.ScrapeConfigs, job)
	}

	k8sJobs, err := nrConfig.Kubernetes.Build(nrConfig.Sharding)
	if err != nil {
		return prometheusConfig, fmt.Errorf("building k8s config: %w", err)
	}

	for _, job := range k8sJobs {
		prometheusConfig.ScrapeConfigs = append(prometheusConfig.ScrapeConfigs, job)
	}

	prometheusConfig.ScrapeConfigs = append(prometheusConfig.ScrapeConfigs, nrConfig.ExtraScrapeConfigs...)

	return prometheusConfig, nil
}

// expand replace some specifics configs that can be defined by env variables.
func expand(config *NrConfig) {
	if licenseKey := os.Getenv(LicenseKeyEnvKey); licenseKey != "" {
		config.RemoteWrite.LicenseKey = licenseKey
	}

	if proxyurl := os.Getenv(ProxyUrlEnvKey); proxyurl != "" {
		config.RemoteWrite.ProxyURL = proxyurl
	}

	chartVersion := os.Getenv(ChartVersionEnvKey)
	if chartVersion != "" {
		config.RemoteWrite.ChartVersion = chartVersion
	}

	dataSourceName := os.Getenv(DataSourceNameEnvKey)
	if dataSourceName != "" {
		config.RemoteWrite.DataSourceName = dataSourceName
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
//	1: newrelic-prometheus-agent-0
//	2: newrelic-prometheus-agent-1
func getIndexFromDataSourceName(dataSourceName string) string {
	parts := strings.Split(dataSourceName, "-")
	if len(parts) <= 1 {
		return ""
	}

	return parts[len(parts)-1]
}
