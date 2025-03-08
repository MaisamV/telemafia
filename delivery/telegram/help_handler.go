package telegram

import (
	"gopkg.in/telebot.v3"
)

// HandleHelp handles the /help command
func (h *BotHandler) HandleHelp(c telebot.Context) error {
	help := `Available commands:
/start - Start the bot
/help - Show this help message
/create_room [name] - Create a new room
/join_room [room_id] - Join a room
/leave_room [room_id] - Leave a room
/list_rooms - List all rooms
/my_rooms - List your rooms`
	return c.Send(help)
} 