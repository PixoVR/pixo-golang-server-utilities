package k8s

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
)

func NewEnvTestClient(namespaces ...string) (kubernetes.Interface, error) {
	testEnv := &envtest.Environment{}
	cfg, err := testEnv.Start()
	if err != nil {
		return nil, err
	}

	k8sClient, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, err
	}

	for _, namespace := range namespaces {
		_, err = k8sClient.CoreV1().Namespaces().Create(context.Background(), &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{Name: namespace},
		}, metav1.CreateOptions{})
		if err != nil {
			return nil, err
		}
	}

	return k8sClient, nil
}
