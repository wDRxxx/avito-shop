package service

import (
	"github.com/pkg/errors"
)

var (
	ErrWrongCredentials    = errors.New("wrong credentials")
	ErrItemNotFound        = errors.New("item not found")
	ErrUserNotFound        = errors.New("user not found")
	ErrInsufficientBalance = errors.New("insufficient balance")
)
