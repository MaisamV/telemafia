package usecase

import (
	"context"
	"telemafia/internal/entity"
)

// AddDescriptionCommand represents the command to add a description to a room
type AddDescriptionCommand struct {
	RoomID          entity.RoomID
	DescriptionName string
	Text            string
}

// AddDescriptionHandler handles adding description to a room
type AddDescriptionHandler struct {
	roomRepo Repository
}

// NewAddDescriptionHandler creates a new AddDescriptionHandler
func NewAddDescriptionHandler(repo Repository) *AddDescriptionHandler {
	return &AddDescriptionHandler{roomRepo: repo}
}

// Handle processes the add description command
func (h *AddDescriptionHandler) Handle(ctx context.Context, cmd AddDescriptionCommand) error {
	room, err := h.roomRepo.GetRoomByID(cmd.RoomID)
	if err != nil {
		return err
	}

	room.SetDescription(cmd.DescriptionName, cmd.Text)
	return nil
}
