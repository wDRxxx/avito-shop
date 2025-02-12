package service

import (
	"context"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"

	"github.com/wDRxxx/avito-shop/internal/config"
	"github.com/wDRxxx/avito-shop/internal/repository"
	"github.com/wDRxxx/avito-shop/internal/service"
	"github.com/wDRxxx/avito-shop/internal/service/converter"
	sm "github.com/wDRxxx/avito-shop/internal/service/models"
	"github.com/wDRxxx/avito-shop/internal/utils"
)

type serv struct {
	repo       repository.Repository
	authConfig *config.AuthConfig
}

func NewService(
	repo repository.Repository,
	authConfig *config.AuthConfig,
) service.Service {
	s := &serv{
		repo:       repo,
		authConfig: authConfig,
	}

	return s
}

func (s *serv) UserToken(ctx context.Context, username string, password string) (string, error) {
	u, err := s.repo.User(ctx, username)
	if err != nil {
		if !errors.Is(err, repository.ErrNotFound) {
			return "", err
		}

		pass, err := bcrypt.GenerateFromPassword([]byte(password), 12)
		if err != nil {
			return "", err
		}
		u, err = s.repo.InsertUser(ctx, username, string(pass))
		if err != nil {
			return "", err
		}
	}
	user := converter.UserFromRepositoryToService(u)

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", service.ErrWrongCredentials
	}

	// generate token
	token, err := utils.GenerateToken(&sm.UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: fmt.Sprint(user.ID),
		}, Username: user.Username,
	}, s.authConfig.TokenSecret(), s.authConfig.TokenTTL())
	if err != nil {
		return "", err
	}

	return token, nil
}

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
		return err
	}

	return nil
}

func (s *serv) SendCoin(ctx context.Context, toUser string, fromUserID int, amount int) error {
	user, err := s.repo.User(ctx, toUser)
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
