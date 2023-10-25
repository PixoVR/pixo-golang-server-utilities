package gcs

type SignedURLPart struct {
	PartNumber int    `json:"partNumber"`
	SignedURL  string `json:"signedUrl"`
}
