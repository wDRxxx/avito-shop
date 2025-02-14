package models

import (
	am "github.com/wDRxxx/avito-shop/internal/models"
)

type ErrorResponse struct {
	Errors string `json:"errors"`
}

type AuthResponse struct {
	am.AuthResponse
	ErrorResponse
}
