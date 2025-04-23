package telegram

import (
	"context"
	"fmt"
	"strings"

	roomEntity "telemafia/internal/domain/room/entity"
	roomCommand "telemafia/internal/domain/room/usecase/command"

	"gopkg.in/telebot.v3"
)

// HandleLeaveRoom handles the /leave_room command (now a function)
func HandleLeaveRoom(h *BotHandler, c telebot.Context) error {
	args := strings.TrimSpace(c.Message().Payload)
	if args == "" {
		return c.Send("Please provide a room ID: /leave_room [room_id]")
	}

	user := ToUser(c.Sender())
	if user == nil {
		return c.Send("Could not identify user.")
	}
	cmd := roomCommand.LeaveRoomCommand{
		RoomID:    roomEntity.RoomID(args),
		Requester: *user,
	}
	if err := h.leaveRoomHandler.Handle(context.Background(), cmd); err != nil {
		return c.Send(fmt.Sprintf("Error leaving room '%s': %v", args, err))
	}

	return c.Send(fmt.Sprintf("Successfully left room %s!", args))
}
