package config

import (
	"os"
	"time"

	"github.com/wDRxxx/avito-shop/internal/utils"
)

type AuthConfig struct {
	tokenSecret string
	tokenTTL    time.Duration
}

func (c *AuthConfig) TokenSecret() string {
	return c.tokenSecret
}

func (c *AuthConfig) TokenTTL() time.Duration {
	return c.tokenTTL
}

func NewAuthConfig() *AuthConfig {
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

	return &AuthConfig{
		tokenSecret: secret,
		tokenTTL:    ttl,
	}
}
