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
)

func TestSendCoin(t *testing.T) {
	t.Parallel()

	type repositoryMockFunc func(mc *minimock.Controller) repository.Repository

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		recipientUsername = gofakeit.Username()
		password          = gofakeit.Password(true, true, true, false, false, 12)
		hashPass, _       = bcrypt.GenerateFromPassword([]byte(password), 12)
		strHashPass       = string(hashPass)

		userID    = int(gofakeit.Uint8())
		recipient = &rm.User{
			ID:        int(gofakeit.Uint8()),
			Username:  recipientUsername,
			Password:  strHashPass,
			Balance:   int(gofakeit.Uint8()),
			CreatedAt: gofakeit.Date(),
			UpdatedAt: gofakeit.Date(),
		}

		amount = recipient.Balance / 2

		repoErr = errors.New("repo err")
	)

	type args struct {
		ctx        context.Context
		toUser     string
		fromUserID int
		Amount     int
	}

	tests := []struct {
		name           string
		err            error
		args           args
		repositoryMock repositoryMockFunc
	}{
		{
			name: "success case",
			err:  nil,
			args: args{
				ctx:        ctx,
				toUser:     recipientUsername,
				fromUserID: userID,
				Amount:     amount,
			},
			repositoryMock: func(mc *minimock.Controller) repository.Repository {
				mock := repoMocks.NewRepositoryMock(mc)
				mock.UserByUsernameMock.Expect(ctx, recipientUsername).Return(recipient, nil)
				mock.SendCoinMock.Expect(ctx, recipient.ID, userID, amount).Return(nil)
				return mock
			},
		},
		{
			name: "UserByUsernameMock repo error case",
			err:  repoErr,
			args: args{
				ctx:        ctx,
				toUser:     recipientUsername,
				fromUserID: userID,
				Amount:     amount,
			},
			repositoryMock: func(mc *minimock.Controller) repository.Repository {
				mock := repoMocks.NewRepositoryMock(mc)
				mock.UserByUsernameMock.Expect(ctx, recipientUsername).Return(nil, repoErr)
				return mock
			},
		},
		{
			name: "SendCoinMock repo error case",
			err:  repoErr,
			args: args{
				ctx:        ctx,
				toUser:     recipientUsername,
				fromUserID: userID,
				Amount:     amount,
			},
			repositoryMock: func(mc *minimock.Controller) repository.Repository {
				mock := repoMocks.NewRepositoryMock(mc)
				mock.UserByUsernameMock.Expect(ctx, recipientUsername).Return(recipient, nil)
				mock.SendCoinMock.Expect(ctx, recipient.ID, userID, amount).Return(repoErr)
				return mock
			},
		},
		{
			name: "user not found error case",
			err:  service.ErrUserNotFound,
			args: args{
				ctx:        ctx,
				toUser:     recipientUsername,
				fromUserID: userID,
				Amount:     amount,
			},
			repositoryMock: func(mc *minimock.Controller) repository.Repository {
				mock := repoMocks.NewRepositoryMock(mc)
				mock.UserByUsernameMock.Expect(ctx, recipientUsername).Return(nil, service.ErrUserNotFound)
				return mock
			},
		},
		{
			name: "insufficient balance error case 1",
			err:  service.ErrInsufficientBalance,
			args: args{
				ctx:        ctx,
				toUser:     recipientUsername,
				fromUserID: userID,
				Amount:     recipient.Balance + 1,
			},
			repositoryMock: func(mc *minimock.Controller) repository.Repository {
				mock := repoMocks.NewRepositoryMock(mc)
				mock.UserByUsernameMock.Expect(ctx, recipientUsername).Return(recipient, nil)
				return mock
			},
		},
		{
			name: "insufficient balance error case 2",
			err:  service.ErrInsufficientBalance,
			args: args{
				ctx:        ctx,
				toUser:     recipientUsername,
				fromUserID: userID,
				Amount:     amount,
			},
			repositoryMock: func(mc *minimock.Controller) repository.Repository {
				mock := repoMocks.NewRepositoryMock(mc)
				mock.UserByUsernameMock.Expect(ctx, recipientUsername).Return(recipient, nil)
				mock.SendCoinMock.Expect(ctx, recipient.ID, userID, amount).Return(repository.ErrNegativeBalance)
				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoMock := tt.repositoryMock(mc)
			service := NewMockService(repoMock, config.NewMockAuthConfig())

			err := service.SendCoin(tt.args.ctx, tt.args.toUser, tt.args.fromUserID, tt.args.Amount)
			require.Equal(t, tt.err, err)
		})
	}
}
