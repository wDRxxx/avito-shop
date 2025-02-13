package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/wDRxxx/avito-shop/internal/repository"
	"github.com/wDRxxx/avito-shop/internal/service"
	"github.com/wDRxxx/avito-shop/internal/service/converter"
	sm "github.com/wDRxxx/avito-shop/internal/service/models"
	"github.com/wDRxxx/avito-shop/internal/utils"
)

func (s *serv) UserToken(ctx context.Context, username string, password string) (string, error) {
	u, err := s.repo.UserByUsername(ctx, username)
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
