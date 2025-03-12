package telegram

import (
	"context"
	"fmt"
	"gopkg.in/telebot.v3"
	"sync"
	flagCommand "telemafia/internal/room/usecase/command"
	"time"
)

type RefreshingMessageType int

const (
	ListRooms RefreshingMessageType = iota
	RoomDetail
)

type RefreshingMessage struct {
	ID          int
	messageType RefreshingMessageType
	data        string
}

type RefreshingChat struct {
	userID  int64
	message RefreshingMessage
}

// Map to store message IDs for each user
var (
	refreshingChats      = make(map[int64]RefreshingMessage)
	refreshingChatsMutex = &sync.RWMutex{}
)

// RefreshRoomsList handles updating the room list for all users
func (h *BotHandler) RefreshRoomsList() {
	updateRoomList := func() {
		refreshingChats := GetRefreshingChats()

		for _, um := range refreshingChats {
			var text string = ""
			var markup *telebot.ReplyMarkup = nil
			var err error = nil
			switch um.message.messageType {
			case ListRooms:
				text, markup, err = h.ListRoomsMessage()
			case RoomDetail:
				text, markup, err = h.RoomDetailMessage(um.message.data)
				if err != nil {
					ChangeRefreshType(um.userID, ListRooms, um.message.data)
					text, markup, err = h.ListRoomsMessage()
				}
			}
			if err != nil {
				fmt.Printf("Error refreshing message: %v", err)
			} else {
				_, err = h.bot.Edit(&telebot.Message{ID: um.message.ID, Chat: &telebot.Chat{ID: um.userID}}, text, markup)
			}
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

func (h *BotHandler) SendMessage(userId int64, text string, markup *telebot.ReplyMarkup, messageType RefreshingMessageType, data string) error {
	refreshingChatsMutex.Lock()
	defer refreshingChatsMutex.Unlock()
	prevMsg, exists := refreshingChats[userId]
	delete(refreshingChats, userId)
	if exists {
		h.bot.Delete(&telebot.Message{ID: prevMsg.ID, Chat: &telebot.Chat{ID: userId}})
	}
	msg, err := h.bot.Send(&telebot.Chat{ID: userId}, text, markup)
	if err == nil {
		refreshingChats[userId] = RefreshingMessage{
			ID:          msg.ID,
			messageType: messageType,
			data:        data,
		}
	}
	return err
}

func (h *BotHandler) UpdateMessage(userId int64, messageID int, text string, markup *telebot.ReplyMarkup) error {
	_, err := h.bot.Edit(&telebot.Message{ID: messageID, Chat: &telebot.Chat{ID: userId}}, text, markup)
	return err
}

func ChangeRefreshType(userId int64, messageType RefreshingMessageType, data string) {
	// Change refresh message type
	refreshingChatsMutex.Lock()
	refreshingChats[userId] = RefreshingMessage{
		ID:          refreshingChats[userId].ID,
		messageType: messageType,
		data:        data,
	}
	refreshingChatsMutex.Unlock()
}

func GetRefreshingChats() []RefreshingChat {
	refreshingChatsMutex.RLock()
	userMessages := make([]RefreshingChat, 0, len(refreshingChats))
	for userID, message := range refreshingChats {
		userMessages = append(userMessages, RefreshingChat{
			userID:  userID,
			message: message,
		})
	}
	refreshingChatsMutex.RUnlock()
	return userMessages
}
