package tgutil

import (
	"strconv"
	"strings"

	sharedEntity "telemafia/internal/shared/entity"

	"gopkg.in/telebot.v4"
)

var adminIds []int64

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
func SetAdminUsers(idList []string) {
	adminIds = make([]int64, len(idList))
	for _, idString := range idList {
		id, err := strconv.ParseInt(idString, 10, 64)
		if err != nil {
			panic("Cannot parse admin IDs")
		}
		adminIds = append(adminIds, id)
	}
}

// IsAdmin checks if a given username is in the configured admin list.
func IsAdmin(id int64) bool {
	for _, admin := range adminIds {
		if id == admin {
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
	isAdmin := IsAdmin(sender.ID)
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
