package query

import (
	"context"
	"telemafia/internal/room/entity"
	"telemafia/internal/room/repo"
)

// GetRoomQuery represents the query to get a specific room
type GetRoomQuery struct {
	RoomID entity.RoomID
}

// GetRoomHandler handles single room queries
type GetRoomHandler struct {
	roomRepo repo.RoomReader
}

// NewGetRoomHandler creates a new GetRoomHandler
func NewGetRoomHandler(repo repo.RoomReader) *GetRoomHandler {
	return &GetRoomHandler{
		roomRepo: repo,
	}
}

// Handle processes the get room query
func (h *GetRoomHandler) Handle(ctx context.Context, query GetRoomQuery) (*entity.Room, error) {
	return h.roomRepo.GetRoomByID(query.RoomID)
}
