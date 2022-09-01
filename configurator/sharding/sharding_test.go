// Copyright 2022 New Relic Corporation. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package sharding_test

import (
	"testing"

	"github.com/newrelic-forks/newrelic-prometheus/configurator/promcfg"
	"github.com/newrelic-forks/newrelic-prometheus/configurator/sharding"
	"github.com/stretchr/testify/assert"
)

func TestConfig_IncludeShardingRules(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name                   string
		config                 sharding.Config
		job                    promcfg.Job
		expectedRelabelConfigs []promcfg.RelabelConfig
		assert                 func(t *testing.T, job promcfg.Job, expectedRelabelConfigs []promcfg.RelabelConfig)
	}{
		{
			name: "NotAddingRulesWhenShardsNotSet",
			job:  promcfg.Job{},
			assert: func(t *testing.T, job promcfg.Job, expectedRelabelConfigs []promcfg.RelabelConfig) {
				assert.Equal(t, 0, len(job.RelabelConfigs))
			},
		},
		{
			name: "AddingStaticTargetsRules",
			config: sharding.Config{
				Kind:             "hash",
				TotalShardsCount: 2,
				ShardIndex:       "1",
			},
			job: promcfg.Job{},
			expectedRelabelConfigs: []promcfg.RelabelConfig{
				{
					SourceLabels: []string{"__address__"},
					Modulus:      2,
					Action:       "hashmod",
					TargetLabel:  "__tmp_hash",
				},
				{
					SourceLabels: []string{"__tmp_hash"},
					Regex:        "^1$",
					Action:       "keep",
				},
			},
			assert: func(t *testing.T, job promcfg.Job, expectedRelabelConfigs []promcfg.RelabelConfig) {
				assert.Equal(t, expectedRelabelConfigs, job.RelabelConfigs)
			},
		},
		{
			name: "Addingk8sEndpointsRules",
			config: sharding.Config{
				Kind:             "hash",
				TotalShardsCount: 2,
				ShardIndex:       "1",
			},
			job: promcfg.Job{
				KubernetesSdConfigs: []promcfg.KubernetesSdConfig{
					{
						Role: "endpoints",
					},
				},
			},
			expectedRelabelConfigs: []promcfg.RelabelConfig{
				{
					SourceLabels: []string{"__address__", "_meta_kubernetes_service_name"},
					Modulus:      2,
					Action:       "hashmod",
					TargetLabel:  "__tmp_hash",
				},
				{
					SourceLabels: []string{"__tmp_hash"},
					Regex:        "^1$",
					Action:       "keep",
				},
			},
			assert: func(t *testing.T, job promcfg.Job, expectedRelabelConfigs []promcfg.RelabelConfig) {
				assert.Equal(t, expectedRelabelConfigs, job.RelabelConfigs)
			},
		},
	}

	for _, c := range testCases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			job := c.config.IncludeShardingRules(c.job)

			c.assert(t, job, c.expectedRelabelConfigs)
		})
	}
}
