package api

import (
	"github.com/pkg/errors"
)

var (
	ErrInternal            = errors.New("internal error, please, try again later")
	ErrSendToYourself      = errors.New("you can't send coins to yourself")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrUserNotFound        = errors.New("user not found")
	ErrInsufficientBalance = errors.New("insufficient balance")
	ErrItemNotFound        = errors.New("item not found")
	ErrWrongCredentials    = errors.New("wrong credentials")
)
