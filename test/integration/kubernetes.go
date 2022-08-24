//go:build integration_test

package integration

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	kubeconfigPath = "../../.kubeconfig-dev"
	defaultPodPort = 8000
)

type k8sEnvironment struct {
	kubeconfigFullPath string
	client             *kubernetes.Clientset
	testNamespace      *corev1.Namespace

	defaultTimeout time.Duration
	defaultBackoff time.Duration
}

// newK8sEnvironment connects to a cluster using kubeconfigPath and creates namespace for the current test.
func newK8sEnvironment(t *testing.T) k8sEnvironment {
	t.Helper()

	// Prometheus needs the full path to read the file.
	kubeconfigFullPath, err := filepath.Abs(kubeconfigPath)
	require.NoError(t, err)

	clientset, err := k8sClient(kubeconfigFullPath)
	require.NoError(t, err)

	namespaceTemplate := corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "newrelic-prometheus-test-",
		},
	}

	testNamespace, err := clientset.CoreV1().Namespaces().Create(context.Background(), &namespaceTemplate, metav1.CreateOptions{})
	require.NoError(t, err)

	t.Cleanup(func() {
		err := clientset.CoreV1().Namespaces().Delete(context.Background(), testNamespace.Name, metav1.DeleteOptions{})
		require.NoError(t, err)
	})

	return k8sEnvironment{
		kubeconfigFullPath: kubeconfigFullPath,
		client:             clientset,
		testNamespace:      testNamespace,
		defaultBackoff:     time.Second,
		defaultTimeout:     time.Second * 20,
	}
}

func (ke *k8sEnvironment) addPod(t *testing.T, pod *corev1.Pod) *corev1.Pod {
	t.Helper()

	p, err := ke.client.CoreV1().Pods(ke.testNamespace.Name).Create(context.Background(), pod, metav1.CreateOptions{})
	require.NoError(t, err)

	return p
}

// addPodAndWaitOnPhase creates the pod and waits until the specified podPhase.
func (ke *k8sEnvironment) addPodAndWaitOnPhase(t *testing.T, pod *corev1.Pod, podPhase corev1.PodPhase) *corev1.Pod {
	t.Helper()

	p := ke.addPod(t, pod)

	err := retryUntilTrue(ke.defaultTimeout, ke.defaultBackoff, func() bool {
		var err error
		// we want to override p with the latest pod retrieved.
		p, err = ke.client.CoreV1().Pods(ke.testNamespace.Name).Get(context.Background(), p.Name, metav1.GetOptions{})
		require.NoError(t, err)

		return p.Status.Phase == podPhase
	})
	require.NoError(t, err)

	return p
}

// addService adds a service using the k8s client.
// It fails in case the service can't be added.
func (ke *k8sEnvironment) addService(t *testing.T, srv *corev1.Service) *corev1.Service {
	t.Helper()

	p, err := ke.client.CoreV1().Services(ke.testNamespace.Name).Create(context.Background(), srv, metav1.CreateOptions{})
	require.NoError(t, err)

	return p
}

//nolint:goerr113
func k8sClient(kubeconfigPath string) (*kubernetes.Clientset, error) {
	conf, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to build config")
	}

	client, err := kubernetes.NewForConfig(conf)
	if err != nil {
		return nil, fmt.Errorf("failed to build client")
	}

	return client, nil
}

func fakePodSpec() corev1.PodSpec {
	return corev1.PodSpec{
		Containers: []corev1.Container{
			{
				Name:    "fake-exporter",
				Image:   "alpine:latest",
				Command: []string{"/bin/sh", "-c", "sleep infinity"},
				Ports: []corev1.ContainerPort{
					{
						ContainerPort: defaultPodPort,
					},
				},
			},
		},
	}
}

func fakePod(namePrefix string, annotations, labels map[string]string) *corev1.Pod {
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: namePrefix,
			Annotations:  annotations,
			Labels:       labels,
		},
		Spec: fakePodSpec(),
	}
}

func fakeService(namePrefix string, selector, annotations, labels map[string]string) *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: namePrefix,
			Annotations:  annotations,
			Labels:       labels,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Port:       80,
					TargetPort: intstr.FromInt(defaultPodPort),
				},
			},
			Selector: selector,
		},
	}
}
