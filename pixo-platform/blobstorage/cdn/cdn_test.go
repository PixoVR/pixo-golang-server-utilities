package cdn_test

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/blobstorage/cdn"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("CDN Signed Prefix", func() {

	var (
		client   *cdn.Client
		keyName  = "test-key"
		keyValue = base64.URLEncoding.EncodeToString([]byte("test-secret-key1"))
		domain   = "https://cdn.example.com"
	)

	BeforeEach(func() {
		var err error
		client, err = cdn.NewClient(cdn.Config{
			KeyName:            keyName,
			KeyValue:           keyValue,
			LoadBalancerDomain: domain,
			BucketName:         "test-cas-bucket",
			DefaultTTL:         12 * time.Hour,
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(client).NotTo(BeNil())
	})

	Describe("NewClient", func() {
		It("should return an error if key name is empty", func() {
			_, err := cdn.NewClient(cdn.Config{
				KeyValue:           keyValue,
				LoadBalancerDomain: domain,
			})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("key name"))
		})

		It("should return an error if key value is empty", func() {
			_, err := cdn.NewClient(cdn.Config{
				KeyName:            keyName,
				LoadBalancerDomain: domain,
			})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("key value"))
		})

		It("should return an error if load balancer domain is empty", func() {
			_, err := cdn.NewClient(cdn.Config{
				KeyName:  keyName,
				KeyValue: keyValue,
			})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("load balancer domain"))
		})

		It("should set the default TTL if not provided", func() {
			c, err := cdn.NewClient(cdn.Config{
				KeyName:            keyName,
				KeyValue:           keyValue,
				LoadBalancerDomain: domain,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(c).NotTo(BeNil())
		})

		It("should strip trailing slash from domain", func() {
			c, err := cdn.NewClient(cdn.Config{
				KeyName:            keyName,
				KeyValue:           keyValue,
				LoadBalancerDomain: "https://cdn.example.com/",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(c.GetLoadBalancerDomain()).To(Equal("https://cdn.example.com"))
		})
	})

	Describe("GenerateSignedPrefix", func() {
		It("should return an error if module ID is empty", func() {
			result, err := client.GenerateSignedPrefix("")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("module ID"))
			Expect(result).To(BeNil())
		})

		It("should generate a valid signed prefix for a module", func() {
			moduleID := "42"
			result, err := client.GenerateSignedPrefix(moduleID)

			Expect(err).NotTo(HaveOccurred())
			Expect(result).NotTo(BeNil())
			Expect(result.QueryString).To(ContainSubstring("URLPrefix="))
			Expect(result.QueryString).To(ContainSubstring("Expires="))
			Expect(result.QueryString).To(ContainSubstring(fmt.Sprintf("KeyName=%s", keyName)))
			Expect(result.QueryString).To(ContainSubstring("Signature="))
			Expect(result.Prefix).To(Equal(fmt.Sprintf("%s/%s/", domain, moduleID)))
			Expect(result.Expiration).To(BeTemporally("~", time.Now().Add(12*time.Hour), 5*time.Second))
		})

		It("should use a custom TTL when provided", func() {
			customTTL := 1 * time.Hour
			result, err := client.GenerateSignedPrefix("42", customTTL)

			Expect(err).NotTo(HaveOccurred())
			Expect(result).NotTo(BeNil())
			Expect(result.Expiration).To(BeTemporally("~", time.Now().Add(customTTL), 5*time.Second))
		})

		It("should produce a verifiable HMAC-SHA1 signature", func() {
			moduleID := "42"
			result, err := client.GenerateSignedPrefix(moduleID)
			Expect(err).NotTo(HaveOccurred())

			parts := strings.Split(result.QueryString, "&Signature=")
			Expect(parts).To(HaveLen(2))

			policyString := parts[0]
			signatureB64 := parts[1]

			decodedKey, err := base64.URLEncoding.DecodeString(keyValue)
			Expect(err).NotTo(HaveOccurred())

			mac := hmac.New(sha1.New, decodedKey)
			mac.Write([]byte(policyString))
			expectedSignature := base64.URLEncoding.EncodeToString(mac.Sum(nil))

			Expect(signatureB64).To(Equal(expectedSignature))
		})

		It("should encode the URL prefix as URL-safe base64 without padding", func() {
			moduleID := "42"
			result, err := client.GenerateSignedPrefix(moduleID)
			Expect(err).NotTo(HaveOccurred())

			prefix := fmt.Sprintf("%s/%s/", domain, moduleID)
			expectedEncoded := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString([]byte(prefix))

			Expect(result.QueryString).To(ContainSubstring(fmt.Sprintf("URLPrefix=%s", expectedEncoded)))
		})

		It("should produce different signatures for different modules", func() {
			result1, err := client.GenerateSignedPrefix("1")
			Expect(err).NotTo(HaveOccurred())

			result2, err := client.GenerateSignedPrefix("2")
			Expect(err).NotTo(HaveOccurred())

			sig1 := strings.Split(result1.QueryString, "&Signature=")[1]
			sig2 := strings.Split(result2.QueryString, "&Signature=")[1]
			Expect(sig1).NotTo(Equal(sig2))
		})
	})

	Describe("BuildAssetURL", func() {
		It("should construct a valid asset URL with the query string", func() {
			queryString := "URLPrefix=abc&Expires=123&KeyName=test-key&Signature=xyz"
			url := client.BuildAssetURL("42", "sha256", "a1b2c3d4", queryString)

			Expect(url).To(Equal("https://cdn.example.com/42/sha256/a1b2c3d4?URLPrefix=abc&Expires=123&KeyName=test-key&Signature=xyz"))
		})
	})

	Describe("GetBucketName", func() {
		It("should return the configured bucket name", func() {
			Expect(client.GetBucketName()).To(Equal("test-cas-bucket"))
		})
	})

	Describe("ConfigFromEnv", func() {
		It("should return an error if CDN_KEY_NAME is not set", func() {
			_, err := cdn.ConfigFromEnv()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring(cdn.EnvCDNKeyName))
		})
	})
})
