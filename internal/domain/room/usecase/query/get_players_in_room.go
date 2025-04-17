package query

import (
	"context"
	roomEntity "telemafia/internal/domain/room/entity"
	roomPort "telemafia/internal/domain/room/port"
	sharedEntity "telemafia/internal/shared/entity"
)

// GetPlayersInRoomQuery represents the query to get players in a room
type GetPlayersInRoomQuery struct {
	RoomID roomEntity.RoomID // Use imported RoomID type
}

// GetPlayersInRoomHandler handles queries for players in a room
type GetPlayersInRoomHandler struct {
	roomRepo roomPort.RoomReader // Use imported RoomReader interface
}

// NewGetPlayersInRoomHandler creates a new GetPlayersInRoomHandler
func NewGetPlayersInRoomHandler(repo roomPort.RoomReader) *GetPlayersInRoomHandler {
	return &GetPlayersInRoomHandler{
		roomRepo: repo,
	}
}

// Handle processes the get players in room query
func (h *GetPlayersInRoomHandler) Handle(ctx context.Context, query GetPlayersInRoomQuery) ([]*sharedEntity.User, error) { // Return type updated
	return h.roomRepo.GetPlayersInRoom(query.RoomID) // Propagates results/errors from repo
}
