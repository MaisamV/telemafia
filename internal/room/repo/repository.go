package repo

import (
	"telemafia/internal/room/entity"
	userEntity "telemafia/internal/user/entity"
)

// RoomReader defines the interface for reading room data
type RoomReader interface {
	// GetRoomByID Get room by ID
	GetRoomByID(id entity.RoomID) (*entity.Room, error)

	// GetRooms Get all rooms
	GetRooms() ([]*entity.Room, error)

	// GetPlayerRooms Get rooms for a player
	GetPlayerRooms(playerID userEntity.UserID) ([]*entity.Room, error)

	// GetPlayersInRoom Get players in a specific room
	GetPlayersInRoom(roomID entity.RoomID) ([]*userEntity.User, error)
}

// RoomWriter defines the interface for writing room data
type RoomWriter interface {
	// CreateRoom Create a new room
	CreateRoom(room *entity.Room) error

	// AddPlayerToRoom Add a player to a room
	AddPlayerToRoom(roomID entity.RoomID, player *userEntity.User) error

	// RemovePlayerFromRoom Remove a player from a room
	RemovePlayerFromRoom(roomID entity.RoomID, playerID userEntity.UserID) error

	// DeleteRoom deletes a room by ID
	DeleteRoom(roomID entity.RoomID) error
}

// Repository defines the interface for room persistence
type Repository interface {
	RoomReader
	RoomWriter
}
