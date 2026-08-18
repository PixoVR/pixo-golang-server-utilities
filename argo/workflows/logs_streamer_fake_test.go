package workflows_test

import (
	"context"
	"errors"
	"time"

	"github.com/PixoVR/pixo-golang-server-utilities/argo/workflows"
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/blobstorage"
	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	argofake "github.com/argoproj/argo-workflows/v3/pkg/client/clientset/versioned/fake"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sfake "k8s.io/client-go/kubernetes/fake"
)

const (
	fakeNamespace = "fake-namespace"
	fakeWorkflow  = "fake-whalesay"
	fakeStepOne   = "fake-step-1"
	fakeStepTwo   = "fake-step-2"
)

// newFakeWorkflow builds a two pod template workflow whose status can be moved
// through the phases the streamer has to cope with.
func newFakeWorkflow(phase v1alpha1.WorkflowPhase, nodePhase v1alpha1.NodePhase, withNodes bool) *v1alpha1.Workflow {
	workflow := &v1alpha1.Workflow{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fakeWorkflow,
			Namespace: fakeNamespace,
		},
		Spec: v1alpha1.WorkflowSpec{
			Entrypoint: fakeWorkflow,
			Templates: []v1alpha1.Template{
				{
					Name: fakeWorkflow,
					DAG: &v1alpha1.DAGTemplate{
						Tasks: []v1alpha1.DAGTask{
							{Name: fakeStepOne, Template: fakeStepOne},
							{Name: fakeStepTwo, Template: fakeStepTwo},
						},
					},
				},
				{
					Name:      fakeStepOne,
					Container: &corev1.Container{Name: fakeStepOne, Image: "docker/whalesay:latest"},
				},
				{
					Name:      fakeStepTwo,
					Container: &corev1.Container{Name: fakeStepTwo, Image: "docker/whalesay:latest"},
				},
			},
		},
		Status: v1alpha1.WorkflowStatus{Phase: phase},
	}

	if withNodes {
		workflow.Status.Nodes = v1alpha1.Nodes{
			"node-1": {
				ID:           "fake-whalesay-1",
				Name:         fakeStepOne,
				BoundaryID:   fakeWorkflow,
				TemplateName: fakeStepOne,
				Type:         v1alpha1.NodeTypePod,
				Phase:        nodePhase,
			},
			"node-2": {
				ID:           "fake-whalesay-2",
				Name:         fakeStepTwo,
				BoundaryID:   fakeWorkflow,
				TemplateName: fakeStepTwo,
				Type:         v1alpha1.NodeTypePod,
				Phase:        nodePhase,
			},
		}
	}

	return workflow
}

var _ = Describe("Stream with fake clients", func() {

	var (
		fakeCtx       context.Context
		cancel        context.CancelFunc
		argoClientset *argofake.Clientset
		streamer      *workflows.LogsStreamer
		storage       *blobstorage.MockStorageClient
	)

	newStreamer := func(workflow *v1alpha1.Workflow) *workflows.LogsStreamer {
		argoClientset = argofake.NewSimpleClientset(workflow)
		storage = blobstorage.NewMockStorageClient()

		newStreamer, err := workflows.NewLogsStreamer(workflows.StreamerConfig{
			K8sClient:     k8sfake.NewSimpleClientset(),
			ArgoClient:    &workflows.Client{Interface: argoClientset},
			StorageClient: storage,
			Namespace:     fakeNamespace,
			WorkflowName:  fakeWorkflow,
		})
		Expect(err).NotTo(HaveOccurred())

		return newStreamer
	}

	updateWorkflow := func(workflow *v1alpha1.Workflow) {
		_, err := argoClientset.ArgoprojV1alpha1().
			Workflows(fakeNamespace).
			Update(fakeCtx, workflow, metav1.UpdateOptions{})
		Expect(err).NotTo(HaveOccurred())
	}

	BeforeEach(func() {
		fakeCtx, cancel = context.WithTimeout(context.Background(), time.Minute)
	})

	AfterEach(func() {
		if streamer != nil {
			Expect(streamer.Close()).To(Succeed())
			streamer = nil
		}
		cancel()
	})

	It("can tail the logs of a running workflow", func() {
		streamer = newStreamer(newFakeWorkflow(v1alpha1.WorkflowRunning, v1alpha1.NodeSucceeded, true))

		stream, err := streamer.Start(fakeCtx)
		Expect(err).NotTo(HaveOccurred())
		Expect(streamer.NumNodes()).To(Equal(2))

		readNLogsAndExpectLinesTo(Not(BeEmpty()), 2, stream)
		Eventually(streamer.IsDone, time.Minute, time.Second).Should(BeTrue())
		Eventually(stream, time.Minute).Should(BeClosed())
	})

	It("can stream the logs of a workflow that has not started yet", func() {
		streamer = newStreamer(newFakeWorkflow("", v1alpha1.NodePending, false))

		stream, err := streamer.Start(fakeCtx)
		Expect(err).NotTo(HaveOccurred())

		Consistently(stream, 2*time.Second).ShouldNot(Receive())
		updateWorkflow(newFakeWorkflow(v1alpha1.WorkflowRunning, v1alpha1.NodeSucceeded, true))

		readNLogsAndExpectLinesTo(Not(BeEmpty()), 2, stream)
	})

	It("can read the archives of a workflow that already completed", func() {
		streamer = newStreamer(newFakeWorkflow(v1alpha1.WorkflowSucceeded, v1alpha1.NodeSucceeded, true))

		stream, err := streamer.Start(fakeCtx)
		Expect(err).NotTo(HaveOccurred())

		readNLogsAndExpectLinesTo(Equal("test"), 2, stream)
		Expect(storage.ReadFileNumTimesCalled).To(Equal(2))
	})

	It("can close the streams of a completed workflow that never created any nodes", func() {
		streamer = newStreamer(newFakeWorkflow(v1alpha1.WorkflowRunning, v1alpha1.NodePending, false))

		stream, err := streamer.Start(fakeCtx)
		Expect(err).NotTo(HaveOccurred())

		updateWorkflow(newFakeWorkflow(v1alpha1.WorkflowSucceeded, v1alpha1.NodePending, false))

		Eventually(streamer.NumDone, time.Minute, time.Second).Should(Equal(2))
		Eventually(stream, time.Minute).Should(BeClosed())
	})

	It("can mark the streams done when the archives are unreadable", func() {
		streamer = newStreamer(newFakeWorkflow(v1alpha1.WorkflowSucceeded, v1alpha1.NodeSucceeded, true))
		storage.ReadFileError = errors.New("bucket is on fire")

		stream, err := streamer.Start(fakeCtx)
		Expect(err).NotTo(HaveOccurred())

		Eventually(streamer.NumDone, time.Minute, time.Second).Should(Equal(2))
		Eventually(stream, time.Minute).Should(BeClosed())
	})

	It("can close the streams of a workflow that never starts", func() {
		streamer = newStreamer(newFakeWorkflow("", v1alpha1.NodePending, false))

		stream, err := streamer.Start(fakeCtx)
		Expect(err).NotTo(HaveOccurred())

		Expect(streamer.Close()).To(Succeed())
		Expect(streamer.Close()).To(Succeed())
		Eventually(stream, time.Minute).Should(BeClosed())
		Expect(streamer.IsDone()).To(BeTrue())
	})

	It("can stop streaming when the context is cancelled", func() {
		streamer = newStreamer(newFakeWorkflow(v1alpha1.WorkflowRunning, v1alpha1.NodePending, false))

		stream, err := streamer.Start(fakeCtx)
		Expect(err).NotTo(HaveOccurred())

		cancel()
		Eventually(stream, time.Minute).Should(BeClosed())
	})

	It("can return an error when the archived node does not exist", func() {
		streamer = newStreamer(newFakeWorkflow(v1alpha1.WorkflowSucceeded, v1alpha1.NodeSucceeded, true))

		archivedLogs, err := streamer.GetArchivedLogsForTemplate(fakeCtx, "nonexistent-template")
		Expect(err).To(MatchError(ContainSubstring("unable to get node")))
		Expect(archivedLogs).To(BeNil())
	})

	It("can return an error when the node has not finished", func() {
		streamer = newStreamer(newFakeWorkflow(v1alpha1.WorkflowRunning, v1alpha1.NodeRunning, true))

		archivedLogs, err := streamer.GetArchivedLogsForTemplate(fakeCtx, fakeStepOne)
		Expect(err).To(MatchError("node is not done"))
		Expect(archivedLogs).To(BeNil())
	})

	It("can return an error when the archives cannot be read", func() {
		streamer = newStreamer(newFakeWorkflow(v1alpha1.WorkflowSucceeded, v1alpha1.NodeSucceeded, true))
		storage.ReadFileError = errors.New("bucket is on fire")

		archivedLogs, err := streamer.GetArchivedLogsForTemplate(fakeCtx, fakeStepOne)
		Expect(err).To(MatchError(ContainSubstring("unable to read archived logs")))
		Expect(archivedLogs).To(BeNil())
	})

	It("can return an error when the workflow does not exist", func() {
		streamer = newStreamer(newFakeWorkflow(v1alpha1.WorkflowSucceeded, v1alpha1.NodeSucceeded, true))
		Expect(argoClientset.ArgoprojV1alpha1().
			Workflows(fakeNamespace).
			Delete(fakeCtx, fakeWorkflow, metav1.DeleteOptions{})).To(Succeed())

		archivedLogs, err := streamer.GetArchivedLogsForTemplate(fakeCtx, fakeStepOne)
		Expect(err).To(MatchError(ContainSubstring("unable to get workflow")))
		Expect(archivedLogs).To(BeNil())
	})

})

var _ = Describe("Archives", func() {

	It("can build the location of an archived log", func() {
		archive := workflows.Archive{
			BucketName:   "fake-bucket",
			WorkflowName: fakeWorkflow,
			PodName:      "fake-pod",
		}

		Expect(archive.GetBucketName()).To(Equal("fake-bucket"))
		Expect(archive.GetFileLocation()).To(Equal("fake-whalesay/fake-pod/main.log"))
		Expect(archive.GetTimestamp()).To(BeZero())
	})

	It("can format an empty pod name for a nil node", func() {
		Expect(workflows.FormatPodName(nil)).To(BeEmpty())
	})

})
