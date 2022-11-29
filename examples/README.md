# Example Instructions

The yaml files in this example folder can be used to test out New Relic's Prometheus-based Kubernetes monitoring quickstart.

1. Deploy the following yamls. Note that if you already have KSM setup in your cluster, not all of these yamls will be necessary.

```
kubectl apply -f cluster-role.yaml
kubectl apply -f ksm-cluster-role-binding.yaml
kubectl apply -f ksm-cluster-role.yaml
kubectl apply -f ksm-deployment.yaml
kubectl apply -f ksm-service-account.yaml
kubectl apply -f ksm-service.yaml
```

2. Modify the `configurator-values.yaml` file to add your cluster name and license key.

3. Install the helm chart:

```
helm repo add newrelic-prometheus https://newrelic.github.io/newrelic-prometheus-configurator
helm upgrade --install newrelic newrelic-prometheus/newrelic-prometheus-agent -f configurator-values.yaml
```

4. Install the Kubernetes Prometheus Quickstart dashboard and watch for data to appear:

# TODO: add link to Quickstart
