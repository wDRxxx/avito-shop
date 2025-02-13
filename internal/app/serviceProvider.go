package app

import (
	"context"
	"log"
	"log/slog"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/wDRxxx/avito-shop/internal/api"
	"github.com/wDRxxx/avito-shop/internal/api/httpServer"
	"github.com/wDRxxx/avito-shop/internal/closer"
	"github.com/wDRxxx/avito-shop/internal/config"
	"github.com/wDRxxx/avito-shop/internal/repository"
	"github.com/wDRxxx/avito-shop/internal/repository/postgres"
	"github.com/wDRxxx/avito-shop/internal/service"
	serviceImpl "github.com/wDRxxx/avito-shop/internal/service/service"
)

type serviceProvider struct {
	httpConfig     config.HttpConfig
	postgresConfig config.PostgresConfig
	authConfig     config.AuthConfig

	repository repository.Repository
	httpServer api.HTTPServer

	service service.Service
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) HttpConfig() config.HttpConfig {
	if s.httpConfig == nil {
		s.httpConfig = config.NewHttpConfig()
	}

	return s.httpConfig
}

func (s *serviceProvider) PostgresConfig() config.PostgresConfig {
	if s.postgresConfig == nil {
		s.postgresConfig = config.NewPostgresConfig()
	}
	return s.postgresConfig
}

func (s *serviceProvider) AuthConfig() config.AuthConfig {
	if s.authConfig == nil {
		s.authConfig = config.NewAuthConfig()
	}

	return s.authConfig
}

func (s *serviceProvider) Repository(ctx context.Context) repository.Repository {
	if s.repository == nil {
		db, err := pgxpool.New(ctx, s.PostgresConfig().ConnectionString())
		if err != nil {
			log.Fatalf("error connecting to database: %v", err)
		}
		closer.Add(2, func() error {
			slog.Info("closing pgxpool")
			db.Close()
			return nil
		})

		err = db.Ping(ctx)
		if err != nil {
			log.Fatalf("error connecting to database: %v", err)
		}

		s.repository = postgres.NewPostgresRepo(db, s.PostgresConfig().Timeout())
	}

	return s.repository
}

func (s *serviceProvider) Service(ctx context.Context) service.Service {
	if s.service == nil {
		s.service = serviceImpl.NewService(s.Repository(ctx), s.AuthConfig())
	}

	return s.service
}

func (s *serviceProvider) HTTPServer(ctx context.Context, wg *sync.WaitGroup) api.HTTPServer {
	if s.httpServer == nil {
		s.httpServer = httpServer.NewHTTPServer(
			s.AuthConfig(),
			s.HttpConfig(),
			s.Service(ctx),
		)
	}

	return s.httpServer
}
