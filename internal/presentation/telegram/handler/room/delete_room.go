package telegram

import (
	"context"
	"fmt"

	roomQuery "telemafia/internal/domain/room/usecase/query"
	tgutil "telemafia/internal/shared/tgutil"

	"gopkg.in/telebot.v3"
)

// HandleDeleteRoom handles the /delete_room command (now a function)
func HandleDeleteRoom(getRoomsHandler *roomQuery.GetRoomsHandler, c telebot.Context) error {
	requester := tgutil.ToUser(c.Sender())
	if requester == nil {
		return c.Send("Could not identify requester.")
	}

	query := roomQuery.GetRoomsQuery{}
	rooms, err := getRoomsHandler.Handle(context.Background(), query)
	if err != nil {
		return c.Send(fmt.Sprintf("Failed to fetch rooms list: %v", err))
	}

	if len(rooms) == 0 {
		return c.Send("No rooms exist to delete.")
	}

	var buttons [][]telebot.InlineButton
	for _, room := range rooms {
		btn := telebot.InlineButton{
			Unique: tgutil.UniqueDeleteRoomSelectRoom,
			Text:   fmt.Sprintf("%s (%s)", room.Name, room.ID),
			Data:   string(room.ID),
		}
		buttons = append(buttons, []telebot.InlineButton{btn})
	}

	markup := &telebot.ReplyMarkup{InlineKeyboard: buttons}
	return c.Send("Select a room to delete:", markup)
}
