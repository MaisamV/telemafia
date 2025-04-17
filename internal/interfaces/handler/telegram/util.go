package telegram

import (
	"telemafia/internal/entity"
	"telemafia/pkg/common"

	"gopkg.in/telebot.v3"
)

var adminUsernames []string

func SetAdminUsers(usernames []string) {
	adminUsernames = make([]string, len(usernames))
	copy(adminUsernames, usernames)
}

func ToUser(sender *telebot.User) *entity.User {
	return &entity.User{
		ID:        entity.UserID(sender.ID),
		FirstName: sender.FirstName,
		LastName:  sender.LastName,
		Username:  sender.Username,
		Admin:     IsAdmin(sender.Username),
	}
}

// IsAdmin checks if a username is in the list of admin usernames
func IsAdmin(username string) bool {
	return common.Contains(adminUsernames, username)
}
