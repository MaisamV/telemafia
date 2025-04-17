package command

import (
	"context"
	roomEntity "telemafia/internal/domain/room/entity"
	roomPort "telemafia/internal/domain/room/port"
	sharedEntity "telemafia/internal/shared/entity"
	sharedEvent "telemafia/internal/shared/event"
	"time"
)

// JoinRoomCommand represents the command to join a room
type JoinRoomCommand struct {
	Requester sharedEntity.User // Use imported User type
	RoomID    roomEntity.RoomID // Use imported RoomID type
}

// JoinRoomHandler handles room joining
type JoinRoomHandler struct {
	roomRepo       roomPort.RoomRepository // Use imported combined Repository interface
	eventPublisher sharedEvent.Publisher
}

// NewJoinRoomHandler creates a new JoinRoomHandler
func NewJoinRoomHandler(repo roomPort.RoomRepository, publisher sharedEvent.Publisher) *JoinRoomHandler {
	return &JoinRoomHandler{
		roomRepo:       repo,
		eventPublisher: publisher,
	}
}

// Handle processes the join room command
func (h *JoinRoomHandler) Handle(ctx context.Context, cmd JoinRoomCommand) error {
	// Get the room first to ensure it exists
	room, err := h.roomRepo.GetRoomByID(cmd.RoomID)
	if err != nil {
		return err // Propagates ErrRoomNotFound etc.
	}

	// Add player to room using the repository method
	if err := h.roomRepo.AddPlayerToRoom(cmd.RoomID, &cmd.Requester); err != nil {
		return err // Propagates potential errors from repo impl (e.g., already exists)
	}

	// Publish domain event
	evt := sharedEvent.PlayerJoinedEvent{ // Use imported event type
		RoomID:   cmd.RoomID,
		PlayerID: cmd.Requester.ID, // Use ID from the User struct
		RoomName: room.Name,        // Get room name from the fetched room
		JoinedAt: time.Now(),
	}

	if err := h.eventPublisher.Publish(evt); err != nil {
		// Log error but don't fail the operation
		// log.Printf("Failed to publish PlayerJoinedEvent: %v", err)
	}

	return nil
}
