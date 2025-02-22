package workflows_test

import (
	"context"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Argo", func() {

	It("can get the list of workflows ", func() {
		workflows, err := workflowsClient.ListWorkflows(context.Background(), namespace)
		Expect(err).NotTo(HaveOccurred())
		Expect(workflows).NotTo(BeNil())
		Expect(len(workflows)).To(BeNumerically(">", 0))
	})

	It("can get the whalesay workflow", func() {
		retrieved, err := workflowsClient.GetWorkflow(ctx, namespace, workflow.Name)
		Expect(err).NotTo(HaveOccurred())
		Expect(retrieved).NotTo(BeNil())
		Expect(retrieved.Name).To(Equal(workflow.Name))
	})

})
