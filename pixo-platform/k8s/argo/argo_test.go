package argo_test

import (
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/k8s/argo"
	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

var _ = Describe("Argo", func() {

	var (
		clientset *argo.Client
		namespace = "dev-multiplayer"
	)

	BeforeEach(func() {
		var err error
		clientset, err = argo.NewLocalArgoClient()
		Expect(err).NotTo(HaveOccurred())
		Expect(clientset).To(Not(BeNil()))
	})

	It("can get the list of workflows ", func() {
		workflows, err := clientset.ListWorkflows(namespace)

		Expect(err).NotTo(HaveOccurred())
		Expect(workflows).NotTo(BeNil())
	})

	It("can create, get, and delete a whalesay workflow", func() {
		name := "whalesay"
		spec := &v1alpha1.Workflow{
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

		workflow, err := clientset.CreateWorkflow(namespace, spec)
		Expect(err).NotTo(HaveOccurred())
		Expect(workflow).NotTo(BeNil())

		retrieved, err := clientset.GetWorkflow(namespace, name)
		Expect(err).NotTo(HaveOccurred())
		Expect(retrieved).NotTo(BeNil())
		Expect(retrieved.Name).To(Equal(name))

		time.Sleep(10 * time.Second)

		err = clientset.DeleteWorkflow(namespace, name)
		Expect(err).NotTo(HaveOccurred())

		retrieved, err = clientset.GetWorkflow(namespace, name)
		Expect(err).To(HaveOccurred())
		Expect(retrieved).To(BeNil())
	})

})
