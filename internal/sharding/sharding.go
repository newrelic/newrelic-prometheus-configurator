package sharding

import (
	"fmt"

	"github.com/newrelic/newrelic-prometheus-configurator/internal/promcfg"
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
	if c.TotalShardsCount <= 1 {
		return job
	}

	job.RelabelConfigs = append(c.RelabelConfigs(), job.RelabelConfigs...)

	return job
}

func (c Config) RelabelConfigs() []promcfg.RelabelConfig {
	return []promcfg.RelabelConfig{
		{
			SourceLabels: []string{"__address__"},
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
