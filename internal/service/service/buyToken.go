package service

import (
	"context"
	"errors"

	"github.com/wDRxxx/avito-shop/internal/repository"
	"github.com/wDRxxx/avito-shop/internal/service"
)

func (s *serv) BuyItem(ctx context.Context, userID int, title string) error {
	item, err := s.repo.Item(ctx, title)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return service.ErrItemNotFound
		}

		return err
	}

	err = s.repo.BuyItem(ctx, userID, item)
	if err != nil {
		if errors.Is(err, repository.ErrNegativeBalance) {
			return service.ErrInsufficientBalance
		}

		return err
	}

	return nil
}
