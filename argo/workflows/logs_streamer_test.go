package workflows_test

import (
	"context"
	"github.com/PixoVR/pixo-golang-server-utilities/argo/workflows"
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/blobstorage/gcs"
	"github.com/alicebob/miniredis/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
	"github.com/redis/go-redis/v9"
	"io"
	"time"
)

var _ = Describe("Stream", func() {

	var (
		storageClient gcs.Client
	)

	BeforeEach(func() {
		var err error
		storageClient, err = gcs.NewClient(gcs.Config{BucketName: bucketName})
		Expect(err).NotTo(HaveOccurred())
		Expect(storageClient).NotTo(BeNil())
	})

	It("can return an error if the namespace is empty", func() {
		invalidStreamer, err := workflows.NewLogsStreamer(workflows.StreamerConfig{
			K8sClient:     k8sClient,
			ArgoClient:    workflowsClient,
			StorageClient: storageClient,
			Namespace:     "",
			WorkflowName:  workflowName,
		})
		Expect(err).To(MatchError("namespace may not be empty"))
		Expect(invalidStreamer).To(BeNil())
	})

	It("can return an error if the workflow name is empty", func() {
		invalidStreamer, err := workflows.NewLogsStreamer(workflows.StreamerConfig{
			K8sClient:     k8sClient,
			ArgoClient:    workflowsClient,
			StorageClient: storageClient,
			Namespace:     namespace,
			WorkflowName:  "",
		})
		Expect(err).To(MatchError("workflowName may not be empty"))
		Expect(invalidStreamer).To(BeNil())
	})

	It("can return an error if the workflow is not found", func() {
		nonexistentStreamer, err := workflows.NewLogsStreamer(workflows.StreamerConfig{
			K8sClient:     k8sClient,
			ArgoClient:    workflowsClient,
			StorageClient: storageClient,
			Namespace:     namespace,
			WorkflowName:  "nonexistent-workflow",
		})
		Expect(err).NotTo(HaveOccurred())

		_, err = nonexistentStreamer.Start(context.Background())
		Expect(err).To(MatchError("workflow not found"))
	})

	Context("after running a whalesay workflow", func() {

		var (
			streamer *workflows.LogsStreamer
			s        *miniredis.Miniredis
			cache    *redis.Client
		)

		BeforeEach(func() {
			var err error
			cache, s, err = workflows.NewMiniCache()
			Expect(err).NotTo(HaveOccurred())
			Expect(s).NotTo(BeNil())
			Expect(cache).NotTo(BeNil())

			streamer, err = workflows.NewLogsStreamer(workflows.StreamerConfig{
				K8sClient:     k8sClient,
				ArgoClient:    workflowsClient,
				StorageClient: storageClient,
				Namespace:     namespace,
				WorkflowName:  workflow.Name,
				LogsCache:     cache,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(streamer).NotTo(BeNil())
		})

		AfterEach(func() {
			s.Close()
		})

		It("can stream logs and read the archives", func() {
			stream, err := streamer.Start(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(stream).NotTo(BeNil())

			Expect(streamer.NumNodes()).To(Equal(2))
			Expect(streamer.NumDone()).To(Equal(0))
			Expect(streamer.IsDone()).To(BeFalse())

			readNLogsAndExpectLinesTo(ContainSubstring("~~~"), 2, stream)

			time.Sleep(10 * time.Second)
			Expect(streamer.NumNodes()).To(Equal(2))
			Expect(streamer.NumDone()).To(Equal(2))
			Expect(streamer.IsDone()).To(BeTrue())

			archivedLogs, err := streamer.GetArchivedLogsForTemplate(context.Background(), templateOneName)
			Expect(err).NotTo(HaveOccurred())
			Expect(archivedLogs).NotTo(BeNil())
			readBytesAndExpectTo(ContainSubstring("~~~"), archivedLogs)
			Expect(archivedLogs.Close()).To(Succeed())

			archivedLogs, err = streamer.GetArchivedLogsForTemplate(context.Background(), templateTwoName)
			Expect(err).NotTo(HaveOccurred())
			Expect(archivedLogs).NotTo(BeNil())
			readBytesAndExpectTo(ContainSubstring("~~~"), archivedLogs)
			Expect(archivedLogs.Close()).To(Succeed())
		})

		It("can wait for the workflow to end and then read archived logs", func() {
			time.Sleep(30 * time.Second)

			stream, err := streamer.Start(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(stream).NotTo(BeNil())

			readLogsUntilDoneAndExpectLinesTo(Not(BeEmpty()), stream)

			time.Sleep(2 * time.Second)
			Expect(streamer.NumDone()).To(Equal(2))
			Expect(streamer.NumNodes()).To(Equal(2))
			Expect(streamer.IsDone()).To(BeTrue())
		})

	})

})

func readNLogsAndExpectLinesTo(matcher types.GomegaMatcher, n int, ch <-chan workflows.Log) {
	for i := 0; i < n; i++ {
		streamLog, ok := <-ch
		if !ok {
			break
		}

		ExpectLinesTo(matcher, &streamLog)
	}
}

func readLogsUntilDoneAndExpectLinesTo(matcher types.GomegaMatcher, ch <-chan workflows.Log) {
	for {
		streamLog, ok := <-ch
		if !ok {
			break
		}

		ExpectLinesTo(matcher, &streamLog)
	}
}

func ExpectLinesTo(matcher types.GomegaMatcher, streamLog *workflows.Log) {
	Expect(streamLog).NotTo(BeNil())
	Expect(streamLog.Step).NotTo(BeEmpty())
	Expect(streamLog.Lines).To(matcher)
}

func readBytesAndExpectTo(matcher types.GomegaMatcher, r io.Reader) {
	logBytes := make([]byte, 1024)
	n, err := r.Read(logBytes)
	Expect(err).NotTo(HaveOccurred())
	Expect(n).To(BeNumerically(">", 0))
	Expect(string(logBytes)).To(matcher)
}
