package auth

import (
	"context"

	platform "github.com/PixoVR/pixo-golang-clients/pixo-platform/primary-api"
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/config"
)

func GetUser(ctx context.Context) *platform.User {
	user, ok := ctx.Value(config.ContextRequestAuthentication.String()).(*platform.User)
	if !ok {
		return nil
	}

	return user
}

func GetCurrentUserID(ctx context.Context) int {
	user := GetUser(ctx)
	if user == nil {
		return 0
	}

	return user.ID
}
