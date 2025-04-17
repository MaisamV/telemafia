package usecase

import (
	"context"
	"telemafia/internal/entity"
)

// GetRoomQuery represents the query to get a specific room
type GetRoomQuery struct {
	RoomID entity.RoomID
}

// GetRoomHandler handles queries for a specific room
type GetRoomHandler struct {
	roomRepo RoomReader
}

// NewGetRoomHandler creates a new GetRoomHandler
func NewGetRoomHandler(repo RoomReader) *GetRoomHandler {
	return &GetRoomHandler{
		roomRepo: repo,
	}
}

// Handle processes the get room query
func (h *GetRoomHandler) Handle(ctx context.Context, query GetRoomQuery) (*entity.Room, error) {
	return h.roomRepo.GetRoomByID(query.RoomID)
}
