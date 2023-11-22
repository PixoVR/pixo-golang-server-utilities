package gcs

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	client "github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/blobstorage"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
)

func (c Client) UploadFile(ctx context.Context, object client.UploadableObject, fileReader io.Reader) (string, error) {

	storageClient, err := storage.NewClient(ctx)
	if err != nil {
		log.Error().Err(err).Msg("unable to create storage client")
		return "", err
	}

	sw := storageClient.Bucket(c.getBucketName(object)).Object(client.GetFullPath(object)).NewWriter(ctx)

	if _, err = io.Copy(sw, fileReader); err != nil {
		log.Error().Err(err).Msg("unable to copy file to bucket")
		return "", err
	}

	if err = sw.Close(); err != nil {
		log.Error().Err(err).Msg("unable to close file writer")
		return "", err
	}

	url, err := c.GetSignedURL(ctx, object)
	if err != nil {
		return "", err
	}

	return url, nil
}

func (c Client) InitResumableUpload(ctx context.Context, object client.UploadableObject) (*client.ResumableUploadResponse, error) {
	log.Debug().Msgf("Initializing resumable upload for %s/%s", c.getBucketName(object), object.GetUploadDestination())

	url := fmt.Sprintf("https://storage.googleapis.com/%s/%s", c.getBucketName(object), client.GetFullPath(object))

	request, err := http.NewRequestWithContext(ctx, "POST", url, nil)
	if err != nil {
		log.Error().Err(err).Msgf("initializing request httpClient failed: %v", err)
		return nil, err
	}

	accessToken, err := getAccessToken()
	if err != nil {
		log.Error().Err(err).Msg("unable to get access token")
		return nil, err
	}

	request.Header.Add("Authorization", "Bearer "+accessToken)
	request.Header.Add("Content-Length", "0")
	request.Header.Add("x-goog-resumable", "start")

	httpClient := &http.Client{}
	response, err := httpClient.Do(request)
	if err != nil {
		log.Error().Err(err).Msg("unable to initiate multipart upload")
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(response.Body)
		bodyString := string(bodyBytes)
		err = fmt.Errorf("unable to initiate multipart upload: %s", bodyString)
		log.Error().Err(err)
		return nil, err
	}

	location := response.Header.Get("Location")
	log.Debug().Msgf("multipart upload url: %s", location)

	return &client.ResumableUploadResponse{
		UploadURL: location,
		Method:    http.MethodPut,
	}, nil
}
