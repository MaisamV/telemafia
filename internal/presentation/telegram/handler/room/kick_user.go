package telegram

import (
	"context"
	"fmt"
	"strings"

	roomEntity "telemafia/internal/domain/room/entity"
	roomCommand "telemafia/internal/domain/room/usecase/command"
	messages "telemafia/internal/presentation/telegram/messages"
	"telemafia/internal/shared/common"
	sharedEntity "telemafia/internal/shared/entity"
	tgutil "telemafia/internal/shared/tgutil"

	"gopkg.in/telebot.v4"
)

// RefreshNotifier is defined in create_room.go (same package)

// HandleKickUser handles the /kick_user command.
func HandleKickUser(
	kickUserHandler *roomCommand.KickUserHandler,
	refreshNotifier RefreshNotifier,
	c telebot.Context,
	msgs *messages.Messages,
) error {
	args := strings.Fields(c.Message().Payload)
	if len(args) != 2 {
		return c.Send(msgs.Room.KickPrompt)
	}
	roomIDStr := args[0]
	userIDStr := args[1]

	userID, err := common.StringToInt64(userIDStr)
	if err != nil {
		return c.Send(msgs.Room.KickInvalidUserID)
	}

	requester := tgutil.ToUser(c.Sender())
	if requester == nil {
		return c.Send(msgs.Common.ErrorIdentifyUser)
	}

	cmd := roomCommand.KickUserCommand{
		Requester: *requester,
		RoomID:    roomEntity.RoomID(roomIDStr),
		PlayerID:  sharedEntity.UserID(userID),
	}

	if err := kickUserHandler.Handle(context.Background(), cmd); err != nil {
		return c.Send(fmt.Sprintf(msgs.Room.KickError, userID, roomIDStr, err))
	}

	refreshNotifier.RaiseRefreshNeeded()
	return c.Send(fmt.Sprintf(msgs.Room.KickSuccess, userID, roomIDStr))
}
