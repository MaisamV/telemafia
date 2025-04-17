package query

import (
	"context"
	"log"
	gameEntity "telemafia/internal/domain/game/entity"
	gamePort "telemafia/internal/domain/game/port"
)

// GetGamesQuery represents the query to get all games
type GetGamesQuery struct {
	// Add filters if needed (e.g., by state, by room)
}

// GetGamesHandler handles queries for all games
type GetGamesHandler struct {
	gameRepo gamePort.GameReader // Use imported GameReader interface
}

// NewGetGamesHandler creates a new GetGamesHandler
func NewGetGamesHandler(repo gamePort.GameReader) *GetGamesHandler {
	return &GetGamesHandler{
		gameRepo: repo,
	}
}

// Handle processes the get all games query
// Note: Query struct is empty, context is unused.
func (h *GetGamesHandler) Handle(ctx context.Context, query GetGamesQuery) ([]*gameEntity.Game, error) {
	log.Printf("Fetching all games")
	games, err := h.gameRepo.GetAllGames()
	if err != nil {
		log.Printf("Error fetching games: %v", err)
		return nil, err // Propagates errors from repo
	}
	log.Printf("Found %d games", len(games))
	return games, nil
}
