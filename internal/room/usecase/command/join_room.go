package command

import (
	"context"
	"telemafia/internal/room/entity"
	"telemafia/internal/room/repo"
	userEntity "telemafia/internal/user/entity"
	"telemafia/pkg/event"
	"time"
)

// JoinRoomCommand represents the command to join a room
type JoinRoomCommand struct {
	RoomID entity.RoomID
	Player userEntity.User
}

// JoinRoomHandler handles room joining
type JoinRoomHandler struct {
	roomRepo       repo.Repository
	eventPublisher event.Publisher
}

// NewJoinRoomHandler creates a new JoinRoomHandler
func NewJoinRoomHandler(repo repo.Repository, publisher event.Publisher) *JoinRoomHandler {
	return &JoinRoomHandler{
		roomRepo:       repo,
		eventPublisher: publisher,
	}
}

// Handle processes the join room command
func (h *JoinRoomHandler) Handle(ctx context.Context, cmd JoinRoomCommand) error {
	// Add user to room
	if err := h.roomRepo.AddPlayerToRoom(cmd.RoomID, &cmd.Player); err != nil {
		return err
	}

	// Publish event with room details
	room, err := h.roomRepo.GetRoomByID(cmd.RoomID)
	if err != nil {
		return err
	}

	event := entity.PlayerJoinedEvent{
		RoomID:   cmd.RoomID,
		PlayerID: cmd.Player.ID,
		RoomName: room.Name, // Assuming RoomName is part of the event
		JoinedAt: time.Now(),
	}

	return h.eventPublisher.Publish(event)
}
