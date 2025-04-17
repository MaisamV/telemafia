package usecase

import (
	"context"
	"telemafia/internal/entity"
	"telemafia/pkg/event"
	"time"
)

// CreateRoomCommand represents the command to create a new room
type CreateRoomCommand struct {
	ID        entity.RoomID
	Name      string
	CreatorID entity.UserID
}

// CreateRoomHandler handles room creation
type CreateRoomHandler struct {
	roomRepo       RoomWriter
	eventPublisher event.Publisher
}

// NewCreateRoomHandler creates a new CreateRoomHandler
func NewCreateRoomHandler(repo RoomWriter, publisher event.Publisher) *CreateRoomHandler {
	return &CreateRoomHandler{
		roomRepo:       repo,
		eventPublisher: publisher,
	}
}

// Handle processes the create room command
func (h *CreateRoomHandler) Handle(ctx context.Context, cmd CreateRoomCommand) (*entity.Room, error) {
	room, err := entity.NewRoom(cmd.ID, cmd.Name)
	if err != nil {
		return nil, err
	}

	if err := h.roomRepo.CreateRoom(room); err != nil {
		return nil, err
	}

	// Publish domain event
	event := entity.RoomCreatedEvent{
		RoomID:    room.ID,
		Name:      room.Name,
		CreatedAt: time.Now(),
	}

	if err := h.eventPublisher.Publish(event); err != nil {
		// Log error but don't fail the operation
		// Consider using a retry mechanism for event publishing
	}

	return room, nil
}
