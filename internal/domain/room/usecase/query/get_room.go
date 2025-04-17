package query

import (
	"context"
	roomEntity "telemafia/internal/domain/room/entity"
	roomPort "telemafia/internal/domain/room/port"
)

// GetRoomQuery represents the query to get a specific room
type GetRoomQuery struct {
	RoomID roomEntity.RoomID // Use imported RoomID type
}

// GetRoomHandler handles queries for a specific room
type GetRoomHandler struct {
	roomRepo roomPort.RoomReader // Use imported RoomReader interface
}

// NewGetRoomHandler creates a new GetRoomHandler
func NewGetRoomHandler(repo roomPort.RoomReader) *GetRoomHandler {
	return &GetRoomHandler{
		roomRepo: repo,
	}
}

// Handle processes the get room query
func (h *GetRoomHandler) Handle(ctx context.Context, query GetRoomQuery) (*roomEntity.Room, error) {
	return h.roomRepo.GetRoomByID(query.RoomID) // Propagates results/errors from repo
}
