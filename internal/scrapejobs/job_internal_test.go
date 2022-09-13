package scrapejobs

import (
	"testing"

	"github.com/newrelic/newrelic-prometheus-configurator/internal/sharding"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIncludeShardingConfig(t *testing.T) {
	t.Parallel()

	cases := []struct {
		Name           string
		Job            Job
		ShardingConfig sharding.Config
		RulesExpected  bool
	}{
		{
			Name:           "Skip sharding",
			Job:            Job{SkipSharding: true},
			ShardingConfig: sharding.Config{TotalShardsCount: 2, ShardIndex: "0"},
			RulesExpected:  false,
		},
		{
			Name:           "Only one shard",
			Job:            Job{SkipSharding: false},
			ShardingConfig: sharding.Config{TotalShardsCount: 1, ShardIndex: "0"},
			RulesExpected:  false,
		},
		{
			Name:           "Include rules",
			Job:            Job{SkipSharding: false},
			ShardingConfig: sharding.Config{TotalShardsCount: 2, ShardIndex: "0"},
			RulesExpected:  true,
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			relabelConfigCount := len(tc.Job.RelabelConfigs)
			promJob := tc.Job.includeShardingRules(tc.ShardingConfig, tc.Job.Job)
			if tc.RulesExpected {
				require.Greater(t, len(promJob.RelabelConfigs), relabelConfigCount)
				assert.Equal(t, tc.ShardingConfig.RelabelConfigs(), promJob.RelabelConfigs[:2])
			} else {
				assert.Len(t, promJob.RelabelConfigs, relabelConfigCount)
			}
		})
	}
}
