# Defaults and Curated Experience

NewRelic provides a list of Dashboards, alerts and entities for several Services. The curated_experience, enabled by default,
configuration allows to scrape only the targets having this experience out of the box.

## What is scraped?
By default, the chart has two jobs configured and curated experience turned on:
- `default` scrapes all targets having `prometheus.io/scrape: true`.
  Since, by default, curated_experience.enabled=true then only targets selected by the curated experience are actually scraped.
- `newrelic` crapes all targets having `newrelic.io/scrape: true`.
  This is useful to extend the 'default-job' allowlisting all services adding the required annotation.

## How to extend the solution
In order to scrape additional targets there are mainly 3 ways (starring from the default configuration):
 - setting `curated_experience.enabled=false` configures k8s to scrape all the targets matching both
`prometheus.io/scrape: true` and `newrelic.io/scrape: true` if the default configuration was not changed.
 - to scrape more targets you can whitelist them one by one adding the `newrelic.io/scrape: true` label to pods or services.
 - `app_values` and `source_labels` can be modified adding or reducing the entries. Note that the dashboards might not be 
working out of the box with custom labels or values

## Upgrading
When upgrading new curated experience might be available. Therefore, the amount of data scraped can increase 
depending on which service are in the cluster after any upgrade involving curated experience.

If a service added by a curated experience was previously scraped by a different Job there is the risk of duplication.
In order to detect targets scraped multiple times by different jobs you can run the following query:
```
FROM Metric select uniqueCount(job)  facet instance, instrumentation.source 
```
