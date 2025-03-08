package query

import (
	"context"
	"telemafia/internal/room/entity"
	"telemafia/internal/room/repo"
)

// GetRoomsQuery represents the query to get all rooms
type GetRoomsQuery struct {
	// Add filters if needed
	Limit  int
	Offset int
}

// GetRoomsHandler handles room queries
type GetRoomsHandler struct {
	roomRepo repo.RoomReader
}

// NewGetRoomsHandler creates a new GetRoomsHandler
func NewGetRoomsHandler(repo repo.RoomReader) *GetRoomsHandler {
	return &GetRoomsHandler{
		roomRepo: repo,
	}
}

// Handle processes the get rooms query
func (h *GetRoomsHandler) Handle(ctx context.Context, query GetRoomsQuery) ([]*entity.Room, error) {
	return h.roomRepo.GetRooms()
}
