package telegram

import (
	"context"
	"fmt"
	"strconv"
	roomQuery "telemafia/internal/room/usecase/query"

	"gopkg.in/telebot.v3"
)

type MinifiedChat struct {
	chatID int64
}

func (r MinifiedChat) Recipient() string {
	return strconv.FormatInt(r.chatID, 10)
}

// HandleListRooms handles the /list_rooms command
func (h *BotHandler) HandleListRooms(c telebot.Context) error {
	text, markup, err := h.ListRoomsMessage()
	if err != nil {
		c.Respond(&telebot.CallbackResponse{Text: fmt.Sprintf("Error listing rooms: %v", err)})
		return err
	}
	return h.SendMessage(c.Sender().ID, text, markup, ListRooms, "")
}

func (h *BotHandler) ListRoomsMessage() (string, *telebot.ReplyMarkup, error) {
	rooms, err := h.getRoomsHandler.Handle(context.Background(), roomQuery.GetRoomsQuery{})
	if err != nil {
		return "", nil, err
	}

	var text string
	var markup *telebot.ReplyMarkup = nil
	if len(rooms) == 0 {
		text = "فعلا بازی در حال شروع شدن نیست."
	} else {
		// Create inline keyboard
		var buttons [][]telebot.InlineButton
		for _, room := range rooms {
			buttonText := fmt.Sprintf("%s (بازیکنان: %d)", room.Name, len(room.Players))
			buttons = append(buttons, []telebot.InlineButton{
				{
					Unique: UniqueJoinSelectRoom,
					Text:   buttonText,
					Data:   string(room.ID),
				},
			})
		}
		text = "Available rooms:"
		markup = &telebot.ReplyMarkup{
			InlineKeyboard: buttons,
		}
	}

	return text, markup, err
}
