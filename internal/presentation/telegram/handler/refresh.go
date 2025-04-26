package telegram

import (
	"context"
	"fmt"
	"log"
	"strings"
	"telemafia/internal/shared/tgutil"
	"time"

	roomQuery "telemafia/internal/domain/room/usecase/query"

	"gopkg.in/telebot.v3"
)

// --- Refresh Logic (potentially move to its own package/file later) ---

// RefreshingMessageType defines the type of content being refreshed.
type RefreshingMessageType int

const (
	ListRooms  RefreshingMessageType = iota
	RoomDetail                       // Placeholder for potential future use
)

// TrackedMessage stores info about a message that needs periodic updates.
type TrackedMessage struct {
	Msg         *telebot.Message
	MessageType RefreshingMessageType
	Data        string // e.g., room ID for RoomDetail
}

// RefreshRoomsList handles updating dynamic messages.
// It periodically checks the RefreshState and updates all active messages if changes occurred.
func (h *BotHandler) RefreshRoomsList() {
	updateMessages := func() {
		// Get a clone of the active messages map from RefreshState
		messagesToUpdate := h.refreshState.GetAllActiveMessages()

		if len(messagesToUpdate) == 0 {
			return
		}

		log.Printf("Refreshing %d messages...", len(messagesToUpdate))
		for chatID, msg := range messagesToUpdate {
			// TODO: Refactor prepareMessageContent to not depend on h or pass necessary parts
			// For now, we assume it needs getRoomsHandler and getPlayersInRoomHandler
			newContent, newMarkup, err := h.prepareListRoomsMessage() // Assuming only ListRooms for now
			if err != nil {
				log.Printf("Error preparing refresh content for chat %d: %v", chatID, err)
				continue
			}

			_, editErr := h.bot.Edit(msg, newContent, newMarkup)
			if editErr != nil {
				if strings.Contains(editErr.Error(), "message to edit not found") ||
					strings.Contains(editErr.Error(), "message is not modified") ||
					strings.Contains(editErr.Error(), "bot was blocked by the user") {
					log.Printf("Removing message for chat %d from refresh list (edit failed: %v)", chatID, editErr)
					// Use RefreshState method to remove
					h.refreshState.RemoveActiveMessage(chatID)
				} else {
					log.Printf("Non-fatal error editing message for chat %d: %v", chatID, editErr)
				}
			}
		}
		log.Println("Finished refreshing messages.")
	}

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		// Use RefreshState method to check and consume the flag
		if h.refreshState.ConsumeRefreshNeeded() {
			updateMessages()
		}
	}
}

// prepareMessageContent generates content based on message type.
func (h *BotHandler) prepareMessageContent(messageType RefreshingMessageType, data string) (string, *telebot.ReplyMarkup, error) {
	switch messageType {
	case ListRooms:
		return h.prepareListRoomsMessage()
	// case RoomDetail: // Example for future use
	// 	 return h.prepareRoomDetailMessage(data)
	default:
		return "", nil, fmt.Errorf("unsupported refreshing message type: %v", messageType)
	}
}

// prepareListRoomsMessage generates the text and markup for the list rooms view
func (h *BotHandler) prepareListRoomsMessage() (string, *telebot.ReplyMarkup, error) {
	rooms, err := h.getRoomsHandler.Handle(context.Background(), roomQuery.GetRoomsQuery{})
	if err != nil {
		return "", nil, fmt.Errorf("error getting rooms: %w", err)
	}

	var response strings.Builder
	markup := &telebot.ReplyMarkup{}
	var rows []telebot.Row

	if len(rooms) == 0 {
		response.WriteString("No rooms available.")
	} else {
		response.WriteString("Available Rooms (refreshed):\n")
		for _, room := range rooms {
			players, _ := h.getPlayersInRoomHandler.Handle(context.Background(), roomQuery.GetPlayersInRoomQuery{RoomID: room.ID})
			playerCount := len(players)
			maxPlayers := 10 // TODO: Define elsewhere
			response.WriteString(fmt.Sprintf("- %s (%s) [%d/%d players]\n", room.Name, room.ID, playerCount, maxPlayers))
			btnJoin := markup.Data(fmt.Sprintf("Join %s", room.Name), tgutil.UniqueJoinRoom, string(room.ID))
			rows = append(rows, markup.Row(btnJoin))
		}
	}
	markup.Inline(rows...)
	return response.String(), markup, nil
}

// SendOrUpdateRefreshingMessage sends a new message and registers it for refreshing, or updates an existing one.
// This now uses the RefreshState manager.
func (h *BotHandler) SendOrUpdateRefreshingMessage(userID int64, messageType RefreshingMessageType, data string) error {
	// TODO: Refactor prepareMessageContent call
	content, markup, err := h.prepareListRoomsMessage() // Assuming ListRooms for now
	if err != nil {
		log.Printf("Error preparing content for user %d type %v: %v", userID, messageType, err)
		_, sendErr := h.bot.Send(telebot.ChatID(userID), "Error preparing dynamic message content.")
		if sendErr != nil {
			log.Printf("Failed to send error message to user %d: %v", userID, sendErr)
		}
		return err
	}

	// Check if we already have an active message for this user
	if existingMsg, ok := h.refreshState.GetActiveMessage(userID); ok {
		updatedMsg, editErr := h.bot.Edit(existingMsg, content, markup)
		if editErr == nil {
			// Update the message pointer in RefreshState
			h.refreshState.AddActiveMessage(userID, updatedMsg)
			log.Printf("Successfully updated refreshing message for user %d", userID)
			return nil
		}
		log.Printf("Failed to edit refreshing message %d for user %d, sending new: %v", existingMsg.ID, userID, editErr)
		// Remove the old one before sending a new one
		h.refreshState.RemoveActiveMessage(userID)
	}

	// Send a new message
	msg, sendErr := h.bot.Send(telebot.ChatID(userID), content, markup)
	if sendErr != nil {
		log.Printf("Error sending new refreshing message to user %d: %v", userID, sendErr)
		return sendErr
	}

	// Add the new message to RefreshState
	h.refreshState.AddActiveMessage(userID, msg)
	log.Printf("Sent and registered new refreshing message %d for user %d", msg.ID, userID)
	return nil
}

// RemoveRefreshingChat removes a chat from the refreshing list via RefreshState.
func (h *BotHandler) RemoveRefreshingChat(userID int64) {
	h.refreshState.RemoveActiveMessage(userID)
	log.Printf("Removed user %d from refreshing messages list.", userID)
}
