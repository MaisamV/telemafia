package telegram

import (
	"context"
	"fmt"
	"log"
	"strings"

	gameEntity "telemafia/internal/domain/game/entity"
	gameCommand "telemafia/internal/domain/game/usecase/command"
	tgutil "telemafia/internal/shared/tgutil"

	"gopkg.in/telebot.v3"
)

// HandleAssignRoles handles the /assign_roles command (now a function)
func HandleAssignRoles(assignRolesHandler *gameCommand.AssignRolesHandler, c telebot.Context) error {
	gameIDStr := strings.TrimSpace(c.Message().Payload)
	if gameIDStr == "" {
		return c.Send("Please provide a game ID: /assign_roles <game_id>")
	}

	requester := tgutil.ToUser(c.Sender())
	if requester == nil {
		return c.Send("Could not identify requester.")
	}
	gameID := gameEntity.GameID(gameIDStr)

	cmd := gameCommand.AssignRolesCommand{
		Requester: *requester,
		GameID:    gameID,
	}

	// Assignments are returned but private sending is handled via callback.
	_, err := assignRolesHandler.Handle(context.Background(), cmd)
	if err != nil {
		log.Printf("Error assigning roles for game '%s': %v", gameID, err)
		return c.Send(fmt.Sprintf("Error assigning roles for game '%s': %v", gameID, err))
	}

	// Send individual messages (potentially sensitive)
	// TODO: Implement private role sending functionality
	// go h.sendAssignmentMessages(assignments)

	// Send confirmation to admin/channel
	// TODO: Re-add "Show My Role" button when callback is implemented
	// markup := &telebot.ReplyMarkup{}
	// btnConfirm := markup.Data("Show My Role", UniqueShowMyRole, gameIDStr) // Button for players
	// markup.Inline(markup.Row(btnConfirm))

	// return c.Send(fmt.Sprintf("Roles assigned for game %s! Players can use the button below to see their role.", gameID), markup)
	return c.Send(fmt.Sprintf("Roles assigned successfully for game %s!", gameID))
}
