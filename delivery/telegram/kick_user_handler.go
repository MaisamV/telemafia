package telegram

import (
	"context"
	"strconv"
	"strings"
	"telemafia/common"
	"telemafia/internal/room/entity"
	roomCommand "telemafia/internal/room/usecase/command"
	roomQuery "telemafia/internal/room/usecase/query"
	userEntity "telemafia/internal/user/entity"

	"gopkg.in/telebot.v3"
)

// HandleKickUser handles the /kick_user command
func (h *BotHandler) HandleKickUser(c telebot.Context) error {
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
			Unique: UniqueKickSelectRoom,
			Text:   room.Name,
			Data:   string(room.ID),
		}
		buttons = append(buttons, []telebot.InlineButton{button})
	}

	markup := &telebot.ReplyMarkup{InlineKeyboard: buttons}
	return c.Send("Select a room:", markup)
}

// HandleKickUserFromRoomCallback handles the callback to kick a specific user from a room
func (h *BotHandler) HandleKickUserFromRoomCallback(c telebot.Context, data string) error {
	// Extract room ID and player ID from callback data
	parts := strings.Split(data, ":")
	if len(parts) != 2 {
		return c.Respond(&telebot.CallbackResponse{
			Text: "Invalid data format.",
		})
	}
	roomID := parts[0]
	playerID := parts[1]

	id, err := common.StringToInt64(playerID)
	if err != nil {
		return c.Respond(&telebot.CallbackResponse{
			Text: "Failed to extract user ID.",
		})
	}
	// Kick user from room
	cmd := roomCommand.KickUserCommand{
		RoomID:   entity.RoomID(roomID),
		PlayerID: userEntity.UserID(id),
	}
	if err := h.kickUserHandler.Handle(context.Background(), cmd); err != nil {
		return c.Respond(&telebot.CallbackResponse{
			Text: "Failed to kick user.",
		})
	}

	return c.Respond(&telebot.CallbackResponse{
		Text: "User successfully kicked from the room!",
	})
}

// HandleKickUserCallback handles the kick user callback
func (h *BotHandler) HandleKickUserCallback(c telebot.Context, data string) error {
	roomID := data
	// Fetch players in the room and send them as inline keyboard buttons
	players, err := h.getPlayersInRoomHandler.Handle(context.Background(), roomQuery.GetPlayersInRoomQuery{RoomID: entity.RoomID(roomID)})
	if err != nil {
		return c.Respond(&telebot.CallbackResponse{
			Text: "Failed to fetch players.",
		})
	}

	var buttons [][]telebot.InlineButton
	for _, player := range players {
		button := telebot.InlineButton{
			Unique: UniqueKickFromRoomSelectPlayer,
			Text:   player.FirstName + " " + player.LastName + " (" + player.Username + ")",
			Data:   roomID + ":" + strconv.FormatInt(int64(player.ID), 10),
		}
		buttons = append(buttons, []telebot.InlineButton{button})
	}

	markup := &telebot.ReplyMarkup{InlineKeyboard: buttons}
	return c.Send("Select a player to kick:", markup)
}
