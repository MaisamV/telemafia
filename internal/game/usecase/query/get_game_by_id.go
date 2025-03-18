package query

import (
	"context"
	"log"
	"telemafia/internal/game/entity"
	"telemafia/internal/game/repo"
)

// GetGameByIDQuery represents the query to get a game by its ID
type GetGameByIDQuery struct {
	ID entity.GameID
}

// GetGameByIDHandler handles queries for a specific game by ID
type GetGameByIDHandler struct {
	gameRepo repo.GameRepository
}

// NewGetGameByIDHandler creates a new GetGameByIDHandler
func NewGetGameByIDHandler(repo repo.GameRepository) *GetGameByIDHandler {
	return &GetGameByIDHandler{
		gameRepo: repo,
	}
}

// Handle processes the get game by ID query
func (h *GetGameByIDHandler) Handle(ctx context.Context, query GetGameByIDQuery) (*entity.Game, error) {
	log.Printf("Fetching game with ID: %s", query.ID)
	game, err := h.gameRepo.GetGameByID(query.ID)
	if err != nil {
		log.Printf("Error fetching game with ID %s: %v", query.ID, err)
		return nil, err
	}
	log.Printf("Successfully retrieved game with ID %s", query.ID)
	return game, nil
}
