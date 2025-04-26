package telegram

import (
	"context"
	"fmt"
	"strings"

	roomQuery "telemafia/internal/domain/room/usecase/query"

	"gopkg.in/telebot.v3"
)

// HandleListRooms handles the /list_rooms command (now a function)
func HandleListRooms(getRoomsHandler *roomQuery.GetRoomsHandler, getPlayersInRoomHandler *roomQuery.GetPlayersInRoomHandler, c telebot.Context) error {
	query := roomQuery.GetRoomsQuery{}
	rooms, err := getRoomsHandler.Handle(context.Background(), query)
	if err != nil {
		return c.Send(fmt.Sprintf("Error getting rooms: %v", err))
	}

	if len(rooms) == 0 {
		return c.Send("No rooms available.")
	}

	var response strings.Builder
	response.WriteString("Available Rooms:\n")
	for _, room := range rooms {
		players, _ := getPlayersInRoomHandler.Handle(context.Background(), roomQuery.GetPlayersInRoomQuery{RoomID: room.ID})
		playerCount := len(players)
		maxPlayers := 10
		response.WriteString(fmt.Sprintf("- %s (%s) [%d/%d players]\n", room.Name, room.ID, playerCount, maxPlayers))
	}

	return c.Send(response.String())
}
