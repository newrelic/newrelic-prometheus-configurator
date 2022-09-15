package scrapejob_test

import (
	"testing"

	"github.com/newrelic/newrelic-prometheus-configurator/internal/promcfg"
	"github.com/newrelic/newrelic-prometheus-configurator/internal/scrapejob"
	"github.com/newrelic/newrelic-prometheus-configurator/internal/sharding"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIncludeShardingRuleRules(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name                         string
		job                          scrapejob.Job
		shardingConfig               sharding.Config
		shardingRulesExpected        bool
		additionalRelabelConfig      []promcfg.RelabelConfig
		expectedRelabelConfigs       []promcfg.RelabelConfig
		expectedMetricRelabelConfigs []promcfg.RelabelConfig
	}{
		{
			name: "Skip sharding",
			job: scrapejob.Job{
				Job:          promcfg.Job{RelabelConfigs: []promcfg.RelabelConfig{{Action: "from-job"}}},
				SkipSharding: true,
			},
			shardingConfig:         sharding.Config{TotalShardsCount: 2, ShardIndex: "0"},
			shardingRulesExpected:  false,
			expectedRelabelConfigs: []promcfg.RelabelConfig{{Action: "from-job"}},
		},
		{
			name: "Only one shard",
			job: scrapejob.Job{
				Job:          promcfg.Job{RelabelConfigs: []promcfg.RelabelConfig{{Action: "from-job"}}},
				SkipSharding: false,
			},
			shardingConfig:         sharding.Config{TotalShardsCount: 1, ShardIndex: "0"},
			shardingRulesExpected:  false,
			expectedRelabelConfigs: []promcfg.RelabelConfig{{Action: "from-job"}},
		},
		{
			name: "Include rules",
			job: scrapejob.Job{
				Job:          promcfg.Job{RelabelConfigs: []promcfg.RelabelConfig{{Action: "from-job"}}},
				SkipSharding: false,
			},
			shardingConfig:        sharding.Config{TotalShardsCount: 2, ShardIndex: "0"},
			shardingRulesExpected: true,
			// Sharding rules are checked separately.
			expectedRelabelConfigs: []promcfg.RelabelConfig{{Action: "from-job"}},
		},
		{
			name: "Rules from all sources",
			job: scrapejob.Job{
				Job: promcfg.Job{
					RelabelConfigs:       []promcfg.RelabelConfig{{Action: "from-job"}},
					MetricRelabelConfigs: []promcfg.RelabelConfig{{Action: "metrics-from-job"}},
				},
				ExtraRelabelConfigs:       []promcfg.RelabelConfig{{Action: "from-extra-1"}, {Action: "from-extra-2"}},
				ExtraMetricRelabelConfigs: []promcfg.RelabelConfig{{Action: "extra-metrics"}},
			},
			additionalRelabelConfig: []promcfg.RelabelConfig{{Action: "from-additional"}},
			expectedRelabelConfigs: []promcfg.RelabelConfig{
				// Sharding rules are checked separately.
				{Action: "from-job"}, {Action: "from-additional"}, {Action: "from-extra-1"}, {Action: "from-extra-2"},
			},
			expectedMetricRelabelConfigs: []promcfg.RelabelConfig{{Action: "metrics-from-job"}, {Action: "extra-metrics"}},
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			job := tc.job.
				WithName("job-name").
				WithRelabelConfigs(tc.additionalRelabelConfig).
				BuildPrometheusJob(tc.shardingConfig)

			assert.Equal(t, "job-name", job.JobName)
			if tc.shardingRulesExpected { // Sharding rules should be included at the beginning when expected.
				shardingRules := tc.shardingConfig.RelabelConfigs()
				shardingRulesCount := len(shardingRules)
				require.GreaterOrEqual(t, len(job.RelabelConfigs), len(shardingRules))
				assert.Equal(t, shardingRules, job.RelabelConfigs[:shardingRulesCount])
				assert.EqualValues(t, tc.expectedRelabelConfigs, job.RelabelConfigs[shardingRulesCount:])
			} else {
				assert.Equal(t, tc.expectedRelabelConfigs, job.RelabelConfigs)
			}
			assert.Equal(t, tc.expectedMetricRelabelConfigs, job.MetricRelabelConfigs)
		})
	}
}
