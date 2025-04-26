package telegram

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	roomEntity "telemafia/internal/domain/room/entity"
	roomCommand "telemafia/internal/domain/room/usecase/command"
	sharedEntity "telemafia/internal/shared/entity"
	tgutil "telemafia/internal/shared/tgutil"

	"gopkg.in/telebot.v3"
)

// RefreshNotifier is defined in create_room.go (same package)

// HandleKickUser handles the /kick_user command (now a function)
func HandleKickUser(kickUserHandler *roomCommand.KickUserHandler, refreshNotifier RefreshNotifier, c telebot.Context) error {
	parts := strings.Fields(c.Message().Payload)
	if len(parts) != 2 {
		return c.Send("Usage: /kick_user <room_id> <user_id>")
	}

	roomID := roomEntity.RoomID(parts[0])
	playerIDStr := parts[1]

	playerID, err := strconv.ParseInt(playerIDStr, 10, 64)
	if err != nil {
		return c.Send("Invalid user ID format.")
	}

	requester := tgutil.ToUser(c.Sender())
	if requester == nil {
		return c.Send("Could not identify requester.")
	}

	cmd := roomCommand.KickUserCommand{
		Requester: *requester,
		RoomID:    roomID,
		PlayerID:  sharedEntity.UserID(playerID),
	}

	if err := kickUserHandler.Handle(context.Background(), cmd); err != nil {
		return c.Send(fmt.Sprintf("Error kicking user %d from room %s: %v", playerID, roomID, err))
	}

	refreshNotifier.RaiseRefreshNeeded()
	return c.Send(fmt.Sprintf("User %d kicked from room %s", playerID, roomID))
}
