package service

import (
	"context"

	sm "github.com/wDRxxx/avito-shop/internal/service/models"
)

type Service interface {
	UserToken(ctx context.Context, username string, password string) (string, error)
	BuyItem(ctx context.Context, userID int, title string) error
	SendCoin(ctx context.Context, toUser string, fromUserID int, amount int) error
	UserInfo(ctx context.Context, userID int) (*sm.UserInfo, error)
}
