package cdn

import (
	"fmt"
	"os"
	"time"
)

const (
	DefaultSignedPrefixTTL = 12 * time.Hour

	EnvCDNKeyName    = "CDN_KEY_NAME"
	EnvCDNKeyValue   = "CDN_KEY_VALUE"
	EnvCDNDomain     = "CDN_LOAD_BALANCER_DOMAIN"
	EnvCASBucketName = "CAS_BUCKET"
)

type Config struct {
	KeyName            string
	KeyValue           string
	LoadBalancerDomain string
	BucketName         string
	DefaultTTL         time.Duration
}

func ConfigFromEnv() (Config, error) {
	cfg := Config{
		KeyName:            os.Getenv(EnvCDNKeyName),
		KeyValue:           os.Getenv(EnvCDNKeyValue),
		LoadBalancerDomain: os.Getenv(EnvCDNDomain),
		BucketName:         os.Getenv(EnvCASBucketName),
		DefaultTTL:         DefaultSignedPrefixTTL,
	}

	if cfg.KeyName == "" {
		return cfg, fmt.Errorf("%s is required", EnvCDNKeyName)
	}
	if cfg.KeyValue == "" {
		return cfg, fmt.Errorf("%s is required", EnvCDNKeyValue)
	}
	if cfg.LoadBalancerDomain == "" {
		return cfg, fmt.Errorf("%s is required", EnvCDNDomain)
	}

	return cfg, nil
}
