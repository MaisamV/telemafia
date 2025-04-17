package command

import (
	"context"
	roomPort "telemafia/internal/domain/room/port"
)

// ResetChangeFlagCommand represents the command to reset (consume) the change flag
// This seems highly specific to the in-memory repo's transient flag.
// Consider if this logic should remain if persistence changes.
type ResetChangeFlagCommand struct{}

// ResetChangeFlagHandler handles the command to reset the change flag
type ResetChangeFlagHandler struct {
	roomRepo roomPort.RoomWriter // Use imported RoomWriter interface
}

// NewResetChangeFlagHandler creates a new handler (Corrected constructor name)
func NewResetChangeFlagHandler(repo roomPort.RoomWriter) *ResetChangeFlagHandler {
	return &ResetChangeFlagHandler{roomRepo: repo}
}

// Handle executes the command
// Note: Command struct is empty, context is unused.
func (h *ResetChangeFlagHandler) Handle(ctx context.Context, cmd ResetChangeFlagCommand) bool {
	return h.roomRepo.ConsumeChangeFlag()
}
