package command

import (
	"context"
	"errors" // Added for permission error
	roomEntity "telemafia/internal/domain/room/entity"
	roomPort "telemafia/internal/domain/room/port"
	sharedEntity "telemafia/internal/shared/entity" // Added for User
)

// DeleteRoomCommand represents the command to delete a room
type DeleteRoomCommand struct {
	Requester sharedEntity.User // The user initiating the delete
	RoomID    roomEntity.RoomID // Use imported RoomID type
}

// DeleteRoomHandler handles room deletion
type DeleteRoomHandler struct {
	roomRepo roomPort.RoomWriter // Use imported RoomWriter interface
}

// NewDeleteRoomHandler creates a new DeleteRoomHandler
func NewDeleteRoomHandler(repo roomPort.RoomWriter) *DeleteRoomHandler {
	return &DeleteRoomHandler{
		roomRepo: repo,
	}
}

// Handle processes the delete room command
func (h *DeleteRoomHandler) Handle(ctx context.Context, cmd DeleteRoomCommand) error {
	// --- Permission Check ---
	if !cmd.Requester.Admin {
		return errors.New("delete room: admin privilege required")
	}
	return h.roomRepo.DeleteRoom(cmd.RoomID) // Propagates errors from repo
}
