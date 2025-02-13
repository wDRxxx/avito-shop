package service

import (
	"context"
	"errors"

	"github.com/wDRxxx/avito-shop/internal/repository"
	"github.com/wDRxxx/avito-shop/internal/service"
	"github.com/wDRxxx/avito-shop/internal/service/converter"
	sm "github.com/wDRxxx/avito-shop/internal/service/models"
)

func (s *serv) UserInfo(ctx context.Context, userID int) (*sm.UserInfo, error) {
	user, err := s.repo.UserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, service.ErrUserNotFound
		}
		return nil, err
	}

	repoInventory, err := s.repo.UserInventory(ctx, userID)
	if err != nil {
		return nil, err
	}
	inventoryItems := converter.InventoryFromRepositoryToService(repoInventory)

	repoIncomingTransactions, err := s.repo.UserIncomingTransactions(ctx, userID)
	if err != nil {
		return nil, err
	}
	incomingTransactions := converter.IncomingTransactionsFromRepositoryToService(repoIncomingTransactions)

	repoOutgoingTransactions, err := s.repo.UserOutgoingTransactions(ctx, userID)
	if err != nil {
		return nil, err
	}
	outgoingTransactions := converter.OutgoingTransactionsFromRepositoryToService(repoOutgoingTransactions)

	return &sm.UserInfo{
		Balance:              user.Balance,
		InventoryItems:       inventoryItems,
		IncomingTransactions: incomingTransactions,
		OutgoingTransactions: outgoingTransactions,
	}, nil
}
