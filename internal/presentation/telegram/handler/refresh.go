package telegram

import (
	"log"
	"strings" // Import messages
	"telemafia/internal/shared/tgutil"
	"time"

	roomHandler "telemafia/internal/presentation/telegram/handler/room"

	"gopkg.in/telebot.v3"
)

func (h *BotHandler) updateMessages(book *tgutil.RefreshingMessageBook, getMessage func(user int64, data string) (string, []interface{}, error)) {
	messagesToUpdate := book.GetAllActiveMessages()
	if len(messagesToUpdate) == 0 {
		return
	}

	log.Printf("Refreshing %d messages...", len(messagesToUpdate))
	for chatID, payload := range messagesToUpdate {
		// Prepare the updated message content using the refactored function
		// Pass the necessary handlers from the BotHandler (h)
		newContent, newMarkup, err := getMessage(chatID, payload.Data)
		if err != nil {
			log.Printf(h.msgs.Refresh.ErrorPrepare, chatID, err) // Use msg
			continue
		}

		_, editErr := h.bot.Edit(&telebot.Message{ID: payload.MessageID, Chat: &telebot.Chat{ID: payload.ChatID}}, newContent, newMarkup...)
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
			h.updateMessages(h.roomListRefreshMessage, func(user int64, data string) (string, []interface{}, error) {
				message, markup, err := roomHandler.PrepareRoomListMessage(
					h.getRoomsHandler,
					h.getPlayersInRoomHandler,
					h.msgs, // Pass messages
				)
				opts := []interface{}{
					markup,
					telebot.NoPreview,
				}
				return message, opts, err
			})
		}
		if h.roomDetailRefreshMessage.ConsumeRefreshNeeded() {
			h.updateMessages(h.roomDetailRefreshMessage, func(user int64, data string) (string, []interface{}, error) {
				message, markup, err := roomHandler.RoomDetailMessage(
					h.getRoomsHandler,
					h.getPlayersInRoomHandler,
					h.msgs, // Pass messages
					tgutil.IsAdmin(user),
					data,
				)
				opts := []interface{}{
					markup,
					telebot.ModeMarkdownV2,
					telebot.NoPreview,
				}
				return message, opts, err
			})
		}
	}
}
