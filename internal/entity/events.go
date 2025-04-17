package entity

import (
	"time"
)

// RoomCreatedEvent is emitted when a new room is created
type RoomCreatedEvent struct {
	RoomID       RoomID
	Name         string
	CreatedAt    time.Time
	ScenarioName string
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
	PlayerID UserID
	RoomName string
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
	PlayerID UserID
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
	PlayerID UserID
	KickedAt time.Time
}

func (e PlayerKickedEvent) EventName() string {
	return "room.player_kicked"
}

func (e PlayerKickedEvent) OccurredAt() time.Time {
	return e.KickedAt
}

// RoomDetailMessage represents the details of a room
type RoomDetailMessage struct {
	RoomID       RoomID
	Name         string
	CreatedAt    time.Time
	PlayerCount  int
	ScenarioName string
	Players      []*User
}

func (m RoomDetailMessage) EventName() string {
	return "room.detail"
}

func (m RoomDetailMessage) OccurredAt() time.Time {
	return m.CreatedAt
}

// NewRoomDetailMessage creates a new RoomDetailMessage from a Room
func NewRoomDetailMessage(room *Room) RoomDetailMessage {
	return RoomDetailMessage{
		RoomID:       room.ID,
		Name:         room.Name,
		CreatedAt:    room.CreatedAt,
		PlayerCount:  len(room.Players),
		ScenarioName: room.ScenarioName,
		Players:      room.Players,
	}
}
