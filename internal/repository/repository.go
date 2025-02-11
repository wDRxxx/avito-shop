package repository

import (
	"context"

	rm "github.com/wDRxxx/avito-shop/internal/repository/models"
)

type Repository interface {
	User(ctx context.Context, username string) (*rm.User, error)
	InsertUser(ctx context.Context, username string, password string) (*rm.User, error)
}
