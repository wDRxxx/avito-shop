package service

import (
	"context"
)

type UsersService interface {
	UserToken(ctx context.Context, username string, password string) (string, error)
}
