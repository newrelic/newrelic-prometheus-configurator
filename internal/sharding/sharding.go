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

// ShouldIncludeShardingRules returns true when additional rules are needed for the current configuration.
func (c Config) ShouldIncludeShardingRules() bool {
	return c.TotalShardsCount > 1
}

func (c Config) RelabelConfigs() []promcfg.RelabelConfig {
	return []promcfg.RelabelConfig{
		{
			SourceLabels: []string{"__address__"},
			Regex:        `(\d{1,3}\.\d{1,3}\.\d{1,3}.\d{1,3})(?::\d+)?`,
			Action:       "replace",
			TargetLabel:  "__tmp_hash",
		},
		{
			SourceLabels: []string{"__tmp_hash"},
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
