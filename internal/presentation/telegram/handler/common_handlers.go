package telegram

import (
	messages "telemafia/internal/presentation/telegram/messages" // Import messages

	"gopkg.in/telebot.v3"
	// room "telemafia/internal/presentation/telegram/handler/room"
)

// HandleHelp provides a simple help message.
func HandleHelp(h *BotHandler, c telebot.Context, msgs *messages.Messages) error {
	// Use help message from config
	return c.Send(msgs.Common.Help, &telebot.SendOptions{DisableWebPagePreview: true})
}

// HandleStart sends a welcome message and then shows the dynamic room list.
func HandleStart(h *BotHandler, c telebot.Context, msgs *messages.Messages) error {
	return h.handleListRooms(c)
}
