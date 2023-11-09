package argo_test

import (
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/k8s/argo"
	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Argo", func() {

	var (
		clientset    *argo.Client
		namespace    = "dev-multiplayer"
		name         = "whalesay"
		whalesaySpec = &v1alpha1.Workflow{
			ObjectMeta: metav1.ObjectMeta{
				Name: name,
			},
			Spec: v1alpha1.WorkflowSpec{
				Entrypoint: name,
				Templates: []v1alpha1.Template{
					{
						Name: name,
						Container: &corev1.Container{
							Image:   "docker/whalesay:latest",
							Command: []string{"cowsay"},
						},
					},
				},
			},
		}
	)

	BeforeEach(func() {
		var err error
		clientset, err = argo.NewLocalArgoClient()
		Expect(err).NotTo(HaveOccurred())
		Expect(clientset).To(Not(BeNil()))

		workflow, err := clientset.CreateWorkflow(namespace, whalesaySpec)
		Expect(err).NotTo(HaveOccurred())
		Expect(workflow).NotTo(BeNil())
	})

	AfterEach(func() {
		err := clientset.DeleteWorkflow(namespace, name)
		Expect(err).NotTo(HaveOccurred())

		retrieved, err := clientset.GetWorkflow(namespace, name)
		Expect(err).To(HaveOccurred())
		Expect(retrieved).To(BeNil())
	})

	It("can get the list of workflows ", func() {
		workflows, err := clientset.ListWorkflows(namespace)

		Expect(err).NotTo(HaveOccurred())
		Expect(workflows).NotTo(BeNil())
	})

	It("can get whalesay workflow", func() {
		retrieved, err := clientset.GetWorkflow(namespace, name)
		Expect(err).NotTo(HaveOccurred())
		Expect(retrieved).NotTo(BeNil())
		Expect(retrieved.Name).To(Equal(name))
	})

})
