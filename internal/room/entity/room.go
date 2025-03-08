package entity

import (
	error2 "telemafia/common/error"
	userEntity "telemafia/internal/user/entity"
	"time"
)

// RoomID represents a unique room identifier
type RoomID string

// Room represents a game room entity
type Room struct {
	ID        RoomID
	Name      string
	CreatedAt time.Time
	Players   []*userEntity.User
}

// NewRoom creates a new Room instance with validation
func NewRoom(id RoomID, name string) (*Room, error) {
	if len(name) < 3 || len(name) > 50 {
		return nil, error2.ErrInvalidRoomName
	}

	return &Room{
		ID:        id,
		Name:      name,
		CreatedAt: time.Now(),
		Players:   make([]*userEntity.User, 0),
	}, nil
}

// AddPlayer adds a player to the room
func (r *Room) AddPlayer(player *userEntity.User) {
	r.Players = append(r.Players, player)
}

// RemovePlayer removes a player from the room
func (r *Room) RemovePlayer(playerID userEntity.UserID) {
	for i, p := range r.Players {
		if p.ID == playerID {
			r.Players = append(r.Players[:i], r.Players[i+1:]...)
			return
		}
	}
}
