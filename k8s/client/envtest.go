package k8s

import (
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
)

func NewEnvTestClient() (kubernetes.Interface, error) {
	testEnv := &envtest.Environment{}
	cfg, err := testEnv.Start()
	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(cfg)
}
