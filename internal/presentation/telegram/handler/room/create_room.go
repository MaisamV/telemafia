package telegram

import (
	"context"
	"fmt"
	"strings"
	"time"

	roomEntity "telemafia/internal/domain/room/entity"
	roomCommand "telemafia/internal/domain/room/usecase/command"

	"gopkg.in/telebot.v3"
)

// HandleCreateRoom handles the /create_room command (now a function)
func HandleCreateRoom(h *BotHandler, c telebot.Context) error {
	args := strings.TrimSpace(c.Message().Payload)
	if args == "" {
		return c.Send("Please provide a room name: /create_room [name]")
	}

	user := ToUser(c.Sender()) // Get user info
	if user == nil {
		return c.Send("Could not identify user.")
	}

	cmd := roomCommand.CreateRoomCommand{
		ID:        roomEntity.RoomID(fmt.Sprintf("room_%d", time.Now().UnixNano())), // Generate unique ID
		Name:      args,
		CreatorID: user.ID, // Pass CreatorID
	}

	createdRoom, err := h.createRoomHandler.Handle(context.Background(), cmd)
	if err != nil {
		return c.Send(fmt.Sprintf("Error creating room: %v", err))
	}

	return c.Send(fmt.Sprintf("Room '%s' created successfully! ID: %s", createdRoom.Name, createdRoom.ID))
}
