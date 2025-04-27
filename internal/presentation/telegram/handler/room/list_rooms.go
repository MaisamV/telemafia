package telegram

import (
	"context"
	"fmt"
	"strings"

	roomQuery "telemafia/internal/domain/room/usecase/query"
	messages "telemafia/internal/presentation/telegram/messages"
	tgutil "telemafia/internal/shared/tgutil"

	"gopkg.in/telebot.v3"
)

// PrepareRoomListMessage fetches rooms and generates the message text and inline button markup.
func PrepareRoomListMessage(
	getRoomsHandler *roomQuery.GetRoomsHandler,
	getPlayersInRoomHandler *roomQuery.GetPlayersInRoomHandler,
	msgs *messages.Messages,
) (string, *telebot.ReplyMarkup, error) {
	rooms, err := getRoomsHandler.Handle(context.Background(), roomQuery.GetRoomsQuery{})
	if err != nil {
		return "", nil, fmt.Errorf(msgs.Room.ListError, err)
	}

	var response strings.Builder
	markup := &telebot.ReplyMarkup{}
	var rows []telebot.Row

	if len(rooms) == 0 {
		response.WriteString(msgs.Room.ListNoRooms)
	} else {
		response.WriteString(msgs.Room.ListTitle)
		for _, room := range rooms {
			players, _ := getPlayersInRoomHandler.Handle(context.Background(), roomQuery.GetPlayersInRoomQuery{RoomID: room.ID})
			playerCount := len(players)

			btnText := fmt.Sprintf(msgs.Room.JoinButtonText, room.Name, playerCount)
			btnJoin := markup.Data(btnText, tgutil.UniqueJoinRoom, string(room.ID))
			rows = append(rows, markup.Row(btnJoin))
		}
	}
	markup.Inline(rows...)
	return response.String(), markup, nil
}

// HandleListRooms handles the /list_rooms command using the new message preparation function.
func HandleListRooms(
	getRoomsHandler *roomQuery.GetRoomsHandler,
	getPlayersInRoomHandler *roomQuery.GetPlayersInRoomHandler,
	bot *telebot.Bot,
	refreshingMessage *tgutil.RefreshingMessageBook,
	c telebot.Context,
	msgs *messages.Messages,
) error {
	chatID := telebot.ChatID(c.Sender().ID)
	content, markup, err := PrepareRoomListMessage(getRoomsHandler, getPlayersInRoomHandler, msgs)
	if err != nil {
		return c.Send(fmt.Sprintf(msgs.Room.ListErrorPrepare, err))
	}

	message, err := bot.Send(chatID, content, markup)
	if err == nil {
		activeMessage, exists := refreshingMessage.GetActiveMessage(c.Sender().ID)
		if exists {
			_ = bot.Delete(activeMessage.Msg)
		}
		refreshingMessage.AddActiveMessage(c.Sender().ID, &tgutil.RefreshingMessage{
			Msg:  message,
			Data: "",
		})
	}
	return err
}
