package tgutil

import (
	"strings"

	sharedEntity "telemafia/internal/shared/entity"

	"gopkg.in/telebot.v3"
)

var adminUsernames []string

// --- Refresh Logic (potentially move to its own package/file later) ---

// RefreshingMessageType defines the type of content being refreshed.
type RefreshingMessageType int

const (
	ListRooms  RefreshingMessageType = iota
	RoomDetail                       // Placeholder for potential future use
)

type Updater interface {
	SendOrUpdateRefreshingMessage(int64, RefreshingMessageType, string) error
}

// SetAdminUsers stores the list of admin usernames.
// NOTE: This uses a package-level variable, which isn't ideal for testing.
// Consider injecting a config or admin checker service instead in a real app.
func SetAdminUsers(usernames []string) {
	adminUsernames = usernames
}

// IsAdmin checks if a given username is in the configured admin list.
func IsAdmin(username string) bool {
	if username == "" {
		return false
	}
	for _, admin := range adminUsernames {
		if strings.EqualFold(username, admin) { // Case-insensitive check for admin username
			return true
		}
	}
	return false
}

// ToUser converts a telebot.User to our internal sharedEntity.User, checking admin status.
func ToUser(sender *telebot.User) *sharedEntity.User {
	if sender == nil {
		return nil
	}
	isAdmin := IsAdmin(sender.Username)
	return &sharedEntity.User{
		ID:         sharedEntity.UserID(sender.ID), // Assuming UserID is based on Telegram ID
		TelegramID: sender.ID,
		FirstName:  sender.FirstName,
		LastName:   sender.LastName,
		Username:   sender.Username,
		Admin:      isAdmin,
	}
}

// SplitCallbackData extracts the unique identifier and payload from callback data.
// Assumes format "unique:payload" or just "unique".
func SplitCallbackData(data string) (unique string, payload string) {
	parts := strings.SplitN(data, "|", 2)
	unique = strings.TrimSpace(parts[0])
	if len(parts) == 2 {
		payload = strings.TrimSpace(parts[1])
	}
	return
}
