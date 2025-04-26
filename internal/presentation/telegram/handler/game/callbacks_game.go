package telegram

import (
	"context"
	"fmt"
	"log"

	gameEntity "telemafia/internal/domain/game/entity"
	gameQuery "telemafia/internal/domain/game/usecase/query"
	tgutil "telemafia/internal/shared/tgutil"

	"gopkg.in/telebot.v3"
)

// HandleConfirmAssignments is triggered when the admin confirms role assignments (maybe)
func HandleConfirmAssignments(getGameByIDHandler *gameQuery.GetGameByIDHandler, c telebot.Context, data string) error {
	gameID := gameEntity.GameID(data)
	// userID := c.Sender().ID // Not needed here

	// Respond immediately to the callback to stop the loading indicator
	_ = c.Respond(&telebot.CallbackResponse{Text: "Sending roles..."})

	// Fetch the game to get assignments
	query := gameQuery.GetGameByIDQuery{ID: gameID} // Use ID field
	game, err := getGameByIDHandler.Handle(context.Background(), query)
	if err != nil {
		log.Printf("[Callback %s] Error fetching game '%s': %v", tgutil.UniqueConfirmAssignments, gameID, err)
		return c.Respond(&telebot.CallbackResponse{Text: fmt.Sprintf("Game '%s' not found: %v", gameID, err), ShowAlert: true})
	}

	if game.Room == nil {
		log.Printf("Game '%s' has no associated room.", gameID)
		return c.Respond(&telebot.CallbackResponse{Text: "Game has no room.", ShowAlert: true})
	}
	roomID := game.Room.ID
	log.Printf("Found game '%s' for room '%s'", game.ID, roomID)

	assignments := game.Assignments
	if len(assignments) == 0 {
		return c.Respond(&telebot.CallbackResponse{Text: "No role assignments found for this game", ShowAlert: true})
	}

	log.Printf("Found %d role assignments for game '%s'", len(assignments), gameID)

	sentCount := 0
	sendErrors := 0
	for playerUserID, role := range assignments {
		// We don't have the BotHandler 'h' anymore, so we can't get the display name easily here.
		// We could potentially fetch all users in the room, but that's inefficient.
		// For now, log with UserID.
		log.Printf("Sending role %s to UserID: %d", role.Name, playerUserID)

		// Get the Telegram User ID from the Game's Room Players list (requires fetching room or having it passed)
		// This is getting complex. Alternative: The AssignRoles use case could return a map[UserID]struct{TelegramID int64; Role Role}
		// For now, let's assume we can get the Telegram ID somehow. (THIS NEEDS REVISITING)
		playerTelegramID := int64(playerUserID) // FIXME: This assumes UserID == TelegramID, which might not be true!

		roleMsg := fmt.Sprintf("🎭 *Your Role Assignment* 🎭\n\nYou have been assigned the role: *%s*\n\nKeep your role secret and follow the game master's instructions!", role.Name)

		// Send the message (error ignored for now)
		_, sendErr := c.Bot().Send(&telebot.User{ID: playerTelegramID}, roleMsg, &telebot.SendOptions{ParseMode: telebot.ModeMarkdown})
		if sendErr != nil {
			log.Printf("[Callback %s] Failed to send role '%s' to user %d for game '%s': %v", tgutil.UniqueConfirmAssignments, role.Name, playerUserID, gameID, sendErr)
			sendErrors++
		}
	}

	_ = c.Edit(fmt.Sprintf("Roles sent to %d players for game %s!", sentCount, gameID))

	return c.Respond(&telebot.CallbackResponse{Text: fmt.Sprintf("Roles sent to %d players!", sentCount)})
}

// handleShowMyRoleCallback could potentially show a user their role again
// func handleShowMyRoleCallback(h *BotHandler, c telebot.Context, data string) error {
// 	// Requires fetching game, finding user's assignment, and sending PM
// 	// ... implementation needed ...
// }
