package command

import (
	"context"
	roomEntity "telemafia/internal/domain/room/entity"
	roomPort "telemafia/internal/domain/room/port"
	sharedEntity "telemafia/internal/shared/entity"
	sharedEvent "telemafia/internal/shared/event"
	"time"
)

// CreateRoomCommand represents the command to create a new room
type CreateRoomCommand struct {
	ID        roomEntity.RoomID
	Name      string
	CreatorID sharedEntity.UserID
}

// CreateRoomHandler handles room creation
type CreateRoomHandler struct {
	roomRepo       roomPort.RoomWriter // Use imported port interface
	eventPublisher sharedEvent.Publisher
}

// NewCreateRoomHandler creates a new CreateRoomHandler
func NewCreateRoomHandler(repo roomPort.RoomWriter, publisher sharedEvent.Publisher) *CreateRoomHandler {
	return &CreateRoomHandler{
		roomRepo:       repo,
		eventPublisher: publisher,
	}
}

// Handle processes the create room command
func (h *CreateRoomHandler) Handle(ctx context.Context, cmd CreateRoomCommand) (*roomEntity.Room, error) {
	room, err := roomEntity.NewRoom(cmd.ID, cmd.Name)
	if err != nil {
		return nil, err
	}

	// Add Creator logic if needed - the entity constructor doesn't take creator anymore
	// creator := sharedEntity.User{ ID: cmd.CreatorID /* Fetch full user? */ }
	// room.AddPlayer(&creator)

	if err := h.roomRepo.CreateRoom(room); err != nil {
		return nil, err
	}

	// Publish domain event
	evt := sharedEvent.RoomCreatedEvent{ // Use imported event type
		RoomID:    room.ID,
		Name:      room.Name,
		CreatedAt: time.Now(),
		// ScenarioName: room.ScenarioName, // Add if needed
	}

	if err := h.eventPublisher.Publish(evt); err != nil {
		// Log error but don't fail the operation
		// Consider using a retry mechanism for event publishing
		// log.Printf("Failed to publish RoomCreatedEvent: %v", err)
	}

	return room, nil
}
