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
	RoomID   entity.RoomID
	PlayerID userEntity.UserID
	Name     string // Added player name
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
	_, err := h.roomRepo.GetRoomByID(cmd.RoomID)
	if err != nil {
		return err
	}

	// Create user entity
	user := &userEntity.User{
		ID:       cmd.PlayerID,
		Username: cmd.Name,
	}

	// Add user to room
	if err := h.roomRepo.AddPlayerToRoom(cmd.RoomID, user); err != nil {
		return err
	}

	event := entity.PlayerJoinedEvent{
		RoomID:   cmd.RoomID,
		PlayerID: cmd.PlayerID,
		JoinedAt: time.Now(),
	}

	return h.eventPublisher.Publish(event)
}
