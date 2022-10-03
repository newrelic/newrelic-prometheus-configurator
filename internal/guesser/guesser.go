package guesser

import (
	"errors"
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/newrelic/newrelic-prometheus-configurator/internal/promapi"
	"github.com/newrelic/newrelic-prometheus-configurator/internal/promcfg"
)

// Guesser fetches metric metadata from Prometheus API and generates relabel configs for
// New Relic remote write to map the correct type of the metric.
// https://docs.newrelic.com/docs/infrastructure/prometheus-integrations/install-configure-remote-write/set-your-prometheus-remote-write-integration/#override-mapping .
type Guesser struct {
	promServer  promapi.PromAPI
	relabelHash string
}

func New(ps promapi.PromAPI) Guesser {
	return Guesser{
		promServer: ps,
	}
}

// UpdatedMetricTypeRules fetches metric metadata from Prometheus API and generates the mapping relabel configs.
// Important: This is an STATEFUL NOT SYNCHRONIZED func which returns configs as long as the result is different from the previous call.
func (g *Guesser) UpdatedMetricTypeRules() ([]promcfg.RelabelConfig, error) {
	m, err := g.promServer.Metadata()
	if err != nil {
		return nil, err
	}

	rc := metricTypeRelabelConfigs(m)

	// TODO improve and abstract the hash thing. The goal is to check if the generated
	// configs are different from the latest generated.
	sort.Slice(rc, func(i, j int) bool {
		return strings.Compare(rc[i].Regex, rc[j].Regex) == -1
	})

	var sorted strings.Builder
	for _, r := range rc {
		sorted.WriteString(r.Regex)
	}

	if strings.EqualFold(g.relabelHash, sorted.String()) {
		return nil, nil
	}

	g.relabelHash = sorted.String()

	return rc, nil
}

var ErrMetadataEndpointFail = errors.New("checking metadata: status not successful")

// MetricTypeRelabelConfigs generates relabel configs needed for the New Relic Prometheus Remote Write
// to use the real metric type instead of the inferred type based on the metric name.
func metricTypeRelabelConfigs(metadata promapi.Metadata) []promcfg.RelabelConfig {
	rc := []promcfg.RelabelConfig{}

	for metricName, metricsMetadata := range metadata.Data {
		// TODO improve this to check if the list contains different types, and continue if that is true.
		// Otherwise keep executing since the relabel will be valid for all metrics.
		// Although we should consider also that an extra metric could appear on the next iteration changing
		// the output of this configurator. Since this we might need to introduce a more complex mechanism that
		// checks for /targets/metadata only for a group of metrics that are ambiguous.
		if len(metricsMetadata) != 1 {
			log.Printf("Metric %s skipped since contains more than 1 metadata", metricName)
			continue
		}

		switch metricsMetadata[0].Type {
		case "counter":
			if !hasCounterSuffix(metricName) {
				rc = append(rc, promcfg.RelabelConfig{
					SourceLabels: []string{"__name__"},
					Regex:        fmt.Sprintf("^%s$", metricName),
					TargetLabel:  "newrelic_metric_type",
					Replacement:  "counter",
					Action:       "replace",
				})
			}

		case "gauge":
			if hasCounterSuffix(metricName) {
				rc = append(rc, promcfg.RelabelConfig{
					SourceLabels: []string{"__name__"},
					Regex:        fmt.Sprintf("^%s$", metricName),
					TargetLabel:  "newrelic_metric_type",
					Replacement:  "gauge",
					Action:       "replace",
				})
			}

			// TODO check behavior of NR with summary. Prometheus adds the _sum and _count suffix automatically to these metrics
			// so by default they will be consider as count. According to doc _sum should be treated as summary.
			// https://docs.newrelic.com/docs/infrastructure/prometheus-integrations/view-query-data/translate-promql-queries-nrql#compare
			//
			// case "summary":
			// 	rc = append(rc, WriteRelabelConfigs{
			// 		SourceLabels: "[__name__]",
			//	    // metricName is the baseName of the metric.
			// 		Regex:        fmt.Sprintf("^%s_sum$", metricName),
			// 		TargetLabel:  "newrelic_metric_type",
			// 		Replacement:  "summary",
			// 		Action:       "replace",
			// 	})
		}
	}

	return rc
}

func hasCounterSuffix(metricName string) bool {
	if strings.HasSuffix(metricName, "_total") {
		return true
	}
	if strings.HasSuffix(metricName, "_count") {
		return true
	}
	if strings.HasSuffix(metricName, "_sum") {
		return true
	}
	if strings.HasSuffix(metricName, "_bucket") {
		return true
	}
	return false
}
