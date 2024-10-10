package workflows

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

	log.Debug().Int("workflows", len(workflows)).Msg("Fetched workflows")
	return workflows, nil
}

func (c Client) GetWorkflow(ctx context.Context, namespace, name string) (*v1alpha1.Workflow, error) {
	workflow, err := c.
		ArgoprojV1alpha1().
		Workflows(namespace).
		Get(ctx, name, metav1.GetOptions{})

	if err != nil {
		log.Err(err).Msg("Error fetching workflows")
		return nil, err
	}

	log.Debug().Str("workflow", workflow.Name).Msg("Fetched workflow")
	return workflow, nil
}

func (c Client) CreateWorkflow(ctx context.Context, namespace string, workflow *v1alpha1.Workflow) (*v1alpha1.Workflow, error) {
	workflow, err := c.
		ArgoprojV1alpha1().
		Workflows(namespace).
		Create(ctx, workflow, metav1.CreateOptions{})

	if err != nil {
		log.Err(err).Msg("Error creating workflow")
		return nil, err
	}

	log.Debug().Str("workflow", workflow.Name).Msg("Created workflow")
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

	log.Debug().Str("workflow", name).Msg("Deleted workflow")
	return nil
}

func (c Client) GetNode(ctx context.Context, workflow *v1alpha1.Workflow, name string) (*v1alpha1.NodeStatus, error) {
	workflow, err := c.GetWorkflow(ctx, workflow.Namespace, workflow.Name)
	if err != nil {
		return nil, err
	}

	var selectedNode *v1alpha1.NodeStatus

	for _, node := range workflow.Status.Nodes {
		if node.Type == v1alpha1.NodeTypePod && node.TemplateName == name {
			selectedNode = &node
			log.Debug().
				Str("workflow", workflow.Name).
				Str("node", node.Name).
				Msg("Node found")
			break
		}
	}

	if selectedNode == nil {
		log.Debug().
			Str("workflow", workflow.Name).
			Msg("Node not found")
		return nil, errors.New("node not found")
	}

	log.Debug().
		Str("workflow", workflow.Name).
		Str("node", selectedNode.Name).
		Msg("Node found")
	return selectedNode, nil
}
