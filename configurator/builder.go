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

// BuildOutput builds the prometheus config output from the provided input, it holds "first level" transformations
// required to obtain a valid prometheus configuration.
func BuildOutput(input *Input) (*Output, error) {
	expand(input)

	if err := validate(input); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	output := &Output{
		RemoteWrite:  []RawPromConfig{input.RemoteWrite.Build(input.DataSourceName)},
		GlobalConfig: input.Common,
	}

	output.RemoteWrite = append(output.RemoteWrite, input.ExtraRemoteWrite...)

	for _, staticTargets := range input.StaticTargets.Build() {
		job := input.Sharding.IncludeShardingRules(staticTargets)
		output.ScrapeConfigs = append(output.ScrapeConfigs, job)
	}

	k8sJobs, err := input.Kubernetes.Build()
	if err != nil {
		return output, fmt.Errorf("building k8s config: %w", err)
	}

	for _, K8sJob := range k8sJobs {
		j := input.Sharding.IncludeShardingRules(K8sJob)
		output.ScrapeConfigs = append(output.ScrapeConfigs, j)
	}

	output.ScrapeConfigs = append(output.ScrapeConfigs, input.ExtraScrapeConfigs...)

	return output, nil
}

// expand replace some specifics configs that can be defined by env variables.
func expand(config *Input) {
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

func validate(config *Input) error {
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
