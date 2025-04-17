package usecase

import (
	"context"
	"telemafia/internal/entity"
)

// GetPlayerRoomsQuery represents the query to get rooms a player is in
type GetPlayerRoomsQuery struct {
	PlayerID entity.UserID
}

// GetPlayerRoomsHandler handles queries for player rooms
type GetPlayerRoomsHandler struct {
	roomRepo RoomReader
}

// NewGetPlayerRoomsHandler creates a new GetPlayerRoomsHandler
func NewGetPlayerRoomsHandler(repo RoomReader) *GetPlayerRoomsHandler {
	return &GetPlayerRoomsHandler{
		roomRepo: repo,
	}
}

// Handle processes the get player rooms query
func (h *GetPlayerRoomsHandler) Handle(ctx context.Context, query GetPlayerRoomsQuery) ([]*entity.Room, error) {
	return h.roomRepo.GetPlayerRooms(query.PlayerID)
}
