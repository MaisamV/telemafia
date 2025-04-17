package usecase

import (
	"context"
	"fmt"
	"telemafia/internal/entity"
	"time"
)

// CreateGameCommand represents the command to create a new game
type CreateGameCommand struct {
	RoomID     entity.RoomID
	ScenarioID string
}

// CreateGameHandler handles game creation
type CreateGameHandler struct {
	gameRepo GameRepository
}

// NewCreateGameHandler creates a new CreateGameHandler
func NewCreateGameHandler(repo GameRepository) *CreateGameHandler {
	return &CreateGameHandler{
		gameRepo: repo,
	}
}

// Handle processes the create game command
func (h *CreateGameHandler) Handle(ctx context.Context, cmd CreateGameCommand) (*entity.Game, error) {
	// Create a new game entity (basic setup, needs room and scenario info)
	game := &entity.Game{
		ID:          entity.GameID(fmt.Sprintf("game_%d", time.Now().UnixNano())),
		Room:        &entity.Room{ID: cmd.RoomID},         // Placeholder, need full room?
		Scenario:    &entity.Scenario{ID: cmd.ScenarioID}, // Placeholder, need full scenario?
		State:       entity.GameStateWaitingForPlayers,
		Assignments: make(map[entity.UserID]entity.Role),
	}

	if err := h.gameRepo.CreateGame(game); err != nil {
		return nil, err
	}

	return game, nil
}
