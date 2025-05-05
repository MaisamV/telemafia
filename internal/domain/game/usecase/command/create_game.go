package command

import (
	"context"
	"errors"
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
	Requester  sharedEntity.User
	RoomID     roomEntity.RoomID // Use imported RoomID type
	ScenarioID string            // Scenario ID remains string based on entity
}

// CreateGameHandler handles game creation
type CreateGameHandler struct {
	gameRepo       gamePort.GameRepository // Use imported GameRepository interface
	roomClient     gamePort.RoomClient     // Use the room client interface
	scenarioClient gamePort.ScenarioClient // Use the scenario client interface
}

// NewCreateGameHandler creates a new CreateGameHandler
func NewCreateGameHandler(repo gamePort.GameRepository, roomClient gamePort.RoomClient, scenarioClient gamePort.ScenarioClient) *CreateGameHandler { // Add client dependencies
	return &CreateGameHandler{
		gameRepo:       repo,
		roomClient:     roomClient,     // Store room client
		scenarioClient: scenarioClient, // Store scenario client
	}
}

// Handle processes the create game command
func (h *CreateGameHandler) Handle(ctx context.Context, cmd CreateGameCommand) (*gameEntity.Game, error) {
	// Fetch the actual Room and Scenario entities using the clients
	room, err := h.roomClient.FetchRoom(cmd.RoomID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch room %s for game creation: %w", cmd.RoomID, err)
	}
	scenario, err := h.scenarioClient.FetchScenario(cmd.ScenarioID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch scenario %s for game creation: %w", cmd.ScenarioID, err)
	}

	// --- Permission Check ---
	// Allow if requester is global admin OR the moderator of this specific room
	isRoomModerator := room.Moderator != nil && room.Moderator.ID == cmd.Requester.ID
	if !cmd.Requester.Admin && !isRoomModerator {
		return nil, errors.New("create game: permission denied (requires admin or room moderator)")
	}

	// Create a new game entity
	game := &gameEntity.Game{
		ID:          gameEntity.GameID(fmt.Sprintf("game_%d", time.Now().UnixNano())), // Simple unique ID generation
		Room:        room,                                                             // Assign fetched room pointer
		Scenario:    scenario,                                                         // Assign fetched scenario pointer
		State:       gameEntity.GameStateWaitingForPlayers,
		Assignments: make(map[sharedEntity.UserID]scenarioEntity.Role), // Use imported types
	}

	if err := h.gameRepo.CreateGame(game); err != nil {
		return nil, err // Propagates errors from repo
	}

	return game, nil
}
