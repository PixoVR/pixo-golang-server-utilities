package workflows_test

import (
	"context"
	"fmt"
	"github.com/PixoVR/pixo-golang-server-utilities/argo/workflows"
	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"os"
	"testing"
	"time"
)

var (
	ctx context.Context
)

func TestArgo(t *testing.T) {
	RegisterFailHandler(Fail)
	_ = os.Setenv("LIFECYCLE", "local")
	RunSpecs(t, "Argo Client Suite")
}

var (
	bucketName string
	namespace  string

	serviceAccountName string

	secretName    string
	secretKeyName string

	workflowName    string
	templateOneName string
	templateTwoName string

	k8sClient       kubernetes.Interface
	workflowsClient *workflows.Client
	workflow        *v1alpha1.Workflow

	archiveLogs = true
)

var _ = BeforeSuite(func() {
	ctx = context.Background()

	var ok bool
	bucketName, ok = os.LookupEnv("GCS_BUCKET_NAME")
	if !ok {
		bucketName = "pixo-test-bucket"
	}

	namespace, ok = os.LookupEnv("NAMESPACE")
	if !ok {
		namespace = "test"
	}

	serviceAccountName, ok = os.LookupEnv("SA_NAME")
	if !ok {
		serviceAccountName = "test-sa"
	}

	secretName, ok = os.LookupEnv("SECRET_NAME")
	if !ok {
		secretName = "google-credentials"
	}

	secretKeyName, ok = os.LookupEnv("SECRET_KEY_NAME")
	if !ok {
		secretKeyName = "credentials"
	}

	var err error
	k8sClient, err = workflows.NewLocalBaseClient()
	Expect(err).NotTo(HaveOccurred())
	Expect(k8sClient).NotTo(BeNil())

	workflowsClient, err = workflows.NewLocalClient()
	Expect(err).NotTo(HaveOccurred())
	Expect(workflowsClient).NotTo(BeNil())

	workflowName = "whalesay-"

	templateOneName = fmt.Sprintf("%s1", workflowName)
	templateTwoName = fmt.Sprintf("%s2", workflowName)

	whalesaySpec := &v1alpha1.Workflow{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: workflowName,
		},
		Spec: v1alpha1.WorkflowSpec{
			Entrypoint:  workflowName,
			ArchiveLogs: &archiveLogs,
			ArtifactRepositoryRef: &v1alpha1.ArtifactRepositoryRef{
				ConfigMap: "artifact-repositories",
				Key:       "gcs-artifact-repository",
			},
			PodGC: &v1alpha1.PodGC{
				Strategy: v1alpha1.PodGCOnWorkflowSuccess,
			},
			Templates: []v1alpha1.Template{
				{
					Name:               workflowName,
					ServiceAccountName: serviceAccountName,
					ArchiveLocation: &v1alpha1.ArtifactLocation{
						GCS: &v1alpha1.GCSArtifact{
							GCSBucket: v1alpha1.GCSBucket{
								Bucket: bucketName,
								ServiceAccountKeySecret: &corev1.SecretKeySelector{
									Optional: &[]bool{false}[0],
									LocalObjectReference: corev1.LocalObjectReference{
										Name: secretName,
									},
									Key: secretKeyName,
								},
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
	workflow, err = workflowsClient.CreateWorkflow(ctx, namespace, whalesaySpec)
	Expect(err).NotTo(HaveOccurred())
	Expect(workflow).NotTo(BeNil())
	time.Sleep(3 * time.Second)
})

var _ = AfterSuite(func() {
	Expect(workflowsClient.DeleteWorkflow(namespace, workflow.Name)).NotTo(HaveOccurred())

	retrieved, err := workflowsClient.GetWorkflow(ctx, namespace, workflow.Name)
	Expect(err).To(HaveOccurred())
	Expect(retrieved).To(BeNil())
})
