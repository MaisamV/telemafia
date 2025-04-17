package telegram

import (
	"context"
	"fmt"
	"strings"
	"telemafia/internal/entity"
	"telemafia/internal/usecase"
	"time"

	"gopkg.in/telebot.v3"
)

// HandleCreateRoom handles the /create_room command
func (h *BotHandler) HandleCreateRoom(c telebot.Context) error {
	args := strings.TrimSpace(c.Message().Payload)
	if args == "" {
		return c.Send("Please provide a room name: /create_room [name]")
	}

	user := ToUser(c.Sender())
	if !IsAdmin(user.Username) {
		return c.Send("Only admins can create rooms")
	}

	// Create room
	cmd := usecase.CreateRoomCommand{
		ID:        entity.RoomID(fmt.Sprintf("room_%d", time.Now().UnixNano())),
		Name:      args,
		CreatorID: user.ID,
	}
	room, err := h.createRoomHandler.Handle(context.Background(), cmd)
	if err != nil {
		// Check for specific known errors from entity package
		if err == entity.ErrRoomAlreadyExists {
			return c.Send(fmt.Sprintf("Error: Room name '%s' already exists.", args))
		} else if err == entity.ErrInvalidRoomName {
			return c.Send(fmt.Sprintf("Error: Invalid room name '%s'. Must be 3-50 characters.", args))
		}
		// Generic error
		return c.Send(fmt.Sprintf("Error creating room: %v", err))
	}

	return c.Send(fmt.Sprintf("Room '%s' created successfully! ID: %s", room.Name, room.ID))
}
