package config

import (
	"context"
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

func GetIPAddress(ctx context.Context) string {
	ipAddress, ok := ctx.Value(ContextRequestHost.String()).(string)
	if !ok {
		return ""
	}

	return ipAddress
}

func GetCurrentUserID(ctx context.Context) int {
	user, ok := ctx.Value(ContextRequestAuthentication).(*User)
	if !ok {
		return 0
	}

	return user.ID
}

func GetAuthorizationEnforcer(ctx context.Context) *casbin.Enforcer {
	enforcer, ok := ctx.Value(ContextRequestAuthorization).(*casbin.Enforcer)
	if !ok {
		return nil
	}

	return enforcer
}
