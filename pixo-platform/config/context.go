package config

import (
	"context"
	"github.com/gin-gonic/gin"
)

const (
	GinContextKey            = "GIN_CONTEXT"
	AuthorizationContextKey  = "ENFORCER_CONTEXT"
	AuthenticationContextKey = "USER_CONTEXT"
	IPAddressContextKey      = "IP_ADDRESS_CONTEXT"

	ContextRequestGin            = ContextRequest(GinContextKey)
	ContextRequestAuthorization  = ContextRequest(AuthorizationContextKey)
	ContextRequestAuthentication = ContextRequest(AuthenticationContextKey)
	ContextRequestIPAddress      = ContextRequest(IPAddressContextKey)
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

type User struct {
	ID int
}

func GetCurrentUserID(userContext context.Context) int {
	user, ok := userContext.Value(ContextRequestAuthentication).(*User)
	if !ok {
		return 0
	}

	return user.ID
}

func GetAuthorization(userContext context.Context) string {
	authorization, ok := userContext.Value(ContextRequestAuthorization).(string)
	if !ok {
		return ""
	}

	return authorization
}

func GetIPAddress(userContext context.Context) string {
	ipAddress, ok := userContext.Value(ContextRequestIPAddress).(string)
	if !ok {
		return "127.0.0.1"
	}

	return ipAddress
}
