package telegram

import (
	"context"
	"fmt"
	"strings"
	"telemafia/delivery/common"
	roomQuery "telemafia/internal/room/usecase/query"

	"gopkg.in/telebot.v3"
)

// HandleMyRooms handles the /my_rooms command
func (h *BotHandler) HandleMyRooms(c telebot.Context) error {
	user := common.ToUser(c.Sender())

	rooms, err := h.getPlayerRoomsHandler.Handle(context.Background(), roomQuery.GetPlayerRoomsQuery{
		PlayerID: user.ID,
	})
	if err != nil {
		return c.Send(fmt.Sprintf("Error getting your rooms: %v", err))
	}

	if len(rooms) == 0 {
		return c.Send("You haven't joined any rooms")
	}

	var sb strings.Builder
	sb.WriteString("Your rooms:\n")
	for _, room := range rooms {
		sb.WriteString(fmt.Sprintf("- %s (ID: %s, Players: %d)\n", room.Name, room.ID, len(room.Players)))
	}

	return c.Send(sb.String())
}
