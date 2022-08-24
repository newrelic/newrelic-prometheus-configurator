{{- /* Return the newrelic-prometheus configuration */ -}}

{{- /* it builds the common configuration from configurator config, cluster name and custom attributes */ -}}
{{- define "newrelic-prometheus.configurator.common" -}}
{{- $tmp := dict "external_labels" (dict "cluster_name" (include "newrelic.common.cluster" . )) -}}

{{- if .Values.config  -}}
  {{- if .Values.config.common -}}
      {{- $tmp := mustMerge $tmp .Values.config.common -}}
  {{- end -}}
{{- end -}}

{{- $tmpCustomAttribute := dict "external_labels" (include "newrelic.common.customAttributes" . | fromYaml ) -}}
{{- $tmp := mustMerge $tmp $tmpCustomAttribute  -}}

common:
{{- $tmp | toYaml | nindent 2 -}}

{{- end -}}


{{- /* it builds the newrelic_remote_write configuration from configurator config */ -}}
{{- define "newrelic-prometheus.configurator.newrelic_remote_write" -}}
{{- $tmp := dict -}}

{{- if (include "newrelic.common.nrStaging" . ) -}}
  {{- $tmp = dict "staging" true  -}}
{{- end -}}

{{- if .Values.config -}}
  {{- if .Values.config.newrelic_remote_write  -}}
    {{- $tmp = mustMerge $tmp .Values.config.newrelic_remote_write  -}}
  {{- end -}}
{{- end -}}

{{- if not (empty $tmp) -}}
  {{- dict "newrelic_remote_write" $tmp | toYaml -}}
{{- end -}}

{{- end -}}

{{- /* it builds the extra_remote_write configuration from configurator config */ -}}
{{- define "newrelic-prometheus.configurator.extra_remote_write" -}}
{{- if .Values.config -}}
  {{- if .Values.config.extra_remote_write  -}}
extra_remote_write:
    {{- .Values.config.extra_remote_write | toYaml | nindent 2 -}}
  {{- end -}}
{{- end -}}
{{- end -}}

{{- define "newrelic-prometheus.configurator.static_targets" -}}
{{- if .Values.config -}}
  {{- if .Values.config.static_targets -}}
static_targets:
    {{- .Values.config.static_targets | toYaml | nindent 2 -}}
  {{- end -}}
{{- end -}}
{{- end -}}

{{- define "newrelic-prometheus.configurator.extra_scrape_configs" -}}
{{- if .Values.config -}}
  {{- if .Values.config.extra_scrape_configs  -}}
extra_scrape_configs:
    {{- .Values.config.extra_scrape_configs | toYaml | nindent 2 -}}
  {{- end -}}
{{- end -}}
{{- end -}}

{{- define "newrelic-prometheus.configurator.kubernetes" -}}
{{- if .Values.config -}}
  {{- if .Values.config.kubernetes  -}}
kubernetes:
  {{- .Values.config.kubernetes | toYaml | nindent 2 -}}
  {{- end -}}
{{- end -}}
{{- end -}}
