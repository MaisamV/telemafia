package command

import (
	"context"
	"telemafia/internal/room/repo"
)

// RaiseChangeFlagCommand represents a command to raise the change flag
type RaiseChangeFlagCommand struct{}

// RaiseChangeFlagHandler handles the raise change flag command
type RaiseChangeFlagHandler struct {
	roomRepo repo.Repository
}

// NewRaiseChangeFlagHandler creates a new raise change flag handler
func NewRaiseChangeFlagHandler(roomRepo repo.Repository) *RaiseChangeFlagHandler {
	return &RaiseChangeFlagHandler{
		roomRepo: roomRepo,
	}
}

// Handle raises the change flag
func (h *RaiseChangeFlagHandler) Handle(ctx context.Context, cmd RaiseChangeFlagCommand) error {
	h.roomRepo.RaiseChangeFlag()
	return nil
}
