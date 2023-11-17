package config

import (
	"context"
	"github.com/gin-gonic/gin"
)

const (
	GinContextKey           = "GIN_CONTEXT"
	ContextRequestUser      = ContextRequest("user")
	ContextRequestIPAddress = ContextRequest("ipAddress")
)

type ContextRequest string

func (c ContextRequest) String() string {
	return string(c)
}

func GetGinContext(ctx context.Context) *gin.Context {
	ginContext, ok := ctx.Value(GinContextKey).(*gin.Context)
	if !ok {
		return nil
	}

	return ginContext
}

type User struct {
	ID int
}

func GetCurrentUserID(userContext context.Context) int {
	user, ok := userContext.Value(ContextRequestUser).(*User)
	if !ok {
		return 0
	}

	return user.ID
}

func GetIPAddress(userContext context.Context) string {
	ipAddress, ok := userContext.Value(ContextRequestIPAddress).(string)
	if !ok {
		return "127.0.0.1"
	}

	return ipAddress
}
