package query

import (
	"context"
	"telemafia/internal/room/entity"
	"telemafia/internal/room/repo"
	userEntity "telemafia/internal/user/entity"
)

// GetPlayerRoomsQuery represents the query to get rooms for a player
type GetPlayerRoomsQuery struct {
	PlayerID userEntity.UserID
}

// GetPlayerRoomsHandler handles player room queries
type GetPlayerRoomsHandler struct {
	roomRepo repo.RoomReader
}

// NewGetPlayerRoomsHandler creates a new GetPlayerRoomsHandler
func NewGetPlayerRoomsHandler(repo repo.RoomReader) *GetPlayerRoomsHandler {
	return &GetPlayerRoomsHandler{
		roomRepo: repo,
	}
}

// Handle processes the get player rooms query
func (h *GetPlayerRoomsHandler) Handle(ctx context.Context, query GetPlayerRoomsQuery) ([]*entity.Room, error) {
	return h.roomRepo.GetPlayerRooms(query.PlayerID)
}
