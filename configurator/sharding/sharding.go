package sharding

import (
	"fmt"
	"github.com/newrelic-forks/newrelic-prometheus/configurator/promcfg"
)

// Config defines all the NewRelic's sharding options.
type Config struct {
	Kind             string `yaml:"kind"`
	TotalShardsCount int    `yaml:"total_shards_count"`
	ShardIndex       string `yaml:"shard_index"`
}

// IncludeShardingRules prepends the proper sharding relabel configs for the given job.
func (c Config) IncludeShardingRules(job promcfg.Job) promcfg.Job {
	// Skip the relabeling if at least there are not 2 shards.
	if c.TotalShardsCount < 2 { //nolint: gomnd
		return job
	}

	var additionalHashModLabels []string

	// In case of a kubernetes job and the role being endpoints, we need to add an extra source
	// label for the hash mod relabel config.
	if len(job.KubernetesSdConfigs) > 0 && job.KubernetesSdConfigs[0].Role == "endpoints" {
		additionalHashModLabels = []string{"_meta_kubernetes_service_name"}
	}

	job.RelabelConfigs = append(c.RelabelConfigs(additionalHashModLabels), job.RelabelConfigs...)

	return job
}

func (c Config) RelabelConfigs(additionalHashModLabels []string) []promcfg.RelabelConfig {
	return []promcfg.RelabelConfig{
		{
			SourceLabels: append([]string{"__address__"}, additionalHashModLabels...),
			Modulus:      c.TotalShardsCount,
			Action:       "hashmod",
			TargetLabel:  "__tmp_hash",
		},
		{
			SourceLabels: []string{"__tmp_hash"},
			Regex:        fmt.Sprintf("^%v$", c.ShardIndex),
			Action:       "keep",
		},
	}
}
