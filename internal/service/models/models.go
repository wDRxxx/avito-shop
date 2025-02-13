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

type InventoryItem struct {
	ID       int
	Title    string
	Quantity int
}

type IncomingTransaction struct {
	SenderUsername string
	Amount         int
}

type OutgoingTransaction struct {
	RecipientUsername string
	Amount            int
}

type UserInfo struct {
	Balance              int
	InventoryItems       []*InventoryItem
	IncomingTransactions []*IncomingTransaction
	OutgoingTransactions []*OutgoingTransaction
}
