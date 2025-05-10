package telegram

import (
	"strings"                                                    // Import strings package
	room "telemafia/internal/presentation/telegram/handler/room" // Import room handlers
	messages "telemafia/internal/presentation/telegram/messages" // Import messages

	"gopkg.in/telebot.v4"
	// room "telemafia/internal/presentation/telegram/handler/room"
)

// HandleHelp provides a simple help message.
func HandleHelp(h *BotHandler, c telebot.Context, msgs *messages.Messages) error {
	// Use help message from config
	return c.Send(msgs.Common.Help, &telebot.SendOptions{DisableWebPagePreview: true})
}

// HandleStart checks for deep link payload for joining rooms, otherwise shows the room list.
func HandleStart(h *BotHandler, c telebot.Context, msgs *messages.Messages) error {
	payload := c.Message().Payload

	if payload != "" {
		// Try parsing payload assuming format "unique-data"
		// unique, data := tgutil.SplitCallbackData(payload) // Don't use this as it uses |
		parts := strings.SplitN(payload, "-", 2)
		if len(parts) == 2 {
			unique := parts[0]
			data := parts[1]

			// Check if it's a join room request
			if unique == "join_room" {
				roomID := data
				// Reuse the existing Join Room callback logic
				return room.HandleJoinRoom(
					h.joinRoomHandler,
					h.getRoomsHandler,
					h.getPlayersInRoomHandler,
					h.roomListRefreshMessage,
					h.roomDetailRefreshMessage,
					c, // Pass the original message context
					roomID,
					h.msgs,
				)
			}
		}
	}

	// Default action: Show the list of rooms
	return h.handleListRooms(c)
}
