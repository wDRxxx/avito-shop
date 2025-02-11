package httpServer

import (
	"net/http"
	"strings"

	"github.com/wDRxxx/avito-shop/internal/api"
	"github.com/wDRxxx/avito-shop/internal/config"
	"github.com/wDRxxx/avito-shop/internal/service"
	sm "github.com/wDRxxx/avito-shop/internal/service/models"
	"github.com/wDRxxx/avito-shop/internal/utils"
)

type server struct {
	mux http.Handler

	authConfig *config.AuthConfig
	httpConfig *config.HttpConfig

	service service.Service
}

func NewHTTPServer(
	authConfig *config.AuthConfig,
	httpConfig *config.HttpConfig,
	service service.Service,
) api.HTTPServer {
	s := &server{
		authConfig: authConfig,
		httpConfig: httpConfig,
		service:    service,
	}

	s.setRoutes()

	return s
}

func (s *server) Handler() http.Handler {
	return s.mux
}

func (s *server) getAndVerifyHeaderToken(r *http.Request) (string, *sm.UserClaims, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", nil, errInvalidAuthHeader
	}

	exploded := strings.Split(authHeader, " ")
	if len(exploded) != 2 || exploded[0] != "Bearer" {
		return "", nil, errInvalidAuthHeader
	}

	token := exploded[1]
	claims, err := utils.VerifyToken(token, s.authConfig.TokenSecret())
	if err != nil {
		return "", nil, errInvalidAuthHeader
	}

	return token, claims, nil
}
