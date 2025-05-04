package telegram

import (
	"context"
	"fmt"
	roomEntity "telemafia/internal/domain/room/entity"
	roomQuery "telemafia/internal/domain/room/usecase/query"
	"telemafia/internal/presentation/telegram/messages"
	"telemafia/internal/shared/tgutil"

	"gopkg.in/telebot.v3"
)

func RoomDetailMessage(
	getRoomsHandler *roomQuery.GetRoomsHandler,
	getPlayersHandler *roomQuery.GetPlayersInRoomHandler,
	msgs *messages.Messages,
	isAdmin bool,
	roomID string,
) (string, *telebot.ReplyMarkup, error) {
	// Fetch players in the room
	players, err := getPlayersHandler.Handle(context.Background(), roomQuery.GetPlayersInRoomQuery{RoomID: roomEntity.RoomID(roomID)})
	if err != nil {
		return "", nil, fmt.Errorf("error fetching players for room detail: %w", err)
	}

	// Fetch room details (using GetRoomByID would be more efficient if available, but using GetRooms for now)
	rooms, err := getRoomsHandler.Handle(context.Background(), roomQuery.GetRoomsQuery{})
	if err != nil {
		return "", nil, fmt.Errorf("error fetching rooms for room detail: %w", err)
	}

	// Find the specific room
	var room *roomEntity.Room
	for _, r := range rooms {
		if r.ID == roomEntity.RoomID(roomID) {
			room = r
			break
		}
	}
	if room == nil {
		return fmt.Sprintf(msgs.Room.RoomNotFound, roomID), nil, nil // Return user-friendly error message
	}

	// Construct player list string
	playerNames := ""

	for i, player := range players {
		playerNames += fmt.Sprintf("%d \\- %s\n", i+1, player.GetProfileLink())
	}

	// Format message text
	var messageText string
	//if room.ScenarioName != "" {
	//	messageText = fmt.Sprintf(msgs.Room.RoomDetailWithScenario,
	//		room.Name,
	//		room.ScenarioName,
	//		playerNames)
	//} else {
	messageText = fmt.Sprintf(msgs.Room.RoomDetail,
		room.Name,
		room.Moderator.GetProfileLink(),
		playerNames)
	//}

	// Create buttons
	markup := &telebot.ReplyMarkup{}

	// Create individual buttons first
	leaveButton := markup.Data(msgs.Room.LeaveButton, tgutil.UniqueLeaveRoomSelectRoom, roomID)
	inviteButton := markup.Data(msgs.Room.InviteLinkButton, tgutil.UniqueGetInviteLink, roomID)

	// Arrange the first row
	firstRow := markup.Row(leaveButton, inviteButton)

	// Prepare admin rows
	adminRows := []telebot.Row{}
	if isAdmin {
		// Admin Action Row (Kick, Change Moderator)
		kickButton := markup.Data(msgs.Room.KickUserButton, tgutil.UniqueKickUserSelect, roomID)
		modButton := markup.Data(msgs.Room.ChangeModeratorButton, tgutil.UniqueChangeModeratorSelect, roomID)
		actionRow := markup.Row(kickButton, modButton) // Add buttons to the same row
		adminRows = append(adminRows, actionRow)

		// Start Game Button Row (Separate Row)
		startButton := markup.Data(msgs.Game.StartButton, tgutil.UniqueCreateGameSelectRoom, roomID)
		startRow := markup.Row(startButton)
		adminRows = append(adminRows, startRow)
	}

	// Combine rows and add to markup
	allRows := []telebot.Row{firstRow}
	allRows = append(allRows, adminRows...)
	markup.Inline(allRows...)

	// No need to manually set markup.InlineKeyboard, markup.Inline handles it
	return messageText, markup, nil
}
