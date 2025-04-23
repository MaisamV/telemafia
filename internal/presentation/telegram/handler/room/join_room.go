package telegram

import (
	"context"
	"fmt"
	"strings"

	roomEntity "telemafia/internal/domain/room/entity"
	roomCommand "telemafia/internal/domain/room/usecase/command"

	"gopkg.in/telebot.v3"
)

// HandleJoinRoom handles the /join_room command (now a function)
func HandleJoinRoom(h *BotHandler, c telebot.Context) error {
	roomIDStr := strings.TrimSpace(c.Message().Payload)
	if roomIDStr == "" {
		return c.Send("Please provide a room ID: /join_room <room_id>")
	}

	roomID := roomEntity.RoomID(roomIDStr)
	user := ToUser(c.Sender())
	if user == nil {
		return c.Send("Could not identify user.")
	}

	cmd := roomCommand.JoinRoomCommand{
		Requester: *user,
		RoomID:    roomID,
	}

	if err := h.joinRoomHandler.Handle(context.Background(), cmd); err != nil {
		return c.Send(fmt.Sprintf("Error joining room '%s': %v", roomID, err))
	}

	markup := &telebot.ReplyMarkup{}
	btnLeave := markup.Data(fmt.Sprintf("Leave Room %s", roomID), UniqueLeaveRoomSelectRoom, string(roomID))
	markup.Inline(markup.Row(btnLeave))

	return c.Send(fmt.Sprintf("Successfully joined room %s", roomID), markup)
}
