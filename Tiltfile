# -*- mode: Python -*-

# Settings and defaults.
project_name = 'newrelic-prometheus'
cluster_context = 'minikube'

# Only use explicitly allowed kubeconfigs as a safety measure.
allow_k8s_contexts(cluster_context)

local_resource('Configurator binary', 'make build-multiarch', deps=[
  './cmd',
  './configurator',
])

# Images are pushed to the docker inside minikube since we use 'eval $(minikube docker-env)'.
docker_build('prometheus-configurator', '.')
docker_build('openmetrics-fake-exporter', './test/openmetrics-fake-exporter/.')

# Deploying Kubernetes resources.
k8s_yaml(
  helm(
    './charts/%s' % project_name,
    name=project_name,
    values=['values-dev.yaml'],
    set=['licenseKey=%s' % os.getenv('NR_PROM_LICENSE_KEY'), 'cluster=%s' % os.getenv('NR_PROM_CLUSTER')],
    ))
k8s_yaml(helm('./charts/internal/test-resources', name='test-resources'))

# Tracking the deployment.
k8s_resource(project_name)
