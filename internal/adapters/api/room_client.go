package api

import (
	gamePort "telemafia/internal/domain/game/port"
	roomEntity "telemafia/internal/domain/room/entity"
	roomPort "telemafia/internal/domain/room/port"
)

// Ensure LocalRoomClient implements the gamePort.RoomClient interface.
var _ gamePort.RoomClient = (*LocalRoomClient)(nil)

// LocalRoomClient implements the RoomClient interface by directly calling
// the Room domain's repository reader within the monolith.
// In a microservice architecture, this would be replaced by an HTTP/gRPC client.
type LocalRoomClient struct {
	roomRepo roomPort.RoomReader
}

// NewLocalRoomClient creates a new LocalRoomClient.
func NewLocalRoomClient(roomRepo roomPort.RoomReader) *LocalRoomClient {
	return &LocalRoomClient{roomRepo: roomRepo}
}

// FetchRoom retrieves a room using the injected RoomReader.
func (c *LocalRoomClient) FetchRoom(id roomEntity.RoomID) (*roomEntity.Room, error) {
	// In a real microservice, this would make an API call to the Room service.
	return c.roomRepo.GetRoomByID(id)
}
