package client

type BasicUploadableObject struct {
	BucketName        string
	UploadDestination string
	Filename          string
}

func (b BasicUploadableObject) GetBucketName() string {
	return b.BucketName
}

func (b BasicUploadableObject) GetUploadDestination() string {
	return b.UploadDestination
}

func (b BasicUploadableObject) GetFileLocation() string {
	return b.Filename
}
