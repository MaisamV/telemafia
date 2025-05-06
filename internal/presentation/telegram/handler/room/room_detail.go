package telegram

import (
	"context"
	"fmt"
	roomEntity "telemafia/internal/domain/room/entity"
	roomQuery "telemafia/internal/domain/room/usecase/query"
	"telemafia/internal/presentation/telegram/messages"
	sharedEntity "telemafia/internal/shared/entity"
	"telemafia/internal/shared/tgutil"

	"gopkg.in/telebot.v3"
)

// RoomDetailMessage generates the content and markup for the room detail view.
// It shows admin buttons if the viewer (identified by chatID) is a global admin
// or the moderator of this specific room.
func RoomDetailMessage(
	getRoomsHandler *roomQuery.GetRoomsHandler, // Consider changing to GetRoomByID handler
	getPlayersHandler *roomQuery.GetPlayersInRoomHandler,
	msgs *messages.Messages,
	requesterID sharedEntity.UserID, // ID of the user viewing the message
	roomID string,
) (string, []interface{}, error) {
	// Fetch players in the room
	players, err := getPlayersHandler.Handle(context.Background(), roomQuery.GetPlayersInRoomQuery{RoomID: roomEntity.RoomID(roomID)})
	if err != nil {
		return "", nil, fmt.Errorf("error fetching players for room detail: %w", err)
	}

	// Fetch room details
	// TODO: Optimize by using a GetRoomByID handler if available/created
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

	// Determine if the viewer has admin privileges for this room
	isRoomAdmin := tgutil.IsAdmin(int64(requesterID)) || (room.Moderator != nil && room.Moderator.ID == requesterID)

	// Format message text
	var messageText string
	moderatorLink := "None"
	if room.Moderator != nil {
		moderatorLink = room.Moderator.GetProfileLink()
	}
	messageText = fmt.Sprintf(msgs.Room.RoomDetail, // Assuming RoomDetail takes Name, Moderator, Players
		room.Name,
		moderatorLink,
		playerNames)

	// Create buttons
	markup := &telebot.ReplyMarkup{}

	// Create individual buttons first
	leaveButton := markup.Data(msgs.Room.LeaveButton, tgutil.UniqueLeaveRoomSelectRoom, roomID)
	inviteButton := markup.Data(msgs.Room.InviteLinkButton, tgutil.UniqueGetInviteLink, roomID)

	// Arrange the first row
	firstRow := markup.Row(leaveButton, inviteButton)

	// Prepare admin rows (if viewer is room admin)
	adminRows := []telebot.Row{}
	if isRoomAdmin {
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

	opts := []interface{}{
		markup,
		telebot.ModeMarkdownV2,
		telebot.NoPreview,
	}
	return messageText, opts, nil
}
