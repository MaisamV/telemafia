package telegram

import (
	"context"
	"fmt"
	"log"

	gameEntity "telemafia/internal/domain/game/entity"
	gameQuery "telemafia/internal/domain/game/usecase/query"
	messages "telemafia/internal/presentation/telegram/messages"

	"gopkg.in/telebot.v3"
)

// HandleConfirmAssignments sends a public confirmation after roles are assigned (placeholder/example)
func HandleConfirmAssignments(
	getGameByIDHandler *gameQuery.GetGameByIDHandler,
	c telebot.Context,
	data string,
	msgs *messages.Messages,
) error {
	gameID := data // Assuming data is the game ID

	// Optional: Fetch game details if needed for the message
	_, err := getGameByIDHandler.Handle(context.Background(), gameQuery.GetGameByIDQuery{ID: gameEntity.GameID(gameID)})
	if err != nil {
		log.Printf("Callback ConfirmAssignments: Error fetching game %s: %v", gameID, err)
		// Use msg
		_ = c.Respond(&telebot.CallbackResponse{Text: fmt.Sprintf(msgs.Common.CallbackErrorGeneric, err), ShowAlert: true})
		return c.Edit(msgs.Common.CallbackFailedEdit) // Use msg
	}

	// Acknowledge the callback
	_ = c.Respond()

	// Edit the original message or send a new one
	// Use msg
	return c.Edit(fmt.Sprintf(msgs.Game.AssignmentsConfirmedResponse, gameID))
}

// handleShowMyRoleCallback could potentially show a user their role again
// func handleShowMyRoleCallback(h *BotHandler, c telebot.Context, data string) error {
// 	// Requires fetching game, finding user's assignment, and sending PM
// 	// ... implementation needed ...
// }
