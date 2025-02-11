package service

import (
	"context"
)

type Service interface {
	UserToken(ctx context.Context, username string, password string) (string, error)
	BuyItem(ctx context.Context, userID int, title string) error
}
