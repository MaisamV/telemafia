package port

import (
	gameEntity "telemafia/internal/domain/game/entity"
	roomEntity "telemafia/internal/domain/room/entity" // Needed for RoomID
)

// GameReader defines the interface for reading game data
type GameReader interface {
	// GetGameByID gets a game by its ID
	GetGameByID(id gameEntity.GameID) (*gameEntity.Game, error)

	// GetGameByRoomID gets a game by room ID
	GetGameByRoomID(roomID roomEntity.RoomID) (*gameEntity.Game, error)

	// GetAllGames gets all games
	GetAllGames() ([]*gameEntity.Game, error)
}

// GameWriter defines the interface for writing game data
type GameWriter interface {
	// CreateGame creates a new game
	CreateGame(game *gameEntity.Game) error

	// UpdateGame updates an existing game
	UpdateGame(game *gameEntity.Game) error

	// DeleteGame deletes a game by ID
	DeleteGame(id gameEntity.GameID) error
}

// GameRepository defines the interface for game persistence
type GameRepository interface {
	GameReader
	GameWriter
}
