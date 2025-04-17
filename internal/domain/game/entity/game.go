package entity

import (
	roomEntity "telemafia/internal/domain/room/entity"
	scenarioEntity "telemafia/internal/domain/scenario/entity"
	sharedEntity "telemafia/internal/shared/entity"
)

// GameID represents a unique game identifier
type GameID string

// Game represents a game entity with a scenario, room, and role assignments
type Game struct {
	ID          GameID
	State       GameState
	Room        *roomEntity.Room                            // Use imported Room type
	Scenario    *scenarioEntity.Scenario                    // Use imported Scenario type
	Assignments map[sharedEntity.UserID]scenarioEntity.Role // Use imported UserID and Role types
}

// GameState represents the current state of a game
type GameState string

const (
	// GameStateWaitingForPlayers means the game has been created but not started
	GameStateWaitingForPlayers GameState = "waiting_for_players"
	// GameStateRolesAssigned means roles have been assigned to players
	GameStateRolesAssigned GameState = "roles_assigned"
	// GameStateInProgress means the game is being played
	GameStateInProgress GameState = "in_progress"
	// GameStateFinished means the game has concluded
	GameStateFinished GameState = "finished"
)

// NewGame construction might be better suited within a use case/factory
// due to dependencies on specific Room and Scenario instances.

// AssignRole assigns a role to a player
func (g *Game) AssignRole(userID sharedEntity.UserID, role scenarioEntity.Role) {
	g.Assignments[userID] = role
}

// SetRolesAssigned updates the game state to roles assigned
func (g *Game) SetRolesAssigned() {
	g.State = GameStateRolesAssigned
}

// StartGame updates the game state to in progress
func (g *Game) StartGame() {
	g.State = GameStateInProgress
}

// FinishGame updates the game state to finished
func (g *Game) FinishGame() {
	g.State = GameStateFinished
}
