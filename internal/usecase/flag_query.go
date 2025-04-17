package usecase

import (
	"context"
	// "telemafia/internal/entity" // Removed unused import
)

// CheckChangeFlagQuery represents the query for the change flag
type CheckChangeFlagQuery struct{}

// CheckChangeFlagHandler handles the query for the change flag
type CheckChangeFlagHandler struct {
	roomRepo RoomReader
}

// NewCheckChangeFlagHandler creates a new handler
func NewCheckChangeFlagHandler(repo RoomReader) *CheckChangeFlagHandler {
	return &CheckChangeFlagHandler{roomRepo: repo}
}

// Handle executes the query
func (h *CheckChangeFlagHandler) Handle(ctx context.Context) bool {
	return h.roomRepo.CheckChangeFlag()
}
