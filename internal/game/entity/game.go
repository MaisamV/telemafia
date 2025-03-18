package entity

import (
	roomEntity "telemafia/internal/room/entity"
	scenarioEntity "telemafia/internal/scenario/entity"
	userEntity "telemafia/internal/user/entity"
)

// GameID represents a unique game identifier
type GameID string

// Game represents a game entity with a scenario, room, and role assignments
type Game struct {
	ID          GameID
	Status      GameStatus
	Room        *roomEntity.Room
	Scenario    *scenarioEntity.Scenario
	Assignments map[userEntity.UserID]scenarioEntity.Role
}

// GameStatus represents the current status of a game
type GameStatus string

const (
	// GameStatusInitialized means the game has been created but roles are not assigned yet
	GameStatusInitialized GameStatus = "initialized"
	// GameStatusRolesAssigned means roles have been assigned to players
	GameStatusRolesAssigned GameStatus = "roles_assigned"
	// GameStatusInProgress means the game is being played
	GameStatusInProgress GameStatus = "in_progress"
	// GameStatusFinished means the game has concluded
	GameStatusFinished GameStatus = "finished"
)

// NewGame creates a new game instance
func NewGame(id GameID, roomID roomEntity.RoomID, scenarioID string, scenarioName string) *Game {
	// Create temporary Room and Scenario objects
	room := &roomEntity.Room{
		ID: roomID,
	}

	scenario := &scenarioEntity.Scenario{
		ID:   scenarioID,
		Name: scenarioName,
	}

	return &Game{
		ID:          id,
		Status:      GameStatusInitialized,
		Room:        room,
		Scenario:    scenario,
		Assignments: make(map[userEntity.UserID]scenarioEntity.Role),
	}
}

// AssignRole assigns a role to a player
func (g *Game) AssignRole(userID userEntity.UserID, role scenarioEntity.Role) {
	g.Assignments[userID] = role
}

// SetRolesAssigned updates the game status to roles assigned
func (g *Game) SetRolesAssigned() {
	g.Status = GameStatusRolesAssigned
}

// StartGame updates the game status to in progress
func (g *Game) StartGame() {
	g.Status = GameStatusInProgress
}

// FinishGame updates the game status to finished
func (g *Game) FinishGame() {
	g.Status = GameStatusFinished
}
