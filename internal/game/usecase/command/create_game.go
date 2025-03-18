package command

import (
	"fmt"
	"telemafia/internal/game/entity"
	"telemafia/internal/game/repo"
	roomEntity "telemafia/internal/room/entity"
	"time"
)

// CreateGameCommand represents a command to create a new game
type CreateGameCommand struct {
	RoomID       roomEntity.RoomID
	ScenarioID   string
	ScenarioName string
}

// CreateGameHandler handles the creation of new games
type CreateGameHandler struct {
	gameRepo repo.GameRepository
}

// NewCreateGameHandler creates a new CreateGameHandler
func NewCreateGameHandler(gameRepo repo.GameRepository) *CreateGameHandler {
	return &CreateGameHandler{
		gameRepo: gameRepo,
	}
}

// Handle processes the create game command
func (h *CreateGameHandler) Handle(cmd CreateGameCommand) (*entity.Game, error) {
	// Generate a unique game ID based on timestamp
	gameID := entity.GameID(fmt.Sprintf("game_%d", time.Now().UnixNano()))

	// Create a new game
	game := entity.NewGame(gameID, cmd.RoomID, cmd.ScenarioID, cmd.ScenarioName)

	// Save the game in the repository
	err := h.gameRepo.CreateGame(game)
	if err != nil {
		return nil, err
	}

	return game, nil
}
