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
	Moderator    *sharedEntity.User // Added Moderator field
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
func NewRoom(id RoomID, name string, creator *sharedEntity.User) (*Room, error) { // Added creator parameter
	if len(name) < 3 || len(name) > 50 {
		return nil, ErrInvalidRoomName
	}
	if creator == nil {
		// Depending on requirements, maybe return an error or assign a default/nil moderator
		return nil, errors.New("room creator cannot be nil")
	}

	return &Room{
		ID:           id,
		Name:         name,
		CreatedAt:    time.Now(),
		Players:      make([]*sharedEntity.User, 0), // Use imported User type
		Description:  make(map[string]string),
		ScenarioName: "",
		Moderator:    creator, // Set the creator as the initial moderator
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

// SetModerator updates the room's moderator, removes the new moderator from the player list if they exist,
// and adds the previous moderator back to the player list.
func (r *Room) SetModerator(newModerator *sharedEntity.User) error {
	if newModerator == nil {
		return errors.New("new moderator cannot be nil")
	}

	previousModerator := r.Moderator // Store the previous moderator

	// Set the new moderator
	r.Moderator = newModerator

	// Check if the new moderator is currently a player and remove them if so
	pModFound := false
	found := false
	playerIndex := -1
	for i, p := range r.Players {
		if p.ID == newModerator.ID {
			found = true
			playerIndex = i
			break
		}
	}

	if found {
		// Remove the new moderator from the players list using the found index
		r.Players = append(r.Players[:playerIndex], r.Players[playerIndex+1:]...)
	}

	// Add the previous moderator back to the players list, if there was one
	if previousModerator != nil {
		for _, p := range r.Players {
			if p.ID == previousModerator.ID {
				pModFound = true
				break
			}
		}
	}

	if previousModerator != nil && !pModFound {
		r.AddPlayer(previousModerator)
	}

	return nil
}
