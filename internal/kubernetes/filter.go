package kubernetes

import (
	"regexp"

	"github.com/newrelic/newrelic-prometheus-configurator/internal/promcfg"
)

const (
	// check kubernetes service discovery metadata docs for more info.
	podMetadata        = "__meta_kubernetes_pod"
	serviceMetadata    = "__meta_kubernetes_service"
	annotationMetadata = "_annotation"
	labelMetadata      = "_label"
	// Prom labels like `__metadata_kubernetes_<role>_<label/annotation>present_` will contain
	// `true` if the label/annotation is present.
	presentSuffix = "present"

	separator = ";"
)

// This regex checks matches characters not supported inside Prometheus metric labels names.
// Copied from https://github.com/prometheus/prometheus/blob/v2.37.0/util/strutil/strconv.go#L41
var invalidPrometheusLabelCharRegex = regexp.MustCompile(`[^a-zA-Z0-9_]`)

// Filter defines the field needed to provided filtering capabilities to a kubernetes scrape job.
type Filter struct {
	Annotations map[string]string `yaml:"annotations,omitempty"`
	Labels      map[string]string `yaml:"labels,omitempty"`
}

// Pod creates a RelabelConfig that will keep only the Pod targets specified in Filter.
func (f Filter) Pod() promcfg.RelabelConfig {
	return f.build(podMetadata)
}

// Endpoints creates a RelabelConfig that will keep only the Endpoints targets specified in Filter.
func (f Filter) Endpoints() promcfg.RelabelConfig {
	return f.build(serviceMetadata)
}

// Valid creates a RelabelConfig that will keep only the Endpoints targets specified in Filter.
func (f Filter) Valid() bool {
	return len(f.Annotations) != 0 || len(f.Labels) != 0
}

// build creates a RelabelConfig that will keep only the targets specified in Filter.
// All conditions are concatenated with 'AND' operation.
// If no value has been specified for the metadata, it will check that exists.
func (f Filter) build(metadataSourcePrefix string) promcfg.RelabelConfig {
	filterCfg := promcfg.RelabelConfig{
		Separator: separator,
		Action:    "keep",
	}

	addConditions(&filterCfg, f.Annotations, metadataSourcePrefix+annotationMetadata)
	addConditions(&filterCfg, f.Labels, metadataSourcePrefix+labelMetadata)

	return filterCfg
}

// addConditions iterates over the metadata and appends the conditions
// to `source_labels` and `regex` of the filter.
func addConditions(relabelConfig *promcfg.RelabelConfig, metadata map[string]string, metadataPrefix string) {
	for prometheusLabels, regex := range metadata {
		// Prometheus sanitize all metadata keys (like kubernetes label/annotations names) to comply
		// with their naming conventions. We have to do the same so we can match in relabel configs.
		// The values in the metadata are not sanitized.
		// e.g: kubernetes label `prometheus.io/scrape` -> `prometheus_io_scrape`
		sanitizedK8sKey := invalidPrometheusLabelCharRegex.ReplaceAllString(prometheusLabels, "_")

		prefix := metadataPrefix
		// If no value has specified for metadata we just check it exist using the
		// `__meta_kubernetes_<role>_<annotation/label>present_<annotation/label name>: true
		if regex == "" {
			prefix += presentSuffix
			regex = "true"
		}

		prometheusSourceLabel := prefix + "_" + sanitizedK8sKey

		// Position on this array really matters since Prometheus will check against the `regex`
		// in order.
		relabelConfig.SourceLabels = append(relabelConfig.SourceLabels, prometheusSourceLabel)

		// Position here also matters since Prometheus parse this regex using the separator and
		// do the match against the same position of the source labels.
		relabelConfig.Regex = appendRegex(relabelConfig.Regex, regex)
	}
}

func appendRegex(regex string, newRegex string) string {
	// avoids to put separator for only one condition. In Prometheus this `regex: true;` doesn't work.
	if regex == "" {
		return newRegex
	}

	return regex + separator + newRegex
}
