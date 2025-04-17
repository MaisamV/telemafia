package usecase

import (
	"context"
	"log"
	"telemafia/internal/entity"
)

// GetGameByIDQuery represents the query to get a game by ID
type GetGameByIDQuery struct {
	ID entity.GameID
}

// GetGameByIDHandler handles queries for a specific game
type GetGameByIDHandler struct {
	gameRepo GameRepository
}

// NewGetGameByIDHandler creates a new GetGameByIDHandler
func NewGetGameByIDHandler(repo GameRepository) *GetGameByIDHandler {
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
