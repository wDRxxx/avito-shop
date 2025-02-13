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

func TestUserToken(t *testing.T) {
	t.Parallel()

	type repositoryMockFunc func(mc *minimock.Controller) repository.Repository

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		username      = gofakeit.Username()
		password      = gofakeit.Password(true, true, true, false, false, 12)
		wrongPassword = gofakeit.Password(true, true, true, false, false, 12)
		hashPass, _   = bcrypt.GenerateFromPassword([]byte(password), 12)
		strHashPass   = string(hashPass)

		user = &rm.User{
			ID:        gofakeit.Int(),
			Username:  username,
			Password:  strHashPass,
			Balance:   int(gofakeit.Uint8()),
			CreatedAt: gofakeit.Date(),
			UpdatedAt: gofakeit.Date(),
		}

		repoErr = errors.New("repo err")
	)

	type args struct {
		ctx      context.Context
		username string
		password string
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
				ctx:      ctx,
				username: username,
				password: password,
			},
			repositoryMock: func(mc *minimock.Controller) repository.Repository {
				mock := repoMocks.NewRepositoryMock(mc)
				mock.UserByUsernameMock.Expect(ctx, username).Return(user, nil)
				return mock
			},
		},
		{
			name: "success case when user doesnt exist",
			err:  nil,
			args: args{
				ctx:      ctx,
				username: username,
				password: password,
			},
			repositoryMock: func(mc *minimock.Controller) repository.Repository {
				mock := repoMocks.NewRepositoryMock(mc)
				mock.UserByUsernameMock.Expect(ctx, username).Return(nil, repository.ErrNotFound)
				mock.InsertUserMock.Set(func(ctx context.Context, usrname, hashedPassword string) (*rm.User, error) {
					if usrname != username {
						return nil, service.ErrWrongCredentials
					}

					if bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) != nil {
						return nil, service.ErrWrongCredentials
					}

					return user, nil
				})
				return mock
			},
		},
		{
			name: "service error case",
			err:  repoErr,
			args: args{
				ctx:      ctx,
				username: username,
				password: password,
			},
			repositoryMock: func(mc *minimock.Controller) repository.Repository {
				mock := repoMocks.NewRepositoryMock(mc)
				mock.UserByUsernameMock.Expect(ctx, username).Return(nil, repoErr)
				return mock
			},
		},
		{
			name: "wrong credentials case",
			err:  service.ErrWrongCredentials,
			args: args{
				ctx:      ctx,
				username: username,
				password: wrongPassword,
			},
			repositoryMock: func(mc *minimock.Controller) repository.Repository {
				mock := repoMocks.NewRepositoryMock(mc)
				mock.UserByUsernameMock.Expect(ctx, username).Return(nil, repository.ErrNotFound)
				mock.InsertUserMock.Set(func(ctx context.Context, usrname, hashedPassword string) (*rm.User, error) {
					if usrname != username {
						return nil, service.ErrWrongCredentials
					}

					if bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) != nil {
						return nil, service.ErrWrongCredentials
					}

					return user, nil
				})
				return mock
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoMock := tt.repositoryMock(mc)
			service := NewMockService(repoMock, config.NewMockAuthConfig())

			_, err := service.UserToken(tt.args.ctx, tt.args.username, tt.args.password)
			require.Equal(t, tt.err, err)
		})
	}
}
