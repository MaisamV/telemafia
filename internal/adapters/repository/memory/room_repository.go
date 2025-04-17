package memory

import (
	"fmt"
	"sync"

	// roomEntity "telemafia/internal/room/entity"
	roomEntity "telemafia/internal/domain/room/entity"
	// roomPort "telemafia/internal/room/port"
	roomPort "telemafia/internal/domain/room/port"
	sharedEntity "telemafia/internal/shared/entity"
)

// Ensure InMemoryRepository implements the roomPort.RoomRepository interface.
var _ roomPort.RoomRepository = (*InMemoryRoomRepository)(nil)

type InMemoryRoomRepository struct {
	rooms map[roomEntity.RoomID]*roomEntity.Room
	mutex sync.RWMutex

	// This change flag is specific to the in-memory implementation and might
	// not belong in the core domain or repository interface if persistence changes.
	changeFlag bool
	// roomToScenario map[roomEntity.RoomID]string // Moved this logic? Room entity now has ScenarioName
}

// NewInMemoryRoomRepository creates a new in-memory room repository
func NewInMemoryRoomRepository() roomPort.RoomRepository { // Return the port interface
	return &InMemoryRoomRepository{
		rooms: make(map[roomEntity.RoomID]*roomEntity.Room),
		// roomToScenario: make(map[roomEntity.RoomID]string),
	}
}

func (r *InMemoryRoomRepository) CreateRoom(room *roomEntity.Room) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.rooms[room.ID]; exists {
		// Consider using a more specific error from roomEntity if defined
		return fmt.Errorf("room with ID %s already exists", room.ID) // Use roomEntity.ErrRoomAlreadyExists?
	}

	r.rooms[room.ID] = room
	r.changeFlag = true // Mark change
	return nil
}

func (r *InMemoryRoomRepository) GetRoomByID(id roomEntity.RoomID) (*roomEntity.Room, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	room, exists := r.rooms[id]
	if !exists {
		return nil, roomEntity.ErrRoomNotFound // Use error from entity package
	}

	return room, nil
}

func (r *InMemoryRoomRepository) GetRooms() ([]*roomEntity.Room, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	rooms := make([]*roomEntity.Room, 0, len(r.rooms))
	for _, room := range r.rooms {
		rooms = append(rooms, room)
	}

	return rooms, nil
}

func (r *InMemoryRoomRepository) GetPlayerRooms(playerID sharedEntity.UserID) ([]*roomEntity.Room, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var playerRooms []*roomEntity.Room
	for _, room := range r.rooms {
		for _, player := range room.Players {
			if player != nil && player.ID == playerID { // Add nil check
				playerRooms = append(playerRooms, room)
				break
			}
		}
	}

	return playerRooms, nil
}

func (r *InMemoryRoomRepository) AddPlayerToRoom(roomID roomEntity.RoomID, player *sharedEntity.User) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	room, exists := r.rooms[roomID]
	if !exists {
		return roomEntity.ErrRoomNotFound
	}

	for _, p := range room.Players {
		if p != nil && p.ID == player.ID { // Add nil check
			// Player already in room, return success or specific error?
			// return roomEntity.ErrPlayerAlreadyInRoom (if defined)
			return nil // Current behavior
		}
	}

	room.Players = append(room.Players, player)
	r.changeFlag = true // Mark change
	return nil
}

func (r *InMemoryRoomRepository) RemovePlayerFromRoom(roomID roomEntity.RoomID, playerID sharedEntity.UserID) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	room, exists := r.rooms[roomID]
	if !exists {
		return roomEntity.ErrRoomNotFound
	}

	found := false
	newPlayers := make([]*sharedEntity.User, 0, len(room.Players))
	for _, p := range room.Players {
		if p != nil && p.ID == playerID { // Add nil check
			found = true
		} else {
			newPlayers = append(newPlayers, p)
		}
	}

	if !found {
		return roomEntity.ErrPlayerNotInRoom // Use error from entity package
	}

	room.Players = newPlayers
	r.changeFlag = true // Mark change
	return nil
}

// GetPlayersInRoom returns the players in a specific room
func (r *InMemoryRoomRepository) GetPlayersInRoom(roomID roomEntity.RoomID) ([]*sharedEntity.User, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	room, exists := r.rooms[roomID]
	if !exists {
		return nil, roomEntity.ErrRoomNotFound
	}

	// Return a copy to prevent external modification?
	playersCopy := make([]*sharedEntity.User, len(room.Players))
	copy(playersCopy, room.Players)
	return playersCopy, nil
}

// DeleteRoom deletes a room by ID
func (r *InMemoryRoomRepository) DeleteRoom(roomID roomEntity.RoomID) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.rooms[roomID]; !exists {
		// Return specific error? The interface doesn't guarantee it for delete.
		return roomEntity.ErrRoomNotFound
	}

	delete(r.rooms, roomID)
	// delete(r.roomToScenario, roomID) // Removed, scenario name is in Room entity
	r.changeFlag = true // Mark change
	return nil
}

// CheckChangeFlag checks the current state of the change flag
func (r *InMemoryRoomRepository) CheckChangeFlag() bool {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	return r.changeFlag
}

// ConsumeChangeFlag checks and resets the change flag
func (r *InMemoryRoomRepository) ConsumeChangeFlag() bool {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	changed := r.changeFlag
	r.changeFlag = false
	return changed
}

// RaiseChangeFlag sets the change flag to true
func (r *InMemoryRoomRepository) RaiseChangeFlag() {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.changeFlag = true
}

// AssignScenarioToRoom assigns a scenario to a room
func (r *InMemoryRoomRepository) AssignScenarioToRoom(roomID roomEntity.RoomID, scenarioName string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	room, exists := r.rooms[roomID]
	if !exists {
		return roomEntity.ErrRoomNotFound
	}
	// Update the room entity directly
	room.ScenarioName = scenarioName
	// r.roomToScenario[roomID] = scenarioName // Removed
	r.changeFlag = true // Mark change
	return nil
}

// GetRoomScenario gets the scenario assigned to a room
func (r *InMemoryRoomRepository) GetRoomScenario(roomID roomEntity.RoomID) (string, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	room, exists := r.rooms[roomID]
	if !exists {
		return "", roomEntity.ErrRoomNotFound
	}
	// scenarioName, ok := r.roomToScenario[roomID] // Removed
	// if !ok {
	// 	return "", fmt.Errorf("no scenario assigned to room %s", roomID)
	// }
	return room.ScenarioName, nil // Get from entity
}
