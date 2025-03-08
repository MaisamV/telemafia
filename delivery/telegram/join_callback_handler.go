package telegram

import (
	"context"
	"telemafia/delivery/common"
	"telemafia/internal/room/entity"
	roomCommand "telemafia/internal/room/usecase/command"

	"gopkg.in/telebot.v3"
)

// HandleJoinRoomCallback handles the join room callback
func (h *BotHandler) HandleJoinRoomCallback(c telebot.Context, data string) error {
	user := common.ToUser(c.Sender())
	roomID := data

	// Join room
	cmd := roomCommand.JoinRoomCommand{
		RoomID:   entity.RoomID(roomID),
		PlayerID: user.ID,
		Name:     user.Username,
	}
	if err := h.joinRoomHandler.Handle(context.Background(), cmd); err != nil {
		return c.Respond(&telebot.CallbackResponse{
			Text: "Error joining room",
		})
	}

	return c.Respond(&telebot.CallbackResponse{
		Text: "Successfully joined the room!",
	})
}
