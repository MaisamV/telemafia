package telegram

import (
	"context"
	"telemafia/common"
	"telemafia/internal/room/entity"
	roomCommand "telemafia/internal/room/usecase/command"
	roomQuery "telemafia/internal/room/usecase/query"

	"gopkg.in/telebot.v3"
)

// HandleDeleteRoom handles the /delete_room command
func (h *BotHandler) HandleDeleteRoom(c telebot.Context) error {
	// Check if the user is an admin
	if !common.Contains(h.adminUsernames, c.Sender().Username) {
		return c.Send("You are not authorized to use this command.")
	}

	// Fetch all rooms and send them as inline keyboard buttons
	rooms, err := h.getRoomsHandler.Handle(context.Background(), roomQuery.GetRoomsQuery{})
	if err != nil {
		return c.Send("Failed to fetch rooms.")
	}

	var buttons [][]telebot.InlineButton
	for _, room := range rooms {
		button := telebot.InlineButton{
			Unique: UniqueDeleteRoomSelectRoom,
			Text:   room.Name,
			Data:   string(room.ID),
		}
		buttons = append(buttons, []telebot.InlineButton{button})
	}

	markup := &telebot.ReplyMarkup{InlineKeyboard: buttons}
	return c.Send("Select a room to delete:", markup)
}

// HandleDeleteRoomCallback handles the callback to delete a specific room
func (h *BotHandler) HandleDeleteRoomCallback(c telebot.Context, roomID string) error {
	cmd := roomCommand.DeleteRoomCommand{
		RoomID: entity.RoomID(roomID),
	}
	if err := h.deleteRoomHandler.Handle(context.Background(), cmd); err != nil {
		return c.Respond(&telebot.CallbackResponse{
			Text: "Failed to delete room.",
		})
	}

	return c.Respond(&telebot.CallbackResponse{
		Text: "Room successfully deleted!",
	})
}
