package util

import (
	"telemafia/common"
	userEntity "telemafia/internal/user/entity"

	"gopkg.in/telebot.v3"
)

var adminUsernames []string = nil

func SetAdminUsers(AdminUsernames []string) {
	adminUsernames = AdminUsernames
}

func ToUser(sender *telebot.User) *userEntity.User {
	return &userEntity.User{
		ID:        userEntity.UserID(sender.ID),
		FirstName: sender.FirstName,
		LastName:  sender.LastName,
		Username:  sender.Username,
		Admin:     common.Contains(adminUsernames, sender.Username),
	}
}

// IsAdmin checks if a username is in the list of admin usernames
func IsAdmin(username string) bool {
	for _, admin := range adminUsernames {
		if admin == username {
			return true
		}
	}
	return false
}
