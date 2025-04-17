package command

import (
	"context"
	roomPort "telemafia/internal/domain/room/port"
)

// RaiseChangeFlagCommand represents the command to raise the change flag
// This seems highly specific to the in-memory repo's transient flag.
// Consider if this logic should remain if persistence changes.
type RaiseChangeFlagCommand struct{}

// RaiseChangeFlagHandler handles the command to raise the change flag
type RaiseChangeFlagHandler struct {
	roomRepo roomPort.RoomWriter // Use imported RoomWriter interface
}

// NewRaiseChangeFlagHandler creates a new handler
func NewRaiseChangeFlagHandler(repo roomPort.RoomWriter) *RaiseChangeFlagHandler {
	return &RaiseChangeFlagHandler{roomRepo: repo}
}

// Handle executes the command
// Note: Command struct is empty, context is unused.
func (h *RaiseChangeFlagHandler) Handle(ctx context.Context, cmd RaiseChangeFlagCommand) {
	h.roomRepo.RaiseChangeFlag()
}
