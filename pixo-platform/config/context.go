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
	ginContext, ok := ctx.Value(ContextRequestGin.String()).(*gin.Context)
	if !ok {
		return nil
	}

	return ginContext
}

func GetIPAddress(ctx context.Context) string {
	ipAddress, ok := ctx.Value(ContextRequestHost.String()).(string)
	if !ok {
		return ""
	}

	return ipAddress
}

func GetCurrentUserID(userContext context.Context) int {
	user, ok := userContext.Value(ContextRequestAuthentication.String()).(*platform.User)
	if !ok {
		return 0
	}

	return user.ID
}

func GetAuthorizationEnforcer(userContext context.Context) *casbin.Enforcer {
	enforcer, ok := userContext.Value(ContextRequestAuthorization.String()).(*casbin.Enforcer)
	if !ok {
		return nil
	}

	return enforcer
}
