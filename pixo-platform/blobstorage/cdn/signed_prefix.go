package cdn

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"strings"
	"time"
)

type SignedPrefixResult struct {
	QueryString string
	Expiration  time.Time
	Prefix      string
}

type Client struct {
	config Config
}

func NewClient(config Config) (*Client, error) {
	if config.KeyName == "" {
		return nil, fmt.Errorf("key name is required")
	}
	if config.KeyValue == "" {
		return nil, fmt.Errorf("key value is required")
	}
	if config.LoadBalancerDomain == "" {
		return nil, fmt.Errorf("load balancer domain is required")
	}
	if config.DefaultTTL == 0 {
		config.DefaultTTL = DefaultSignedPrefixTTL
	}

	config.LoadBalancerDomain = strings.TrimRight(config.LoadBalancerDomain, "/")

	return &Client{config: config}, nil
}

func (c *Client) GenerateSignedPrefix(moduleID string, ttl ...time.Duration) (*SignedPrefixResult, error) {
	if moduleID == "" {
		return nil, fmt.Errorf("module ID is required")
	}

	duration := c.config.DefaultTTL
	if len(ttl) > 0 && ttl[0] > 0 {
		duration = ttl[0]
	}

	expiration := time.Now().Add(duration)
	expirationUnix := expiration.Unix()

	prefix := fmt.Sprintf("%s/%s/", c.config.LoadBalancerDomain, moduleID)

	encodedPrefix := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString([]byte(prefix))

	policyString := fmt.Sprintf("URLPrefix=%s&Expires=%d&KeyName=%s", encodedPrefix, expirationUnix, c.config.KeyName)

	decodedKey, err := base64.URLEncoding.DecodeString(c.config.KeyValue)
	if err != nil {
		decodedKey = []byte(c.config.KeyValue)
	}

	mac := hmac.New(sha1.New, decodedKey)
	mac.Write([]byte(policyString))
	signature := base64.URLEncoding.EncodeToString(mac.Sum(nil))

	queryString := fmt.Sprintf("%s&Signature=%s", policyString, signature)

	return &SignedPrefixResult{
		QueryString: queryString,
		Expiration:  expiration,
		Prefix:      prefix,
	}, nil
}

func (c *Client) BuildAssetURL(moduleID, algorithm, hash, queryString string) string {
	return fmt.Sprintf("%s/%s/%s/%s?%s",
		c.config.LoadBalancerDomain,
		moduleID,
		algorithm,
		hash,
		queryString,
	)
}

func (c *Client) GetBucketName() string {
	return c.config.BucketName
}

func (c *Client) GetLoadBalancerDomain() string {
	return c.config.LoadBalancerDomain
}
