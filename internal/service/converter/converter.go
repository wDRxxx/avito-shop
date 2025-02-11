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
