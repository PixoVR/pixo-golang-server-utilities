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
		streamer = argo.NewLogsStreamer(k8sClient, workflowName)
		Expect(streamer).NotTo(BeNil())
	})

	It("can stream logs from a workflow", func() {
		workflow, err := argoClient.GetWorkflow(namespace, workflowName)
		Expect(err).NotTo(HaveOccurred())
		Expect(workflow).NotTo(BeNil())
		node := workflow.Status.Nodes[workflow.Name]
		Expect(node).NotTo(BeNil())

		stream, err := streamer.Tail(context.Background(), namespace, &node)
		Expect(err).To(BeNil())
		Expect(stream).NotTo(BeNil())
		log := <-stream
		Expect(log.Step).To(Equal(workflowName))
	})

})
