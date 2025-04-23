package telegram

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	roomEntity "telemafia/internal/domain/room/entity"
	roomCommand "telemafia/internal/domain/room/usecase/command"
	sharedEntity "telemafia/internal/shared/entity"

	"gopkg.in/telebot.v3"
)

// HandleKickUser handles the /kick_user command (now a function)
func HandleKickUser(h *BotHandler, c telebot.Context) error {
	args := strings.Fields(c.Message().Payload)
	if len(args) != 2 {
		return c.Send("Usage: /kick_user <room_id> <user_id>")
	}

	roomID := roomEntity.RoomID(args[0])
	playerIDStr := args[1]

	playerID, err := strconv.ParseInt(playerIDStr, 10, 64)
	if err != nil {
		return c.Send("Invalid user ID format.")
	}

	requester := ToUser(c.Sender())
	if requester == nil {
		return c.Send("Could not identify user.")
	}

	cmd := roomCommand.KickUserCommand{
		Requester: *requester,
		RoomID:    roomID,
		PlayerID:  sharedEntity.UserID(playerID),
	}

	if err := h.kickUserHandler.Handle(context.Background(), cmd); err != nil {
		return c.Send(fmt.Sprintf("Error kicking user %d from room %s: %v", playerID, roomID, err))
	}

	return c.Send(fmt.Sprintf("User %d kicked from room %s", playerID, roomID))
}
