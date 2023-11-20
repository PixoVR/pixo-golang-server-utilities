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
		invalidStreamer, err := argo.NewLogsStreamer(k8sClient, argoClient, "", workflowName)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("not be empty"))
		Expect(invalidStreamer).To(BeNil())
	})

	It("can return an error if the workflow name is empty", func() {
		invalidStreamer, err := argo.NewLogsStreamer(k8sClient, argoClient, namespace, "")
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("not be empty"))
		Expect(invalidStreamer).To(BeNil())
	})

	It("can return an error if the workflow is not found", func() {
		nonexistentStreamer, err := argo.NewLogsStreamer(k8sClient, argoClient, namespace, "nonexistent-workflow")
		Expect(err).NotTo(HaveOccurred())

		_, err = nonexistentStreamer.TailAll(context.Background())

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("not found"))
	})

	It("can tail logs for all templates in a workflow", func() {
		streamer, err := argo.NewLogsStreamer(k8sClient, argoClient, namespace, workflow.Name)
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
	})

	It("can tail the logs with a combined stream then pull the archived logs from cloud storage", func() {
		newWorkflow, err := argoClient.CreateWorkflow(namespace, whalesaySpec)
		Expect(err).NotTo(HaveOccurred())
		Expect(newWorkflow).NotTo(BeNil())

		newStreamer, err := argo.NewLogsStreamer(k8sClient, argoClient, namespace, newWorkflow.Name)
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

		// Pull the archived logs from cloud storage
		archivedLogs, err := newStreamer.GetArchivedLogs(context.Background())
		Expect(err).NotTo(HaveOccurred())
		Expect(archivedLogs).NotTo(BeEmpty())
	})

})

func readFromStreamAndExpectLinesTo(matcher types.GomegaMatcher, streamer *argo.LogsStreamer, templateName string) {
	logTwo := streamer.ReadFromStream(templateName)
	Expect(logTwo).NotTo(BeNil())
	Expect(logTwo.Step).To(Equal(templateName))
	Expect(logTwo.Lines).To(matcher)
}
