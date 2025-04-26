package telegram

import (
	"context"
	"fmt"
	"strings"
	"time"

	roomEntity "telemafia/internal/domain/room/entity"
	roomCommand "telemafia/internal/domain/room/usecase/command"
	tgutil "telemafia/internal/shared/tgutil"

	"gopkg.in/telebot.v3"
)

// RefreshNotifier defines an interface for triggering a refresh.
// *tgutil.RefreshState satisfies this interface.
type RefreshNotifier interface {
	RaiseRefreshNeeded()
}

// HandleCreateRoom handles the /create_room command (now a function)
func HandleCreateRoom(createRoomHandler *roomCommand.CreateRoomHandler, refreshNotifier RefreshNotifier, c telebot.Context) error {
	args := strings.TrimSpace(c.Message().Payload)
	if args == "" {
		return c.Send("Please provide a room name: /create_room [name]")
	}

	user := tgutil.ToUser(c.Sender())
	if user == nil {
		return c.Send("Could not identify user.")
	}

	cmd := roomCommand.CreateRoomCommand{
		ID:        roomEntity.RoomID(fmt.Sprintf("room_%d", time.Now().UnixNano())), // Generate unique ID
		Name:      args,
		CreatorID: user.ID, // Pass CreatorID
	}

	createdRoom, err := createRoomHandler.Handle(context.Background(), cmd)
	if err != nil {
		return c.Send(fmt.Sprintf("Error creating room: %v", err))
	}

	refreshNotifier.RaiseRefreshNeeded() // Raise flag on success
	return c.Send(fmt.Sprintf("Room '%s' created successfully! ID: %s", createdRoom.Name, createdRoom.ID))
}
