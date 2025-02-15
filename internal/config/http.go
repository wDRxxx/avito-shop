package config

import (
	"net"
	"os"
	"time"
)

type HttpConfig interface {
	Address() string
	ReadHeaderTimeout() time.Duration
}

type httpConfig struct {
	host              string
	port              string
	readHeaderTimeout time.Duration
}

func (c *httpConfig) ReadHeaderTimeout() time.Duration {
	return c.readHeaderTimeout
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

	readHeaderTimeout := os.Getenv("HTTP_READ_HEADER_TIMEOUT")
	if readHeaderTimeout == "" {
		panic("HTTP_READ_HEADER_TIMEOUT environment variable is empty")
	}

	t, err := time.ParseDuration(readHeaderTimeout)
	if err != nil {
		panic("HTTP_READ_HEADER_TIMEOUT environment variable has wrong format")
	}

	return &httpConfig{
		host:              host,
		port:              port,
		readHeaderTimeout: t,
	}
}

func NewMockHttpConfig() HttpConfig {
	return &httpConfig{
		host: "::",
		port: "8080",
	}
}
