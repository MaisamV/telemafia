package command

import (
	"context"
	"fmt"

	// gameEntity "telemafia/internal/game/entity"
	gameEntity "telemafia/internal/domain/game/entity"
	// gamePort "telemafia/internal/game/port"
	gamePort "telemafia/internal/domain/game/port"
	// roomEntity "telemafia/internal/room/entity"
	roomEntity "telemafia/internal/domain/room/entity"
	// scenarioEntity "telemafia/internal/scenario/entity"
	scenarioEntity "telemafia/internal/domain/scenario/entity"
	sharedEntity "telemafia/internal/shared/entity"
	"time"
)

// CreateGameCommand represents the command to create a new game
type CreateGameCommand struct {
	RoomID     roomEntity.RoomID // Use imported RoomID type
	ScenarioID string            // Scenario ID remains string based on entity
	// We might need actual Room and Scenario objects/pointers here
	// instead of just IDs to properly initialize the Game entity.
	// This depends on whether the Game entity *needs* the full objects
	// or just their references/IDs at creation time.
}

// CreateGameHandler handles game creation
type CreateGameHandler struct {
	gameRepo gamePort.GameRepository // Use imported GameRepository interface
	// May need RoomReader and ScenarioReader ports injected here
	// to fetch the actual Room and Scenario entities based on IDs.
}

// NewCreateGameHandler creates a new CreateGameHandler
func NewCreateGameHandler(repo gamePort.GameRepository) *CreateGameHandler {
	return &CreateGameHandler{
		gameRepo: repo,
	}
}

// Handle processes the create game command
func (h *CreateGameHandler) Handle(ctx context.Context, cmd CreateGameCommand) (*gameEntity.Game, error) {
	// TODO: Fetch the actual Room and Scenario entities using their repositories
	// room, err := h.roomRepo.GetRoomByID(cmd.RoomID)
	// scenario, err := h.scenarioRepo.GetScenarioByID(cmd.ScenarioID)
	// Handle errors...

	// Create a new game entity
	game := &gameEntity.Game{
		ID: gameEntity.GameID(fmt.Sprintf("game_%d", time.Now().UnixNano())), // Simple unique ID generation
		// Room:        room,     // Assign fetched room
		// Scenario:    scenario, // Assign fetched scenario
		Room:        &roomEntity.Room{ID: cmd.RoomID},             // TEMPORARY: Using placeholder ID
		Scenario:    &scenarioEntity.Scenario{ID: cmd.ScenarioID}, // TEMPORARY: Using placeholder ID
		State:       gameEntity.GameStateWaitingForPlayers,
		Assignments: make(map[sharedEntity.UserID]scenarioEntity.Role), // Use imported types
	}

	if err := h.gameRepo.CreateGame(game); err != nil {
		return nil, err // Propagates errors from repo
	}

	return game, nil
}
