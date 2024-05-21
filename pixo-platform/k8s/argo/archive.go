package argo

import (
	"fmt"
)

type Archive struct {
	BucketName   string
	WorkflowName string `json:"workflowName"`
	PodName      string `json:"podName"`
}

func (a Archive) GetBucketName() string {
	return a.BucketName
}

func (a Archive) GetFileLocation() string {
	return fmt.Sprintf("%s/%s/main.log", a.WorkflowName, a.PodName)
}

func (a Archive) GetTimestamp() int64 {
	return 0
}
