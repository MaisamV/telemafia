package port

import (
	roomEntity "telemafia/internal/domain/room/entity"
)

// RoomClient defines an interface for the Game domain to fetch Room data.
// Implementations could be local (monolith) or remote (microservice).
type RoomClient interface {
	FetchRoom(id roomEntity.RoomID) (*roomEntity.Room, error)
}
