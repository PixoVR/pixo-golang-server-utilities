package servicetest

import (
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"net/http"
	"net/http/httptest"
)

type ServerTestFeature struct {
	Engine   *gin.Engine
	Recorder *httptest.ResponseRecorder

	Client         *resty.Client
	Response       *http.Response
	ResponseString string
	StatusCode     int

	Token  string
	ID     string
	UserID int
	APIKey string
	Err    error

	DirectoryFilePath string
	GraphQLOperation  string
	SendFileKey       string
	SendFile          string
	GraphQLResponse   map[string]interface{}
}

func (s *ServerTestFeature) Reset() {
	s.Recorder = httptest.NewRecorder()

	s.Client = resty.New()
	s.Response = nil
	s.ResponseString = ""
	s.StatusCode = 0

	s.Token = ""
	s.ID = ""
	s.UserID = 0
	s.APIKey = ""
	s.Err = nil

	s.GraphQLResponse = make(map[string]interface{})
	s.DirectoryFilePath = ""
	s.SendFileKey = ""
	s.SendFileKey = ""
}
