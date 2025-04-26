package telegram

import (
	"log"
	game "telemafia/internal/presentation/telegram/handler/game"
	room "telemafia/internal/presentation/telegram/handler/room"

	tgutil "telemafia/internal/shared/tgutil"

	"gopkg.in/telebot.v3"
	// Import necessary command/query handlers for callback routing
)

// HandleCallback routes all callback queries from inline buttons (now a method)
func (h *BotHandler) handleCallback(c telebot.Context) error {
	callback := c.Callback()
	if callback == nil {
		log.Println("Received update that is not a callback")
		return nil
	}

	// SplitCallbackData defined in util.go (same package)
	unique, data := tgutil.SplitCallbackData(callback.Data)
	userID := c.Sender().ID
	log.Printf("Callback received: User=%d, Unique=%s, Data=%s", userID, unique, data)

	switch unique {
	// Room Callbacks - Call functions, passing required handlers
	case tgutil.UniqueJoinRoom:
		return room.HandleJoinRoomCallback(h.joinRoomHandler, c, data)
	case tgutil.UniqueDeleteRoomSelectRoom:
		return room.HandleDeleteRoomSelectCallback(h.getRoomHandler, c, data)
	case tgutil.UniqueDeleteRoomConfirm:
		return room.HandleDeleteRoomConfirmCallback(h.deleteRoomHandler, c, data)
	case tgutil.UniqueLeaveRoomSelectRoom:
		// This callback likely needs access to GetPlayerRoomsQuery
		return room.HandleLeaveRoomSelectCallback(h.getPlayerRoomsHandler, c, data)
	case tgutil.UniqueLeaveRoomConfirm:
		return room.HandleLeaveRoomConfirmCallback(h.leaveRoomHandler, c, data)

	// Game Callbacks - Call functions, passing required handlers
	case tgutil.UniqueConfirmAssignments:
		return game.HandleConfirmAssignments(h.getGameByIDHandler, c, data)
	// case UniqueShowMyRole:
	// 	 return handleShowMyRoleCallback(h, c, data) // Implement if needed

	// General Callbacks
	case tgutil.UniqueCancel:
		log.Printf("User %d cancelled operation.", userID)
		_ = c.Respond(&telebot.CallbackResponse{Text: "Operation cancelled."})
		_ = c.Delete()
		return nil

	default:
		log.Printf("Unknown callback unique identifier: %s", unique)
		return c.Respond(&telebot.CallbackResponse{Text: "Unknown action."})
	}
}
