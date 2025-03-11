package telegram

import (
	"context"
	"telemafia/common"
	roomQuery "telemafia/internal/room/usecase/query"

	"gopkg.in/telebot.v3"
)

// HandleKickUser handles the /kick_user command
func (h *BotHandler) HandleKickUser(c telebot.Context) error {
	// Check if the user is an admin
	if !common.Contains(h.adminUsernames, c.Sender().Username) {
		return c.Send("You are not authorized to use this command.")
	}

	// Fetch all rooms and send them as inline keyboard buttons
	rooms, err := h.getRoomsHandler.Handle(context.Background(), roomQuery.GetRoomsQuery{})
	if err != nil {
		return c.Send("Failed to fetch rooms.")
	}

	var buttons [][]telebot.InlineButton
	for _, room := range rooms {
		button := telebot.InlineButton{
			Unique: UniqueKickSelectRoom,
			Text:   room.Name,
			Data:   string(room.ID),
		}
		buttons = append(buttons, []telebot.InlineButton{button})
	}

	markup := &telebot.ReplyMarkup{InlineKeyboard: buttons}
	return c.Send("Select a room:", markup)
}
