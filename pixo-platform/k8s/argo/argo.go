package argo

import (
	"context"
	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"github.com/rs/zerolog/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c Client) ListWorkflows() (*v1alpha1.WorkflowList, error) {
	workflows, err := c.Clientset.ArgoprojV1alpha1().Workflows(c.Namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Err(err).Msg("Error fetching workflows")
		return nil, err
	}

	return workflows, nil
}

func (c Client) GetWorkflow(name string) (*v1alpha1.Workflow, error) {
	workflow, err := c.Clientset.ArgoprojV1alpha1().Workflows(c.Namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		log.Err(err).Msg("Error fetching workflows")
		return nil, err
	}

	return workflow, nil
}

func (c Client) CreateWorkflow(workflow *v1alpha1.Workflow) (*v1alpha1.Workflow, error) {
	workflow, err := c.Clientset.ArgoprojV1alpha1().Workflows(c.Namespace).Create(context.Background(), workflow, metav1.CreateOptions{})
	if err != nil {
		log.Err(err).Msg("Error creating workflow")
		return nil, err
	}

	return workflow, nil
}

func (c Client) DeleteWorkflow(name string) error {
	err := c.Clientset.ArgoprojV1alpha1().Workflows(c.Namespace).Delete(context.Background(), name, metav1.DeleteOptions{})
	if err != nil {
		log.Err(err).Msg("Error deleting workflow")
		return err
	}

	return nil
}
