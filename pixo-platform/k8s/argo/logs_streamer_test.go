package argo_test

import (
	"context"
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/k8s/argo"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Stream", Ordered, func() {

	var (
		streamer *argo.LogsStreamer
	)

	BeforeAll(func() {
		streamer = argo.NewLogsStreamer(k8sClient, argoClient, workflowName)
		Expect(streamer).NotTo(BeNil())
	})

	It("can return an error if the workflow name is empty", func() {
		wf, err := streamer.TailAll(context.Background(), namespace, "")
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("not be empty"))
		Expect(wf).To(BeNil())
	})

	It("can return an error if the workflow is not found", func() {
		wf, err := streamer.TailAll(context.Background(), namespace, "nonexistent-workflow")
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("not found"))
		Expect(wf).To(BeNil())
	})

	It("can tail logs for all templates in a workflow", func() {
		logStreams, err := streamer.TailAll(context.Background(), namespace, workflowName)
		Expect(err).To(BeNil())
		Expect(logStreams).NotTo(BeNil())

		Expect(streamer.NumStreams()).To(Equal(2))
		Expect(streamer.IsDone()).To(BeFalse())

		logOne := streamer.ReadFromStream(templateOneName)
		Expect(logOne).NotTo(BeNil())
		Expect(logOne.Step).To(Equal(templateOneName))
		Expect(logOne.Lines).NotTo(BeEmpty())
		Expect(streamer.NumClosed()).To(Equal(1))

		emptyLog := streamer.ReadFromStream(templateOneName)
		Expect(emptyLog).To(BeNil())

		logTwo := streamer.ReadFromStream(templateTwoName)
		Expect(logTwo).NotTo(BeNil())
		Expect(logTwo.Step).To(Equal(templateTwoName))
		Expect(logTwo.Lines).NotTo(BeEmpty())
		Expect(streamer.NumClosed()).To(Equal(2))

		emptyLog = streamer.ReadFromStream(templateTwoName)
		Expect(emptyLog).To(BeNil())

		Expect(streamer.NumStreams()).To(Equal(2))
		Expect(streamer.IsDone()).To(BeTrue())
	})

})
