package config

import (
	"os"
	"time"

	"github.com/wDRxxx/avito-shop/internal/utils"
)

type AuthConfig interface {
	TokenSecret() string
	TokenTTL() time.Duration
}

type authConfig struct {
	tokenSecret string
	tokenTTL    time.Duration
}

func (c *authConfig) TokenSecret() string {
	return c.tokenSecret
}

func (c *authConfig) TokenTTL() time.Duration {
	return c.tokenTTL
}

func NewAuthConfig() AuthConfig {
	secret := os.Getenv("TOKEN_SECRET")
	if secret == "" {
		panic("TOKEN_SECRET environment variable is empty")
	}

	ttlStr := os.Getenv("TOKEN_TTL")
	if ttlStr == "" {
		panic("TOKEN_TTL environment variable is empty")
	}
	ttl, err := utils.ParseCustomDuration(ttlStr)
	if err != nil {
		panic("TOKEN_TTL environment variable is invalid")
	}

	return &authConfig{
		tokenSecret: secret,
		tokenTTL:    ttl,
	}
}

func NewMockAuthConfig() AuthConfig {
	return &authConfig{
		tokenSecret: "mock",
		tokenTTL:    1 * time.Minute,
	}
}
