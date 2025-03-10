package command

import (
	"context"
	"telemafia/internal/room/repo"
)

// ResetChangeFlagCommand represents the command to reset the change flag
type ResetChangeFlagCommand struct{}

type ResetChangeFlagHandler struct {
	roomRepo repo.RoomWriter
}

// NewResetChangeFlagCommand creates a new ResetChangeFlagHandler
func NewResetChangeFlagCommand(repo repo.RoomWriter) *ResetChangeFlagHandler {
	return &ResetChangeFlagHandler{
		roomRepo: repo,
	}
}

// Handle processes the delete room command
func (h *ResetChangeFlagHandler) Handle(ctx context.Context, cmd ResetChangeFlagCommand) bool {
	return h.roomRepo.ConsumeChangeFlag()
}
