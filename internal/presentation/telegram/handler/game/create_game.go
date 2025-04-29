package telegram

import (
	"context"
	"fmt"
	"log"

	// Import room entity
	roomQuery "telemafia/internal/domain/room/usecase/query"
	messages "telemafia/internal/presentation/telegram/messages"
	"telemafia/internal/shared/tgutil"

	"gopkg.in/telebot.v3"
)

// HandleCreateGame initiates the interactive game creation process.
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

	if !requester.Admin {
		return c.Send(msgs.Common.ErrorPermissionDenied) // Direct access
	}

	// Fetch available rooms
	rooms, err := getRoomsHandler.Handle(context.Background(), roomQuery.GetRoomsQuery{}) // Pass value type
	if err != nil {
		errMsg := fmt.Sprintf(msgs.Game.CreateGameErrorFetchRooms, err) // Use correct field
		log.Printf("Error fetching rooms for /create_game: %v", err)
		return c.Send(errMsg)
	}

	if len(rooms) == 0 {
		return c.Send(msgs.Room.ListNoRooms) // Correct path
	}

	// Create inline keyboard with room buttons
	markup := &telebot.ReplyMarkup{}
	var rows []telebot.Row
	for _, room := range rooms {
		// Get player count safely
		playerCount := 0
		if room.Players != nil {
			playerCount = len(room.Players)
		}
		/*btn := &telebot.InlineButton{
			Unique: tgutil.UniqueCreateGameSelectRoom,
			Text:   fmt.Sprintf("%s (%d players)", room.Name, playerCount),
			Data:   string(room.ID),
		}*/
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
