package telegram

import (
	"fmt"
	"strings"

	"gopkg.in/telebot.v3"
)

// Unique identifiers for inline buttons
const (
	UniqueJoinToRoom         = "join"
	UniqueKickFromRoom       = "kick"
	UniqueKickPlayerFromRoom = "kick_player"
	UniqueDeleteRoom         = "delete_room"
)

func (h *BotHandler) HandleCallback(c telebot.Context) error {
	// Extract room ID and player ID from callback data
	parts := strings.Split(c.Callback().Data, "|")
	if len(parts) != 2 {
		return c.Respond(&telebot.CallbackResponse{
			Text: "Invalid data format.",
		})
	}
	unique := strings.TrimSpace(parts[0])
	data := parts[1]
	fmt.Println(unique)
	fmt.Println(data)
	if unique == UniqueJoinToRoom {
		return h.HandleJoinRoomCallback(c, data)
	} else if unique == UniqueKickFromRoom {
		return h.HandleKickUserCallback(c, data)
	} else if unique == UniqueKickPlayerFromRoom {
		return h.HandleKickUserFromRoomCallback(c, data)
	} else if unique == UniqueDeleteRoom {
		return h.HandleDeleteRoomCallback(c, data)
	}
	fmt.Println("button command not found")
	return c.Respond()
}
