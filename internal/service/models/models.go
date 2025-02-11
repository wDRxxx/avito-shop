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

type Item struct {
	ID    int
	Title string
	Price int
}

type Transaction struct {
	ID          int
	Type        bool
	SenderID    int
	RecipientID int
	Amount      int
}

type Inventory struct {
	ID       int
	UserID   int
	ItemID   int
	Quantity int
}
