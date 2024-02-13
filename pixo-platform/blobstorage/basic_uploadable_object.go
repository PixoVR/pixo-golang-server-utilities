package client

type BasicUploadable struct {
	BucketName        string
	UploadDestination string
	Filename          string
}

func (b BasicUploadable) GetBucketName() string {
	return b.BucketName
}

func (b BasicUploadable) GetFileLocation() string {
	return b.UploadDestination + "/" + b.Filename
}

type PathUploadable struct {
	BucketName string
	Filepath   string
}

func (p PathUploadable) GetBucketName() string {
	return p.BucketName
}

func (p PathUploadable) GetFileLocation() string {
	return p.Filepath
}
