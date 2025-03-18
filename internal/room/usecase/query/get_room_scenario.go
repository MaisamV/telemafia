package query

import (
	"context"
	"telemafia/internal/room/entity"
	"telemafia/internal/room/repo"
)

// GetRoomScenarioQuery represents a query to get the scenario assigned to a room
type GetRoomScenarioQuery struct {
	RoomID entity.RoomID
}

// GetRoomScenarioHandler handles the get room scenario query
type GetRoomScenarioHandler struct {
	roomRepo repo.Repository
}

// NewGetRoomScenarioHandler creates a new get room scenario handler
func NewGetRoomScenarioHandler(roomRepo repo.Repository) *GetRoomScenarioHandler {
	return &GetRoomScenarioHandler{
		roomRepo: roomRepo,
	}
}

// Handle gets the scenario assigned to a room
func (h *GetRoomScenarioHandler) Handle(ctx context.Context, query GetRoomScenarioQuery) (string, error) {
	return h.roomRepo.GetRoomScenario(query.RoomID)
}
