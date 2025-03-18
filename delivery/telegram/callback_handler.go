package telegram

import (
	"log"
	"strings"

	"gopkg.in/telebot.v3"
)

// Unique identifiers for inline buttons
const (
	UniqueJoinSelectRoom           = "join_selectRoom"
	UniqueKickSelectRoom           = "kick_selectRoom"
	UniqueKickFromRoomSelectPlayer = "kickFromRoom_selectPlayer"
	UniqueDeleteRoomSelectRoom     = "deleteRoom_selectRoom"
	UniqueLeaveRoomSelectRoom      = "leaveRoom_selectRoom"
	UniqueConfirm                  = "confirm_assignments"
)

// HandleCallback handles all callbacks from inline buttons
func (h *BotHandler) HandleCallback(c telebot.Context) error {
	// Extract room ID and player ID from callback data
	data := c.Callback().Data

	log.Printf("Received callback: %s", data)

	// Process other callbacks with the original format
	parts := strings.Split(data, "|")
	if len(parts) != 2 {
		log.Printf("Invalid callback data format: %s", data)
		return c.Respond(&telebot.CallbackResponse{
			Text: "Invalid data format.",
		})
	}

	uniqueFromData := strings.TrimSpace(parts[0])
	callbackData := parts[1]

	log.Printf("Parsed callback: unique=%s, data=%s", uniqueFromData, callbackData)

	if uniqueFromData == UniqueConfirm {
		log.Printf("Routing to HandleConfirmAssignments with data: %s", data)
		return h.HandleConfirmAssignments(c, callbackData)
	} else if uniqueFromData == UniqueJoinSelectRoom {
		return h.HandleJoinRoomCallback(c, callbackData)
	} else if uniqueFromData == UniqueKickSelectRoom {
		return h.HandleKickUserCallback(c, callbackData)
	} else if uniqueFromData == UniqueKickFromRoomSelectPlayer {
		return h.HandleKickUserFromRoomCallback(c, callbackData)
	} else if uniqueFromData == UniqueDeleteRoomSelectRoom {
		return h.HandleDeleteRoomCallback(c, callbackData)
	} else if uniqueFromData == UniqueLeaveRoomSelectRoom {
		return h.HandleLeaveRoomCallback(c, callbackData)
	}

	log.Printf("Button command not found: %s", uniqueFromData)
	return c.Respond(&telebot.CallbackResponse{Text: "Unknown callback!"})
}
