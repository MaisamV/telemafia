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

// HandleJoinRoom handles the /join_room command (now a function)
func HandleJoinRoom(joinRoomHandler *roomCommand.JoinRoomHandler, refreshNotifier RefreshNotifier, c telebot.Context) error {
	roomIDStr := strings.TrimSpace(c.Message().Payload)
	if roomIDStr == "" {
		return c.Send("Please provide a room ID: /join_room <room_id>")
	}

	roomID := roomEntity.RoomID(roomIDStr)
	user := tgutil.ToUser(c.Sender())
	if user == nil {
		return c.Send("Could not identify user.")
	}

	cmd := roomCommand.JoinRoomCommand{
		Requester: *user,
		RoomID:    roomID,
	}

	err := joinRoomHandler.Handle(context.Background(), cmd)
	if err != nil {
		return c.Send(fmt.Sprintf("Error joining room '%s': %v", roomID, err))
	}

	refreshNotifier.RaiseRefreshNeeded()

	markup := &telebot.ReplyMarkup{}
	btnLeave := markup.Data(fmt.Sprintf("Leave Room %s", roomID), tgutil.UniqueLeaveRoomSelectRoom, string(roomID))
	markup.Inline(markup.Row(btnLeave))

	return c.Send(fmt.Sprintf("Successfully joined room %s", roomID), markup)
}
