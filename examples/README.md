# Example Instructions

The yaml files in this example folder can be used to test out New Relic's Prometheus observability solution.

1. Deploy the following yamls. Note that if you already have KSM setup in your cluster, not all of these yamls will be necessary.

```
kubectl apply -f cluster-role.yaml
kubectl apply -f ksm-cluster-role-binding.yaml
kubectl apply -f ksm-cluster-role.yaml
kubectl apply -f ksm-deployment.yaml
kubectl apply -f ksm-service-account.yaml
kubectl apply -f ksm-service.yaml
```

2. Install the helm chart:

```
helm repo add newrelic-prometheus https://newrelic.github.io/newrelic-prometheus-configurator
helm upgrade --install newrelic newrelic-prometheus/newrelic-prometheus-agent -f configurator-values.yaml
```

3. Watch for data to appear in the following dashboards:

# TODO: add link to Quickstart
