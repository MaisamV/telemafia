package command

import (
	"context"
	"telemafia/internal/room/entity"
	"telemafia/internal/room/repo"
)

// DeleteRoomCommand represents the command to delete a room
type DeleteRoomCommand struct {
	RoomID entity.RoomID
}

// DeleteRoomHandler handles deleting a room
type DeleteRoomHandler struct {
	roomRepo repo.RoomWriter
}

// NewDeleteRoomHandler creates a new DeleteRoomHandler
func NewDeleteRoomHandler(repo repo.RoomWriter) *DeleteRoomHandler {
	return &DeleteRoomHandler{
		roomRepo: repo,
	}
}

// Handle processes the delete room command
func (h *DeleteRoomHandler) Handle(ctx context.Context, cmd DeleteRoomCommand) error {
	return h.roomRepo.DeleteRoom(cmd.RoomID)
}
