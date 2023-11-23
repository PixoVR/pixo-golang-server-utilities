package argo_test

import (
	"context"
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/k8s/argo"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
	"time"
)

var _ = Describe("Stream", func() {

	It("can return an error if the namespace is empty", func() {
		invalidStreamer, err := argo.NewLogsStreamer(k8sClient, argoClient, "", workflowName, bucketName)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("not be empty"))
		Expect(invalidStreamer).To(BeNil())
	})

	It("can return an error if the workflow name is empty", func() {
		invalidStreamer, err := argo.NewLogsStreamer(k8sClient, argoClient, namespace, "", bucketName)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("not be empty"))
		Expect(invalidStreamer).To(BeNil())
	})

	It("can return an error if the workflow is not found", func() {
		nonexistentStreamer, err := argo.NewLogsStreamer(k8sClient, argoClient, namespace, "nonexistent-workflow", bucketName)
		Expect(err).NotTo(HaveOccurred())

		_, err = nonexistentStreamer.Start(context.Background())

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("not found"))
	})

	Context("after running a whalesay workflow", func() {

		var (
			ctx      context.Context
			streamer *argo.LogsStreamer
		)

		BeforeEach(func() {
			var err error
			streamer, err = argo.NewLogsStreamer(k8sClient, argoClient, namespace, workflow.Name, bucketName)
			Expect(err).NotTo(HaveOccurred())
			Expect(streamer).NotTo(BeNil())

			ctx = context.Background()
		})

		It("can stream logs and read the archives", func() {
			stream, err := streamer.Start(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(stream).NotTo(BeNil())

			Expect(streamer.NumNodes()).To(Equal(2))
			Expect(streamer.NumDone()).To(Equal(0))
			Expect(streamer.IsDone()).To(BeFalse())

			readNLogsFromChannelAndExpectLinesTo(ContainSubstring("~~~"), 2, stream)

			time.Sleep(3 * time.Second)
			Expect(streamer.NumDone()).To(Equal(2))
			Expect(streamer.NumNodes()).To(Equal(2))
			Expect(streamer.IsDone()).To(BeTrue())
			Expect(argo.IsClosed(stream)).To(BeTrue())

			archivedLogs, err := streamer.GetArchivedLogsForTemplate(context.Background(), templateOneName)
			Expect(err).NotTo(HaveOccurred())
			Expect(archivedLogs).NotTo(BeNil())

			logBytes := make([]byte, 1024)
			n, err := archivedLogs.Read(logBytes)
			Expect(err).NotTo(HaveOccurred())
			Expect(n).To(BeNumerically(">", 0))
			Expect(string(logBytes)).To(ContainSubstring("~~~"))
			Expect(archivedLogs.Close()).To(Succeed())

			archivedLogs, err = streamer.GetArchivedLogsForTemplate(context.Background(), templateTwoName)
			Expect(err).NotTo(HaveOccurred())
			Expect(archivedLogs).NotTo(BeNil())

			logBytes = make([]byte, 1024)
			n, err = archivedLogs.Read(logBytes)
			Expect(err).NotTo(HaveOccurred())
			Expect(n).To(BeNumerically(">", 0))
			Expect(logBytes).To(ContainSubstring("~~~"))
			Expect(archivedLogs.Close()).To(Succeed())
		})

		It("can wait for the workflow to end and then read archived logs", func() {
			time.Sleep(30 * time.Second)

			stream, err := streamer.Start(ctx)
			Expect(err).NotTo(HaveOccurred())
			Expect(stream).NotTo(BeNil())

			readNLogsFromChannelAndExpectLinesTo(ContainSubstring("~~~"), 2, stream)

			time.Sleep(5 * time.Second)
			Expect(streamer.NumDone()).To(Equal(2))
			Expect(streamer.NumNodes()).To(Equal(2))
			Expect(streamer.IsDone()).To(BeTrue())
			Expect(argo.IsClosed(stream)).To(BeTrue())
		})

	})

})

func readNLogsFromChannelAndExpectLinesTo(matcher types.GomegaMatcher, n int, ch <-chan argo.Log) {
	for i := 0; i < n; i++ {
		streamLog := <-ch
		ExpectLinesTo(matcher, &streamLog)
	}
}

func ExpectLinesTo(matcher types.GomegaMatcher, streamLog *argo.Log) {
	Expect(streamLog).NotTo(BeNil())
	Expect(streamLog.Lines).To(matcher)
}
