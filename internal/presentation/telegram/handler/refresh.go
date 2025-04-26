package telegram

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"telemafia/internal/shared/tgutil"
	"time"

	roomCommand "telemafia/internal/domain/room/usecase/command"
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

var (
	refreshingMessages      = make(map[int64]TrackedMessage)
	refreshingMessagesMutex sync.RWMutex
)

// RefreshRoomsList handles updating dynamic messages.
// It periodically checks a flag (set by commands that modify rooms)
// and updates all registered messages if changes occurred.
func (h *BotHandler) RefreshRoomsList() {
	updateMessages := func() {
		refreshingMessagesMutex.RLock()
		messagesToUpdate := make(map[int64]TrackedMessage)
		for userID, trackedMsg := range refreshingMessages {
			messagesToUpdate[userID] = trackedMsg
		}
		refreshingMessagesMutex.RUnlock()

		if len(messagesToUpdate) == 0 {
			return
		}

		log.Printf("Refreshing %d messages...", len(messagesToUpdate))
		for userID, trackedMsg := range messagesToUpdate {
			newContent, newMarkup, err := h.prepareMessageContent(trackedMsg.MessageType, trackedMsg.Data)
			if err != nil {
				log.Printf("Error preparing refresh content for user %d (type %v, data %s): %v", userID, trackedMsg.MessageType, trackedMsg.Data, err)
				continue
			}

			_, editErr := h.bot.Edit(trackedMsg.Msg, newContent, newMarkup)
			if editErr != nil {
				if strings.Contains(editErr.Error(), "message to edit not found") ||
					strings.Contains(editErr.Error(), "message is not modified") ||
					strings.Contains(editErr.Error(), "bot was blocked by the user") {
					log.Printf("Removing message for user %d from refresh list (edit failed: %v)", userID, editErr)
					RemoveRefreshingChat(userID)
				} else {
					log.Printf("Non-fatal error editing message for user %d: %v", userID, editErr)
				}
			}
		}
		log.Println("Finished refreshing messages.")
	}

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		if h.checkRefreshHandler.Handle(context.Background(), roomQuery.CheckChangeFlagQuery{}) {
			h.resetRefreshHandler.Handle(context.Background(), roomCommand.ResetChangeFlagCommand{})
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
func (h *BotHandler) SendOrUpdateRefreshingMessage(userID int64, messageType RefreshingMessageType, data string) error {
	content, markup, err := h.prepareMessageContent(messageType, data)
	if err != nil {
		log.Printf("Error preparing content for user %d type %v: %v", userID, messageType, err)
		_, sendErr := h.bot.Send(telebot.ChatID(userID), "Error preparing dynamic message content.")
		if sendErr != nil {
			log.Printf("Failed to send error message to user %d: %v", userID, sendErr)
		}
		return err
	}

	refreshingMessagesMutex.Lock()
	defer refreshingMessagesMutex.Unlock()

	if existingTrackedMsg, ok := refreshingMessages[userID]; ok {
		updatedMsg, editErr := h.bot.Edit(existingTrackedMsg.Msg, content, markup)
		if editErr == nil {
			refreshingMessages[userID] = TrackedMessage{
				Msg:         updatedMsg,
				MessageType: messageType,
				Data:        data,
			}
			log.Printf("Successfully updated refreshing message for user %d", userID)
			return nil
		}
		log.Printf("Failed to edit refreshing message %d for user %d, sending new: %v", existingTrackedMsg.Msg.ID, userID, editErr)
		delete(refreshingMessages, userID)
	}

	msg, sendErr := h.bot.Send(telebot.ChatID(userID), content, markup)
	if sendErr != nil {
		log.Printf("Error sending new refreshing message to user %d: %v", userID, sendErr)
		return sendErr
	}

	refreshingMessages[userID] = TrackedMessage{
		Msg:         msg,
		MessageType: messageType,
		Data:        data,
	}
	log.Printf("Sent and registered new refreshing message %d for user %d", msg.ID, userID)
	return nil
}

// ChangeRefreshType updates the type of message being refreshed for a user.
func ChangeRefreshType(userID int64, messageType RefreshingMessageType, data string) {
	refreshingMessagesMutex.Lock()
	defer refreshingMessagesMutex.Unlock()
	if trackedMsg, exists := refreshingMessages[userID]; exists {
		trackedMsg.MessageType = messageType
		trackedMsg.Data = data
		refreshingMessages[userID] = trackedMsg
		log.Printf("Changed refresh type for user %d to %v (data: %s)", userID, messageType, data)
	} else {
		log.Printf("Attempted to change refresh type for user %d, but no existing message found.", userID)
	}
}

// GetRefreshingChats returns a snapshot of the chats being refreshed.
// The returned slice contains copies, safe for concurrent reading.
type RefreshingChat struct {
	userID  int64
	message TrackedMessage
}

func GetRefreshingChats() []RefreshingChat {
	refreshingMessagesMutex.RLock()
	defer refreshingMessagesMutex.RUnlock()
	userMessages := make([]RefreshingChat, 0, len(refreshingMessages))
	for userID, trackedMsg := range refreshingMessages {
		userMessages = append(userMessages, RefreshingChat{userID: userID, message: trackedMsg})
	}
	return userMessages
}

// RemoveRefreshingChat removes a chat from the refreshing list.
func RemoveRefreshingChat(userID int64) {
	refreshingMessagesMutex.Lock()
	defer refreshingMessagesMutex.Unlock()
	delete(refreshingMessages, userID)
	log.Printf("Removed user %d from refreshing messages list.", userID)
}
