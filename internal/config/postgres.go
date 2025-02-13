package config

import (
	"fmt"
	"os"
	"time"
)

type PostgresConfig interface {
	Timeout() time.Duration
	ConnectionString() string
}

type postgresConfig struct {
	host     string
	port     string
	user     string
	password string
	database string
	timeout  time.Duration
}

func (c *postgresConfig) Timeout() time.Duration {
	return c.timeout
}

func (c *postgresConfig) ConnectionString() string {
	return fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		c.host,
		c.port,
		c.database,
		c.user,
		c.password,
	)
}

func NewPostgresConfig() PostgresConfig {
	host := os.Getenv("POSTGRES_HOST")
	if host == "" {
		panic("POSTGRES_HOST environment variable is empty")
	}

	port := os.Getenv("POSTGRES_PORT")
	if port == "" {
		panic("POSTGRES_PORT environment variable is empty")
	}

	user := os.Getenv("POSTGRES_USER")
	if user == "" {
		panic("POSTGRES_USER environment variable is empty")
	}

	password := os.Getenv("POSTGRES_PASSWORD")
	if password == "" {
		panic("POSTGRES_PASSWORD environment variable is empty")
	}

	database := os.Getenv("POSTGRES_DB")
	if database == "" {
		panic("POSTGRES_DB environment variable is empty")
	}

	timeout := os.Getenv("POSTGRES_TIMEOUT")
	if timeout == "" {
		panic("POSTGRES_TIMEOUT environment variable is empty")
	}

	t, err := time.ParseDuration(timeout)
	if err != nil {
		panic("POSTGRES_TIMEOUT must be an integer")
	}

	return &postgresConfig{
		host:     host,
		port:     port,
		user:     user,
		password: password,
		database: database,
		timeout:  t * time.Second,
	}
}
