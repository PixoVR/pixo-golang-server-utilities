package auth

import (
	"context"
	platform "github.com/PixoVR/pixo-golang-clients/pixo-platform/primary-api"
)

func GetUser(ctx context.Context) *platform.User {
	user, ok := ctx.Value(UserKey).(*platform.User)
	if !ok {
		return nil
	}

	return user
}
