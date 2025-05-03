package tgutil

import (
	"sync"
)

// RefreshingMessageBook manages the state for dynamic message refreshing.
// It tracks whether a refresh is needed and which messages are actively being refreshed.
type RefreshingMessageBook struct {
	mutex        sync.RWMutex
	needsRefresh bool
	// activeMessages maps chatID to the message currently showing the dynamic list.
	activeMessages map[int64]*RefreshingMessage
}

type RefreshingMessage struct {
	ChatID    int64
	MessageID int
	Data      string
}

// NewRefreshState creates a new RefreshingMessageBook manager.
func NewRefreshState() *RefreshingMessageBook {
	return &RefreshingMessageBook{
		activeMessages: make(map[int64]*RefreshingMessage),
	}
}

// RaiseRefreshNeeded sets the flag indicating the list needs refreshing.
func (rs *RefreshingMessageBook) RaiseRefreshNeeded() {
	rs.mutex.Lock()
	defer rs.mutex.Unlock()
	rs.needsRefresh = true
}

// CheckRefreshNeeded checks if a refresh is needed without resetting the flag.
func (rs *RefreshingMessageBook) CheckRefreshNeeded() bool {
	rs.mutex.RLock()
	defer rs.mutex.RUnlock()
	return rs.needsRefresh
}

// ConsumeRefreshNeeded checks if a refresh is needed and resets the flag if true.
func (rs *RefreshingMessageBook) ConsumeRefreshNeeded() bool {
	rs.mutex.Lock()
	defer rs.mutex.Unlock()
	if rs.needsRefresh {
		rs.needsRefresh = false
		return true
	}
	return false
}

// AddActiveMessage adds or updates the active message for a given chat ID.
func (rs *RefreshingMessageBook) AddActiveMessage(chatID int64, msg *RefreshingMessage) {
	rs.mutex.Lock()
	defer rs.mutex.Unlock()
	rs.activeMessages[chatID] = msg
}

// RemoveActiveMessage removes the active message tracking for a given chat ID.
func (rs *RefreshingMessageBook) RemoveActiveMessage(chatID int64) {
	rs.mutex.Lock()
	defer rs.mutex.Unlock()
	delete(rs.activeMessages, chatID)
}

// GetActiveMessage retrieves the active message for a specific chat ID.
func (rs *RefreshingMessageBook) GetActiveMessage(chatID int64) (*RefreshingMessage, bool) {
	rs.mutex.RLock()
	defer rs.mutex.RUnlock()
	msg, exists := rs.activeMessages[chatID]
	return msg, exists
}

// GetAllActiveMessages returns a *copy* of the map of active messages.
func (rs *RefreshingMessageBook) GetAllActiveMessages() map[int64]*RefreshingMessage {
	rs.mutex.RLock()
	defer rs.mutex.RUnlock()
	// Create a copy to avoid race conditions on the map elsewhere
	clone := make(map[int64]*RefreshingMessage, len(rs.activeMessages))
	for k, v := range rs.activeMessages {
		clone[k] = v
	}
	return clone
}
