package configurator

import (
	"regexp"
)

const (
	podAnnotationPrefix     = "__meta_kubernetes_pod_annotation_"
	podLabelPrefix          = "__meta_kubernetes_pod_label_"
	serviceAnnotationPrefix = "__meta_kubernetes_service_annotation_"
	serviceLabelPrefix      = "__meta_kubernetes_service_label_"
	separator               = ";"
)

// Copied from https://github.com/prometheus/prometheus/blob/v2.37.0/util/strutil/strconv.go#L41
var invalidLabelCharRegex = regexp.MustCompile(`[^a-zA-Z0-9_]`)

// Filter defines the field needed to provided filtering capabilities to a kubernetes scrape job.
type Filter struct {
	Annotations map[string]string
	Labels      map[string]string
}

func podFilterSettingsBuilder(jobOutput JobOutput, k8sJob KubernetesJob) JobOutput {
	return buildFilter(jobOutput, k8sJob, podAnnotationPrefix, podLabelPrefix)
}

func endpointsFilterSettingsBuilder(jobOutput JobOutput, k8sJob KubernetesJob) JobOutput {
	return buildFilter(jobOutput, k8sJob, serviceAnnotationPrefix, serviceLabelPrefix)
}

func buildFilter(jobOutput JobOutput, k8sJob KubernetesJob, annotationPrefix, labelPrefix string) JobOutput {
	if k8sJob.TargetDiscovery.Filter == nil {
		// If no filter is added all discovered targets will be scrape
		return jobOutput
	}

	filterCfg := RelabelConfig{
		Separator: separator,
		Action:    "keep",
	}

	for annotation, val := range k8sJob.TargetDiscovery.Filter.Annotations {
		sanitizedAnnotation := annotationPrefix + invalidLabelCharRegex.ReplaceAllString(annotation, "_")

		filterCfg.SourceLabels = append(filterCfg.SourceLabels, sanitizedAnnotation)

		filterCfg.Regex = appendRegex(filterCfg.Regex, val)
	}

	for label, val := range k8sJob.TargetDiscovery.Filter.Labels {
		sanitizedLabel := labelPrefix + invalidLabelCharRegex.ReplaceAllString(label, "_")

		filterCfg.SourceLabels = append(filterCfg.SourceLabels, sanitizedLabel)

		filterCfg.Regex = appendRegex(filterCfg.Regex, val)
	}

	jobOutput.RelabelConfigs = append(jobOutput.RelabelConfigs, filterCfg)

	return jobOutput
}

func appendRegex(regex string, val string) string {
	if regex == "" {
		return val
	}

	return regex + separator + val
}
