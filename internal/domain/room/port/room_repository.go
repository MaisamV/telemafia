package port

import (
	roomEntity "telemafia/internal/domain/room/entity"
	sharedEntity "telemafia/internal/shared/entity"
)

// RoomReader defines the interface for reading room data
type RoomReader interface {
	// GetRoomByID Get room by ID
	GetRoomByID(id roomEntity.RoomID) (*roomEntity.Room, error)

	// GetRooms Get all rooms
	GetRooms() ([]*roomEntity.Room, error)

	// GetPlayerRooms Get rooms for a player
	GetPlayerRooms(playerID sharedEntity.UserID) ([]*roomEntity.Room, error)

	// GetPlayersInRoom Get players in a specific room
	GetPlayersInRoom(roomID roomEntity.RoomID) ([]*sharedEntity.User, error)
}

// RoomWriter defines the interface for writing room data
type RoomWriter interface {
	// CreateRoom Create a new room
	CreateRoom(room *roomEntity.Room) error

	// UpdateRoom updates a room
	UpdateRoom(room *roomEntity.Room) error

	// AddPlayerToRoom Add a player to a room
	AddPlayerToRoom(roomID roomEntity.RoomID, player *sharedEntity.User) error

	// RemovePlayerFromRoom Remove a player from a room
	RemovePlayerFromRoom(roomID roomEntity.RoomID, playerID sharedEntity.UserID) error

	// DeleteRoom deletes a room by ID
	DeleteRoom(roomID roomEntity.RoomID) error
}

// RoomRepository defines the combined interface for room persistence
type RoomRepository interface {
	RoomReader
	RoomWriter
}
