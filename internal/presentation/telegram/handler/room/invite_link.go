package telegram

import (
	"fmt"

	messages "telemafia/internal/presentation/telegram/messages"

	"gopkg.in/telebot.v3"
)

// HandleGetInviteLinkCallback responds with the room's invite link.
func HandleGetInviteLinkCallback(
	bot *telebot.Bot, // Need bot instance to get username
	c telebot.Context,
	data string, // roomID
	msgs *messages.Messages,
) error {
	roomID := data
	botUsername := bot.Me.Username

	// Construct the deep link
	// Format: tg://resolve?domain=YourBotUsername&start=join_room-ROOMID
	payload := fmt.Sprintf("join_room-%s", roomID)
	inviteLink := fmt.Sprintf("https://t.me/%s?start=%s", botUsername, payload)

	// Respond to the callback privately with the link
	responseText := fmt.Sprintf(msgs.Room.InviteLinkResponse, inviteLink)

	// Acknowledge the callback first (silently) to remove loading state
	_ = c.Respond()

	// Send the link as a new message in the chat
	return c.Send(responseText, &telebot.SendOptions{})
}
