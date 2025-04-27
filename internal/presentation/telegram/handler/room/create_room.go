package telegram

import (
	"context"
	"fmt"
	"strings"
	"time"

	roomEntity "telemafia/internal/domain/room/entity"
	roomCommand "telemafia/internal/domain/room/usecase/command"
	messages "telemafia/internal/presentation/telegram/messages"
	tgutil "telemafia/internal/shared/tgutil"

	"gopkg.in/telebot.v3"
)

// HandleCreateRoom handles the /create_room command (now a function)
func HandleCreateRoom(
	createRoomHandler *roomCommand.CreateRoomHandler,
	refreshNotifier RefreshNotifier,
	c telebot.Context,
	msgs *messages.Messages,
) error {
	args := strings.TrimSpace(c.Message().Payload)
	if args == "" {
		return c.Send(msgs.Room.CreatePrompt)
	}

	user := tgutil.ToUser(c.Sender())
	if user == nil {
		return c.Send(msgs.Common.ErrorIdentifyUser)
	}

	cmd := roomCommand.CreateRoomCommand{
		ID:        roomEntity.RoomID(fmt.Sprintf("room_%d", time.Now().UnixNano())), // Generate unique ID
		Name:      args,
		CreatorID: user.ID, // Pass CreatorID
	}

	createdRoom, err := createRoomHandler.Handle(context.Background(), cmd)
	if err != nil {
		return c.Send(fmt.Sprintf(msgs.Room.CreateError, err))
	}

	refreshNotifier.RaiseRefreshNeeded() // Raise flag on success
	return c.Send(fmt.Sprintf(msgs.Room.CreateSuccess, createdRoom.Name, createdRoom.ID))
}
