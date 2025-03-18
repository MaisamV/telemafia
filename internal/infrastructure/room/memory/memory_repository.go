package memory

import (
	"sync"
	error2 "telemafia/common/error"
	"telemafia/internal/room/entity"
	"telemafia/internal/room/repo"
	userEntity "telemafia/internal/user/entity"
)

type InMemoryRepository struct {
	rooms map[entity.RoomID]*entity.Room
	mutex sync.RWMutex

	changeFlag      bool
	changeFlagMutex sync.Mutex
	roomToScenario  map[entity.RoomID]string
}

// NewInMemoryRepository creates a new in-memory repository
func NewInMemoryRepository() repo.Repository {
	return &InMemoryRepository{
		rooms:          make(map[entity.RoomID]*entity.Room),
		roomToScenario: make(map[entity.RoomID]string),
	}
}

func (r *InMemoryRepository) CreateRoom(room *entity.Room) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.rooms[room.ID]; exists {
		return error2.ErrRoomAlreadyExists
	}

	r.rooms[room.ID] = room
	r.changeFlagMutex.Lock()
	r.changeFlag = true
	r.changeFlagMutex.Unlock()
	return nil
}

func (r *InMemoryRepository) GetRoomByID(id entity.RoomID) (*entity.Room, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	room, exists := r.rooms[id]
	if !exists {
		return nil, error2.ErrRoomNotFound
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

func (r *InMemoryRepository) GetPlayerRooms(playerID userEntity.UserID) ([]*entity.Room, error) {
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

func (r *InMemoryRepository) AddPlayerToRoom(roomID entity.RoomID, player *userEntity.User) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	room, exists := r.rooms[roomID]
	if !exists {
		return error2.ErrRoomNotFound
	}

	// Check if player is already in the room
	for _, p := range room.Players {
		if p.ID == player.ID {
			return nil
		}
	}

	room.Players = append(room.Players, player)
	r.changeFlagMutex.Lock()
	r.changeFlag = true
	r.changeFlagMutex.Unlock()
	return nil
}

func (r *InMemoryRepository) RemovePlayerFromRoom(roomID entity.RoomID, playerID userEntity.UserID) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	room, exists := r.rooms[roomID]
	if !exists {
		return error2.ErrRoomNotFound
	}

	for i, player := range room.Players {
		if player.ID == playerID {
			// Remove player by swapping with last element and truncating
			room.Players[i] = room.Players[len(room.Players)-1]
			room.Players = room.Players[:len(room.Players)-1]
			r.changeFlagMutex.Lock()
			r.changeFlag = true
			r.changeFlagMutex.Unlock()
			return nil
		}
	}

	return error2.ErrPlayerNotInRoom
}

// GetPlayersInRoom returns the players in a specific room
func (r *InMemoryRepository) GetPlayersInRoom(roomID entity.RoomID) ([]*userEntity.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	room, exists := r.rooms[roomID]
	if !exists {
		return nil, error2.ErrRoomNotFound
	}

	return room.Players, nil
}

// DeleteRoom deletes a room by ID
func (r *InMemoryRepository) DeleteRoom(roomID entity.RoomID) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.rooms[roomID]; !exists {
		return error2.ErrRoomNotFound
	}

	delete(r.rooms, roomID)
	r.changeFlagMutex.Lock()
	r.changeFlag = true
	r.changeFlagMutex.Unlock()
	return nil
}

// CheckChangeFlag checks the current state of the change flag
func (r *InMemoryRepository) CheckChangeFlag() bool {
	r.changeFlagMutex.Lock()
	defer r.changeFlagMutex.Unlock()

	return r.changeFlag
}

// ConsumeChangeFlag checks and resets the change flag
func (r *InMemoryRepository) ConsumeChangeFlag() bool {
	r.changeFlagMutex.Lock()
	defer r.changeFlagMutex.Unlock()

	changed := r.changeFlag
	r.changeFlag = false
	return changed
}

// RaiseChangeFlag sets the change flag to true
func (r *InMemoryRepository) RaiseChangeFlag() {
	r.changeFlagMutex.Lock()
	r.changeFlag = true
	r.changeFlagMutex.Unlock()
}

// AssignScenarioToRoom assigns a scenario to a room
func (r *InMemoryRepository) AssignScenarioToRoom(roomID entity.RoomID, scenarioName string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.rooms[roomID]; !exists {
		return error2.ErrRoomNotFound
	}
	r.roomToScenario[roomID] = scenarioName
	r.RaiseChangeFlag()
	return nil
}

// GetRoomScenario gets the scenario assigned to a room
func (r *InMemoryRepository) GetRoomScenario(roomID entity.RoomID) (string, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	if _, exists := r.rooms[roomID]; !exists {
		return "", error2.ErrRoomNotFound
	}
	return r.roomToScenario[roomID], nil
}
