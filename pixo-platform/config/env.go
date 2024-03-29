package config

import (
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"os"
	"path/filepath"
	"strings"
)

func init() {
	LoadEnvVars()
}

func LoadEnvVars(differential ...string) {
	envPath := filepath.Join(GetProjectRoot(differential...), ".env")

	_ = godotenv.Load(envPath)
}

func GetEnvOrReturn(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func GetEnvOrCrash(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		log.Fatal().Msgf("Missing required environment variable: %s", key)
	}

	return value
}

func GetLifecycle() string {
	return strings.ToLower(GetEnvOrReturn("LIFECYCLE", "local"))
}

func GetDomain() string {
	return GetEnvOrReturn("DOMAIN", "localhost")
}

func GetRegion() string {
	region := strings.ToLower(GetEnvOrReturn("REGION", "us-central1"))
	if strings.Contains(region, "me-central") || strings.Contains(region, "saudi") {
		return "saudi"
	}

	return "us-central1"
}
