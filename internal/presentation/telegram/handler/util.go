package telegram

import (
	"gopkg.in/telebot.v3"
	"telemafia/internal/shared/common"
	// "telemafia/internal/entity" // Old import
	sharedEntity "telemafia/internal/shared/entity" // Use shared entity
)

var adminUsernames []string

// SetAdminUsers stores the list of admin usernames locally.
func SetAdminUsers(usernames []string) {
	adminUsernames = make([]string, len(usernames))
	copy(adminUsernames, usernames)
}

// ToUser converts a telebot.User to our internal sharedEntity.User.
func ToUser(sender *telebot.User) *sharedEntity.User {
	if sender == nil {
		return nil // Handle nil sender
	}
	return &sharedEntity.User{
		ID:         sharedEntity.UserID(sender.ID),
		TelegramID: sender.ID, // Store Telegram specific ID
		FirstName:  sender.FirstName,
		LastName:   sender.LastName,
		Username:   sender.Username,
		Admin:      IsAdmin(sender.Username),
	}
}

// IsAdmin checks if a username is in the list of admin usernames.
func IsAdmin(username string) bool {
	return common.Contains(adminUsernames, username)
}
