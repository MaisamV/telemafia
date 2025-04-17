package usecase

import (
	"context"
	"telemafia/internal/entity"
)

// DeleteRoomCommand represents the command to delete a room
type DeleteRoomCommand struct {
	RoomID entity.RoomID
}

// DeleteRoomHandler handles room deletion
type DeleteRoomHandler struct {
	roomRepo RoomWriter
}

// NewDeleteRoomHandler creates a new DeleteRoomHandler
func NewDeleteRoomHandler(repo RoomWriter) *DeleteRoomHandler {
	return &DeleteRoomHandler{
		roomRepo: repo,
	}
}

// Handle processes the delete room command
func (h *DeleteRoomHandler) Handle(ctx context.Context, cmd DeleteRoomCommand) error {
	return h.roomRepo.DeleteRoom(cmd.RoomID)
}
