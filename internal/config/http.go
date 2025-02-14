package config

import (
	"net"
	"os"
)

type HttpConfig interface {
	Address() string
}

type httpConfig struct {
	host string
	port string
}

func (c *httpConfig) Address() string {
	return net.JoinHostPort(c.host, c.port)
}

func NewHttpConfig() HttpConfig {
	host := os.Getenv("HTTP_HOST")
	if host == "" {
		panic("HTTP_HOST environment variable is empty")
	}

	port := os.Getenv("HTTP_PORT")
	if port == "" {
		panic("HTTP_PORT environment variable is empty")
	}

	return &httpConfig{
		host: host,
		port: port,
	}
}

func NewMockHttpConfig() HttpConfig {
	return &httpConfig{
		host: "::",
		port: "8080",
	}
}
