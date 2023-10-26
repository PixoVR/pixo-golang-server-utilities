package blob

import "fmt"

type Blob struct {
	bucketName        string
	uploadDestination string
	filename          string
}

func (b *Blob) FullFilepath() string {
	return fmt.Sprintf("%s/%s", b.bucketName, b.Filepath())
}

func (b *Blob) Filepath() string {
	if b.uploadDestination == "" {
		return b.filename
	}

	return fmt.Sprintf("%s/%s", b.uploadDestination, b.filename)
}

func (b *Blob) BucketName() string {
	return b.bucketName
}

func (b *Blob) UploadDestination() string {
	return b.uploadDestination
}

func (b *Blob) Filename() string {
	return b.filename
}

func (b *Blob) SetBucketName(bucketName string) {
	b.bucketName = bucketName
}

func (b *Blob) SetUploadDestination(uploadDestination string) {
	b.uploadDestination = uploadDestination
}

func (b *Blob) SetFilename(filename string) {
	b.filename = filename
}
