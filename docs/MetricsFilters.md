# Metrics and label transformations.

They can be applied on two levels, per job (`static_targets` or `kubernetes`) or per remote write level. If configured per job, the filtering only applies to the metrics of targets scraped by that job, and if they are applied at `newrelic_remote_write` level, the filters apply to all metrics that are being sent to New Relic.
The metric filter process happens after these has been scraped from the targets.
The `extra_metric_relabel_config` parameter can be used to apply the filters, which adds entries of [metric_relabel_config](https://prometheus.io/docs/prometheus/latest/configuration/configuration/#relabel_config). This parameter is present at `static_targets.jobs`, `kubernetes.jobs` and the `extra_write_relabel_configs` parameter for `newrelic_remote_write`.

``` yaml
static_targets:
- name: self-metrics
  urls:
    - 'http://static-service:8181'
  extra_metric_relabel_config:
  # Drop metrics with prefix 'go_' for this target.
  - source_labels: [__name__]
    regex: 'go_.+'
    action: drop

newrelic_remote_write:
  extra_write_relabel_configs:
  # Drop all metrics with the specified name before sent to New Relic.
  - source_labels: [__name__]
    regex: 'metric_name'
    action: drop
```


## Keep/drop metrics examples

``` yaml
# Drops metrics staring with 'prefix_'.
- source_labels: [__name__]
  regex: 'prefix_.+'
  action: drop

# Drops metrics with specific kubernetes label.
- source_labels: [k8s_io_app]
  regex: 'appLabelValue'
  action: drop

# Drops metrics staring with 'prefix_' that also contains a kubernetes label.
- source_labels: [__name__,k8s_io_app]
  regex: 'prefix_.+;appLabelValue'
  action: drop

# Keeps only metrics staring with 'prefix_'.
- source_labels: [__name__]
  regex: 'prefix_.+'
  action: keep
```

## Add or Drop Metric labels

Note: Metric Labels names must comply with [Prometheus DataModel](https://prometheus.io/docs/concepts/data_model/#metric-names-and-labels).

``` yaml
# Add labels to metrics with prefix.
- source_labels: [__name__]
  regex: 'prefix_.+'
  target_label: new_label
  action: replace
  replacement: newLabelValue

# Drop any label that matches the regex from all metrics.
# Care must be taken if removing identifying labels to ensure correct metrics aggregations are obtained.
- regex: 'label_name'
  action: labeldrop
```
