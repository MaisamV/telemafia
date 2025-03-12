package telegram

import (
	"context"
	"fmt"
	"gopkg.in/telebot.v3"
	flagCommand "telemafia/internal/room/usecase/command"
	"time"
)

// RefreshRoomsList handles updating the room list for all users
func (h *BotHandler) RefreshRoomsList() {
	updateRoomList := func() {
		userMessages := GetUpdatableMessages()

		for _, um := range userMessages {
			h.UpdateRoomListMessage(um)
		}
	}

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		if h.resetRefreshHandler.Handle(context.Background(), flagCommand.ResetChangeFlagCommand{}) {
			updateRoomList()
		}
	}
}

func (h *BotHandler) UpdateRoomListMessage(um UserMessage) {
	text, markup, err := h.ListRoomsMessage()
	if err != nil {
		fmt.Printf("Error refreshing message: %v", err)
	} else {
		_, err = h.bot.Edit(&telebot.Message{ID: um.messageID, Chat: &telebot.Chat{ID: um.userID}}, text, markup)
	}
}

func GetUpdatableMessages() []UserMessage {
	messageIDsMutex.RLock()
	userMessages := make([]UserMessage, 0, len(messageIDs))
	for userID, messageID := range messageIDs {
		userMessages = append(userMessages, UserMessage{
			userID:    userID,
			messageID: messageID,
		})
	}
	messageIDsMutex.RUnlock()
	return userMessages
}
