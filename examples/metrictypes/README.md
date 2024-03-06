## Override Metric Type Mappings

If you have metrics that don't follow Prometheus naming conventions, you can configure remote write to tag the metric with a newrelic_metric_type label that indicates the metric type. This label is stripped when received by New Relic.

Example: You have a counter metric named my_counter, which does not have our naming convention suffix of _bucket, _count or _total. 

In this situation, your metric would be identified as a gauge rather than a counter. To correct this, add the following relabel configuration to your prometheus.yml:

### Prometheus Server Example
```
- url: https://metric-api.newrelic.com/prometheus/v1/write?X-License-Key=...
  write_relabel_configs:
  - source_labels: [__name__]
    regex: ^my_counter$
    target_label: newrelic_metric_type
    replacement: "counter"
    action: replace
```

### Prometheus Configurator Example
```
config:
  newrelic_remote_write:
    extra_write_relabel_configs:
      - source_labels: [__name__]
        regex: ^my_counter$
        target_label: newrelic_metric_type
        replacement: "counter"
        action: replace
```

This rule matches any metric with the name my_counter and adds a newrelic_metric_type label that identifies it as a counter. You can use the following (case sensitive!) values as the replacement value:

- counter
- gauge
- summary

When a newrelic_metric_type label is present on a metric received and set to one of the valid values, New Relic will assign the indicated type to the metric (and strip the label) before downstream consumption in the data pipeline. If you have multiple metrics that don't follow the above naming conventions, you can add multiple rules with each rule matching different source labels.