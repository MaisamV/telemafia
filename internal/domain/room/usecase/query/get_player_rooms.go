package query

import (
	"context"
	roomEntity "telemafia/internal/domain/room/entity"
	roomPort "telemafia/internal/domain/room/port"
	sharedEntity "telemafia/internal/shared/entity"
)

// GetPlayerRoomsQuery represents the query to get rooms a player is in
type GetPlayerRoomsQuery struct {
	PlayerID sharedEntity.UserID // Use imported UserID type
}

// GetPlayerRoomsHandler handles queries for player rooms
type GetPlayerRoomsHandler struct {
	roomRepo roomPort.RoomReader // Use imported RoomReader interface
}

// NewGetPlayerRoomsHandler creates a new GetPlayerRoomsHandler
func NewGetPlayerRoomsHandler(repo roomPort.RoomReader) *GetPlayerRoomsHandler {
	return &GetPlayerRoomsHandler{
		roomRepo: repo,
	}
}

// Handle processes the get player rooms query
func (h *GetPlayerRoomsHandler) Handle(ctx context.Context, query GetPlayerRoomsQuery) ([]*roomEntity.Room, error) {
	return h.roomRepo.GetPlayerRooms(query.PlayerID) // Propagates results/errors from repo
}
