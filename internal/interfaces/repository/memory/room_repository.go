package memory

import (
	"fmt"
	"sync"
	"telemafia/internal/entity"
	"telemafia/internal/usecase"
)

// Ensure InMemoryRepository implements the usecase.Repository interface.
var _ usecase.Repository = (*InMemoryRepository)(nil)

type InMemoryRepository struct {
	rooms map[entity.RoomID]*entity.Room
	mutex sync.RWMutex

	changeFlag     bool
	roomToScenario map[entity.RoomID]string
}

// NewInMemoryRepository creates a new in-memory repository
func NewInMemoryRepository() usecase.Repository {
	return &InMemoryRepository{
		rooms:          make(map[entity.RoomID]*entity.Room),
		roomToScenario: make(map[entity.RoomID]string),
	}
}

func (r *InMemoryRepository) CreateRoom(room *entity.Room) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.rooms[room.ID]; exists {
		return entity.ErrRoomAlreadyExists
	}

	r.rooms[room.ID] = room
	r.changeFlag = true
	return nil
}

func (r *InMemoryRepository) GetRoomByID(id entity.RoomID) (*entity.Room, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	room, exists := r.rooms[id]
	if !exists {
		return nil, entity.ErrRoomNotFound
	}

	return room, nil
}

func (r *InMemoryRepository) GetRooms() ([]*entity.Room, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	rooms := make([]*entity.Room, 0, len(r.rooms))
	for _, room := range r.rooms {
		rooms = append(rooms, room)
	}

	return rooms, nil
}

func (r *InMemoryRepository) GetPlayerRooms(playerID entity.UserID) ([]*entity.Room, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var playerRooms []*entity.Room
	for _, room := range r.rooms {
		for _, player := range room.Players {
			if player.ID == playerID {
				playerRooms = append(playerRooms, room)
				break
			}
		}
	}

	return playerRooms, nil
}

func (r *InMemoryRepository) AddPlayerToRoom(roomID entity.RoomID, player *entity.User) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	room, exists := r.rooms[roomID]
	if !exists {
		return entity.ErrRoomNotFound
	}

	for _, p := range room.Players {
		if p.ID == player.ID {
			return nil
		}
	}

	room.Players = append(room.Players, player)
	r.changeFlag = true
	return nil
}

func (r *InMemoryRepository) RemovePlayerFromRoom(roomID entity.RoomID, playerID entity.UserID) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	room, exists := r.rooms[roomID]
	if !exists {
		return entity.ErrRoomNotFound
	}

	for i, player := range room.Players {
		if player.ID == playerID {
			room.Players[i] = room.Players[len(room.Players)-1]
			room.Players = room.Players[:len(room.Players)-1]
			r.changeFlag = true
			return nil
		}
	}

	return entity.ErrPlayerNotInRoom
}

// GetPlayersInRoom returns the players in a specific room
func (r *InMemoryRepository) GetPlayersInRoom(roomID entity.RoomID) ([]*entity.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	room, exists := r.rooms[roomID]
	if !exists {
		return nil, entity.ErrRoomNotFound
	}

	return room.Players, nil
}

// DeleteRoom deletes a room by ID
func (r *InMemoryRepository) DeleteRoom(roomID entity.RoomID) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.rooms[roomID]; !exists {
		return entity.ErrRoomNotFound
	}

	delete(r.rooms, roomID)
	delete(r.roomToScenario, roomID)
	r.changeFlag = true
	return nil
}

// CheckChangeFlag checks the current state of the change flag
func (r *InMemoryRepository) CheckChangeFlag() bool {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	return r.changeFlag
}

// ConsumeChangeFlag checks and resets the change flag
func (r *InMemoryRepository) ConsumeChangeFlag() bool {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	changed := r.changeFlag
	r.changeFlag = false
	return changed
}

// RaiseChangeFlag sets the change flag to true
func (r *InMemoryRepository) RaiseChangeFlag() {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.changeFlag = true
}

// AssignScenarioToRoom assigns a scenario to a room
func (r *InMemoryRepository) AssignScenarioToRoom(roomID entity.RoomID, scenarioName string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if _, exists := r.rooms[roomID]; !exists {
		return entity.ErrRoomNotFound
	}
	r.roomToScenario[roomID] = scenarioName
	r.changeFlag = true
	return nil
}

// GetRoomScenario gets the scenario assigned to a room
func (r *InMemoryRepository) GetRoomScenario(roomID entity.RoomID) (string, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	if _, exists := r.rooms[roomID]; !exists {
		return "", entity.ErrRoomNotFound
	}
	scenarioName, ok := r.roomToScenario[roomID]
	if !ok {
		return "", fmt.Errorf("no scenario assigned to room %s", roomID)
	}
	return scenarioName, nil
}
