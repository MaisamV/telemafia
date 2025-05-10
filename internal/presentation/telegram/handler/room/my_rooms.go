package telegram

import (
	"context"
	"fmt"
	"strings"

	roomQuery "telemafia/internal/domain/room/usecase/query"
	messages "telemafia/internal/presentation/telegram/messages"
	tgutil "telemafia/internal/shared/tgutil"

	"gopkg.in/telebot.v4"
)

// HandleMyRooms handles the /my_rooms command.
func HandleMyRooms(
	getPlayerRoomsHandler *roomQuery.GetPlayerRoomsHandler,
	c telebot.Context,
	msgs *messages.Messages,
) error {
	user := tgutil.ToUser(c.Sender())
	if user == nil {
		return c.Send(msgs.Common.ErrorIdentifyUser)
	}

	query := roomQuery.GetPlayerRoomsQuery{PlayerID: user.ID}
	rooms, err := getPlayerRoomsHandler.Handle(context.Background(), query)
	if err != nil {
		return c.Send(fmt.Sprintf(msgs.Room.MyRoomsError, err))
	}

	if len(rooms) == 0 {
		return c.Send(msgs.Room.MyRoomsNone)
	}

	var response strings.Builder
	for _, room := range rooms {
		response.WriteString(fmt.Sprintf(msgs.Room.MyRoomsTitle, room.Name, room.ID))
	}

	return c.Send(response.String())
}
