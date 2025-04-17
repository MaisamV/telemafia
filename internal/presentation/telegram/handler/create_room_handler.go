package telegram

import (
	"context"
	"fmt"
	"strings"

	// "telemafia/internal/room/usecase" // Old path
	"telemafia/internal/domain/room/usecase/command" // New path (assuming CreateRoomCommand is here)
	"time"

	"gopkg.in/telebot.v3"

	// roomEntity "telemafia/internal/room/entity" // Old path
	roomEntity "telemafia/internal/domain/room/entity" // New path
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
	cmd := command.CreateRoomCommand{ // Use imported 'command' package
		ID:        roomEntity.RoomID(fmt.Sprintf("room_%d", time.Now().UnixNano())),
		Name:      args,
		CreatorID: user.ID,
	}
	// room, err := h.createRoomHandler.Handle(context.Background(), cmd) // Assuming createRoomHandler is now a specific command handler
	// Need to potentially update how the handler is accessed if it moved, e.g.:
	// room, err := h.roomCommands.CreateRoom.Handle(context.Background(), cmd)
	// For now, just changing the import and command struct usage.
	// The user might need to adjust the handler invocation later depending on how DI is set up.
	room, err := h.createRoomHandler.Handle(context.Background(), cmd)
	if err != nil {
		// Check for specific known errors from entity package
		if err == roomEntity.ErrRoomAlreadyExists {
			return c.Send(fmt.Sprintf("Error: Room name '%s' already exists.", args))
		} else if err == roomEntity.ErrInvalidRoomName {
			return c.Send(fmt.Sprintf("Error: Invalid room name '%s'. Must be 3-50 characters.", args))
		}
		// Generic error
		return c.Send(fmt.Sprintf("Error creating room: %v", err))
	}

	return c.Send(fmt.Sprintf("Room '%s' created successfully! ID: %s", room.Name, room.ID))
}
