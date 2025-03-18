package repo

import (
	"fmt"
	"telemafia/internal/room/entity"
	userEntity "telemafia/internal/user/entity"
)

// RoomReader defines the interface for reading room data
type RoomReader interface {
	// GetRoomByID Get room by ID
	GetRoomByID(id entity.RoomID) (*entity.Room, error)

	// GetRooms Get all rooms
	GetRooms() ([]*entity.Room, error)

	// GetPlayerRooms Get rooms for a player
	GetPlayerRooms(playerID userEntity.UserID) ([]*entity.Room, error)

	// GetPlayersInRoom Get players in a specific room
	GetPlayersInRoom(roomID entity.RoomID) ([]*userEntity.User, error)

	// CheckChangeFlag checks the current state of the change flag
	CheckChangeFlag() bool

	// GetRoomScenario gets the scenario assigned to a room
	GetRoomScenario(roomID entity.RoomID) (string, error)
}

// RoomWriter defines the interface for writing room data
type RoomWriter interface {
	// CreateRoom Create a new room
	CreateRoom(room *entity.Room) error

	// AddPlayerToRoom Add a player to a room
	AddPlayerToRoom(roomID entity.RoomID, player *userEntity.User) error

	// RemovePlayerFromRoom Remove a player from a room
	RemovePlayerFromRoom(roomID entity.RoomID, playerID userEntity.UserID) error

	// DeleteRoom deletes a room by ID
	DeleteRoom(roomID entity.RoomID) error

	// ConsumeChangeFlag checks and resets the change flag
	ConsumeChangeFlag() bool

	// RaiseChangeFlag sets the change flag to true
	RaiseChangeFlag()

	// AssignScenarioToRoom assigns a scenario to a room
	AssignScenarioToRoom(roomID entity.RoomID, scenarioName string) error
}

// Repository defines the interface for room persistence
type Repository interface {
	RoomReader
	RoomWriter
}

// InMemoryRepository provides an in-memory implementation of Repository
type InMemoryRepository struct {
	rooms          map[entity.RoomID]*entity.Room
	changeFlag     bool
	roomToScenario map[entity.RoomID]string
}

// NewInMemoryRepository creates a new in-memory repository
func NewInMemoryRepository() *InMemoryRepository {
	return &InMemoryRepository{
		rooms:          make(map[entity.RoomID]*entity.Room),
		changeFlag:     false,
		roomToScenario: make(map[entity.RoomID]string),
	}
}

// AssignScenarioToRoom assigns a scenario to a room
func (r *InMemoryRepository) AssignScenarioToRoom(roomID entity.RoomID, scenarioName string) error {
	if _, exists := r.rooms[roomID]; !exists {
		return fmt.Errorf("room not found")
	}
	r.roomToScenario[roomID] = scenarioName
	r.RaiseChangeFlag()
	return nil
}

// GetRoomScenario gets the scenario assigned to a room
func (r *InMemoryRepository) GetRoomScenario(roomID entity.RoomID) (string, error) {
	if _, exists := r.rooms[roomID]; !exists {
		return "", fmt.Errorf("room not found")
	}
	return r.roomToScenario[roomID], nil
}

// RaiseChangeFlag sets the change flag to true
func (r *InMemoryRepository) RaiseChangeFlag() {
	r.changeFlag = true
}
