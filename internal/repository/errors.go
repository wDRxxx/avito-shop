package repository

import (
	"github.com/pkg/errors"
)

var (
	ErrNegativeBalance = errors.New("balance can't be negative")
	ErrNotFound        = errors.New("not found")
)
