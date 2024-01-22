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

func (a Archive) GetUploadDestination() string {
	return fmt.Sprintf("%s/%s", a.WorkflowName, a.PodName)
}

func (a Archive) GetFilename() string {
	return "main.log"
}
