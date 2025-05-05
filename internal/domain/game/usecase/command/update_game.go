package command

import (
	"context"
	"fmt"
	gameEntity "telemafia/internal/domain/game/entity"
	gamePort "telemafia/internal/domain/game/port"
)

// UpdateGameCommand represents the command to update a game entity.
// It includes the full game object to allow updating any field.
type UpdateGameCommand struct {
	Game *gameEntity.Game
	// Requester sharedEntity.User // Optional: Add requester if permission needed
}

// UpdateGameHandler handles updating a game.
type UpdateGameHandler struct {
	gameRepo gamePort.GameRepository
}

// NewUpdateGameHandler creates a new UpdateGameHandler.
func NewUpdateGameHandler(repo gamePort.GameRepository) *UpdateGameHandler {
	return &UpdateGameHandler{
		gameRepo: repo,
	}
}

// Handle processes the update game command.
func (h *UpdateGameHandler) Handle(ctx context.Context, cmd UpdateGameCommand) error {
	if cmd.Game == nil {
		return fmt.Errorf("update game: game data cannot be nil")
	}

	// Optional: Add permission checks here if needed based on the Requester.
	// For now, assume permission is checked upstream or not required for simple updates.

	if err := h.gameRepo.UpdateGame(cmd.Game); err != nil {
		return fmt.Errorf("update game: failed to update game %s: %w", cmd.Game.ID, err)
	}
	return nil
}
