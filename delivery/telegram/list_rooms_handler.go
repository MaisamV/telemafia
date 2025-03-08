package telegram

import (
	"context"
	"fmt"
	roomQuery "telemafia/internal/room/usecase/query"

	"gopkg.in/telebot.v3"
)

// HandleListRooms handles the /list_rooms command
func (h *BotHandler) HandleListRooms(c telebot.Context) error {
	rooms, err := h.getRoomsHandler.Handle(context.Background(), roomQuery.GetRoomsQuery{})
	if err != nil {
		return c.Send(fmt.Sprintf("Error getting rooms: %v", err))
	}

	if len(rooms) == 0 {
		return c.Send("فعلا بازی در حال شروع شدن نیست.")
	}

	// Create inline keyboard
	var buttons [][]telebot.InlineButton
	for _, room := range rooms {
		buttonText := fmt.Sprintf("%s (بازیکنان: %d)", room.Name, len(room.Players))
		buttons = append(buttons, []telebot.InlineButton{
			{
				Unique: UniqueJoinToRoom,
				Text:   buttonText,
				Data:   string(room.ID),
			},
		})
	}

	return c.Send("Available rooms:", &telebot.ReplyMarkup{
		InlineKeyboard: buttons,
	})
}
