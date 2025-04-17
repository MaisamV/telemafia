package query

import (
	"context"
	"log"
	gameEntity "telemafia/internal/domain/game/entity"
	gamePort "telemafia/internal/domain/game/port"
)

// GetGameByIDQuery represents the query to get a game by ID
type GetGameByIDQuery struct {
	ID gameEntity.GameID // Use imported GameID type
}

// GetGameByIDHandler handles queries for a specific game
type GetGameByIDHandler struct {
	gameRepo gamePort.GameReader // Use imported GameReader interface
}

// NewGetGameByIDHandler creates a new GetGameByIDHandler
func NewGetGameByIDHandler(repo gamePort.GameReader) *GetGameByIDHandler {
	return &GetGameByIDHandler{
		gameRepo: repo,
	}
}

// Handle processes the get game by ID query
func (h *GetGameByIDHandler) Handle(ctx context.Context, query GetGameByIDQuery) (*gameEntity.Game, error) {
	log.Printf("Fetching game with ID: %s", query.ID)
	game, err := h.gameRepo.GetGameByID(query.ID)
	if err != nil {
		log.Printf("Error fetching game with ID %s: %v", query.ID, err)
		return nil, err // Propagates errors from repo
	}
	log.Printf("Successfully retrieved game with ID %s", query.ID)
	return game, nil
}
