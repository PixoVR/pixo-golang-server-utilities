package servicetest

import (
	"encoding/json"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/abstract"
	"github.com/PixoVR/pixo-golang-clients/pixo-platform/platform"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"github.com/gorilla/websocket"
	"net"
	"net/http"
	"net/http/httptest"
	"time"
)

type ServerTestFeature struct {
	staticSubstitutions  map[string]string
	dynamicSubstitutions map[string]SubstitutionFunc

	BeforeRequest func(body []byte)
	Engine        *gin.Engine
	Recorder      *httptest.ResponseRecorder
	Client        *resty.Client

	UDPConn              *net.UDPConn
	WebsocketConn        *websocket.Conn
	WebsocketReadTimeout *time.Duration

	ServiceClient  abstract.AbstractClient
	PlatformClient platform.Client

	HTTPResponse   *http.Response
	ResponseString string
	StatusCode     int

	RandomInt int
	UUID      string

	SecretKey string
	Token     string
	ID        string
	UserID    int
	APIKey    string
	Message   string
	Err       error

	DirectoryFilePath string
	GraphQLOperation  string
	FilesToSend       []FileToSend
	GraphQLResponse   map[string]interface{}
}

type FileToSend struct {
	Key  string
	Path string
}

type SubstitutionFunc func(data []byte) string

func NewServerTestFeature() *ServerTestFeature {
	return &ServerTestFeature{
		Client:               resty.New(),
		staticSubstitutions:  make(map[string]string),
		dynamicSubstitutions: make(map[string]SubstitutionFunc),
	}
}

func (s *ServerTestFeature) AddStaticSubstitution(key, value string) {
	s.staticSubstitutions[key] = value
}

func (s *ServerTestFeature) AddDynamicSubstitution(key string, value SubstitutionFunc) {
	s.dynamicSubstitutions[key] = value
}

// UnwrapResponseString returns the inner data from a re-wrapped GraphQL response.
// After MakeGraphQLRequest, ResponseString is formatted as {"operationName": <data>}
// for XPath query compatibility. This method strips the operation wrapper and returns
// the raw inner data for use in functions that do direct JSON unmarshalling.
func (s *ServerTestFeature) UnwrapResponseString() string {
	if s.GraphQLOperation == "" {
		return s.ResponseString
	}
	var obj map[string]json.RawMessage
	if err := json.Unmarshal([]byte(s.ResponseString), &obj); err != nil {
		return s.ResponseString
	}
	if inner, ok := obj[s.GraphQLOperation]; ok && len(obj) == 1 {
		return string(inner)
	}
	return s.ResponseString
}

func (s *ServerTestFeature) Reset(interface{}) {
	s.Recorder = httptest.NewRecorder()

	if s.ServiceClient != nil {
		s.ServiceClient.SetToken("")
		s.ServiceClient.SetAPIKey("")
	}

	s.HTTPResponse = nil
	s.ResponseString = ""
	s.StatusCode = 0

	s.Token = ""
	s.SecretKey = ""
	s.APIKey = ""

	s.ID = ""
	s.UserID = 0
	s.Err = nil

	s.GraphQLOperation = ""
	s.GraphQLResponse = make(map[string]interface{})
	s.DirectoryFilePath = ""
	s.FilesToSend = nil

	s.UUID = ""
	s.RandomInt = 0
}
