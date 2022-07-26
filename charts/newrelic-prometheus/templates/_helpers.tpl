{{- /* Return the newrelic-prometheus configuration */ -}}
{{- define "newrelic-prometheus.configurator.config" -}}

{{- /* TODO: we should consider using an external label to set the cluster name instead */ -}}
data_source_name: {{ include "newrelic.common.cluster" . }}
{{ include "newrelic-prometheus.configurator.config._remoteWrite" . }}
{{ include "newrelic-prometheus.configurator.config._common" . }}

{{- end -}}

{{- /* Internal use: it builds the common configuration from configurator config, cluster name and custom attributes */ -}}
{{- define "newrelic-prometheus.configurator.config._common" -}}
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


{{- /* Internal use: it builds the remote_write configuration from configurator config */ -}}
{{- define "newrelic-prometheus.configurator.config._remoteWrite" -}}

newrelic_remote_write:
{{- if (include "newrelic.common.nrStaging" . ) }}
  staging: true
{{- end -}}

{{- if .Values.config -}}

{{- if .Values.config.remote_write  -}}
{{- .Values.config.remote_write | toYaml | nindent 4 -}}
{{- end -}}

{{- if .Values.config.extra_remote_write }}
extra_remote_write:
  {{- .Values.config.extra_remote_write | toYaml | nindent 4 -}}
{{- end -}}

{{- if .Values.config.static_targets  }}
static_targets:
{{- .Values.config.static_targets | toYaml | nindent 2 -}}
{{- end -}}

{{- if .Values.config.extra_scrape_configs }}
extra_scrape_configs:
  {{- .Values.config.extra_scrape_configs | toYaml | nindent 2 -}}
{{- end -}}

{{- end -}}

{{- end -}}
