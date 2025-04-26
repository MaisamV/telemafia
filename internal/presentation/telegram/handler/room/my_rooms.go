package telegram

import (
	"context"
	"fmt"
	"strings"

	roomQuery "telemafia/internal/domain/room/usecase/query"
	tgutil "telemafia/internal/shared/tgutil"

	"gopkg.in/telebot.v3"
)

// HandleMyRooms handles the /my_rooms command (now a function)
func HandleMyRooms(getPlayerRoomsHandler *roomQuery.GetPlayerRoomsHandler, c telebot.Context) error {
	user := tgutil.ToUser(c.Sender())
	if user == nil {
		return c.Send("Could not identify user.")
	}
	query := roomQuery.GetPlayerRoomsQuery{PlayerID: user.ID}
	rooms, err := getPlayerRoomsHandler.Handle(context.Background(), query)
	if err != nil {
		return c.Send(fmt.Sprintf("Error getting your rooms: %v", err))
	}

	if len(rooms) == 0 {
		return c.Send("You are not in any rooms.")
	}

	var response strings.Builder
	response.WriteString("Rooms you are in:\n")
	for _, room := range rooms {
		response.WriteString(fmt.Sprintf("- %s (%s)\n", room.Name, room.ID))
	}

	return c.Send(response.String())
}
