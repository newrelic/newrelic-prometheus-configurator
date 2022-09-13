package scrapejob

import (
	"testing"

	"github.com/newrelic/newrelic-prometheus-configurator/internal/sharding"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIncludeShardingRuleRules(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name           string
		job            Job
		shardingConfig sharding.Config
		rulesExpected  bool
	}{
		{
			name:           "Skip sharding",
			job:            Job{SkipSharding: true},
			shardingConfig: sharding.Config{TotalShardsCount: 2, ShardIndex: "0"},
			rulesExpected:  false,
		},
		{
			name:           "Only one shard",
			job:            Job{SkipSharding: false},
			shardingConfig: sharding.Config{TotalShardsCount: 1, ShardIndex: "0"},
			rulesExpected:  false,
		},
		{
			name:           "Include rules",
			job:            Job{SkipSharding: false},
			shardingConfig: sharding.Config{TotalShardsCount: 2, ShardIndex: "0"},
			rulesExpected:  true,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			relabelConfigCount := len(tc.job.RelabelConfigs)
			promJob := tc.job.includeShardingRules(tc.shardingConfig, tc.job.Job)
			if tc.rulesExpected {
				require.Greater(t, len(promJob.RelabelConfigs), relabelConfigCount)
				assert.Equal(t, tc.shardingConfig.RelabelConfigs(), promJob.RelabelConfigs[:2])
			} else {
				assert.Len(t, promJob.RelabelConfigs, relabelConfigCount)
			}
		})
	}
}
