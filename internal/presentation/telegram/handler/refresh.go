package telegram

import (
	"log"
	"strings" // Import messages
	"telemafia/internal/shared/tgutil"
	"time"

	"gopkg.in/telebot.v4"
)

func (h *BotHandler) RefreshMessages(book *tgutil.RefreshingMessageBook) {
	messagesToUpdate := book.GetAllActiveMessages()
	if len(messagesToUpdate) == 0 {
		return
	}

	log.Printf("Refreshing %d messages...", len(messagesToUpdate))
	for chatID, payload := range messagesToUpdate {
		// Prepare the updated message content using the refactored function
		// Pass the necessary handlers from the BotHandler (h)
		newContent, newMarkup, err := book.GetMessage(chatID, payload.Data)
		if err != nil {
			log.Printf(h.msgs.Refresh.ErrorPrepare, chatID, err) // Use msg
			continue
		}

		_, editErr := h.bot.Edit(&telebot.Message{ID: payload.MessageID, Chat: &telebot.Chat{ID: payload.ChatID}}, newContent, newMarkup...)
		if editErr != nil {
			if !strings.Contains(editErr.Error(), "message is not modified") {
				book.RemoveActiveMessage(chatID)
			} else {
				log.Printf(h.msgs.Refresh.ErrorEditRemoving, chatID, editErr)
			}
		}
	}
	log.Println("Finished refreshing messages.")
}

// StartRefreshTimer handles updating dynamic messages.
// It periodically checks the RefreshingMessageBook and updates all active messages if changes occurred.
func (h *BotHandler) StartRefreshTimer() {
	ticker := time.NewTicker(1 * time.Second) // Reduced interval for faster feedback during testing
	defer ticker.Stop()
	for range ticker.C {

		// --- Game Admin Assignment Tracker Refresh ---
		h.adminRefreshMutex.RLock() // Lock for reading the map
		for gameID, book := range h.adminAssignmentTrackers {
			if book.ConsumeRefreshNeeded() {
				log.Printf("Refresh needed for Admin Tracker Game ID: %s", gameID)
				h.RefreshMessages(book)
			}
		}
		h.adminRefreshMutex.RUnlock()

		// --- Game Player Role Choice Refresh ---
		h.playerRefreshMutex.RLock() // Lock for reading the map
		for gameID, book := range h.playerRoleChoiceRefreshers {
			if book.ConsumeRefreshNeeded() {
				log.Printf("Refresh needed for Player Choice Game ID: %s", gameID)
				h.RefreshMessages(book)
			}
		}
		h.playerRefreshMutex.RUnlock()

		// --- Room List Refresh ---
		if h.roomListRefreshMessage.ConsumeRefreshNeeded() {
			h.RefreshMessages(h.roomListRefreshMessage)
		}

		// --- Room Detail Refresh ---
		if h.roomDetailRefreshMessage.ConsumeRefreshNeeded() {
			h.RefreshMessages(h.roomDetailRefreshMessage)
		}
	}
}
