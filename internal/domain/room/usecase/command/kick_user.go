package command

import (
	"context"
	"errors"
	"fmt" // Import fmt for error formatting
	roomEntity "telemafia/internal/domain/room/entity"
	roomPort "telemafia/internal/domain/room/port"
	sharedEntity "telemafia/internal/shared/entity"
	sharedEvent "telemafia/internal/shared/event"
	"time"
)

// KickUserCommand represents the command to kick a user from a room
type KickUserCommand struct {
	Requester sharedEntity.User   // The user initiating the kick
	RoomID    roomEntity.RoomID   // Use imported RoomID type
	PlayerID  sharedEntity.UserID // Use imported UserID type
}

// KickUserHandler handles kicking a user from a room
type KickUserHandler struct {
	roomRepo       roomPort.RoomRepository // Use imported Repository interface
	eventPublisher sharedEvent.Publisher
}

// NewKickUserHandler creates a new KickUserHandler
func NewKickUserHandler(repo roomPort.RoomRepository, publisher sharedEvent.Publisher) *KickUserHandler {
	return &KickUserHandler{
		roomRepo:       repo,
		eventPublisher: publisher,
	}
}

// Handle processes the kick user command
func (h *KickUserHandler) Handle(ctx context.Context, cmd KickUserCommand) error {
	// Fetch the room first to check moderator status
	room, err := h.roomRepo.GetRoomByID(cmd.RoomID)
	if err != nil {
		// Handle specific error like RoomNotFound if needed, otherwise return generic error
		return fmt.Errorf("kick user: could not find room %s: %w", cmd.RoomID, err)
	}

	// --- Permission Check ---
	// Allow if requester is global admin OR the moderator of this specific room
	isRoomModerator := room.Moderator != nil && room.Moderator.ID == cmd.Requester.ID
	if !cmd.Requester.Admin && !isRoomModerator {
		return errors.New("kick user: permission denied (requires admin or room moderator)")
	}

	// Remove player from room
	if err := h.roomRepo.RemovePlayerFromRoom(cmd.RoomID, cmd.PlayerID); err != nil {
		return fmt.Errorf("kick user: failed to remove player: %w", err) // Propagates ErrPlayerNotInRoom etc.
	}

	// Publish domain event
	evt := sharedEvent.PlayerKickedEvent{ // Use imported event type
		RoomID:   cmd.RoomID,
		PlayerID: cmd.PlayerID,
		KickedAt: time.Now(),
	}

	if err := h.eventPublisher.Publish(evt); err != nil {
		// Log error but don't fail the operation
		// log.Printf("Failed to publish PlayerKickedEvent: %v", err)
	}

	return nil
}
