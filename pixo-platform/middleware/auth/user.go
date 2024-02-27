package auth

import (
	"context"
	platform "github.com/PixoVR/pixo-golang-clients/pixo-platform/primary-api"
	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/config"
)

func GetUser(ctx context.Context) *platform.User {
	user, ok := ctx.Value(config.ContextRequestAuthentication).(*platform.User)
	if !ok {
		return nil
	}

	return user
}
