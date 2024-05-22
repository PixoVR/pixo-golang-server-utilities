package servicetest

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"math/rand"
	"os"
	"strings"
)

func (s *ServerTestFeature) UseSecretKey() error {
	envKey := "SECRET_KEY"

	currentLifecycle := viper.GetString("lifecycle")
	if currentLifecycle != "" && currentLifecycle != "local" && currentLifecycle != "internal" {
		envKey = strings.ToUpper(fmt.Sprintf("%s_%s", viper.GetString("lifecycle"), "SECRET_KEY"))
	}

	s.SecretKey = os.Getenv(envKey)
	if s.SecretKey == "" {
		return errors.New("secret key not found")
	}

	s.ServiceClient.SetToken(s.SecretKey)

	return nil
}

func (s *ServerTestFeature) CreateRandomInt() {
	s.RandomInt = rand.Int()
}

func (s *ServerTestFeature) SignedInAsA(role string) error {
	username := getUsername(role)
	if username == "" {
		return errors.New("username not found")
	}

	password := getPassword(role)
	if password == "" {
		return errors.New("password not found")
	}

	if s.PlatformClient == nil {
		return errors.New("platform client not initialized")
	}

	if err := s.PlatformClient.Login(username, password); err != nil {
		return fmt.Errorf("failed to login: %w", err)
	}

	s.Token = s.PlatformClient.GetToken()
	if s.Token == "" {
		return errors.New("token not found")
	}

	s.ServiceClient.SetToken(s.Token)

	s.UserID = s.PlatformClient.ActiveUserID()
	if s.UserID == 0 {
		return errors.New("user id not found")
	}

	return nil
}

func getUsername(role string) string {
	return os.Getenv(strings.ToUpper(role) + "_USERNAME")
}

func getPassword(role string) string {
	return os.Getenv(strings.ToUpper(role) + "_PASSWORD")
}
