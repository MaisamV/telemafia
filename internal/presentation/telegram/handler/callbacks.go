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

	rawData := callback.Data
	unique, data := tgutil.SplitCallbackData(rawData)
	userID := c.Sender().ID
	log.Printf("Callback received: User=%d, Data=%s", userID, rawData)
	log.Printf("unique=%s, data=%s", unique, data)

	switch unique {
	// Game Creation Callbacks
	case tgutil.UniqueCreateGameSelectRoom:
		return game.HandleSelectRoomForCreateGame(h.getAllScenariosHandler, c, data, h.msgs)
	case tgutil.UniqueCreateGameSelectScenario:
		roomID, scenarioID := tgutil.SplitCallbackData(data)
		if roomID == "" || scenarioID == "" {
			log.Printf("Invalid creategame_scen callback data: %s", rawData)
			return c.Respond(&telebot.CallbackResponse{Text: "Invalid callback data.", ShowAlert: true})
		}
		return game.HandleSelectScenarioForCreateGame(h.createGameHandler, h.getPlayersInRoomHandler, h.getScenarioByIDHandler, c, roomID, scenarioID, h.msgs)
	case tgutil.UniqueStartGame:
		return game.HandleStartCreatedGame(h.assignRolesHandler, h.bot, c, data, h.msgs)
	case tgutil.UniqueCancelGame:
		return game.HandleCancelCreateGame(c, h.msgs, data)

	// Existing Room Callbacks (assuming tgutil still defines these constants)
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
	case tgutil.UniqueGetInviteLink:
		return room.HandleGetInviteLinkCallback(h.bot, c, data, h.msgs)

	// Kick User Flow Callbacks
	case tgutil.UniqueKickUserSelect:
		return room.HandleKickUserSelectCallback(h.getPlayersInRoomHandler, c, data, h.msgs)
	case tgutil.UniqueKickUserConfirm:
		return room.HandleKickUserConfirmCallback(h.kickUserHandler, h.getRoomsHandler, h.getPlayersInRoomHandler, h.roomListRefreshMessage, h.roomDetailRefreshMessage, c, data, h.msgs)

	// Existing Game Callbacks
	case tgutil.UniqueConfirmAssignments:
		return game.HandleConfirmAssignments(h.getGameByIDHandler, c, data, h.msgs)

	// Existing General Callbacks
	case tgutil.UniqueCancel:
		_ = c.Respond(&telebot.CallbackResponse{Text: h.msgs.Common.CallbackCancelled})
		return c.Delete()

	default:
		log.Printf("Unknown callback data format: %s", rawData)
		return c.Respond(&telebot.CallbackResponse{Text: fmt.Sprintf("Unknown action: %s", rawData), ShowAlert: true})
	}
}
