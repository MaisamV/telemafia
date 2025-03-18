package query

import (
	"context"
	"log"
	"telemafia/internal/game/entity"
	"telemafia/internal/game/repo"
)

// GetGamesQuery represents the query to get all games
type GetGamesQuery struct{}

// GetGamesHandler handles queries for all games
type GetGamesHandler struct {
	gameRepo repo.GameRepository
}

// NewGetGamesHandler creates a new GetGamesHandler
func NewGetGamesHandler(repo repo.GameRepository) *GetGamesHandler {
	return &GetGamesHandler{
		gameRepo: repo,
	}
}

// Handle processes the get all games query
func (h *GetGamesHandler) Handle(ctx context.Context) ([]*entity.Game, error) {
	log.Printf("Fetching all games")
	games, err := h.gameRepo.GetAllGames()
	if err != nil {
		log.Printf("Error fetching games: %v", err)
		return nil, err
	}
	log.Printf("Found %d games", len(games))
	return games, nil
}
