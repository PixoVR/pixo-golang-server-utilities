package workflows

import (
	"fmt"
)

type Archive struct {
	WorkflowName string `json:"workflowName"`
	PodName      string `json:"podName"`
	BucketName   string `json:"-"`
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
