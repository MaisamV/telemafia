package query

import (
	"context"
	roomEntity "telemafia/internal/domain/room/entity"
	roomPort "telemafia/internal/domain/room/port"
)

// GetRoomsQuery represents the query to get all rooms
type GetRoomsQuery struct {
	// Add filters if needed, e.g., by name, status
	// Limit  int
	// Offset int
}

// GetRoomsHandler handles queries for all rooms
type GetRoomsHandler struct {
	roomRepo roomPort.RoomReader // Use imported RoomReader interface
}

// NewGetRoomsHandler creates a new GetRoomsHandler
func NewGetRoomsHandler(repo roomPort.RoomReader) *GetRoomsHandler {
	return &GetRoomsHandler{
		roomRepo: repo,
	}
}

// Handle processes the get all rooms query
// Note: The query parameter is currently unused, but kept for consistency
func (h *GetRoomsHandler) Handle(ctx context.Context, query GetRoomsQuery) ([]*roomEntity.Room, error) {
	return h.roomRepo.GetRooms() // Propagates results/errors from repo
}
