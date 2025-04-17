package query

import (
	"context"
	roomPort "telemafia/internal/domain/room/port"
)

// CheckChangeFlagQuery represents the query for the change flag
// This seems highly specific to the in-memory repo's transient flag.
// Consider if this logic should remain if persistence changes.
type CheckChangeFlagQuery struct{}

// CheckChangeFlagHandler handles the query for the change flag
type CheckChangeFlagHandler struct {
	roomRepo roomPort.RoomReader // Use imported RoomReader interface
}

// NewCheckChangeFlagHandler creates a new handler
func NewCheckChangeFlagHandler(repo roomPort.RoomReader) *CheckChangeFlagHandler {
	return &CheckChangeFlagHandler{roomRepo: repo}
}

// Handle executes the query
// Note: Query struct is empty, context is unused.
func (h *CheckChangeFlagHandler) Handle(ctx context.Context, query CheckChangeFlagQuery) bool {
	return h.roomRepo.CheckChangeFlag()
}
