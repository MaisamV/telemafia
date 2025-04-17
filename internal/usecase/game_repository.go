package usecase

import (
	"telemafia/internal/entity"
)

// GameReader defines the interface for reading game data
type GameReader interface {
	// GetGameByID gets a game by its ID
	GetGameByID(id entity.GameID) (*entity.Game, error)

	// GetGameByRoomID gets a game by room ID
	GetGameByRoomID(roomID entity.RoomID) (*entity.Game, error)

	// GetAllGames gets all games
	GetAllGames() ([]*entity.Game, error)
}

// GameWriter defines the interface for writing game data
type GameWriter interface {
	// CreateGame creates a new game
	CreateGame(game *entity.Game) error

	// UpdateGame updates an existing game
	UpdateGame(game *entity.Game) error

	// DeleteGame deletes a game by ID
	DeleteGame(id entity.GameID) error
}

// GameRepository defines the interface for game persistence
type GameRepository interface {
	GameReader
	GameWriter
}
