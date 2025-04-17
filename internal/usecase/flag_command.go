package usecase

import (
	"context"
)

// ResetChangeFlagCommand represents the command to reset the change flag
type ResetChangeFlagCommand struct{}

// ResetChangeFlagHandler handles the command to reset the change flag
type ResetChangeFlagHandler struct {
	roomRepo RoomWriter
}

// NewResetChangeFlagCommand creates a new handler
func NewResetChangeFlagCommand(repo RoomWriter) *ResetChangeFlagHandler {
	return &ResetChangeFlagHandler{roomRepo: repo}
}

// Handle executes the command
func (h *ResetChangeFlagHandler) Handle(ctx context.Context) bool {
	return h.roomRepo.ConsumeChangeFlag()
}
