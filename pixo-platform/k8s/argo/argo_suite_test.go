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
	bucketName         = "dev-multiplayer-allocator-build-logs"
	namespace          = "dev-multiplayer"
	serviceAccountName = "multiplayer-workload"
	workflowName       = "whalesay-"
	templateOneName    = fmt.Sprintf("%s1", workflowName)
	templateTwoName    = fmt.Sprintf("%s2", workflowName)

	k8sClient  base.Client
	argoClient argo.Client
	workflow   *v1alpha1.Workflow

	archiveLogs = true

	whalesaySpec = &v1alpha1.Workflow{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: workflowName,
		},
		Spec: v1alpha1.WorkflowSpec{
			Entrypoint:  workflowName,
			ArchiveLogs: &archiveLogs,
			PodGC: &v1alpha1.PodGC{
				Strategy: v1alpha1.PodGCOnWorkflowCompletion,
			},
			Templates: []v1alpha1.Template{
				{
					Name:               workflowName,
					ServiceAccountName: serviceAccountName,
					ArchiveLocation: &v1alpha1.ArtifactLocation{
						GCS: &v1alpha1.GCSArtifact{
							GCSBucket: v1alpha1.GCSBucket{
								Bucket: bucketName,
							},
							Key: "whalesay-{{workflow.name}}-{{pod.name}}-{{time}}",
						},
					},
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
					Name:               templateOneName,
					ServiceAccountName: serviceAccountName,
					Container: &corev1.Container{
						Name:    templateOneName,
						Image:   "docker/whalesay:latest",
						Command: []string{"cowsay"},
					},
				},
				{
					Name:               templateTwoName,
					ServiceAccountName: serviceAccountName,
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
	err := argoClient.DeleteWorkflow(namespace, workflow.Name)
	Expect(err).NotTo(HaveOccurred())

	retrieved, err := argoClient.GetWorkflow(namespace, workflow.Name)
	Expect(err).To(HaveOccurred())
	Expect(retrieved).To(BeNil())
})
