package telegram

import (
	"context"
	"fmt"
	"gopkg.in/telebot.v3"
	"telemafia/internal/domain/room/entity"
	roomQuery "telemafia/internal/domain/room/usecase/query"
	"telemafia/internal/presentation/telegram/messages"
	"telemafia/internal/shared/tgutil"
)

func RoomDetailMessage(getRoomsHandler *roomQuery.GetRoomsHandler,
	getPlayersInRoomHandler *roomQuery.GetPlayersInRoomHandler,
	msgs *messages.Messages,
	roomID string) (string, *telebot.ReplyMarkup, error) {
	// Fetch players in the room
	players, err := getPlayersInRoomHandler.Handle(context.Background(), roomQuery.GetPlayersInRoomQuery{RoomID: entity.RoomID(roomID)})
	if err != nil {
		return "", nil, err
	}

	// Fetch rooms
	rooms, err := getRoomsHandler.Handle(context.Background(), roomQuery.GetRoomsQuery{})
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
		playerNames += fmt.Sprintf("@%s\n", player.Username)
	}

	// Format message with scenario information if available
	var messageText string
	if room.ScenarioName != "" {
		messageText = fmt.Sprintf(msgs.Room.RoomDetailWithScenario,
			room.Name,
			room.ScenarioName,
			playerNames)
	} else {
		messageText = fmt.Sprintf(msgs.Room.RoomDetail,
			room.Name,
			playerNames)
	}

	// Create a "Leave this room" button
	leaveButton := telebot.InlineButton{
		Unique: tgutil.UniqueLeaveRoomSelectRoom,
		Text:   msgs.Room.LeaveButton,
		Data:   roomID,
	}
	markup := &telebot.ReplyMarkup{InlineKeyboard: [][]telebot.InlineButton{{leaveButton}}}
	return messageText, markup, nil
}
