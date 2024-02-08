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

	sw := storageClient.Bucket(c.getBucketName(object)).Object(client.GetFullPath(object)).NewWriter(ctx)

	if _, err = io.Copy(sw, fileReader); err != nil {
		return "", err
	}

	if err = sw.Close(); err != nil {
		return "", err
	}

	url, err := c.GetSignedURL(ctx, object)
	if err != nil {
		return "", err
	}

	return url, nil
}

func (c Client) InitResumableUpload(ctx context.Context, object client.UploadableObject) (*client.ResumableUploadResponse, error) {
	url := fmt.Sprintf("https://storage.googleapis.com/%s/%s", c.getBucketName(object), client.GetFullPath(object))

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
