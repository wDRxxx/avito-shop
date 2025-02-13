package converter

import (
	"github.com/wDRxxx/avito-shop/internal/models"
	sm "github.com/wDRxxx/avito-shop/internal/service/models"
)

func UserInfoFromServiceToApi(info *sm.UserInfo) *models.InfoResponse {
	return &models.InfoResponse{
		Coins:       info.Balance,
		Inventory:   InventoryFromServiceToApi(info.InventoryItems),
		CoinHistory: CoinHistoryFromServiceToApi(info.IncomingTransactions, info.OutgoingTransactions),
	}
}

func InventoryFromServiceToApi(inventory []*sm.InventoryItem) []*models.InventoryItem {
	var res []*models.InventoryItem
	for _, item := range inventory {
		res = append(res, &models.InventoryItem{
			Type:     item.Title,
			Quantity: item.Quantity,
		})
	}
	return res
}

func CoinHistoryFromServiceToApi(incoming []*sm.IncomingTransaction, outgoing []*sm.OutgoingTransaction) *models.CoinHistory {
	return &models.CoinHistory{
		Received: ReceivedCoinsFromServiceToApi(incoming),
		Sent:     SentCoinsFromServiceToApi(outgoing),
	}
}

func ReceivedCoinsFromServiceToApi(incoming []*sm.IncomingTransaction) []models.ReceivedCoinsItem {
	var res []models.ReceivedCoinsItem
	for _, trans := range incoming {
		res = append(res, ReceivedCoinsItemFromServiceToApi(trans))
	}
	return res
}

func ReceivedCoinsItemFromServiceToApi(incoming *sm.IncomingTransaction) models.ReceivedCoinsItem {
	return models.ReceivedCoinsItem{
		FromUser: incoming.SenderUsername,
		Amount:   incoming.Amount,
	}
}

func SentCoinsFromServiceToApi(outgoing []*sm.OutgoingTransaction) []models.SentCoinsItem {
	var res []models.SentCoinsItem
	for _, trans := range outgoing {
		res = append(res, SentCoinsItemFromServiceToApi(trans))
	}
	return res
}

func SentCoinsItemFromServiceToApi(outgoing *sm.OutgoingTransaction) models.SentCoinsItem {
	return models.SentCoinsItem{
		ToUser: outgoing.RecipientUsername,
		Amount: outgoing.Amount,
	}
}
