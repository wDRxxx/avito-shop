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
