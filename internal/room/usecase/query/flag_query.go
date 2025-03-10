package query

import (
	"context"
	"telemafia/internal/room/repo"
)

// CheckChangeFlagQuery represents the query to check the change flag state
type CheckChangeFlagQuery struct{}

type CheckChangeFlagHandler struct {
	roomRepo repo.RoomReader
}

// NewCheckChangeFlagHandler handles checking the change flag state
func NewCheckChangeFlagHandler(repo repo.RoomReader) *CheckChangeFlagHandler {
	return &CheckChangeFlagHandler{
		roomRepo: repo,
	}
}

// Handle handles checking the change flag state
func (h *CheckChangeFlagHandler) Handle(ctx context.Context, cmd CheckChangeFlagQuery) bool {
	return h.roomRepo.CheckChangeFlag()
}
