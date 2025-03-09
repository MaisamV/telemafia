package telegram

import (
	"context"
	"telemafia/internal/room/entity"
	roomCommand "telemafia/internal/room/usecase/command"

	"gopkg.in/telebot.v3"
)

// HandleDeleteRoomCallback handles the callback to delete a specific room
func (h *BotHandler) HandleDeleteRoomCallback(c telebot.Context, roomID string) error {
	cmd := roomCommand.DeleteRoomCommand{
		RoomID: entity.RoomID(roomID),
	}
	if err := h.deleteRoomHandler.Handle(context.Background(), cmd); err != nil {
		return c.Respond(&telebot.CallbackResponse{
			Text: "Failed to delete room.",
		})
	}

	return c.Respond(&telebot.CallbackResponse{
		Text: "Room successfully deleted!",
	})
}
