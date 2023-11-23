package argo

import (
	"context"
	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
)

func (s *LogsStreamer) readFromCache(ctx context.Context, workflow *v1alpha1.Workflow) error {
	return s.logsCache.Get(ctx, workflow.Name).Err()
}

func (s *LogsStreamer) writeToCache(ctx context.Context, workflow *v1alpha1.Workflow) error {
	return s.logsCache.Set(ctx, workflow.Name, workflow.Status.Nodes, 0).Err()
}
