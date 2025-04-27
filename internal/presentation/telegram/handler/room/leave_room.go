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

// HandleLeaveRoom handles the /leave_room command (now a function)
func HandleLeaveRoom(
	leaveRoomHandler *roomCommand.LeaveRoomHandler,
	refreshNotifier RefreshNotifier,
	c telebot.Context,
	msgs *messages.Messages,
) error {
	roomIDStr := strings.TrimSpace(c.Message().Payload)
	if roomIDStr == "" {
		return c.Send(msgs.Room.LeavePrompt)
	}

	user := tgutil.ToUser(c.Sender())
	if user == nil {
		return c.Send(msgs.Common.ErrorIdentifyUser)
	}
	cmd := roomCommand.LeaveRoomCommand{
		RoomID:    roomEntity.RoomID(roomIDStr),
		Requester: *user,
	}
	if err := leaveRoomHandler.Handle(context.Background(), cmd); err != nil {
		return c.Send(fmt.Sprintf(msgs.Room.LeaveError, roomIDStr, err))
	}

	refreshNotifier.RaiseRefreshNeeded()
	return c.Send(fmt.Sprintf(msgs.Room.LeaveSuccess, roomIDStr))
}
