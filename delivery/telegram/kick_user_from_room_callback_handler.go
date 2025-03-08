package telegram

import (
	"context"
	"strings"
	"telemafia/common"
	"telemafia/internal/room/entity"
	roomCommand "telemafia/internal/room/usecase/command"
	userEntity "telemafia/internal/user/entity"

	"gopkg.in/telebot.v3"
)

// HandleKickUserFromRoomCallback handles the callback to kick a specific user from a room
func (h *BotHandler) HandleKickUserFromRoomCallback(c telebot.Context, data string) error {
	// Extract room ID and player ID from callback data
	parts := strings.Split(data, ":")
	if len(parts) != 2 {
		return c.Respond(&telebot.CallbackResponse{
			Text: "Invalid data format.",
		})
	}
	roomID := parts[0]
	playerID := parts[1]

	id, err := common.StringToInt64(playerID)
	if err != nil {
		return c.Respond(&telebot.CallbackResponse{
			Text: "Failed to extract user ID.",
		})
	}
	// Kick user from room
	cmd := roomCommand.KickUserCommand{
		RoomID:   entity.RoomID(roomID),
		PlayerID: userEntity.UserID(id),
	}
	if err := h.kickUserHandler.Handle(context.Background(), cmd); err != nil {
		return c.Respond(&telebot.CallbackResponse{
			Text: "Failed to kick user.",
		})
	}

	return c.Respond(&telebot.CallbackResponse{
		Text: "User successfully kicked from the room!",
	})
}
