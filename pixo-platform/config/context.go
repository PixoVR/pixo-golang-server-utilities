package config

import (
	"context"
	"github.com/gin-gonic/gin"
)

const (
	GinContextKey = "GIN_CONTEXT"
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
