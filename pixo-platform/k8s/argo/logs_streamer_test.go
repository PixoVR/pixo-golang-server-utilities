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

		_, err = nonexistentStreamer.TailAll(context.Background())

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("not found"))
	})

	It("can tail logs for all templates in a workflow separately and read from the archives", func() {
		streamer, err := argo.NewLogsStreamer(k8sClient, argoClient, namespace, workflow.Name, bucketName)
		Expect(err).NotTo(HaveOccurred())
		Expect(streamer).NotTo(BeNil())

		logStreams, err := streamer.TailAll(context.Background())
		Expect(err).To(BeNil())
		Expect(logStreams).NotTo(BeNil())

		Expect(streamer.NumTotalStreams()).To(Equal(2))
		Expect(streamer.IsDone()).To(BeFalse())

		readFromStreamAndExpectLinesTo(Not(BeEmpty()), streamer, templateOneName)
		time.Sleep(3 * time.Second)
		Expect(streamer.NumClosed()).To(Equal(1))
		Expect(streamer.IsDone()).To(BeFalse())

		emptyLog := streamer.ReadFromStream(templateOneName)
		Expect(emptyLog).To(BeNil())

		readFromStreamAndExpectLinesTo(Not(BeEmpty()), streamer, templateTwoName)
		time.Sleep(3 * time.Second)
		Expect(streamer.NumClosed()).To(Equal(2))
		Expect(streamer.ReadFromStream(templateTwoName)).To(BeNil())
		Expect(streamer.NumTotalStreams()).To(Equal(2))
		Expect(streamer.IsDone()).To(BeTrue())

		archivedLogs, err := streamer.GetArchivedLogs(context.Background(), templateOneName)
		Expect(err).NotTo(HaveOccurred())
		Expect(archivedLogs).NotTo(BeNil())
		logBytes := make([]byte, 1024)
		n, err := archivedLogs.Read(logBytes)
		Expect(err).NotTo(HaveOccurred())
		Expect(n).To(BeNumerically(">", 0))
		Expect(archivedLogs.Close()).To(Succeed())
	})

	It("can tail the logs with a combined stream while running", func() {
		newWorkflow, err := argoClient.CreateWorkflow(namespace, whalesaySpec)
		Expect(err).NotTo(HaveOccurred())
		Expect(newWorkflow).NotTo(BeNil())

		newStreamer, err := argo.NewLogsStreamer(k8sClient, argoClient, namespace, newWorkflow.Name, bucketName)
		Expect(err).NotTo(HaveOccurred())
		Expect(newStreamer).NotTo(BeNil())

		time.Sleep(3 * time.Second)
		combinedStream := newStreamer.GetLogsStream()
		Expect(combinedStream).NotTo(BeNil())

		logOne := <-combinedStream
		Expect(logOne).NotTo(BeNil())
		Expect(logOne.Step).To(Equal(templateOneName))
		Expect(logOne.Lines).NotTo(BeEmpty())

		logTwo := <-combinedStream
		Expect(logTwo).NotTo(BeNil())
		Expect(logTwo.Step).To(Equal(templateTwoName))
		Expect(logTwo.Lines).NotTo(BeEmpty())

		time.Sleep(3 * time.Second)
		Expect(newStreamer.IsDone()).To(BeTrue())
	})

	It("can wait for the workflow to end and then read archived logs", func() {

		time.Sleep(30 * time.Second)

		newStreamer, err := argo.NewLogsStreamer(k8sClient, argoClient, namespace, workflow.Name, bucketName)
		Expect(err).NotTo(HaveOccurred())
		Expect(newStreamer).NotTo(BeNil())

		combinedStream := newStreamer.GetLogsStream()
		Expect(combinedStream).NotTo(BeNil())

		logOne := <-combinedStream
		Expect(logOne).NotTo(BeNil())
		Expect(logOne.Lines).NotTo(BeEmpty())

		logTwo := <-combinedStream
		Expect(logTwo).NotTo(BeNil())
		Expect(logTwo.Lines).NotTo(BeEmpty())

		time.Sleep(3 * time.Second)
		Expect(newStreamer.IsDone()).To(BeTrue())
	})

})

func readFromStreamAndExpectLinesTo(matcher types.GomegaMatcher, streamer *argo.LogsStreamer, templateName string) {
	logTwo := streamer.ReadFromStream(templateName)
	Expect(logTwo).NotTo(BeNil())
	Expect(logTwo.Step).To(Equal(templateName))
	Expect(logTwo.Lines).To(matcher)
}
