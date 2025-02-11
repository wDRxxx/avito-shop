package config

import (
	"net"
	"os"
)

type HttpConfig struct {
	host string
	port string
}

func NewHttpConfig() *HttpConfig {
	host := os.Getenv("HTTP_HOST")
	if host == "" {
		panic("HTTP_HOST environment variable is empty")
	}

	port := os.Getenv("HTTP_PORT")
	if port == "" {
		panic("HTTP_PORT environment variable is empty")
	}

	return &HttpConfig{
		host: host,
		port: port,
	}
}

func (c *HttpConfig) Address() string {
	return net.JoinHostPort(c.host, c.port)
}
