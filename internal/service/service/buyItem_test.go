package service

import (
	"context"
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/require"

	"github.com/wDRxxx/avito-shop/internal/config"
	"github.com/wDRxxx/avito-shop/internal/repository"
	repoMocks "github.com/wDRxxx/avito-shop/internal/repository/mocks"
	rm "github.com/wDRxxx/avito-shop/internal/repository/models"
	"github.com/wDRxxx/avito-shop/internal/service"
)

func TestBuyItem(t *testing.T) {
	t.Parallel()

	type repositoryMockFunc func(mc *minimock.Controller) repository.Repository

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		userID = int(gofakeit.Uint8())
		title  = gofakeit.BeerName()

		item = &rm.Item{
			ID:    int(gofakeit.Uint()),
			Title: title,
			Price: int(gofakeit.Uint8()),
		}

		repoErr = errors.New("repo err")
	)

	type args struct {
		ctx    context.Context
		userID int
		title  string
	}

	tests := []struct {
		name           string
		err            error
		args           args
		repositoryMock repositoryMockFunc
	}{
		{
			name: "success case when user exists",
			err:  nil,
			args: args{
				ctx:    ctx,
				userID: userID,
				title:  title,
			},
			repositoryMock: func(mc *minimock.Controller) repository.Repository {
				mock := repoMocks.NewRepositoryMock(mc)
				mock.ItemMock.Expect(ctx, title).Return(item, nil)
				mock.BuyItemMock.Expect(ctx, userID, item).Return(nil)
				return mock
			},
		},
		{
			name: "ItemMock repo error case",
			err:  repoErr,
			args: args{
				ctx:    ctx,
				userID: userID,
				title:  title,
			},
			repositoryMock: func(mc *minimock.Controller) repository.Repository {
				mock := repoMocks.NewRepositoryMock(mc)
				mock.ItemMock.Expect(ctx, title).Return(nil, repoErr)
				return mock
			},
		},
		{
			name: "BuyItemMock repo error case",
			err:  repoErr,
			args: args{
				ctx:    ctx,
				userID: userID,
				title:  title,
			},
			repositoryMock: func(mc *minimock.Controller) repository.Repository {
				mock := repoMocks.NewRepositoryMock(mc)
				mock.ItemMock.Expect(ctx, title).Return(item, nil)
				mock.BuyItemMock.Expect(ctx, userID, item).Return(repoErr)
				return mock
			},
		},
		{
			name: "item not found error case",
			err:  service.ErrItemNotFound,
			args: args{
				ctx:    ctx,
				userID: userID,
				title:  title,
			},
			repositoryMock: func(mc *minimock.Controller) repository.Repository {
				mock := repoMocks.NewRepositoryMock(mc)
				mock.ItemMock.Expect(ctx, title).Return(nil, repository.ErrNotFound)
				return mock
			},
		},
		{
			name: "insufficient balance error case",
			err:  service.ErrInsufficientBalance,
			args: args{
				ctx:    ctx,
				userID: userID,
				title:  title,
			},
			repositoryMock: func(mc *minimock.Controller) repository.Repository {
				mock := repoMocks.NewRepositoryMock(mc)
				mock.ItemMock.Expect(ctx, title).Return(item, nil)
				mock.BuyItemMock.Expect(ctx, userID, item).Return(repository.ErrNegativeBalance)
				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repoMock := tt.repositoryMock(mc)
			service := NewMockService(repoMock, config.NewMockAuthConfig())

			err := service.BuyItem(tt.args.ctx, tt.args.userID, tt.args.title)
			require.Equal(t, tt.err, err)
		})
	}
}
