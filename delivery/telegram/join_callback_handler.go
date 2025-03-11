package telegram

import (
	"context"
	"fmt"
	"telemafia/delivery/util"
	"telemafia/internal/room/entity"
	roomCommand "telemafia/internal/room/usecase/command"
	roomQuery "telemafia/internal/room/usecase/query"

	"gopkg.in/telebot.v3"
)

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

	// Fetch players in the room
	players, err := h.getPlayersInRoomHandler.Handle(context.Background(), roomQuery.GetPlayersInRoomQuery{RoomID: entity.RoomID(roomID)})
	if err != nil {
		return c.Respond(&telebot.CallbackResponse{
			Text: "Failed to fetch players in the room.",
		})
	}

	// Fetch rooms
	rooms, err := h.getRoomsHandler.Handle(context.Background(), roomQuery.GetRoomsQuery{})
	if err != nil {
		return c.Respond(&telebot.CallbackResponse{
			Text: "Failed to fetch rooms.",
		})
	}

	// Construct the message with room name and player list
	roomName := ""
	for _, room := range rooms {
		if room.ID == entity.RoomID(roomID) {
			roomName = room.Name
			break
		}
	}
	playerNames := ""
	for _, player := range players {
		playerNames += fmt.Sprintf("- %s\n", player.FirstName)
	}
	messageText := fmt.Sprintf("You have joined the room [%s].\nPlayers in the room:\n%s", roomName, playerNames)

	// Remove user from refresh list
	messageIDsMutex.Lock()
	delete(messageIDs, int64(user.ID))
	messageIDsMutex.Unlock()

	// Create a "Leave this room" button
	leaveButton := telebot.InlineButton{
		Unique: UniqueLeaveRoomSelectRoom,
		Text:   "Leave this room",
		Data:   string(roomID),
	}
	markup := &telebot.ReplyMarkup{InlineKeyboard: [][]telebot.InlineButton{{leaveButton}}}

	// Edit the original message with the room and player list
	return c.Edit(messageText, markup)
}
