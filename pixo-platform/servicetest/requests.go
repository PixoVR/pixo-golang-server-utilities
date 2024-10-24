package servicetest

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/urlfinder"
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/middleware/auth"
	"github.com/antchfx/jsonquery"
	"github.com/cucumber/godog"
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"os"
	"path/filepath"
	"strings"
)

func (s *ServerTestFeature) MakeRequest(method, tenant, service, endpoint string, body *godog.DocString, paramsMap map[string]string) error {
	var bodyContent []byte
	if body != nil {
		bodyContent = []byte(body.Content)
	}

	bodyContent = s.PerformSubstitutions(bodyContent)

	if err := s.PerformRequest(method, tenant, service, endpoint, bodyContent, paramsMap); err != nil {
		return err
	}

	log.Debug().
		Str("endpoint", endpoint).
		Str("response", s.ResponseString).
		Int("status_code", s.StatusCode).
		Msg("HTTP RESPONSE")

	var responseBody map[string]interface{}
	if err := json.Unmarshal([]byte(s.ResponseString), &responseBody); err == nil {
		if value, ok := responseBody["message"]; ok {
			s.Message = value.(string)
		}
	}

	return nil
}

func (s *ServerTestFeature) PerformRequest(method, tenant, service, endpoint string, body []byte, paramsMap map[string]string) error {

	if s.BeforeRequest != nil {
		s.BeforeRequest(body)
	}

	currentLifecycle := strings.ToLower(viper.GetString("lifecycle"))

	if service == "" && (currentLifecycle == "" || currentLifecycle == "internal") {
		url := fmt.Sprintf("/%s%s", s.ServiceClient.Path(), endpoint)
		log.Debug().Msgf("URL: %s: - Method %s", url, method)

		req, err := http.NewRequest(method, url, bytes.NewReader(body))
		if err != nil {
			return fmt.Errorf("failed to create request: %s", err)
		}

		if s.Token != "" {
			req.Header.Set(auth.AuthorizationHeader, fmt.Sprintf("Bearer %s", s.Token))
		}

		if s.SecretKey != "" {
			req.Header.Set(auth.SecretKeyHeader, s.SecretKey)
		}

		if s.APIKey != "" {
			req.Header.Set(auth.APIKeyHeader, s.APIKey)
		}

		if method == http.MethodPost || method == http.MethodPatch {
			log.Debug().
				Str("url", url).
				Str("method", method).
				Str("body", string(body))
			req.Header.Set("Content-Type", "application/json")
			req.Body = io.NopCloser(bytes.NewReader(body))
		}

		for key, value := range paramsMap {
			paramsMap[key] = string(s.PerformSubstitutions([]byte(value)))
		}

		if len(paramsMap) > 0 {
			q := req.URL.Query()
			for key, value := range paramsMap {
				q.Add(key, value)
			}
			req.URL.RawQuery = q.Encode()
		}

		if s.Engine == nil {
			return errors.New("engine is nil")
		}

		s.Engine.ServeHTTP(s.Recorder, req)

		s.HTTPResponse = s.Recorder.Result()
		s.ResponseString = s.Recorder.Body.String()
		s.StatusCode = s.Recorder.Code
		s.Err = nil

	} else {
		url := s.ServiceClient.GetURL() + endpoint
		url = string(s.PerformSubstitutions([]byte(url)))
		req := s.Client.R()

		if s.Token != "" {
			req.SetAuthToken(s.Token)
		}

		if s.SecretKey != "" {
			req.SetHeader(auth.SecretKeyHeader, s.SecretKey)
		}

		if s.APIKey != "" {
			req.SetHeader(auth.APIKeyHeader, s.APIKey)
		}

		if paramsMap != nil {
			for key, value := range paramsMap {
				paramsMap[key] = string(s.PerformSubstitutions([]byte(value)))
			}
			req.SetQueryParams(paramsMap)
		}

		if service != "" {
			serviceConfig := urlfinder.ServiceConfig{
				Tenant:    tenant,
				Service:   service,
				Region:    viper.GetString("region"),
				Lifecycle: currentLifecycle,
			}
			url = serviceConfig.FormatURL() + endpoint
		}

		log.Debug().
			Str("url", url).
			Str("method", method).
			Msg("Making request")

		var res *resty.Response
		var err error
		switch method {
		case http.MethodGet:
			log.Debug().
				Str("url", url).
				Str("method", method)
			res, err = req.Get(url)
		case http.MethodDelete:
			log.Debug().
				Str("url", url).
				Str("method", method)
			res, err = req.Delete(url)
		case http.MethodPost:
			log.Debug().
				Str("url", url).
				Str("method", method).
				Str("body", string(body)).
				Msg("Making POST request")
			if len(s.FilesToSend) == 0 {
				req.SetHeader("Content-Type", "application/json").
					SetBody(body)
			} else {
				for _, value := range s.FilesToSend {
					log.Debug().
						Str("key", value.Key).
						Str("path", value.Path).
						Msgf("Uploading file")
					req.SetFile(value.Key, value.Path)
				}

				formBodyMap := make(map[string]interface{})
				if err = json.Unmarshal(body, &formBodyMap); err != nil {
					return fmt.Errorf("failed to unmarshal body: %s", err)
				}

				bodyFormData := make(map[string]string)
				for key, value := range formBodyMap {
					bodyFormData[key] = fmt.Sprint(value)
				}

				req.SetFormData(bodyFormData)
			}
			res, err = req.Post(url)
		case http.MethodPatch:
			log.Debug().
				Str("url", url).
				Str("method", method).
				Str("body", string(body)).
				Msg("Making POST request")
			if len(s.FilesToSend) == 0 {
				req.SetHeader("Content-Type", "application/json").
					SetBody(body)
			} else {
				for _, value := range s.FilesToSend {
					log.Debug().
						Str("key", value.Key).
						Str("path", value.Path).
						Msgf("Uploading file")
					req.SetFile(value.Key, value.Path)
				}

				formBodyMap := make(map[string]interface{})
				if err = json.Unmarshal(body, &formBodyMap); err != nil {
					return fmt.Errorf("failed to unmarshal body: %s", err)
				}

				bodyFormData := make(map[string]string)
				for key, value := range formBodyMap {
					bodyFormData[key] = fmt.Sprint(value)
				}

				req.SetFormData(bodyFormData)
			}
			res, err = req.Patch(url)
		}

		s.Err = err
		if err != nil {
			return fmt.Errorf("failed to make request: %s", err)
		}

		if res == nil {
			return errors.New("response is nil")
		}

		s.HTTPResponse = res.RawResponse
		s.ResponseString = string(res.Body())
		s.StatusCode = res.StatusCode()
	}

	s.FilesToSend = nil
	return nil
}

func (s *ServerTestFeature) MakeGraphQLRequest(endpoint, serviceName, body string) error {
	req := s.Client.R()

	body = string(s.PerformSubstitutions([]byte(body)))

	if len(s.FilesToSend) > 0 {
		payload := &bytes.Buffer{}
		writer := multipart.NewWriter(payload)
		_ = writer.WriteField("operations", body)

		mapData := map[string][]string{}
		for i, upload := range s.FilesToSend {
			mapData[fmt.Sprint(i)] = []string{fmt.Sprintf(`variables.%s`, upload.Key)}
		}
		jsonData, _ := json.Marshal(mapData)

		_ = writer.WriteField("map", string(jsonData))

		for i, upload := range s.FilesToSend {
			log.Debug().
				Str("key", upload.Key).
				Str("path", upload.Path).
				Msgf("Uploading file")
			file, err := os.Open(upload.Path)
			if err != nil {
				log.Err(err).
					Str("path", upload.Path).
					Msg("Failed to open file")
				return err
			}

			part, err := createFormFile(writer, fmt.Sprint(i), filepath.Base(upload.Path))
			if err != nil {
				log.Err(err).
					Str("path", upload.Path).
					Msg("Failed to create form file")
				return err
			}
			_, _ = io.Copy(part, file)
		}
		if err := writer.Close(); err != nil {
			return err
		}
		req.SetHeader("Content-Type", writer.FormDataContentType())
		req.SetBody(payload)
	} else {
		req.SetHeader("Content-Type", "application/json")
		req.SetBody(body)
	}

	log.Debug().Msgf("GraphQL request body: %s", body)

	if s.Token != "" {
		req.SetAuthToken(s.Token)
	}

	if s.SecretKey != "" {
		req.SetHeader(auth.SecretKeyHeader, s.SecretKey)
	}

	if s.APIKey != "" {
		req.SetHeader(auth.APIKeyHeader, s.APIKey)
	}

	method := "POST"
	serviceConfig := urlfinder.ServiceConfig{
		ServiceName: serviceName,
		Region:      viper.GetString("region"),
		Lifecycle:   viper.GetString("lifecycle"),
	}
	url := serviceConfig.FormatURL() + endpoint
	url = string(s.PerformSubstitutions([]byte(url)))

	log.Debug().Msgf("URL: %s: - Method %s", url, method)

	response, err := req.Post(url)
	if err != nil {
		return err
	}

	s.ResponseString = response.String()

	log.Debug().Msgf("RESPONSE: %v", response)
	doc, err := jsonquery.Parse(strings.NewReader(response.String()))
	if err != nil {
		return err
	}

	errorValue := jsonquery.FindOne(doc, "//errors")
	if errorValue != nil {
		errorBytes, _ := json.Marshal(errorValue.Value())
		s.Err = errors.New(string(errorBytes))
	} else {
		extractedValue := jsonquery.FindOne(doc, fmt.Sprintf("//data/%s", s.GraphQLOperation))
		if extractedValue != nil {
			responseBytes, _ := json.Marshal(extractedValue.Value())
			s.ResponseString = string(responseBytes)
		}
	}

	s.HTTPResponse = response.RawResponse
	s.StatusCode = response.StatusCode()

	s.FilesToSend = nil
	return nil
}

func createFormFile(w *multipart.Writer, fieldName, filename string) (io.Writer, error) {
	fileContentType := mime.TypeByExtension(filepath.Ext(filename))
	if fileContentType == "" {
		log.Debug().Msgf("File content type is empty, setting to application/octet-stream")
		fileContentType = "application/octet-stream"
	}
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, fieldName, filename))
	h.Set("Content-Type", fileContentType)
	return w.CreatePart(h)
}
