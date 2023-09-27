package config

import (
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"os"
	"path/filepath"
)

func init() {
	LoadEnvVars()
}

func LoadEnvVars(differential ...string) {
	envPath := filepath.Join(GetProjectRoot(differential...), ".env")

	if err := godotenv.Load(envPath); err != nil {
		log.Warn().Msgf("No .env file loaded: %s", err)
	}
}

func GetEnvOrReturn(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func GetLifecycle() string {
	return GetEnvOrReturn("LIFECYCLE", "dev")
}

func GetDomain() string {
	return GetEnvOrReturn("DOMAIN", "localhost")
}
