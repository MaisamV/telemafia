package telegram

import (
	"gopkg.in/telebot.v3"
)

// HandleStart handles the /start command
func (h *BotHandler) HandleStart(c telebot.Context) error {
	return h.HandleListRooms(c)
}
