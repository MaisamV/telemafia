package telegram

import (
	"context"
	"fmt"
	"log"

	roomQuery "telemafia/internal/domain/room/usecase/query"

	"gopkg.in/telebot.v3"
)

// HandleDeleteRoom handles the /delete_room command (now a function)
func HandleDeleteRoom(h *BotHandler, c telebot.Context) error {
	rooms, err := h.getRoomsHandler.Handle(context.Background(), roomQuery.GetRoomsQuery{})
	if err != nil {
		log.Printf("Error fetching rooms for deletion selection: %v", err)
		return c.Send("Failed to fetch rooms list.")
	}

	if len(rooms) == 0 {
		return c.Send("No rooms exist to delete.")
	}

	var buttons [][]telebot.InlineButton
	for _, room := range rooms {
		button := telebot.InlineButton{
			Unique: UniqueDeleteRoomSelectRoom,
			Text:   fmt.Sprintf("%s (%s)", room.Name, room.ID),
			Data:   string(room.ID),
		}
		buttons = append(buttons, []telebot.InlineButton{button})
	}

	markup := &telebot.ReplyMarkup{InlineKeyboard: buttons}
	return c.Send("Select a room to delete:", markup)
}
