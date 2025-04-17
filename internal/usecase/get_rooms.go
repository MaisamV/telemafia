package usecase

import (
	"context"
	"telemafia/internal/entity"
)

// GetRoomsQuery represents the query to get all rooms
type GetRoomsQuery struct {
	// Add filters if needed
	Limit  int
	Offset int
}

// GetRoomsHandler handles queries for all rooms
type GetRoomsHandler struct {
	roomRepo RoomReader
}

// NewGetRoomsHandler creates a new GetRoomsHandler
func NewGetRoomsHandler(repo RoomReader) *GetRoomsHandler {
	return &GetRoomsHandler{
		roomRepo: repo,
	}
}

// Handle processes the get all rooms query
func (h *GetRoomsHandler) Handle(ctx context.Context) ([]*entity.Room, error) {
	return h.roomRepo.GetRooms()
}
