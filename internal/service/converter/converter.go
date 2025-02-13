package converter

import (
	rm "github.com/wDRxxx/avito-shop/internal/repository/models"
	sm "github.com/wDRxxx/avito-shop/internal/service/models"
)

func UserFromRepositoryToService(user *rm.User) *sm.User {
	return &sm.User{
		ID:       user.ID,
		Username: user.Username,
		Password: user.Password,
		Balance:  user.Balance,
	}
}

func ItemFromRepositoryToService(item *rm.Item) *sm.Item {
	return &sm.Item{
		ID:    item.ID,
		Title: item.Title,
		Price: item.Price,
	}
}
func InventoryFromRepositoryToService(inventoryItems []*rm.InventoryItem) []*sm.InventoryItem {
	var items []*sm.InventoryItem
	for _, item := range inventoryItems {
		items = append(items, InventoryItemFromRepositoryToService(item))
	}

	return items
}

func InventoryItemFromRepositoryToService(inventoryItem *rm.InventoryItem) *sm.InventoryItem {
	return &sm.InventoryItem{
		ID:       inventoryItem.ID,
		Title:    inventoryItem.Item.Title,
		Quantity: inventoryItem.Quantity,
	}
}

func IncomingTransactionFromRepositoryToService(transaction *rm.Transaction) *sm.IncomingTransaction {
	return &sm.IncomingTransaction{
		SenderUsername: transaction.Sender.Username,
		Amount:         transaction.Amount,
	}
}

func IncomingTransactionsFromRepositoryToService(transactions []*rm.Transaction) []*sm.IncomingTransaction {
	var res []*sm.IncomingTransaction
	for _, transaction := range transactions {
		res = append(res, IncomingTransactionFromRepositoryToService(transaction))
	}
	return res
}

func OutgoingTransactionFromRepositoryToService(transaction *rm.Transaction) *sm.OutgoingTransaction {
	return &sm.OutgoingTransaction{
		RecipientUsername: transaction.Recipient.Username,
		Amount:            transaction.Amount,
	}
}

func OutgoingTransactionsFromRepositoryToService(transactions []*rm.Transaction) []*sm.OutgoingTransaction {
	var res []*sm.OutgoingTransaction
	for _, transaction := range transactions {
		res = append(res, OutgoingTransactionFromRepositoryToService(transaction))
	}
	return res
}
