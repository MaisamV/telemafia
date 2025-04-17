package entity

import (
	"errors"
	sharedEntity "telemafia/internal/shared/entity" // Updated import for User
	"time"
)

// RoomID represents a unique room identifier
// Keeping RoomID definition here as it seems specific to the Room entity context
type RoomID string

// Room represents a game room entity
type Room struct {
	ID           RoomID
	Name         string
	CreatedAt    time.Time
	Players      []*sharedEntity.User // Use imported User type
	Description  map[string]string
	ScenarioName string
}

// Predefined error variables (using standard errors)
var (
	ErrInvalidRoomName   = errors.New("invalid room name")
	ErrRoomNotFound      = errors.New("room not found") // Note: Potentially duplicate error in old entity/room.go? Consolidate later if needed.
	ErrRoomAlreadyExists = errors.New("room already exists")
	ErrPlayerNotInRoom   = errors.New("player not in room")
	// Add other common room-related errors here if needed
)

// NewRoom creates a new Room instance with validation
func NewRoom(id RoomID, name string) (*Room, error) {
	if len(name) < 3 || len(name) > 50 {
		return nil, ErrInvalidRoomName
	}

	return &Room{
		ID:           id,
		Name:         name,
		CreatedAt:    time.Now(),
		Players:      make([]*sharedEntity.User, 0), // Use imported User type
		Description:  make(map[string]string),
		ScenarioName: "",
	}, nil
}

// AddPlayer adds a player to the room
func (r *Room) AddPlayer(player *sharedEntity.User) { // Use imported User type
	r.Players = append(r.Players, player)
}

// RemovePlayer removes a player from the room
func (r *Room) RemovePlayer(playerID sharedEntity.UserID) { // Use imported UserID type
	for i, p := range r.Players {
		if p.ID == playerID {
			r.Players = append(r.Players[:i], r.Players[i+1:]...)
			return
		}
	}
}

func (r *Room) SetDescription(descriptionName string, text string) {
	r.Description[descriptionName] = text
}
