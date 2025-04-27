package telegram

import (
	"context"
	"fmt"
	"strings"
	roomEntity "telemafia/internal/domain/room/entity"
	roomCommand "telemafia/internal/domain/room/usecase/command"
	messages "telemafia/internal/presentation/telegram/messages"
	tgutil "telemafia/internal/shared/tgutil"

	"gopkg.in/telebot.v3"
)

// RefreshNotifier is defined in create_room.go (same package)

// HandleJoinRoom handles the /join_room command (now a function)
func HandleJoinRoom(
	joinRoomHandler *roomCommand.JoinRoomHandler,
	roomListRefreshNotifier RefreshNotifier,
	roomDetailRefreshNotifier RefreshNotifier,
	c telebot.Context,
	msgs *messages.Messages,
) error {
	roomIDStr := strings.TrimSpace(c.Message().Payload)
	if roomIDStr == "" {
		return c.Send(msgs.Room.JoinPrompt)
	}

	roomID := roomEntity.RoomID(roomIDStr)
	user := tgutil.ToUser(c.Sender())
	if user == nil {
		return c.Send(msgs.Common.ErrorIdentifyUser)
	}

	cmd := roomCommand.JoinRoomCommand{
		Requester: *user,
		RoomID:    roomID,
	}

	err := joinRoomHandler.Handle(context.Background(), cmd)
	if err != nil {
		return c.Send(fmt.Sprintf(msgs.Room.JoinError, roomID, err))
	}

	roomListRefreshNotifier.RaiseRefreshNeeded()
	roomDetailRefreshNotifier.RaiseRefreshNeeded()

	markup := &telebot.ReplyMarkup{}
	btnLeave := markup.Data(msgs.Room.LeaveConfirmButton, tgutil.UniqueLeaveRoomSelectRoom, string(roomID))
	markup.Inline(markup.Row(btnLeave))

	return c.Send(fmt.Sprintf(msgs.Room.JoinSuccess, roomID), markup)
}
