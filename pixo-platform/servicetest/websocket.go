package servicetest

import (
	"errors"
	"fmt"
	"github.com/cucumber/godog"
	"github.com/rs/zerolog/log"
	"strconv"
	"strings"
	"time"
)

func (s *ServerTestFeature) OpenWebsocket(endpoint string) error {
	log.Debug().Str("endpoint", endpoint).Msg("Opening websocket connection")

	s.WebsocketConn, s.httpResponse, s.Err = s.ServiceClient.DialWebsocket(endpoint)
	if s.Err != nil {
		log.Info().Err(s.Err).Msg("Error opening websocket")
	}

	if s.httpResponse != nil {
		log.Debug().Str("status", s.httpResponse.Status).Msg("Websocket connection opened")
		s.StatusCode = s.httpResponse.StatusCode
	}

	return nil
}

func (s *ServerTestFeature) WebsocketIsConnected() error {
	if s.WebsocketConn == nil {
		return errors.New("not connected to websocket")
	}

	return nil
}

func (s *ServerTestFeature) WebsocketIsNotConnected() error {
	if s.WebsocketConn != nil {
		return errors.New("websocket is connected")
	}

	return nil
}

func (s *ServerTestFeature) SetWebsocketReadTimeout(seconds string) error {
	timeout, err := strconv.Atoi(seconds)
	if err != nil {
		return err
	}

	if timeout <= 0 {
		return errors.New("timeout must be greater than 0")
	}

	timeoutDuration := time.Duration(timeout) * time.Second
	s.WebsocketReadTimeout = &timeoutDuration

	if err = s.WebsocketConn.SetReadDeadline(time.Now().Add(*s.WebsocketReadTimeout)); err != nil {
		return fmt.Errorf("error setting read deadline: %w", err)
	}

	return nil
}

func (s *ServerTestFeature) SendWebsocketMessage(body *godog.DocString) error {
	if err := s.WebsocketIsConnected(); err != nil {
		return err
	}

	msg := s.PerformSubstitutions([]byte(body.Content))
	if err := s.ServiceClient.WriteToWebsocket(msg); err != nil {
		return fmt.Errorf("error sending message to websocket: %w", err)
	}

	log.Debug().Str("message", string(msg)).Msg("Sent message to websocket")
	return nil
}

func (s *ServerTestFeature) GetWebsocketMessage() error {
	if err := s.WebsocketIsConnected(); err != nil {
		return err
	}

	s.Message = ""

	_, msg, err := s.ServiceClient.ReadFromWebsocket()
	if err != nil {
		return fmt.Errorf("error reading from websocket: %w", err)
	}

	s.Message = string(msg)
	log.Debug().Str("message", s.Message).Msg("Received message from websocket")
	return nil
}

func (s *ServerTestFeature) CheckMessageNotEmpty() error {
	if s.Message == "" {
		return errors.New("message is empty")
	}
	return nil
}

func (s *ServerTestFeature) TheMessageShouldContainA(key string) error {
	key = string(s.PerformSubstitutions([]byte(key)))
	actual := TrimString(s.Message)
	if !strings.Contains(actual, key) {
		return fmt.Errorf("expected message to contain %s, but got %s", key, actual)
	}

	return nil
}

func (s *ServerTestFeature) TheMessageShouldNotContainA(key string) error {
	key = string(s.PerformSubstitutions([]byte(key)))
	actual := TrimString(s.Message)
	if strings.Contains(actual, key) {
		return fmt.Errorf("expected message not to contain %s, but got %s", key, actual)
	}

	return nil
}
