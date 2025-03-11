package telegram

import (
	"context"
	"strings"
	errorHandler "telemafia/common/error"
	"telemafia/delivery/util"
	"telemafia/internal/room/entity"
	roomCommand "telemafia/internal/room/usecase/command"
	userEntity "telemafia/internal/user/entity"

	"gopkg.in/telebot.v3"
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

// HandleLeaveRoomCallback handles the leave room callback
func (h *BotHandler) HandleLeaveRoomCallback(c telebot.Context, roomID string) error {
	user := c.Sender()

	// Create leave room command
	cmd := roomCommand.LeaveRoomCommand{
		RoomID:   entity.RoomID(roomID),
		PlayerID: userEntity.UserID(user.ID),
	}

	// Execute leave room command
	if err := h.leaveRoomHandler.Handle(context.Background(), cmd); err != nil {
		return c.Respond(&telebot.CallbackResponse{
			Text: "Error leaving room",
		})
	}

	// Notify user
	return c.Respond(&telebot.CallbackResponse{
		Text: "You have left the room.",
	})
}
