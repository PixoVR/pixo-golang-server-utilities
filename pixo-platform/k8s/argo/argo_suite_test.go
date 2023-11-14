package argo_test

import (
	"fmt"
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/k8s/argo"
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/k8s/base"
	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestArgo(t *testing.T) {
	RegisterFailHandler(Fail)
	_ = os.Setenv("IS_LOCAL", "true")
	RunSpecs(t, "Argo Client Suite")
}

var (
	namespace       = "dev-multiplayer"
	workflowName    = "whalesay"
	templateOneName = fmt.Sprintf("%s-1", workflowName)
	templateTwoName = fmt.Sprintf("%s-2", templateOneName)

	k8sClient  *base.Client
	argoClient *argo.Client
	workflow   *v1alpha1.Workflow

	whalesaySpec = &v1alpha1.Workflow{
		ObjectMeta: metav1.ObjectMeta{
			Name: workflowName,
		},
		Spec: v1alpha1.WorkflowSpec{
			Entrypoint: workflowName,
			Templates: []v1alpha1.Template{
				{
					Name: workflowName,
					DAG: &v1alpha1.DAGTemplate{
						Tasks: []v1alpha1.DAGTask{
							{
								Name:     templateOneName,
								Template: templateOneName,
							},
							{
								Name:         templateTwoName,
								Template:     templateTwoName,
								Dependencies: []string{templateOneName},
							},
						},
					},
				},
				{
					Name: templateOneName,
					Container: &corev1.Container{
						Name:    templateOneName,
						Image:   "docker/whalesay:latest",
						Command: []string{"cowsay"},
					},
				},
				{
					Name: templateTwoName,
					Container: &corev1.Container{
						Name:    templateTwoName,
						Image:   "docker/whalesay:latest",
						Command: []string{"cowsay"},
					},
				},
			},
		},
	}
)

var _ = BeforeSuite(func() {
	var err error
	k8sClient, err = base.NewLocalClient()
	Expect(err).NotTo(HaveOccurred())
	Expect(k8sClient).To(Not(BeNil()))

	argoClient, err = argo.NewLocalArgoClient()
	Expect(err).NotTo(HaveOccurred())
	Expect(argoClient).To(Not(BeNil()))

	workflow, err = argoClient.CreateWorkflow(namespace, whalesaySpec)
	Expect(err).NotTo(HaveOccurred())
	Expect(workflow).NotTo(BeNil())
	time.Sleep(3 * time.Second)
})

var _ = AfterSuite(func() {
	err := argoClient.DeleteWorkflow(namespace, workflowName)
	Expect(err).NotTo(HaveOccurred())

	retrieved, err := argoClient.GetWorkflow(namespace, workflowName)
	Expect(err).To(HaveOccurred())
	Expect(retrieved).To(BeNil())
})
