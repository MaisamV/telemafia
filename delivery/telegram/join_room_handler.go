package telegram

import (
	"context"
	"gopkg.in/telebot.v3"
	"strings"
	errorHandler "telemafia/common/error"
	"telemafia/delivery/util"
	"telemafia/internal/room/entity"
	roomCommand "telemafia/internal/room/usecase/command"
)

// HandleJoinRoom handles the /join_room command
func (h *BotHandler) HandleJoinRoom(c telebot.Context) error {
	args := strings.TrimSpace(c.Message().Payload)
	if args == "" {
		return c.Send("Please provide a room ID: /join_room [room_id]")
	}

	user := util.ToUser(c.Sender())
	// Join room
	cmd := roomCommand.JoinRoomCommand{
		RoomID: entity.RoomID(args),
		Player: *user,
	}
	if err := h.joinRoomHandler.Handle(context.Background(), cmd); err != nil {
		return c.Send(errorHandler.HandleError(err, "Error joining room"))
	}

	return c.Send("Successfully joined the room!")
}
