package command

import (
	"context" // Added context parameter to Handle method
	"errors"
	"fmt"
	"log"
	"sort"
	gameEntity "telemafia/internal/domain/game/entity"
	gamePort "telemafia/internal/domain/game/port"
	roomPort "telemafia/internal/domain/room/port" // Use imported roomPort
	scenarioEntity "telemafia/internal/domain/scenario/entity"
	scenarioPort "telemafia/internal/domain/scenario/port"
	sharedEntity "telemafia/internal/shared/entity"
)

// AssignRolesCommand represents a command to assign roles to players in a game
type AssignRolesCommand struct {
	Requester sharedEntity.User // Added
	GameID    gameEntity.GameID // Use imported GameID type
	// Users and Roles might be fetched within the handler based on GameID
}

// AssignRolesHandler handles role assignments
type AssignRolesHandler struct {
	gameRepo     gamePort.GameRepository     // Use imported GameRepository interface
	scenarioRepo scenarioPort.ScenarioReader // Use imported ScenarioReader interface
	roomRepo     roomPort.RoomReader         // Use imported roomPort
}

// NewAssignRolesHandler creates a new AssignRolesHandler
func NewAssignRolesHandler(gameRepo gamePort.GameRepository, scenarioRepo scenarioPort.ScenarioReader, roomRepo roomPort.RoomReader) *AssignRolesHandler {
	return &AssignRolesHandler{
		gameRepo:     gameRepo,
		scenarioRepo: scenarioRepo,
		roomRepo:     roomRepo,
	}
}

// Handle processes the assign roles command
func (h *AssignRolesHandler) Handle(ctx context.Context, cmd AssignRolesCommand) (map[sharedEntity.User]scenarioEntity.Role, error) { // Updated return type
	// Get the game by ID
	game, err := h.gameRepo.GetGameByID(cmd.GameID)
	if err != nil {
		log.Printf("Error getting game with ID '%s': %v", cmd.GameID, err)
		return nil, fmt.Errorf("game '%s' not found: %w", cmd.GameID, err)
	}

	log.Printf("Found game '%s' for room '%s'", game.ID, game.Room.ID)

	// --- Permission Check ---
	// Must have room info to check moderator status
	if game.Room == nil {
		return nil, errors.New("assign roles: game has no associated room")
	}
	isRoomModerator := game.Room.Moderator != nil && game.Room.Moderator.ID == cmd.Requester.ID
	if !cmd.Requester.Admin && !isRoomModerator {
		return nil, errors.New("assign roles: permission denied (requires admin or room moderator)")
	}

	// Fetch the scenario to get its roles
	if game.Scenario == nil {
		log.Printf("Game '%s' has no scenario assigned", game.ID)
		return nil, errors.New("game has no scenario assigned")
	}
	scenario, err := h.scenarioRepo.GetScenarioByID(game.Scenario.ID)
	if err != nil {
		log.Printf("Error fetching scenario '%s': %v", game.Scenario.ID, err)
		return nil, fmt.Errorf("error fetching scenario '%s': %w", game.Scenario.ID, err)
	}

	// Fetch players from the game's room
	// Room presence already checked above
	players, err := h.roomRepo.GetPlayersInRoom(game.Room.ID)
	if err != nil {
		log.Printf("Error fetching players for room '%s': %v", game.Room.ID, err)
		return nil, fmt.Errorf("error fetching players for room '%s': %w", game.Room.ID, err)
	}
	// Convert []*User to []User if needed by sorting/assignment logic
	// Assuming the repository returns []*sharedEntity.User
	users := make([]sharedEntity.User, 0, len(players))
	for _, p := range players {
		if p != nil { // Add nil check for safety
			users = append(users, *p)
		}
	}

	sort.Slice(users, func(i, j int) bool {
		return users[i].ID < users[j].ID
	})
	log.Printf("Found %d players in room '%s'", len(users), game.Room.ID)

	rolesToAssign := scenario.GetShuffledRoles(len(users))
	// Ensure we have enough roles for the players
	if len(rolesToAssign) != len(users) {
		msg := fmt.Sprintf("role count (%d) does not match player count (%d) for game '%s'", len(rolesToAssign), len(users), game.ID)
		log.Println(msg)
		return nil, errors.New(msg)
	}

	// Store assignments in the game entity and prepare response map
	assignments := make(map[sharedEntity.User]scenarioEntity.Role)
	for i, user := range users {
		if i >= len(rolesToAssign) {
			break // Safety check
		}
		game.AssignRole(user.ID, rolesToAssign[i])
		assignments[user] = rolesToAssign[i]
		log.Printf("Assigned role '%s' to user ID %d", rolesToAssign[i].Name, user.ID)
	}

	// Update game status
	game.SetRolesAssigned()
	log.Printf("Game '%s' status updated to RolesAssigned", game.ID)

	// Update the game in the repository
	err = h.gameRepo.UpdateGame(game)
	if err != nil {
		log.Printf("Error updating game '%s': %v", game.ID, err)
		return nil, fmt.Errorf("error updating game '%s': %w", game.ID, err)
	}

	log.Printf("Successfully assigned %d roles in game '%s'", len(assignments), game.ID)
	return assignments, nil
}

// --- Methods GetAssignments and GetAssignmentsByRoomID removed ---
// These seemed like query logic (reading assignments) rather than part of the
// AssignRoles command handler. They should likely be separate query handlers.
// // GetAssignments gets the role assignments for a game by ID
// func (h *AssignRolesHandler) GetAssignments(gameID string) (map[string]scenarioEntity.Role, error) {
// 	...
// }
// // GetAssignmentsByRoomID gets the role assignments for a room
// func (h *AssignRolesHandler) GetAssignmentsByRoomID(roomID string) (map[string]scenarioEntity.Role, error) {
// 	...
// }
