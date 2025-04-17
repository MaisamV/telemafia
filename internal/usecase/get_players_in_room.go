package usecase

import (
	"context"
	"telemafia/internal/entity"
)

// GetPlayersInRoomQuery represents the query to get players in a room
type GetPlayersInRoomQuery struct {
	RoomID entity.RoomID
}

// GetPlayersInRoomHandler handles queries for players in a room
type GetPlayersInRoomHandler struct {
	roomRepo RoomReader
}

// NewGetPlayersInRoomHandler creates a new GetPlayersInRoomHandler
func NewGetPlayersInRoomHandler(repo RoomReader) *GetPlayersInRoomHandler {
	return &GetPlayersInRoomHandler{
		roomRepo: repo,
	}
}

// Handle processes the get players in room query
func (h *GetPlayersInRoomHandler) Handle(ctx context.Context, query GetPlayersInRoomQuery) ([]*entity.User, error) {
	return h.roomRepo.GetPlayersInRoom(query.RoomID)
}
