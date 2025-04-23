package telegram

import (
	"log"

	"gopkg.in/telebot.v3"
)

// Unique identifiers for inline buttons
const (
	// Join/Leave related
	UniqueJoinRoom            = "join_room"
	UniqueLeaveRoomSelectRoom = "leave_room_select"
	UniqueLeaveRoomConfirm    = "leave_room_confirm"

	// Delete Room related
	UniqueDeleteRoomSelectRoom = "delete_room_select"
	UniqueDeleteRoomConfirm    = "delete_room_confirm"

	// Game/Assignment related
	UniqueConfirmAssignments = "confirm_assignments"
	// UniqueShowMyRole          = "show_my_role" // Placeholder if needed later

	// Generic Cancel (might need context)
	UniqueCancel = "cancel"
)

// HandleCallback routes all callback queries from inline buttons (now a method)
func (h *BotHandler) handleCallback(c telebot.Context) error {
	callback := c.Callback()
	if callback == nil {
		log.Println("Received update that is not a callback")
		return nil
	}

	// SplitCallbackData defined in util.go (same package)
	unique, data := SplitCallbackData(callback.Data)
	userID := c.Sender().ID
	log.Printf("Callback received: User=%d, Unique=%s, Data=%s", userID, unique, data)

	switch unique {
	// Room Callbacks - Call functions, passing h
	case UniqueJoinRoom:
		return handleJoinRoomCallback(h, c, data)
	case UniqueDeleteRoomSelectRoom:
		return handleDeleteRoomSelectCallback(h, c, data)
	case UniqueDeleteRoomConfirm:
		return handleDeleteRoomConfirmCallback(h, c, data)
	case UniqueLeaveRoomSelectRoom:
		return handleLeaveRoomSelectCallback(h, c, data)
	case UniqueLeaveRoomConfirm:
		return handleLeaveRoomConfirmCallback(h, c, data)

	// Game Callbacks - Call functions, passing h
	case UniqueConfirmAssignments:
		return HandleConfirmAssignments(h, c, data)
	// case UniqueShowMyRole:
	// 	 return handleShowMyRoleCallback(h, c, data) // Implement if needed

	// General Callbacks
	case UniqueCancel:
		log.Printf("User %d cancelled operation.", userID)
		_ = c.Respond(&telebot.CallbackResponse{Text: "Operation cancelled."})
		_ = c.Delete()
		return nil

	default:
		log.Printf("Unknown callback unique identifier: %s", unique)
		return c.Respond(&telebot.CallbackResponse{Text: "Unknown action."})
	}
}
