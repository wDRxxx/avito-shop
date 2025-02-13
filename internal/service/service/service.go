package service

import (
	"github.com/wDRxxx/avito-shop/internal/config"
	"github.com/wDRxxx/avito-shop/internal/repository"
	"github.com/wDRxxx/avito-shop/internal/service"
)

type serv struct {
	repo       repository.Repository
	authConfig config.AuthConfig
}

func NewService(
	repo repository.Repository,
	authConfig config.AuthConfig,
) service.Service {
	s := &serv{
		repo:       repo,
		authConfig: authConfig,
	}

	return s
}
