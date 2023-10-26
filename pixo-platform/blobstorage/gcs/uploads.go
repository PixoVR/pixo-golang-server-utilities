package gcs

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/blobstorage/blob"
	"github.com/rs/zerolog/log"
	"io"
	"mime/multipart"
	"net/http"
)

func UploadFile(bucketName, objectName string, file multipart.File) (string, error) {
	ctx := context.Background()
	storageClient, err := storage.NewClient(ctx)
	if err != nil {
		log.Error().Err(err).Msg("unable to create storage client")
		return "", err
	}

	sw := storageClient.Bucket(bucketName).Object(objectName).NewWriter(ctx)

	if _, err = io.Copy(sw, file); err != nil {
		log.Error().Err(err).Msg("unable to copy file to bucket")
		return "", err
	}

	if err = sw.Close(); err != nil {
		log.Error().Err(err).Msg("unable to close file writer")
		return "", err
	}

	url, err := GetSignedURL(bucketName, objectName)
	if err != nil {
		return "", err
	}

	return url, nil
}

func InitResumableUpload(bucketName, objectName string) (blob.ResumableUploadResponse, error) {
	log.Debug().Msgf("Initializing resumable upload for %s/%s", bucketName, objectName)

	url := fmt.Sprintf("https://storage.googleapis.com/%s/%s", bucketName, objectName)

	request, err := http.NewRequest("POST", url, nil)
	if err != nil {
		log.Error().Err(err).Msgf("initializing request client failed: %v", err)
		return blob.ResumableUploadResponse{}, err
	}

	accessToken, err := getAccessToken()
	if err != nil {
		log.Error().Err(err).Msg("unable to get access token")
		return blob.ResumableUploadResponse{}, err
	}

	request.Header.Add("Authorization", "Bearer "+accessToken)
	request.Header.Add("Content-Length", "0")
	request.Header.Add("x-goog-resumable", "start")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Error().Err(err).Msg("unable to initiate multipart upload")
		return blob.ResumableUploadResponse{}, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(response.Body)
		bodyString := string(bodyBytes)
		err = fmt.Errorf("unable to initiate multipart upload: %s", bodyString)
		log.Error().Err(err)
		return blob.ResumableUploadResponse{}, err
	}

	location := response.Header.Get("Location")
	log.Debug().Msgf("multipart upload url: %s", location)

	return blob.ResumableUploadResponse{
		UploadURL: location,
		Method:    http.MethodPut,
	}, nil
}
