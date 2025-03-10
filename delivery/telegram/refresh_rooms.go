package telegram

import (
	"context"
	"fmt"
	flagCommand "telemafia/internal/room/usecase/command"
	roomQuery "telemafia/internal/room/usecase/query"
	"time"

	"gopkg.in/telebot.v3"
)

// RefreshRoomsList handles updating the room list for all users
func (h *BotHandler) RefreshRoomsList() {
	updateRoomList := func() {
		fmt.Println("Refreshing room list")
		rooms, err := h.getRoomsHandler.Handle(context.Background(), roomQuery.GetRoomsQuery{})
		if err != nil {
			fmt.Println("RefreshRoomsList: Error getting rooms.")
			return
		}

		messageIDsMutex.RLock()
		userMessages := make([]UserMessage, 0, len(messageIDs))
		for userID, messageID := range messageIDs {
			userMessages = append(userMessages, UserMessage{
				userID:    userID,
				messageID: messageID,
			})
		}
		messageIDsMutex.RUnlock()

		if len(rooms) == 0 {
			for _, um := range userMessages {
				fmt.Println(fmt.Sprintf("Refreshing message for user %d with message ID %d", um.userID, um.messageID))
				h.bot.Edit(&telebot.Message{ID: um.messageID, Chat: &telebot.Chat{ID: um.userID}}, "فعلا بازی در حال شروع شدن نیست.")
			}
			return
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

		markup := &telebot.ReplyMarkup{InlineKeyboard: buttons}

		for _, um := range userMessages {
			fmt.Println(fmt.Sprintf("Refreshing message for user %d with message ID %d", um.userID, um.messageID))
			h.bot.Edit(&telebot.Message{ID: um.messageID, Chat: &telebot.Chat{ID: um.userID}}, "Available rooms:", markup)
		}
	}

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		if h.resetRefreshHandler.Handle(context.Background(), flagCommand.ResetChangeFlagCommand{}) {
			updateRoomList()
		}
	}
}
