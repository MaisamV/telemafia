package entity

// GameID represents a unique game identifier
type GameID string

// Game represents a game entity with a scenario, room, and role assignments
type Game struct {
	ID          GameID
	State       GameState
	Room        *Room
	Scenario    *Scenario
	Assignments map[UserID]Role
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

// NewGame creates a new game instance
// Note: Removed - constructor logic seems complex here, better handled in use case?
// func NewGame(id GameID, roomID RoomID, scenarioID string, scenarioName string) *Game {
// 	 // Create temporary Room and Scenario objects?
// 	 room := &Room{ID: roomID}
// 	 scenario := &Scenario{ID: scenarioID, Name: scenarioName}
//
// 	 return &Game{
// 		 ID:          id,
// 		 State:       GameStateWaitingForPlayers,
// 		 Room:        room,
// 		 Scenario:    scenario,
// 		 Assignments: make(map[UserID]Role),
// 	 }
// }

// AssignRole assigns a role to a player
func (g *Game) AssignRole(userID UserID, role Role) {
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
