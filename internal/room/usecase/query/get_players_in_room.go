package query

import (
	"context"
	"telemafia/internal/room/entity"
	"telemafia/internal/room/repo"
	userEntity "telemafia/internal/user/entity"
)

// GetPlayersInRoomQuery represents the query to get players in a specific room
type GetPlayersInRoomQuery struct {
	RoomID entity.RoomID
}

// GetPlayersInRoomHandler handles queries to get players in a room
type GetPlayersInRoomHandler struct {
	roomRepo repo.RoomReader
}

// NewGetPlayersInRoomHandler creates a new GetPlayersInRoomHandler
func NewGetPlayersInRoomHandler(repo repo.RoomReader) *GetPlayersInRoomHandler {
	return &GetPlayersInRoomHandler{
		roomRepo: repo,
	}
}

// Handle processes the get players in room query
func (h *GetPlayersInRoomHandler) Handle(ctx context.Context, query GetPlayersInRoomQuery) ([]*userEntity.User, error) {
	return h.roomRepo.GetPlayersInRoom(query.RoomID)
}
