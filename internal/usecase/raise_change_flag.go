package usecase

import (
	"context"
	// "telemafia/internal/entity" // Removed unused import
)

// RaiseChangeFlagCommand represents the command to raise the change flag
type RaiseChangeFlagCommand struct{}

// RaiseChangeFlagHandler handles the command to raise the change flag
type RaiseChangeFlagHandler struct {
	roomRepo RoomWriter
}

// NewRaiseChangeFlagHandler creates a new handler
func NewRaiseChangeFlagHandler(repo RoomWriter) *RaiseChangeFlagHandler {
	return &RaiseChangeFlagHandler{roomRepo: repo}
}

// Handle executes the command
func (h *RaiseChangeFlagHandler) Handle(ctx context.Context) {
	h.roomRepo.RaiseChangeFlag()
}
