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

// HandleLeaveRoom handles the /leave_room command
func (h *BotHandler) HandleLeaveRoom(c telebot.Context) error {
	args := strings.TrimSpace(c.Message().Payload)
	if args == "" {
		return c.Send("Please provide a room ID: /leave_room [room_id]")
	}

	user := util.ToUser(c.Sender())
	// Leave room
	cmd := roomCommand.LeaveRoomCommand{
		RoomID:   entity.RoomID(args),
		PlayerID: user.ID,
	}
	if err := h.leaveRoomHandler.Handle(context.Background(), cmd); err != nil {
		return c.Send(errorHandler.HandleError(err, "Error leaving room"))
	}

	return c.Send("Successfully left the room!")
}
