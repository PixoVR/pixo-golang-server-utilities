package auth

import (
	"context"

	"github.com/PixoVR/pixo-golang-server-utilities/pixo-platform/config"
)

func GetUser(ctx context.Context) *User {
	user, ok := ctx.Value(config.ContextRequestAuthentication.String()).(*User)
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
