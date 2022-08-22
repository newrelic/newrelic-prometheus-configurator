This is just a simple example to see how docker service discovery works.

A more real POC involving the configurator could be done creating the `docker` configuration, similar to `kubernetes`.

```yaml
docker:
  jobs:
    - job_name_prefix: default
      target_discovery:
        filter:
          labels:
            prometheus.io/scrape: true
```

Also in case of supporting the docker use case we might consider baking an image based on Prometheus image that adds and runs the configurator since `init containers` will not be an option on a plain docker environment.
