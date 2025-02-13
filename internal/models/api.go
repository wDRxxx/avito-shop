package models

type ErrorResponse struct {
	Errors string `json:"errors"`
}

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

type SendCoinRequest struct {
	ToUser string `json:"toUser"`
	Amount int    `json:"amount"`
}

type InfoResponse struct {
	Coins       int              `json:"coins"`
	Inventory   []*InventoryItem `json:"inventory"`
	CoinHistory *CoinHistory     `json:"coinHistory"`
}

type InventoryItem struct {
	Type     string `json:"type"`
	Quantity int    `json:"quantity"`
}

type CoinHistory struct {
	Received []ReceivedCoinsItem `json:"received"`
	Sent     []SentCoinsItem     `json:"sent"`
}

type ReceivedCoinsItem struct {
	FromUser string `json:"fromUser"`
	Amount   int    `json:"amount"`
}

type SentCoinsItem struct {
	ToUser string `json:"toUser"`
	Amount int    `json:"amount"`
}
