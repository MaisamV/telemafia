package usecase

import (
	// "fmt" // Removed unused import
	"telemafia/internal/entity"
)

// RoomReader defines the interface for reading room data
type RoomReader interface {
	// GetRoomByID Get room by ID
	GetRoomByID(id entity.RoomID) (*entity.Room, error)

	// GetRooms Get all rooms
	GetRooms() ([]*entity.Room, error)

	// GetPlayerRooms Get rooms for a player
	GetPlayerRooms(playerID entity.UserID) ([]*entity.Room, error)

	// GetPlayersInRoom Get players in a specific room
	GetPlayersInRoom(roomID entity.RoomID) ([]*entity.User, error)

	// CheckChangeFlag checks the current state of the change flag
	CheckChangeFlag() bool
}

// RoomWriter defines the interface for writing room data
type RoomWriter interface {
	// CreateRoom Create a new room
	CreateRoom(room *entity.Room) error

	// AddPlayerToRoom Add a player to a room
	AddPlayerToRoom(roomID entity.RoomID, player *entity.User) error

	// RemovePlayerFromRoom Remove a player from a room
	RemovePlayerFromRoom(roomID entity.RoomID, playerID entity.UserID) error

	// DeleteRoom deletes a room by ID
	DeleteRoom(roomID entity.RoomID) error

	// ConsumeChangeFlag checks and resets the change flag
	ConsumeChangeFlag() bool

	// RaiseChangeFlag sets the change flag to true
	RaiseChangeFlag()

	// AssignScenarioToRoom assigns a scenario to a room
	AssignScenarioToRoom(roomID entity.RoomID, scenarioName string) error

	// GetRoomScenario gets the scenario assigned to a room
	GetRoomScenario(roomID entity.RoomID) (string, error)
}

// Repository defines the interface for room persistence
type Repository interface {
	RoomReader
	RoomWriter
}
