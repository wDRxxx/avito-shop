package repository

import (
	"context"

	rm "github.com/wDRxxx/avito-shop/internal/repository/models"
)

type Repository interface {
	UserByUsername(ctx context.Context, username string) (*rm.User, error)
	UserByID(ctx context.Context, userID int) (*rm.User, error)
	InsertUser(ctx context.Context, username string, password string) (*rm.User, error)
	Item(ctx context.Context, title string) (*rm.Item, error)
	BuyItem(ctx context.Context, userID int, item *rm.Item) error
	SendCoin(ctx context.Context, toUserID int, fromUserID int, amount int) error
	UserInventory(ctx context.Context, userID int) ([]*rm.InventoryItem, error)
	UserIncomingTransactions(ctx context.Context, userID int) ([]*rm.Transaction, error)
	UserOutgoingTransactions(ctx context.Context, userID int) ([]*rm.Transaction, error)
}
