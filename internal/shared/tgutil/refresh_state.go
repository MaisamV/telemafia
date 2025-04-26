package tgutil

import (
	"sync"

	"gopkg.in/telebot.v3"
)

// RefreshState manages the state for dynamic message refreshing.
// It tracks whether a refresh is needed and which messages are actively being refreshed.
type RefreshState struct {
	mutex        sync.RWMutex
	needsRefresh bool
	// activeMessages maps chatID to the message currently showing the dynamic list.
	activeMessages map[int64]*telebot.Message
}

// NewRefreshState creates a new RefreshState manager.
func NewRefreshState() *RefreshState {
	return &RefreshState{
		activeMessages: make(map[int64]*telebot.Message),
	}
}

// RaiseRefreshNeeded sets the flag indicating the list needs refreshing.
func (rs *RefreshState) RaiseRefreshNeeded() {
	rs.mutex.Lock()
	defer rs.mutex.Unlock()
	rs.needsRefresh = true
}

// CheckRefreshNeeded checks if a refresh is needed without resetting the flag.
func (rs *RefreshState) CheckRefreshNeeded() bool {
	rs.mutex.RLock()
	defer rs.mutex.RUnlock()
	return rs.needsRefresh
}

// ConsumeRefreshNeeded checks if a refresh is needed and resets the flag if true.
func (rs *RefreshState) ConsumeRefreshNeeded() bool {
	rs.mutex.Lock()
	defer rs.mutex.Unlock()
	if rs.needsRefresh {
		rs.needsRefresh = false
		return true
	}
	return false
}

// AddActiveMessage adds or updates the active message for a given chat ID.
func (rs *RefreshState) AddActiveMessage(chatID int64, msg *telebot.Message) {
	rs.mutex.Lock()
	defer rs.mutex.Unlock()
	rs.activeMessages[chatID] = msg
}

// RemoveActiveMessage removes the active message tracking for a given chat ID.
func (rs *RefreshState) RemoveActiveMessage(chatID int64) {
	rs.mutex.Lock()
	defer rs.mutex.Unlock()
	delete(rs.activeMessages, chatID)
}

// GetActiveMessage retrieves the active message for a specific chat ID.
func (rs *RefreshState) GetActiveMessage(chatID int64) (*telebot.Message, bool) {
	rs.mutex.RLock()
	defer rs.mutex.RUnlock()
	msg, exists := rs.activeMessages[chatID]
	return msg, exists
}

// GetAllActiveMessages returns a *copy* of the map of active messages.
func (rs *RefreshState) GetAllActiveMessages() map[int64]*telebot.Message {
	rs.mutex.RLock()
	defer rs.mutex.RUnlock()
	// Create a copy to avoid race conditions on the map elsewhere
	clone := make(map[int64]*telebot.Message, len(rs.activeMessages))
	for k, v := range rs.activeMessages {
		clone[k] = v
	}
	return clone
}
