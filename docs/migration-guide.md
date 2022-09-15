# Migration guide from POMI (nri-prometheus) to newrelic-prometheus

What have changed:
POMI (nri-prometheus) was a custom solution build by New Relic with some of the Prometheus Server capabilities. The `newrelic-prometheus` solution is actually running the Prometheus Server in Agent mode. This allows to use any scrape related feature of the Prometheus Server like authorization methods and relabel configs.
Also we have built a `prometheus-configurator` on top of that which improved the configuration experience of the solution.

## Metadata

Found details [here](./MetricLabels.md) about the labels added by `newrelic-prometheus`.

This is a list of POMI metadata and it replacement in the new solution.

| POMI | newrelic-prometheus | Note |
|------|---------------------|------|
| `namespaceName`| `namespace` | Renamed |
| `nodeName`| `node` | Renamed |
| `podName`| `pod` | Renamed |
| `serviceName`| `service` | Renamed |
| `deploymentName`| - | Removed |
| `clusterName` | `cluster_name` | Renamed |
| `integrationName` | - | Removed |
| `integrationVersion` | - | Removed |
| `metricName` | - | Removed |
| `nrMetricType` | - | Removed |
| `promMetricType` | - | Removed |
| `label.<kubernetesLabel>`| `<kubernetesLabel>` | Renamed - The `label.` prefix has been removed and the label name sanitized to Prometheus naming standards. |
| `targetName` | - | Removed |
| `scrapedTargetKind` | `job` | Replaced - The `job` can be used to identify the scrape kind |
| `scrapedTargetName` | - | Removed |
| `scrapedTargetURL` | `instance` | Replaced - `instance` contains the `host:port` for the target |

## Kubernetes target discovery

The target discovery configuration has improved with the introduction of the Jobs and can be easily configured as explained [here](./KuberntesTargetFilter.md).

There are some default behaviors that differs on the solutions.
By default POMI scrapes Pods and Services that has `prometheus.io/scrape=true` Label or Annotation. But `newrelic-prometheus` scrapes Pods and Endpoints from Services with `prometheus.io/scrape=true` Annotation. 

## Metrics types

POMI converts Prometheus to New Relic metrics before send them, applying the [mappings](https://docs.newrelic.com/docs/infrastructure/prometheus-integrations/view-query-data/translate-promql-queries-nrql#compare) according to the metric metadata type read from the scrape data( `# TYPE <metric_type>`).

On the newrelic-prometheus new solution the metrics are sent directly to the New Relic Remote Write endpoint which takes care of this [conversion based on the metric name](https://docs.newrelic.com/docs/infrastructure/prometheus-integrations/install-configure-remote-write/set-your-prometheus-remote-write-integration#mapping). So is possible that some metrics that are transformed correctly by POMI will not be converted by the Remote Write on the following two cases:
- **Prometheus Counters** metrics **not having** the name suffix `total`, `count`, `sum`, `bucket` will be consider as **Gauge** type.
- On the other side, **Prometheus Gauges** metrics **having** the name suffix `total`, `count`, `sum`, `bucket` will be consider as **Counter** type.

Metric [type mappings](https://docs.newrelic.com/docs/infrastructure/prometheus-integrations/install-configure-remote-write/set-your-prometheus-remote-write-integration#override-mapping) relabel configs can be added in order to correct this behavior.

## Transformations

The POMI transformations are replaced with Prometheus relabel configs.
[Here](./MetricsFilters.md) you can find a list of examples of relabels configs to cover different use cases that were covered by the POMI Transformations. 


## Self-instrumentation

As well as POMI, newrelic-prometheus selfs-scrapes internal metrics. This metrics have the prefix `prometheus_` and can be used to observe the status of the Prometheus instance.
By default only a set of these metrics are sent. You can find the list on the default `values.yaml` of the chart under `config.static_targets.jobs` with the `self-metrics` job name.
