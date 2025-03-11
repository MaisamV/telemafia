package telegram

import (
	"context"
	"fmt"
	"telemafia/delivery/util"
	"telemafia/internal/room/entity"
	roomCommand "telemafia/internal/room/usecase/command"
	roomQuery "telemafia/internal/room/usecase/query"

	"gopkg.in/telebot.v3"
)

// HandleJoinRoomCallback handles the join room callback
func (h *BotHandler) HandleJoinRoomCallback(c telebot.Context, roomID string) error {
	user := util.ToUser(c.Sender())

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

	// Fetch updated room list
	rooms, err := h.getRoomsHandler.Handle(context.Background(), roomQuery.GetRoomsQuery{})
	if err != nil {
		return c.Respond(&telebot.CallbackResponse{
			Text: "Failed to fetch updated rooms.",
		})
	}

	// Create inline keyboard with updated player counts
	var buttons [][]telebot.InlineButton
	for _, room := range rooms {
		buttonText := fmt.Sprintf("%s (بازیکنان: %d)", room.Name, len(room.Players))
		button := telebot.InlineButton{
			Unique: UniqueJoinToRoom,
			Text:   buttonText,
			Data:   string(room.ID),
		}
		buttons = append(buttons, []telebot.InlineButton{button})
	}

	markup := &telebot.ReplyMarkup{InlineKeyboard: buttons}

	c.Respond(&telebot.CallbackResponse{
		Text: "Successfully joined the room!",
	})
	// Edit the original message with updated room list
	return c.Edit("Available rooms:", markup)
}
