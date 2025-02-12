package api

import (
	"github.com/pkg/errors"
)

var (
	ErrInternal       = errors.New("internal error, please, try again later")
	ErrSendToYourself = errors.New("you can't send coins to yourself")
)
