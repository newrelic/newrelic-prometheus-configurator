{{- /* Return the newrelic-prometheus configuration */ -}}
{{- define "newrelic-prometheus.configurator.config" -}}

{{- /* TODO: we should consider using an external label to set the cluster name instead */ -}}
data_source_name: {{ .Values.cluster }}
{{ include "nerelic-prometheus.configurator.config._remoteWrite" . }}

{{- end -}}


{{- /* Internal use: it builds the remote_write configuration from configurator config */ -}}
{{- define "nerelic-prometheus.configurator.config._remoteWrite" -}}

newrelic_remote_write:
  staging: {{ include "newrelic.common.nrStaging.value" . }}

{{- if .Values.config -}}

{{- $remoteWrite := .Values.config.remote_write -}}
{{- if $remoteWrite  -}}
{{- if $remoteWrite.extra_write_relabel_configs }}
  extra_write_relabel_configs:
    {{- $remoteWrite.extra_write_relabel_configs | toYaml | nindent 6 -}}
{{- end -}}
{{- if $remoteWrite.proxy_url }}
  proxy_url: {{ $remoteWrite.proxy_url }}
{{- end -}}
{{- if $remoteWrite.remote_timeout }}
  remote_timeout: {{ $remoteWrite.remote_timeout }}
{{- end -}}
{{- if $remoteWrite.tls_config }}
  tls_config:
    {{- $remoteWrite.tls_config | toYaml | nindent 6 -}}
{{- end -}}
{{- if $remoteWrite.queue_config }}
  queue_config:
    {{- $remoteWrite.queue_config | toYaml | nindent 6 -}}
{{- end -}}
{{- end -}} {{- /* $remoteWrite */ -}}

{{- if .Values.config.extra_remote_write }}
extra_remote_write:
  {{- .Values.config.extra_remote_write | toYaml | nindent 4 -}}
{{- end -}}

{{- end -}}{{- /* .Values.config */ -}}

{{- end -}}
