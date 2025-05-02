package telegram

import (
	"fmt"
	"log"
	"strings" // Import messages
	"telemafia/internal/shared/tgutil"
	"time"

	roomHandler "telemafia/internal/presentation/telegram/handler/room"

	"gopkg.in/telebot.v3"
)

// TrackedMessage stores info about a message that needs periodic updates.
type TrackedMessage struct {
	Msg         *telebot.Message
	MessageType tgutil.RefreshingMessageType
	Data        string // e.g., room ID for RoomDetail
}

func (h *BotHandler) updateMessages(book *tgutil.RefreshingMessageBook, getMessage func(data string) (string, *telebot.ReplyMarkup, error)) {
	messagesToUpdate := book.GetAllActiveMessages()
	if len(messagesToUpdate) == 0 {
		return
	}

	log.Printf("Refreshing %d messages...", len(messagesToUpdate))
	for chatID, payload := range messagesToUpdate {
		// Prepare the updated message content using the refactored function
		// Pass the necessary handlers from the BotHandler (h)
		newContent, newMarkup, err := getMessage(payload.Data)
		if err != nil {
			log.Printf(h.msgs.Refresh.ErrorPrepare, chatID, err) // Use msg
			continue
		}

		_, editErr := h.bot.Edit(payload.Msg, newContent, newMarkup, telebot.NoPreview)
		if editErr != nil {
			if strings.Contains(editErr.Error(), "message to edit not found") ||
				strings.Contains(editErr.Error(), "bot was blocked by the user") {
				log.Printf(h.msgs.Refresh.ErrorEditRemoving, chatID, editErr) // Use msg
				book.RemoveActiveMessage(chatID)
			} else {
				log.Printf(h.msgs.Refresh.ErrorEdit, chatID, editErr) // Use msg
			}
		}
	}
	log.Println("Finished refreshing messages.")
}

// StartRefreshTimer handles updating dynamic messages.
// It periodically checks the RefreshingMessageBook and updates all active messages if changes occurred.
func (h *BotHandler) StartRefreshTimer() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		if h.roomListRefreshMessage.ConsumeRefreshNeeded() {
			h.updateMessages(h.roomListRefreshMessage, func(data string) (string, *telebot.ReplyMarkup, error) {
				return roomHandler.PrepareRoomListMessage(
					h.getRoomsHandler,
					h.getPlayersInRoomHandler,
					h.msgs, // Pass messages
				)
			})
		}
		if h.roomDetailRefreshMessage.ConsumeRefreshNeeded() {
			h.updateMessages(h.roomDetailRefreshMessage, func(data string) (string, *telebot.ReplyMarkup, error) {
				return roomHandler.RoomDetailMessage(
					h.getRoomsHandler,
					h.getPlayersInRoomHandler,
					h.msgs, // Pass messages
					data,
				)
			})
		}
	}
}

// prepareMessageContent generates content based on message type.
func (h *BotHandler) prepareMessageContent(messageType tgutil.RefreshingMessageType, data string) (string, *telebot.ReplyMarkup, error) {
	switch messageType {
	case tgutil.ListRooms:
		// Call the refactored function, passing messages
		return roomHandler.PrepareRoomListMessage(h.getRoomsHandler, h.getPlayersInRoomHandler, h.msgs)
	default:
		return "", nil, fmt.Errorf("unsupported refreshing message type: %v", messageType)
	}
}

// SendOrUpdateRefreshingMessage sends a new message and registers it for refreshing, or updates an existing one.
func (h *BotHandler) SendOrUpdateRefreshingMessage(userID int64, messageType tgutil.RefreshingMessageType, data string) error {
	content, markup, err := h.prepareMessageContent(messageType, data)
	if err != nil {
		log.Printf(h.msgs.Refresh.ErrorPrepare, userID, messageType, err)                     // Use msg
		_, sendErr := h.bot.Send(telebot.ChatID(userID), h.msgs.Common.ErrorPreparingContent) // Use msg
		if sendErr != nil {
			log.Printf("Failed to send error message to user %d: %v", userID, sendErr)
		}
		return err
	}

	if existingMsg, ok := h.roomListRefreshMessage.GetActiveMessage(userID); ok {
		updatedMsg, editErr := h.bot.Edit(existingMsg.Msg, content, markup)
		if editErr == nil {
			log.Printf(fmt.Sprintf("%d : %d", existingMsg.Msg.ID, updatedMsg.ID))
			h.roomListRefreshMessage.AddActiveMessage(userID, &tgutil.RefreshingMessage{
				Msg:  updatedMsg,
				Data: data,
			})
			log.Printf(h.msgs.Refresh.LogUpdateSuccess, userID) // Use msg
			return nil
		}
		log.Printf(h.msgs.Refresh.LogEditFailSendingNew, existingMsg.Msg.ID, userID, editErr) // Use msg
		h.roomListRefreshMessage.RemoveActiveMessage(userID)
	}

	msg, sendErr := h.bot.Send(telebot.ChatID(userID), content, markup)
	if sendErr != nil {
		log.Printf(h.msgs.Refresh.ErrorSendNew, userID, sendErr) // Use msg
		return sendErr
	}

	h.roomListRefreshMessage.AddActiveMessage(userID, &tgutil.RefreshingMessage{
		Msg:  msg,
		Data: data,
	})
	log.Printf(h.msgs.Refresh.LogSendNewSuccess, msg.ID, userID) // Use msg
	return nil
}

// RemoveRefreshingChat removes a chat from the refreshing list via RefreshingMessageBook.
func (h *BotHandler) RemoveRefreshingChat(userID int64) {
	h.roomListRefreshMessage.RemoveActiveMessage(userID)
	log.Printf(h.msgs.Refresh.LogRemovedUser, userID) // Use msg
}
