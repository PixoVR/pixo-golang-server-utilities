package gcs

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	client "github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/blobstorage"
	"io"
	"net/http"
)

func (c Client) UploadFile(ctx context.Context, object client.UploadableObject, fileReader io.Reader) (string, error) {

	storageClient, err := storage.NewClient(ctx)
	if err != nil {
		return "", err
	}

	bucketName := c.getBucketName(object)
	fileLocation := object.GetFileLocation()

	sanitizedFileLocation := c.SanitizeFilename(fileLocation)

	sw := storageClient.
		Bucket(bucketName).
		Object(sanitizedFileLocation).
		NewWriter(ctx)

	if _, err = io.Copy(sw, fileReader); err != nil {
		return "", err
	}

	if err = sw.Close(); err != nil {
		return "", err
	}

	return sanitizedFileLocation, nil
}

func (c Client) InitResumableUpload(ctx context.Context, object client.UploadableObject) (*client.ResumableUploadResponse, error) {
	url := fmt.Sprintf("https://storage.googleapis.com/%s/%s", c.getBucketName(object), object.GetFileLocation())

	request, err := http.NewRequestWithContext(ctx, "POST", url, nil)
	if err != nil {
		return nil, err
	}

	accessToken, err := getAccessToken()
	if err != nil {
		return nil, err
	}

	request.Header.Add("Authorization", "Bearer "+accessToken)
	request.Header.Add("Content-Length", "0")
	request.Header.Add("x-goog-resumable", "start")

	httpClient := &http.Client{}
	response, err := httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(response.Body)
		bodyString := string(bodyBytes)
		return nil, fmt.Errorf("unable to initiate multipart upload: %s", bodyString)
	}

	return &client.ResumableUploadResponse{
		UploadURL: response.Header.Get("Location"),
		Method:    http.MethodPut,
	}, nil
}
