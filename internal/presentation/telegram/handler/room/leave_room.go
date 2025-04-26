package telegram

import (
	"context"
	"fmt"
	"strings"

	roomEntity "telemafia/internal/domain/room/entity"
	roomCommand "telemafia/internal/domain/room/usecase/command"
	tgutil "telemafia/internal/shared/tgutil"

	"gopkg.in/telebot.v3"
)

// RefreshNotifier is defined in create_room.go (same package)

// HandleLeaveRoom handles the /leave_room command (now a function)
func HandleLeaveRoom(leaveRoomHandler *roomCommand.LeaveRoomHandler, refreshNotifier RefreshNotifier, c telebot.Context) error {
	roomIDStr := strings.TrimSpace(c.Message().Payload)
	if roomIDStr == "" {
		return c.Send("Please provide a room ID: /leave_room [room_id]")
	}

	user := tgutil.ToUser(c.Sender())
	if user == nil {
		return c.Send("Could not identify user.")
	}
	cmd := roomCommand.LeaveRoomCommand{
		RoomID:    roomEntity.RoomID(roomIDStr),
		Requester: *user,
	}
	if err := leaveRoomHandler.Handle(context.Background(), cmd); err != nil {
		return c.Send(fmt.Sprintf("Error leaving room '%s': %v", roomIDStr, err))
	}

	refreshNotifier.RaiseRefreshNeeded()
	return c.Send(fmt.Sprintf("Successfully left room %s!", roomIDStr))
}
