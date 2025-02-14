package service

import (
	"context"
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	"github.com/wDRxxx/avito-shop/internal/config"
	"github.com/wDRxxx/avito-shop/internal/repository"
	repoMocks "github.com/wDRxxx/avito-shop/internal/repository/mocks"
	rm "github.com/wDRxxx/avito-shop/internal/repository/models"
	"github.com/wDRxxx/avito-shop/internal/service"
	"github.com/wDRxxx/avito-shop/internal/service/converter"
	sm "github.com/wDRxxx/avito-shop/internal/service/models"
)

func TestUserInfo(t *testing.T) {
	t.Parallel()

	type repositoryMockFunc func(mc *minimock.Controller) repository.Repository

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		userID      = int(gofakeit.Uint8())
		password    = gofakeit.Password(true, true, true, false, false, 12)
		hashPass, _ = bcrypt.GenerateFromPassword([]byte(password), 12)
		balance     = int(gofakeit.Uint8())
		user        = &rm.User{
			ID:        userID,
			Username:  gofakeit.Username(),
			Password:  string(hashPass),
			Balance:   balance,
			CreatedAt: gofakeit.Date(),
			UpdatedAt: gofakeit.Date(),
		}

		itemID         = int(gofakeit.Uint8())
		inventoryItems = []*rm.InventoryItem{
			{
				UserID: userID,
				ItemID: itemID,
				Item: rm.Item{
					ID:    itemID,
					Title: gofakeit.BeerName(),
					Price: int(gofakeit.Uint8()),
				},
			},
		}

		incomingTransactions = []*rm.Transaction{
			{
				Amount: int(gofakeit.Uint8()),
				Sender: rm.User{
					Username: gofakeit.Username(),
				},
			},
		}
		outgoingTransactions = []*rm.Transaction{
			{
				Amount: int(gofakeit.Uint8()),
				Recipient: rm.User{
					Username: gofakeit.Username(),
				},
			},
		}

		userInfo = &sm.UserInfo{
			Balance:              balance,
			InventoryItems:       converter.InventoryFromRepositoryToService(inventoryItems),
			IncomingTransactions: converter.IncomingTransactionsFromRepositoryToService(incomingTransactions),
			OutgoingTransactions: converter.OutgoingTransactionsFromRepositoryToService(outgoingTransactions),
		}

		repoErr = errors.New("repo err")
	)

	type args struct {
		ctx    context.Context
		userID int
	}

	tests := []struct {
		name           string
		want           *sm.UserInfo
		err            error
		args           args
		repositoryMock repositoryMockFunc
	}{
		{
			name: "success case",
			err:  nil,
			want: userInfo,
			args: args{
				ctx:    ctx,
				userID: userID,
			},
			repositoryMock: func(mc *minimock.Controller) repository.Repository {
				mock := repoMocks.NewRepositoryMock(mc)
				mock.UserByIDMock.Expect(ctx, userID).Return(user, nil)
				mock.UserInventoryMock.Expect(ctx, userID).Return(inventoryItems, nil)
				mock.UserIncomingTransactionsMock.Expect(ctx, userID).Return(incomingTransactions, nil)
				mock.UserOutgoingTransactionsMock.Expect(ctx, userID).Return(outgoingTransactions, nil)
				return mock
			},
		},
		{
			name: "UserByIDMock repo error case",
			err:  repoErr,
			want: nil,
			args: args{
				ctx:    ctx,
				userID: userID,
			},
			repositoryMock: func(mc *minimock.Controller) repository.Repository {
				mock := repoMocks.NewRepositoryMock(mc)
				mock.UserByIDMock.Expect(ctx, userID).Return(nil, repoErr)
				return mock
			},
		},
		{
			name: "UserInventoryMock repo error case",
			err:  repoErr,
			want: nil,
			args: args{
				ctx:    ctx,
				userID: userID,
			},
			repositoryMock: func(mc *minimock.Controller) repository.Repository {
				mock := repoMocks.NewRepositoryMock(mc)
				mock.UserByIDMock.Expect(ctx, userID).Return(user, nil)
				mock.UserInventoryMock.Expect(ctx, userID).Return(nil, repoErr)
				return mock
			},
		},
		{
			name: "UserIncomingTransactionsMock repo error case",
			err:  repoErr,
			want: nil,
			args: args{
				ctx:    ctx,
				userID: userID,
			},
			repositoryMock: func(mc *minimock.Controller) repository.Repository {
				mock := repoMocks.NewRepositoryMock(mc)
				mock.UserByIDMock.Expect(ctx, userID).Return(user, nil)
				mock.UserInventoryMock.Expect(ctx, userID).Return(inventoryItems, nil)
				mock.UserIncomingTransactionsMock.Expect(ctx, userID).Return(nil, repoErr)
				return mock
			},
		},
		{
			name: "UserOutgoingTransactionsMock repo error case",
			err:  repoErr,
			want: nil,
			args: args{
				ctx:    ctx,
				userID: userID,
			},
			repositoryMock: func(mc *minimock.Controller) repository.Repository {
				mock := repoMocks.NewRepositoryMock(mc)
				mock.UserByIDMock.Expect(ctx, userID).Return(user, nil)
				mock.UserInventoryMock.Expect(ctx, userID).Return(inventoryItems, nil)
				mock.UserIncomingTransactionsMock.Expect(ctx, userID).Return(incomingTransactions, nil)
				mock.UserOutgoingTransactionsMock.Expect(ctx, userID).Return(nil, repoErr)
				return mock
			},
		},
		{
			name: "user not found error case",
			err:  service.ErrUserNotFound,
			want: nil,
			args: args{
				ctx:    ctx,
				userID: userID,
			},
			repositoryMock: func(mc *minimock.Controller) repository.Repository {
				mock := repoMocks.NewRepositoryMock(mc)
				mock.UserByIDMock.Expect(ctx, userID).Return(user, repository.ErrNotFound)
				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repoMock := tt.repositoryMock(mc)
			service := NewMockService(repoMock, config.NewMockAuthConfig())

			info, err := service.UserInfo(tt.args.ctx, tt.args.userID)
			require.Equal(t, tt.want, info)
			require.Equal(t, tt.err, err)
		})
	}
}
