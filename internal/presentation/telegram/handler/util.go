package telegram

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"telemafia/internal/shared/common"

	"gopkg.in/telebot.v3"

	roomQuery "telemafia/internal/domain/room/usecase/query"
	scenarioEntity "telemafia/internal/domain/scenario/entity"
	sharedEntity "telemafia/internal/shared/entity"
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
		Admin:      IsAdmin(sender.Username), // Call IsAdmin locally
	}
}

// IsAdmin checks if a username is in the list of admin usernames.
func IsAdmin(username string) bool {
	return common.Contains(adminUsernames, username)
}

// SplitCallbackData helper function to parse callback data (unique:payload)
func SplitCallbackData(data string) (unique string, payload string) {
	parts := strings.SplitN(data, ":", 2)
	if len(parts) == 2 {
		return parts[0], parts[1]
	} else if len(parts) == 1 {
		switch parts[0] {
		case UniqueCancel:
			return parts[0], ""
		default:
			// log.Printf("Warning: Callback data '%s' has no colon, assuming empty payload.", data)
			return parts[0], ""
		}
	}
	// log.Printf("Warning: Could not split callback data '%s' correctly.", data)
	return "", data
}

// getUserDisplayName is a helper to get a display name (now a method on BotHandler)
func (h *BotHandler) getUserDisplayName(userID sharedEntity.UserID) string {
	// Inefficiently search all rooms for the user - Needs optimization
	// Consider adding a GetUserByID method to a shared user repository/service.
	allRooms, _ := h.getRoomsHandler.Handle(context.Background(), roomQuery.GetRoomsQuery{})
	for _, room := range allRooms {
		players, _ := h.getPlayersInRoomHandler.Handle(context.Background(), roomQuery.GetPlayersInRoomQuery{RoomID: room.ID})
		for _, player := range players {
			if player != nil && player.ID == userID {
				if player.Username != "" {
					return "@" + player.Username
				} else if player.FirstName != "" {
					name := player.FirstName
					if player.LastName != "" {
						name += " " + player.LastName
					}
					return name
				}
			}
		}
	}
	return fmt.Sprintf("User %d", userID)
}

// formatAssignments formats the role assignments map for display (now a method on BotHandler)
func (h *BotHandler) formatAssignments(assignments map[sharedEntity.UserID]scenarioEntity.Role) string {
	var b strings.Builder
	ids := make([]int64, 0, len(assignments))
	for uid := range assignments {
		ids = append(ids, int64(uid))
	}
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })

	for _, id := range ids {
		uid := sharedEntity.UserID(id)
		role := assignments[uid]
		userName := h.getUserDisplayName(uid) // Call method on h

		b.WriteString(fmt.Sprintf("%s: %s\n", userName, role.Name))
	}
	return b.String()
}
