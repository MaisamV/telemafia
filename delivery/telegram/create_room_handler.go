package telegram

import (
	"context"
	"fmt"
	"strings"
	errorHandler "telemafia/common/error"
	"telemafia/delivery/common"
	"telemafia/internal/room/entity"
	roomCommand "telemafia/internal/room/usecase/command"
	"time"

	"gopkg.in/telebot.v3"
)

// HandleCreateRoom handles the /create_room command
func (h *BotHandler) HandleCreateRoom(c telebot.Context) error {
	args := strings.TrimSpace(c.Message().Payload)
	if args == "" {
		return c.Send("Please provide a room name: /create_room [name]")
	}

	user := common.ToUser(c.Sender())
	if !user.CanCreateRoom() {
		return c.Send("Only admins can create rooms")
	}

	// Create room
	cmd := roomCommand.CreateRoomCommand{
		ID:        entity.RoomID(fmt.Sprintf("room_%d", time.Now().UnixNano())),
		Name:      args,
		CreatorID: user.ID,
	}
	room, err := h.createRoomHandler.Handle(context.Background(), cmd)
	if err != nil {
		return c.Send(errorHandler.HandleError(err, fmt.Sprintf("Error creating room: %v", err)))
	}

	return c.Send(fmt.Sprintf("Room created successfully! ID: %s", room.ID))
}
