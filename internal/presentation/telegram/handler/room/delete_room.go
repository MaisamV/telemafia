package telegram

import (
	"context"
	"fmt"

	roomQuery "telemafia/internal/domain/room/usecase/query"
	messages "telemafia/internal/presentation/telegram/messages"
	tgutil "telemafia/internal/shared/tgutil"

	"gopkg.in/telebot.v4"
)

// HandleDeleteRoom handles the first step of /delete_room, showing the selection.
func HandleDeleteRoom(
	getRoomsHandler *roomQuery.GetRoomsHandler,
	c telebot.Context,
	msgs *messages.Messages,
) error {
	user := tgutil.ToUser(c.Sender())
	if user == nil || !user.Admin {
		return c.Send(msgs.Common.ErrorPermissionDenied)
	}

	rooms, err := getRoomsHandler.Handle(context.Background(), roomQuery.GetRoomsQuery{})
	if err != nil {
		return c.Send(msgs.Room.DeleteErrorFetch)
	}

	if len(rooms) == 0 {
		return c.Send(msgs.Room.DeleteNoRooms)
	}

	markup := &telebot.ReplyMarkup{}
	var rows []telebot.Row
	for _, room := range rooms {
		btn := markup.Data(fmt.Sprintf("%s (%s)", room.Name, room.ID), tgutil.UniqueDeleteRoomSelectRoom, string(room.ID))
		rows = append(rows, markup.Row(btn))
	}
	markup.Inline(rows...)

	return c.Send(msgs.Room.DeletePromptSelect, markup)
}
