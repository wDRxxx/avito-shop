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
	Type        bool
	SenderID    int
	RecipientID int
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

type Inventory struct {
	ID       int
	UserID   int
	ItemID   int
	Quantity int

	CreatedAt time.Time
	UpdatedAt time.Time
}
