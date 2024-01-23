package config

import (
	"context"
	platform "github.com/PixoVR/pixo-golang-clients/pixo-platform/primary-api"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
)

const (
	ginContextKey            = "GIN_CONTEXT"
	authorizationContextKey  = "ENFORCER_CONTEXT"
	authenticationContextKey = "USER_CONTEXT"
	ipAddressContextKey      = "IP_ADDRESS_CONTEXT"
	customContextKey         = "CUSTOM_CONTEXT"

	ContextRequestGin            = ContextRequest(ginContextKey)
	ContextRequestAuthorization  = ContextRequest(authorizationContextKey)
	ContextRequestAuthentication = ContextRequest(authenticationContextKey)
	ContextRequestHost           = ContextRequest(ipAddressContextKey)
	ContextRequestCustom         = ContextRequest(customContextKey)
)

type ContextRequest string

func (c ContextRequest) String() string {
	return string(c)
}

func GetGinContext(ctx context.Context) *gin.Context {
	ginContext, ok := ctx.Value(ContextRequestGin).(*gin.Context)
	if !ok {
		return nil
	}

	return ginContext
}

func GetIPAddress(userContext context.Context) string {
	ipAddress, ok := userContext.Value(ContextRequestHost).(string)
	if !ok {
		return "127.0.0.1"
	}

	return ipAddress
}

func GetCurrentUserID(userContext context.Context) int {
	user, ok := userContext.Value(ContextRequestAuthentication).(*platform.User)
	if !ok {
		return 0
	}

	return user.ID
}

func GetAuthorizationEnforcer(userContext context.Context) *casbin.Enforcer {
	enforcer, ok := userContext.Value(ContextRequestAuthorization).(*casbin.Enforcer)
	if !ok {
		return nil
	}

	return enforcer
}
