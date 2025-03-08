package common

import (
	"telemafia/common"
	"telemafia/config"
	userEntity "telemafia/internal/user/entity"

	"gopkg.in/telebot.v3"
)

func ToUser(sender *telebot.User) *userEntity.User {
	return &userEntity.User{
		ID:        userEntity.UserID(sender.ID),
		FirstName: sender.FirstName,
		LastName:  sender.LastName,
		Username:  sender.Username,
		Admin:     common.Contains(config.GetGlobalConfig().AdminUsernames, sender.Username),
	}
}
