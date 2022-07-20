{{- /* Return the newrelic-prometheus configuration */ -}}
{{- define "newrelic-prometheus.configurator.config" -}}

{{- /* TODO: we should consider using an external label to set the cluster name instead */ -}}
data_source_name: {{ .Values.cluster }}
{{ include "nerelic-prometheus.configurator.config._remoteWrite" . }}

{{- end -}}


{{- /* Internal use: it builds the remote_write configuration from configurator config */ -}}
{{- define "nerelic-prometheus.configurator.config._remoteWrite" -}}

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

{{- end -}}

{{- end -}}
