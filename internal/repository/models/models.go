package models

import (
	"time"
)

type User struct {
	ID       int
	Username string
	Password string
	Balance  int

	CreatedAt time.Time
	UpdatedAt time.Time
}

type Transaction struct {
	ID          int
	SenderID    int
	Sender      User
	RecipientID int
	Recipient   User
	Amount      int

	CreatedAt time.Time
	UpdatedAt time.Time
}

type Item struct {
	ID    int
	Title string
	Price int

	CreatedAt time.Time
	UpdatedAt time.Time
}

type InventoryItem struct {
	ID       int
	UserID   int
	ItemID   int
	Item     Item
	Quantity int

	CreatedAt time.Time
	UpdatedAt time.Time
}
