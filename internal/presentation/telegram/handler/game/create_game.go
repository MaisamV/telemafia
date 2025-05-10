package telegram

import (
	"context"
	"fmt"
	"log"

	// Import room entity
	roomEntity "telemafia/internal/domain/room/entity"
	roomQuery "telemafia/internal/domain/room/usecase/query"
	messages "telemafia/internal/presentation/telegram/messages"
	"telemafia/internal/shared/tgutil"

	"gopkg.in/telebot.v4"
)

// HandleCreateGame initiates the interactive game creation process.
// Global admins see all rooms. Room moderators see only the rooms they moderate.
func HandleCreateGame(
	getRoomsHandler *roomQuery.GetRoomsHandler,
	// other handlers needed for callbacks will be passed to HandleCallback
	c telebot.Context,
	msgs *messages.Messages,
) error {
	requester := tgutil.ToUser(c.Sender()) // Correct way to get user
	if requester == nil {
		log.Println(msgs.Common.ErrorIdentifyRequester) // Direct access
		return c.Send(msgs.Common.ErrorIdentifyRequester)
	}

	// Fetch available rooms
	allRooms, err := getRoomsHandler.Handle(context.Background(), roomQuery.GetRoomsQuery{}) // Pass value type
	if err != nil {
		errMsg := fmt.Sprintf(msgs.Game.CreateGameErrorFetchRooms, err) // Use correct field
		log.Printf("Error fetching rooms for /create_game: %v", err)
		return c.Send(errMsg)
	}

	// Filter rooms based on permission (Global Admin or Room Moderator)
	isGlobalAdmin := requester.Admin
	var roomsToShow []*roomEntity.Room
	for _, room := range allRooms {
		isRoomModerator := room.Moderator != nil && room.Moderator.ID == requester.ID
		if isGlobalAdmin || isRoomModerator {
			roomsToShow = append(roomsToShow, room)
		}
	}

	if len(roomsToShow) == 0 {
		// User is neither admin nor moderator of any room
		return c.Send(msgs.Common.ErrorPermissionDenied) // Or a more specific message
	}

	// Create inline keyboard with allowed room buttons
	markup := &telebot.ReplyMarkup{}
	var rows []telebot.Row
	for _, room := range roomsToShow {
		// Get player count safely
		playerCount := 0
		if room.Players != nil {
			playerCount = len(room.Players)
		}
		btn := markup.Data(
			fmt.Sprintf("%s (%d players)", room.Name, playerCount), // Show player count
			tgutil.UniqueCreateGameSelectRoom,
			string(room.ID),
		)
		rows = append(rows, markup.Row(btn))
	}
	// Add a cancel button
	rows = append(rows, markup.Row(markup.Data(msgs.Game.CreateGameCancelButton, tgutil.UniqueCancel))) // Use correct field
	markup.Inline(rows...)

	// Send message asking to select a room
	promptMsg := msgs.Game.CreateGameSelectRoomPrompt // Use correct field
	return c.Send(promptMsg, markup)
}

// Note: The rest of the flow (scenario selection, confirmation, starting) will be handled
// via callbacks, likely in a shared callback handler file.
