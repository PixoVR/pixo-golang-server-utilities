package servicetest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/urlfinder"
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/middleware/auth"
	"github.com/cucumber/godog"
	"github.com/go-resty/resty/v2"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"strings"
)

func (s *ServerTestFeature) MakeRequest(method string, endpoint string, body *godog.DocString, paramsMap map[string]string) error {
	var bodyContent []byte
	if body != nil {
		bodyContent = []byte(body.Content)
	}

	bodyContent = s.PerformSubstitutions(bodyContent)

	if err := s.PerformRequest(method, endpoint, bodyContent, paramsMap); err != nil {
		return err
	}

	var responseBody map[string]interface{}
	if err := json.Unmarshal([]byte(s.Recorder.Body.String()), &responseBody); err != nil {
		return fmt.Errorf("failed to unmarshal response body: %s", err)
	}

	if value, ok := responseBody["message"]; ok {
		s.Message = value.(string)
	}

	return nil
}

func (s *ServerTestFeature) PerformRequest(method, endpoint string, body []byte, paramsMap map[string]string) error {

	body = s.PerformSubstitutions(body)

	if s.BeforeRequest != nil {
		s.BeforeRequest(body)
	}

	currentLifecycle := strings.ToLower(viper.GetString("lifecycle"))

	if currentLifecycle == "" || currentLifecycle == "local" {
		url := fmt.Sprintf("/%s%s", s.ServiceClient.Path(), endpoint)
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

		if method == "POST" {
			req.Header.Set("Content-Type", "application/json")
			req.Body = io.NopCloser(bytes.NewReader(body))
		}

		if paramsMap != nil {
			for key, value := range paramsMap {
				req.URL.Query().Add(key, value)
			}
		}

		s.Engine.ServeHTTP(s.Recorder, req)

		s.Response = s.Recorder.Result()
		s.ResponseString = s.Recorder.Body.String()
		s.StatusCode = s.Recorder.Code
		s.Err = nil
	} else {
		req := s.Client.R()

		if s.Token != "" {
			req.Header.Set(auth.AuthorizationHeader, fmt.Sprintf("Bearer %s", s.Token))
		}

		if s.SecretKey != "" {
			req.Header.Set(auth.SecretKeyHeader, s.SecretKey)
		}

		if paramsMap != nil {
			req.SetQueryParams(paramsMap)
		}

		url := fmt.Sprintf("%s%s", s.ServiceClient.GetURL(), endpoint)

		var res *resty.Response
		var err error
		switch method {
		case "GET":
			res, err = req.Get(url)
		case "POST":
			req.SetBody(body)
			res, err = req.Post(url)
		}

		Expect(err).NotTo(HaveOccurred())
		Expect(res).NotTo(BeNil())

		log.Debug().
			Str("endpoint", endpoint).
			Str("response", string(res.Body())).
			Int("status_code", res.StatusCode()).
			Msg("HTTP RESPONSE")

		s.Response = res.RawResponse
		s.ResponseString = string(res.Body())
		s.StatusCode = res.StatusCode()
		s.Err = err
	}

	return nil
}

func (s *ServerTestFeature) makeGraphQLRequest(endpoint, serviceName, body string) error {
	req := s.Client.R()

	body = string(s.PerformSubstitutions([]byte(body)))

	if s.SendFileKey != "" && s.SendFile != "" {
		req.FormData.Add("operations", body)
		req.FormData.Add("map", fmt.Sprintf(`{"0": ["variables.%s"]}`, s.SendFileKey))
		req.SetFiles(map[string]string{"0": s.SendFile})
	} else {
		req.SetHeader("Content-Type", "application/json").
			SetBody(body)
	}

	if s.Token != "" {
		req.SetAuthToken(s.Token)
	}

	method := "POST"
	serviceConfig := urlfinder.ServiceConfig{
		ServiceName: serviceName,
		Region:      viper.GetString("region"),
		Lifecycle:   viper.GetString("lifecycle"),
	}
	url := serviceConfig.FormatURL() + endpoint
	url = string(s.PerformSubstitutions([]byte(url)))

	log.Info().Msgf("URL: %s: - Method %s", url, method)

	response, err := req.Post(url)
	if err != nil {
		return err
	}

	log.Debug().Msgf("RESPONSE: %v", response)
	s.Response = response.RawResponse
	s.ResponseString = response.String()
	s.StatusCode = response.StatusCode()

	// reset so it is not used automatically for the next request
	s.SendFileKey = ""
	s.SendFileKey = ""

	return nil
}
