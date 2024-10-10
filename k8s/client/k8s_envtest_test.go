package client_test

import (
	"context"
	"github.com/PixoVR/pixo-golang-server-utilities/k8s/client"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var _ = Describe("EnvTest k8s", Ordered, func() {

	var (
		ctx        context.Context
		baseClient kubernetes.Interface
		podName    string
	)

	BeforeAll(func() {
		ctx = context.Background()

		var err error
		baseClient, err = client.NewEnvTestClient()
		Expect(err).NotTo(HaveOccurred())
		Expect(baseClient).NotTo(BeNil())

		podName = "nginx"
		_, err = baseClient.CoreV1().
			Pods(namespace).
			Create(ctx, &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{Name: podName},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Name:  "nginx",
						Image: "nginx",
						Ports: []corev1.ContainerPort{{ContainerPort: 80}},
					}},
				},
			}, metav1.CreateOptions{})
	})

	AfterAll(func() {
		err := baseClient.
			CoreV1().
			Pods(namespace).
			Delete(ctx, podName, metav1.DeleteOptions{})
		Expect(err).NotTo(HaveOccurred())
	})

	It("can get pod by name", func() {
		pod, err := baseClient.
			CoreV1().
			Pods(namespace).
			Get(ctx, podName, metav1.GetOptions{})

		Expect(err).NotTo(HaveOccurred())
		Expect(pod).NotTo(BeNil())
		Expect(pod.Name).To(Equal(podName))
	})

})
