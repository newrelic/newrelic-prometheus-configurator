package guesser_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/newrelic/newrelic-prometheus-configurator/internal/guesser"
	"github.com/newrelic/newrelic-prometheus-configurator/internal/promapi"
	"github.com/newrelic/newrelic-prometheus-configurator/internal/promcfg"
)

type promAPIMock struct {
	response promapi.Metadata
	fail     bool
}

func (pam promAPIMock) Metadata() (promapi.Metadata, error) {
	if pam.fail {
		return promapi.Metadata{}, fmt.Errorf("fail")
	}

	return pam.response, nil
}

func Test_UpdatedMetricTypeRules(t *testing.T) {
	t.Parallel()

	type args struct {
		metadata  promapi.Metadata
		fail      bool
		unsuccess bool
	}
	tests := []struct {
		name string
		args args
		want []promcfg.RelabelConfig
	}{
		{
			name: "count metric without prefix",
			args: args{
				metadata: promapi.Metadata{
					Status: "success",
					Data: map[string][]promapi.MetricMetadata{
						"count_metric": {
							{Type: "counter"},
						},
					},
				},
			},
			want: []promcfg.RelabelConfig{
				{
					SourceLabels: []string{"__name__"},
					Regex:        "^count_metric$",
					TargetLabel:  "newrelic_metric_type",
					Replacement:  "counter",
					Action:       "replace",
				},
			},
		},
		{
			name: "gauge metrics with count sufix",
			args: args{
				metadata: promapi.Metadata{
					Status: "success",
					Data: map[string][]promapi.MetricMetadata{
						"gauge_metric_total": {
							{Type: "gauge"},
						},
					},
				},
			},
			want: []promcfg.RelabelConfig{
				{
					SourceLabels: []string{"__name__"},
					Regex:        "^gauge_metric_total$",
					TargetLabel:  "newrelic_metric_type",
					Replacement:  "gauge",
					Action:       "replace",
				},
			},
		},
	}

	for _, tt := range tests {
		tc := tt
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			g := guesser.New(promAPIMock{response: tc.args.metadata, fail: tc.args.fail})

			got, _ := g.UpdatedMetricTypeRules()

			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("MetricTypeRelabelConfigs() = %v, want %v", got, tc.want)
			}
		})
	}
}
