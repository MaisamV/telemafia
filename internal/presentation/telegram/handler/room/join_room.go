package telegram

import (
	"context"
	"fmt"
	"log"
	roomEntity "telemafia/internal/domain/room/entity"
	roomCommand "telemafia/internal/domain/room/usecase/command"
	roomQuery "telemafia/internal/domain/room/usecase/query"
	messages "telemafia/internal/presentation/telegram/messages"
	tgutil "telemafia/internal/shared/tgutil"

	"gopkg.in/telebot.v3"
)

// RefreshNotifier is defined in create_room.go (same package)

// HandleJoinRoom handles the /join_room command (now a function)
func HandleJoinRoom(
	joinRoomHandler *roomCommand.JoinRoomHandler,
	getRoomsHandler *roomQuery.GetRoomsHandler,
	getPlayersInRoomHandler *roomQuery.GetPlayersInRoomHandler,
	roomList RefreshNotifier,
	roomDetail RefreshNotifier,
	c telebot.Context,
	data string,
	msgs *messages.Messages,
) error {
	if data == "" {
		return c.Send(msgs.Room.JoinPrompt)
	}
	roomID := roomEntity.RoomID(data)
	user := tgutil.ToUser(c.Sender())
	if user == nil {
		return c.Respond(&telebot.CallbackResponse{Text: msgs.Common.ErrorIdentifyUser, ShowAlert: true})
	}

	cmd := roomCommand.JoinRoomCommand{
		Requester: *user,
		RoomID:    roomID,
	}

	if err := joinRoomHandler.Handle(context.Background(), cmd); err != nil {
		log.Printf("Error handling join room callback for room '%s': %v", roomID, err)
		return c.Respond(&telebot.CallbackResponse{Text: fmt.Sprintf(msgs.Common.CallbackErrorGeneric, err), ShowAlert: true})
	}
	chatID := c.Sender().ID
	listMessage, listExists := roomList.GetActiveMessage(chatID)
	roomMessage, roomExists := roomList.GetActiveMessage(chatID)
	roomList.RemoveActiveMessage(chatID)
	roomDetail.AddActiveMessage(chatID, &tgutil.RefreshingMessage{
		Msg:  c.Message(),
		Data: string(roomID),
	})
	roomList.RaiseRefreshNeeded()
	roomDetail.RaiseRefreshNeeded()
	message, markup, err := RoomDetailMessage(getRoomsHandler, getPlayersInRoomHandler, msgs, data)
	if err != nil {
		return err
	}
	msg, err := c.Bot().Send(c.Sender(), message, markup, telebot.NoPreview)
	roomDetail.AddActiveMessage(chatID, &tgutil.RefreshingMessage{
		Msg:  msg,
		Data: string(roomID),
	})
	if listExists {
		_ = c.Bot().Delete(listMessage.Msg)
	}
	if roomExists {
		_ = c.Bot().Delete(roomMessage.Msg)
	}
	return err
}
