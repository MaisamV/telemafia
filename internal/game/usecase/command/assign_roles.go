package command

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"sort"
	"strconv"
	"telemafia/common"
	"telemafia/internal/game/entity"
	"telemafia/internal/game/repo"
	roomEntity "telemafia/internal/room/entity"
	scenarioEntity "telemafia/internal/scenario/entity"
	scenarioRepo "telemafia/internal/scenario/repo"
	userEntity "telemafia/internal/user/entity"
	"time"
)

// AssignRolesCommand represents a command to assign roles to players
type AssignRolesCommand struct {
	GameID entity.GameID         // Now using GameID instead of RoomID
	Users  []userEntity.User     // Optional; if empty, will get players from the game's room
	Roles  []scenarioEntity.Role // Optional; if empty, will use roles from game's scenario
}

// AssignRolesHandler handles role assignments
type AssignRolesHandler struct {
	gameRepo     repo.GameRepository
	scenarioRepo scenarioRepo.ScenarioReader
}

// NewAssignRolesHandler creates a new AssignRolesHandler
func NewAssignRolesHandler(gameRepo repo.GameRepository, scenarioRepo scenarioRepo.ScenarioReader) *AssignRolesHandler {
	return &AssignRolesHandler{
		gameRepo:     gameRepo,
		scenarioRepo: scenarioRepo,
	}
}

// Handle processes the assign roles command
func (h *AssignRolesHandler) Handle(cmd AssignRolesCommand) (map[userEntity.UserID]scenarioEntity.Role, error) {
	// Get the game by ID
	game, err := h.gameRepo.GetGameByID(cmd.GameID)
	if err != nil {
		log.Printf("Error getting game with ID '%s': %v", cmd.GameID, err)
		return nil, fmt.Errorf("game not found: %w", err)
	}

	log.Printf("Found game '%s' for room '%s'", game.ID, game.Room.ID)

	// Get roles - either use provided roles or fetch from scenario
	var roles []scenarioEntity.Role
	if len(cmd.Roles) > 0 {
		roles = cmd.Roles
	} else {
		// Fetch the scenario to get its roles
		scenario, err := h.scenarioRepo.GetScenarioByID(game.Scenario.ID)
		if err != nil {
			log.Printf("Error fetching scenario '%s': %v", game.Scenario.ID, err)
			return nil, fmt.Errorf("error fetching scenario: %w", err)
		}
		roles = scenario.Roles
		// Sort roles by their hash
		log.Printf("Using %d roles from scenario '%s'", len(roles), game.Scenario.ID)
	}
	sort.Slice(roles, func(i, j int) bool {
		return common.Hash(roles[i].Name) < common.Hash(roles[j].Name)
	})

	// Get users - either use provided users or get from the game's room
	var users []userEntity.User
	if len(cmd.Users) > 0 {
		users = cmd.Users
		// Sort users by ID
		sort.Slice(users, func(i, j int) bool {
			return users[i].ID < users[j].ID
		})
	}

	// Ensure we have enough roles for the players
	if len(roles) != len(users) {
		log.Printf("Not enough roles (%d) for all players (%d)", len(roles), len(users))
		return nil, errors.New("not enough roles for all players")
	}

	// Randomly assign roles to players
	rolesToAssign := make([]scenarioEntity.Role, len(roles))
	copy(rolesToAssign, roles)
	rand.Seed(time.Now().UTC().UnixNano())
	rand.Shuffle(len(rolesToAssign), func(i, j int) {
		rolesToAssign[i], rolesToAssign[j] = rolesToAssign[j], rolesToAssign[i]
	})

	// Store assignments in the game entity and prepare response map
	assignments := make(map[userEntity.UserID]scenarioEntity.Role)
	for i, user := range users {
		if i >= len(rolesToAssign) {
			break // Safety check
		}
		game.AssignRole(user.ID, rolesToAssign[i])
		assignments[user.ID] = rolesToAssign[i]
		log.Printf("Assigned role '%s' to user ID %d", rolesToAssign[i].Name, user.ID)
	}

	// Update game status
	game.SetRolesAssigned()
	log.Printf("Game '%s' status updated to RolesAssigned", game.ID)

	// Update the game in the repository
	err = h.gameRepo.UpdateGame(game)
	if err != nil {
		log.Printf("Error updating game '%s': %v", game.ID, err)
		return nil, fmt.Errorf("error updating game: %w", err)
	}

	log.Printf("Successfully assigned %d roles in game '%s'", len(assignments), game.ID)
	return assignments, nil
}

// GetAssignments gets the role assignments for a game by ID
func (h *AssignRolesHandler) GetAssignments(gameID string) (map[string]scenarioEntity.Role, error) {
	log.Printf("Getting assignments for game: '%s'", gameID)

	// Get the game by ID
	game, err := h.gameRepo.GetGameByID(entity.GameID(gameID))
	if err != nil {
		log.Printf("Error getting game with ID '%s': %v", gameID, err)
		return nil, fmt.Errorf("game not found: %w", err)
	}

	log.Printf("Found game '%s' with %d assignments", game.ID, len(game.Assignments))

	// Convert assignments to the expected format
	assignments := make(map[string]scenarioEntity.Role)
	for userID, role := range game.Assignments {
		assignments[strconv.FormatInt(int64(userID), 10)] = role
		log.Printf("Assignment for user %d: %s", userID, role.Name)
	}

	return assignments, nil
}

// GetAssignmentsByRoomID gets the role assignments for a room
func (h *AssignRolesHandler) GetAssignmentsByRoomID(roomID string) (map[string]scenarioEntity.Role, error) {
	log.Printf("Getting assignments for room: '%s'", roomID)

	// Get the game for this room
	game, err := h.gameRepo.GetGameByRoomID(roomEntity.RoomID(roomID))
	if err != nil {
		log.Printf("Error getting game for room '%s': %v", roomID, err)
		return nil, err
	}

	return h.GetAssignments(string(game.ID))
}
