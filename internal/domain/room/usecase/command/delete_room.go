package command

import (
	"context"
	roomEntity "telemafia/internal/domain/room/entity"
	roomPort "telemafia/internal/domain/room/port"
)

// DeleteRoomCommand represents the command to delete a room
type DeleteRoomCommand struct {
	RoomID roomEntity.RoomID // Use imported RoomID type
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
	return h.roomRepo.DeleteRoom(cmd.RoomID) // Propagates errors from repo
}
