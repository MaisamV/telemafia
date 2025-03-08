package entity

import (
	userEntity "telemafia/internal/user/entity"
	"time"
)

// RoomCreatedEvent is emitted when a new room is created
type RoomCreatedEvent struct {
	RoomID    RoomID
	Name      string
	CreatedAt time.Time
}

func (e RoomCreatedEvent) EventName() string {
	return "room.created"
}

func (e RoomCreatedEvent) OccurredAt() time.Time {
	return e.CreatedAt
}

// PlayerJoinedEvent is emitted when a player joins a room
type PlayerJoinedEvent struct {
	RoomID   RoomID
	PlayerID userEntity.UserID
	JoinedAt time.Time
}

func (e PlayerJoinedEvent) EventName() string {
	return "room.player_joined"
}

func (e PlayerJoinedEvent) OccurredAt() time.Time {
	return e.JoinedAt
}

// PlayerLeftEvent is emitted when a player leaves a room
type PlayerLeftEvent struct {
	RoomID   RoomID
	PlayerID userEntity.UserID
	LeftAt   time.Time
}

func (e PlayerLeftEvent) EventName() string {
	return "room.player_left"
}

func (e PlayerLeftEvent) OccurredAt() time.Time {
	return e.LeftAt
}

// PlayerKickedEvent is emitted when a player is kicked from a room
type PlayerKickedEvent struct {
	RoomID   RoomID
	PlayerID userEntity.UserID
	KickedAt time.Time
}

func (e PlayerKickedEvent) EventName() string {
	return "room.player_kicked"
}

func (e PlayerKickedEvent) OccurredAt() time.Time {
	return e.KickedAt
}
