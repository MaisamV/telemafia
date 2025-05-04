package event

import (
	// roomEntity "telemafia/internal/room/entity"
	roomEntity "telemafia/internal/domain/room/entity"
	sharedEntity "telemafia/internal/shared/entity"
	"time"
)

// Base event interface (from pkg/event/event.go - needs consolidation)
// type Event interface {
// 	EventName() string
// 	OccurredAt() time.Time
// }

// RoomCreatedEvent is emitted when a new room is created
type RoomCreatedEvent struct {
	RoomID       roomEntity.RoomID // Updated type
	Name         string
	CreatorID    sharedEntity.UserID // Add CreatorID
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
	RoomID   roomEntity.RoomID   // Updated type
	PlayerID sharedEntity.UserID // Updated type
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
	RoomID   roomEntity.RoomID   // Updated type
	PlayerID sharedEntity.UserID // Updated type
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
	RoomID   roomEntity.RoomID   // Updated type
	PlayerID sharedEntity.UserID // Updated type
	KickedAt time.Time
}

func (e PlayerKickedEvent) EventName() string {
	return "room.player_kicked"
}

func (e PlayerKickedEvent) OccurredAt() time.Time {
	return e.KickedAt
}

// RoomDetailMessage represents the details of a room (Consider if this is an event or a DTO)
type RoomDetailMessage struct {
	RoomID       roomEntity.RoomID // Updated type
	Name         string
	CreatedAt    time.Time
	PlayerCount  int
	ScenarioName string
	Players      []*sharedEntity.User // Updated type
}

func (m RoomDetailMessage) EventName() string {
	return "room.detail"
}

func (m RoomDetailMessage) OccurredAt() time.Time {
	// Using CreatedAt, might need a different timestamp if it's a query result
	return m.CreatedAt
}

// NewRoomDetailMessage creates a new RoomDetailMessage from a Room
func NewRoomDetailMessage(room *roomEntity.Room) RoomDetailMessage { // Updated type
	return RoomDetailMessage{
		RoomID:       room.ID,
		Name:         room.Name,
		CreatedAt:    room.CreatedAt,
		PlayerCount:  len(room.Players),
		ScenarioName: room.ScenarioName,
		Players:      room.Players,
	}
}
