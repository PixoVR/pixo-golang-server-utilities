package argo

import (
	"context"
	"errors"
	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	"github.com/rs/zerolog/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sort"
)

func (c Client) ListWorkflows(ctx context.Context, namespace string) ([]v1alpha1.Workflow, error) {

	workflowList, err := c.
		ArgoprojV1alpha1().
		Workflows(namespace).
		List(ctx, metav1.ListOptions{})

	if err != nil {
		log.Err(err).Msg("Error fetching workflows")
		return nil, err
	}

	workflows := workflowList.Items

	sort.Slice(workflows, func(i, j int) bool {
		return workflows[j].CreationTimestamp.Time.Before(workflows[i].CreationTimestamp.Time)
	})

	return workflows, nil
}

func (c Client) GetWorkflow(namespace, name string) (*v1alpha1.Workflow, error) {

	workflow, err := c.
		ArgoprojV1alpha1().
		Workflows(namespace).
		Get(context.Background(), name, metav1.GetOptions{})

	if err != nil {
		log.Err(err).Msg("Error fetching workflows")
		return nil, err
	}

	return workflow, nil
}

func (c Client) CreateWorkflow(namespace string, workflow *v1alpha1.Workflow) (*v1alpha1.Workflow, error) {

	workflow, err := c.
		ArgoprojV1alpha1().
		Workflows(namespace).
		Create(context.Background(), workflow, metav1.CreateOptions{})

	if err != nil {
		log.Err(err).Msg("Error creating workflow")
		return nil, err
	}

	return workflow, nil
}

func (c Client) DeleteWorkflow(namespace, name string) error {

	err := c.
		ArgoprojV1alpha1().
		Workflows(namespace).
		Delete(context.Background(), name, metav1.DeleteOptions{})

	if err != nil {
		log.Err(err).Msg("Error deleting workflow")
		return err
	}

	return nil
}

func (c Client) GetNode(workflow *v1alpha1.Workflow, name string) (*v1alpha1.NodeStatus, error) {
	var selectedNode *v1alpha1.NodeStatus

	for _, node := range workflow.Status.Nodes {
		if node.Type == v1alpha1.NodeTypePod && node.TemplateName == name {
			selectedNode = &node
			break
		}
	}

	if selectedNode == nil {
		return nil, errors.New("node not found")
	}

	return selectedNode, nil
}
