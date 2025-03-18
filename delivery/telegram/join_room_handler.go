package telegram

import (
	"context"
	"fmt"
	"strings"
	errorHandler "telemafia/common/error"
	"telemafia/delivery/util"
	"telemafia/internal/room/entity"
	roomCommand "telemafia/internal/room/usecase/command"
	roomQuery "telemafia/internal/room/usecase/query"

	"gopkg.in/telebot.v3"
)

// HandleJoinRoom handles the /join_room command
func (h *BotHandler) HandleJoinRoom(c telebot.Context) error {
	args := strings.TrimSpace(c.Message().Payload)
	if args == "" {
		return c.Send("Please provide a room ID: /join_room [room_id]")
	}

	return h.SendJoinRoom(c, args)
}

// HandleJoinRoomCallback handles the join room callback
func (h *BotHandler) HandleJoinRoomCallback(c telebot.Context, roomID string) error {
	user := util.ToUser(c.Sender())

	// Join room
	cmd := roomCommand.JoinRoomCommand{
		RoomID: entity.RoomID(roomID),
		Player: *user,
	}
	if err := h.joinRoomHandler.Handle(context.Background(), cmd); err != nil {
		return c.Respond(&telebot.CallbackResponse{
			Text: "Error joining room",
		})
	}

	// Change refresh message type
	ChangeRefreshType(int64(user.ID), RoomDetail, roomID)

	messageText, markup, err2 := h.RoomDetailMessage(roomID)
	if err2 != nil {
		return c.Respond(&telebot.CallbackResponse{
			Text: fmt.Sprintf("Failed to fetch room info: %v", err2),
		})
	}

	// Edit the original message with the room and player list
	return h.UpdateMessage(c.Sender().ID, c.Message().ID, messageText, markup)
}

func (h *BotHandler) SendJoinRoom(c telebot.Context, roomID string) error {
	user := util.ToUser(c.Sender())
	// Join room
	cmd := roomCommand.JoinRoomCommand{
		RoomID: entity.RoomID(roomID),
		Player: *user,
	}
	if err := h.joinRoomHandler.Handle(context.Background(), cmd); err != nil {
		return c.Respond(&telebot.CallbackResponse{
			Text: errorHandler.HandleError(err, "Error joining room"),
		})
	}

	c.Respond(&telebot.CallbackResponse{
		Text: "Successfully joined the room!",
	})

	messageText, markup, err2 := h.RoomDetailMessage(roomID)
	if err2 != nil {
		fmt.Printf("Failed to fetch room info: %v", err2)
	}
	return h.SendMessage(c.Sender().ID, messageText, markup, RoomDetail, roomID)
}

func (h *BotHandler) RoomDetailMessage(roomID string) (string, *telebot.ReplyMarkup, error) {
	// Fetch players in the room
	players, err := h.getPlayersInRoomHandler.Handle(context.Background(), roomQuery.GetPlayersInRoomQuery{RoomID: entity.RoomID(roomID)})
	if err != nil {
		return "", nil, err
	}

	// Fetch rooms
	rooms, err := h.getRoomsHandler.Handle(context.Background(), roomQuery.GetRoomsQuery{})
	if err != nil {
		return "", nil, err
	}

	// Construct the message with room name and player list
	var room *entity.Room
	for _, r := range rooms {
		if r.ID == entity.RoomID(roomID) {
			room = r
			break
		}
	}
	if room == nil {
		return "", nil, fmt.Errorf("room not found")
	}

	playerNames := ""
	for _, player := range players {
		playerNames += fmt.Sprintf("- %s\n", player.FirstName)
	}

	// Format message with scenario information if available
	var messageText string
	if room.ScenarioName != "" {
		messageText = fmt.Sprintf("Room: %s\nScenario: %s\nPlayers in the room:\n%s",
			room.Name,
			room.ScenarioName,
			playerNames)
	} else {
		messageText = fmt.Sprintf("Room: %s\nPlayers in the room:\n%s",
			room.Name,
			playerNames)
	}

	// Create a "Leave this room" button
	leaveButton := telebot.InlineButton{
		Unique: UniqueLeaveRoomSelectRoom,
		Text:   "Leave this room",
		Data:   roomID,
	}
	markup := &telebot.ReplyMarkup{InlineKeyboard: [][]telebot.InlineButton{{leaveButton}}}
	return messageText, markup, nil
}
