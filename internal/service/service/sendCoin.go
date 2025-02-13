package service

import (
	"context"
	"errors"

	"github.com/wDRxxx/avito-shop/internal/repository"
	"github.com/wDRxxx/avito-shop/internal/service"
)

func (s *serv) SendCoin(ctx context.Context, toUser string, fromUserID int, amount int) error {
	user, err := s.repo.UserByUsername(ctx, toUser)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return service.ErrUserNotFound
		}

		return err
	}

	err = s.repo.SendCoin(ctx, user.ID, fromUserID, amount)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return service.ErrItemNotFound
		}
		if errors.Is(err, repository.ErrNegativeBalance) {
			return service.ErrInsufficientBalance
		}

		return err
	}

	return nil
}
