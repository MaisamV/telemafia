package telegram

import (
	"context"
	"fmt"
	"log"

	gameEntity "telemafia/internal/domain/game/entity"
	gameQuery "telemafia/internal/domain/game/usecase/query"

	"gopkg.in/telebot.v3"
)

// HandleConfirmAssignments handles the callback to send roles privately (now a function)
func HandleConfirmAssignments(h *BotHandler, c telebot.Context, gameIDStr string) error {
	gameID := gameEntity.GameID(gameIDStr)
	if gameID == "" {
		return c.Respond(&telebot.CallbackResponse{Text: "Error: Missing game ID", ShowAlert: true})
	}
	log.Printf("Confirming assignments for game: '%s'", gameID)

	targetGame, err := h.getGameByIDHandler.Handle(context.Background(), gameQuery.GetGameByIDQuery{ID: gameID})
	if err != nil {
		log.Printf("Error fetching game with ID '%s': %v", gameID, err)
		return c.Respond(&telebot.CallbackResponse{Text: fmt.Sprintf("Game '%s' not found: %v", gameID, err), ShowAlert: true})
	}

	if targetGame.Room == nil {
		log.Printf("Game '%s' has no associated room.", gameID)
		return c.Respond(&telebot.CallbackResponse{Text: "Game has no room.", ShowAlert: true})
	}
	roomID := targetGame.Room.ID
	log.Printf("Found game '%s' for room '%s'", targetGame.ID, roomID)

	assignments := targetGame.Assignments
	if len(assignments) == 0 {
		return c.Respond(&telebot.CallbackResponse{Text: "No role assignments found for this game", ShowAlert: true})
	}

	log.Printf("Found %d role assignments for game '%s'", len(assignments), gameID)

	successCount := 0
	for userID, role := range assignments {
		userChat := &telebot.Chat{ID: int64(userID)}
		userName := h.getUserDisplayName(userID)

		log.Printf("Sending role %s to %s (ID: %d)", role.Name, userName, userID)
		message := fmt.Sprintf("🎭 *Your Role Assignment* 🎭\n\nYou have been assigned the role: *%s*\n\nKeep your role secret and follow the game master's instructions!", role.Name)
		_, err = h.bot.Send(userChat, message, &telebot.SendOptions{ParseMode: telebot.ModeMarkdown})

		if err != nil {
			log.Printf("Failed to send role to %s (ID: %d): %v", userName, userID, err)
		} else {
			log.Printf("Successfully sent role to %s (ID: %d)", userName, userID)
			successCount++
		}
	}

	_ = c.Edit(fmt.Sprintf("Roles sent privately to %d players for game %s.", successCount, gameID))

	return c.Respond(&telebot.CallbackResponse{Text: fmt.Sprintf("Roles sent to %d players!", successCount)})
}
