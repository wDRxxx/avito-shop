package usersService

import (
	"context"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"

	"github.com/wDRxxx/avito-shop/internal/config"
	"github.com/wDRxxx/avito-shop/internal/repository"
	"github.com/wDRxxx/avito-shop/internal/service"
	"github.com/wDRxxx/avito-shop/internal/service/converter"
	sm "github.com/wDRxxx/avito-shop/internal/service/models"
	"github.com/wDRxxx/avito-shop/internal/utils"
)

type usersServ struct {
	repo       repository.Repository
	authConfig *config.AuthConfig
}

func NewUsersService(
	repo repository.Repository,
	authConfig *config.AuthConfig,
) service.UsersService {
	s := &usersServ{
		repo:       repo,
		authConfig: authConfig,
	}

	return s
}

func (s *usersServ) UserToken(ctx context.Context, username string, password string) (string, error) {
	u, err := s.repo.User(ctx, username)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
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
