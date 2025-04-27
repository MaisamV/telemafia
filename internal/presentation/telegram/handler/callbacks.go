package telegram

import (
	"fmt"
	"log"
	game "telemafia/internal/presentation/telegram/handler/game"
	room "telemafia/internal/presentation/telegram/handler/room"

	// Import messages
	tgutil "telemafia/internal/shared/tgutil"

	"gopkg.in/telebot.v3"
	// Import necessary command/query handlers for callback routing
)

// HandleCallback routes all callback queries from inline buttons (now a method)
func (h *BotHandler) handleCallback(c telebot.Context) error {
	callback := c.Callback()
	if callback == nil {
		// This shouldn't happen for telebot.OnCallback, maybe log differently
		log.Println("Received update that is not a callback")
		return nil
	}

	// SplitCallbackData defined in util.go (same package)
	unique, data := tgutil.SplitCallbackData(callback.Data)
	userID := c.Sender().ID
	log.Printf("Callback received: User=%d, Unique=%s, Data=%s", userID, unique, data)

	switch unique {
	// Room Callbacks - Pass messages struct
	case tgutil.UniqueJoinRoom:
		return room.HandleJoinRoomCallback(h.joinRoomHandler, h.getRoomsHandler, h.getPlayersInRoomHandler, h.roomListRefreshMessage, h.roomDetailRefreshMessage, c, data, h.msgs)
	case tgutil.UniqueDeleteRoomSelectRoom:
		return room.HandleDeleteRoomSelectCallback(h.getRoomHandler, c, data, h.msgs)
	case tgutil.UniqueDeleteRoomConfirm:
		return room.HandleDeleteRoomConfirmCallback(h.deleteRoomHandler, h.roomListRefreshMessage, c, data, h.msgs)
	case tgutil.UniqueLeaveRoomSelectRoom:
		return room.HandleLeaveRoomSelectCallback(h.leaveRoomHandler, h.getRoomsHandler, h.getPlayersInRoomHandler, h.roomListRefreshMessage, h.roomDetailRefreshMessage, c, data, h.msgs)
	case tgutil.UniqueLeaveRoomConfirm:
		return room.HandleLeaveRoomConfirmCallback(h.leaveRoomHandler, h.roomListRefreshMessage, c, data, h.msgs)

	// Game Callbacks - Pass messages struct
	case tgutil.UniqueConfirmAssignments:
		return game.HandleConfirmAssignments(h.getGameByIDHandler, c, data, h.msgs)
	// case UniqueShowMyRole:
	// 	 return handleShowMyRoleCallback(h, c, data) // Implement if needed

	// General Callbacks
	case tgutil.UniqueCancel:
		// Use msg
		_ = c.Respond(&telebot.CallbackResponse{Text: h.msgs.Common.CallbackCancelled})
		return c.Delete()

	default:
		log.Printf("Unknown callback unique identifier: %s", unique)
		// Use msg
		return c.Respond(&telebot.CallbackResponse{Text: fmt.Sprintf("Unknown action: %s", unique), ShowAlert: true})
	}
}
