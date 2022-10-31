# Defaults and Integration Filters

NewRelic provides a list of Dashboards, alerts and entities for several Services. The integrations_filter configuration
, enabled by default, allows to scrape only the targets having this experience out of the box.

## What is scraped?
By default, the chart has two jobs configured and integration filters turned on:
- `default` scrapes all targets having `prometheus.io/scrape: true`.
  Since, by default, `integrations_filter.enabled=true` then only targets selected by the integration filters are actually scraped.
- `newrelic` scrapes all targets having `newrelic.io/scrape: true`.
  This is useful to extend the 'default-job' allowlisting all services adding the required annotation.
  This job overrides the `integrations_filter` default configuration setting it to false.

## How to extend the solution
In order to scrape additional targets with the default job configuration there are mainly 3 ways:
 - setting `integrations_filter.enabled=false` configures k8s to scrape all the targets matching both
`prometheus.io/scrape: true` and `newrelic.io/scrape: true`.
 - to scrape more targets you can allowlist them one by one adding the `newrelic.io/scrape: true` label to pods or services.
 - `app_values` and `source_labels` can be modified adding or reducing the entries. Note that the dashboards might not be 
working out of the box with custom labels or values.

## Upgrading
When upgrading, new integration filters might be available. Therefore, the amount of data scraped can increase 
depending on which service are in the cluster after any upgrade involving integration filters.

In case, you do not want to risk to increase the amount of target scrapes, you can save in your `values.yaml` a fixed list for
`app_values` and `source_labels`.

If a service added by a new default of integration filters was previously scraped by a different Job there is the risk of duplication.
In order to detect targets scraped multiple times by different jobs you can run the following query:
```
FROM Metric select uniqueCount(job) facet instance, instrumentation.source, cluster_name
```
