package models

import (
	"github.com/golang-jwt/jwt/v5"
)

type User struct {
	ID       int
	Username string
	Password string
	Balance  int
}

type UserClaims struct {
	jwt.RegisteredClaims
	Username string
}
