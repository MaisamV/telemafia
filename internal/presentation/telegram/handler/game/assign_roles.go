package telegram

import (
	"context"
	"fmt"
	"log"
	"strings"

	gameEntity "telemafia/internal/domain/game/entity"
	gameCommand "telemafia/internal/domain/game/usecase/command"

	"gopkg.in/telebot.v3"
)

// HandleAssignRoles handles the /assign_roles command (now a function)
func HandleAssignRoles(h *BotHandler, c telebot.Context) error {
	gameIDStr := strings.TrimSpace(c.Message().Payload)
	if gameIDStr == "" {
		return c.Send("Usage: /assign_roles <game_id>")
	}

	requester := ToUser(c.Sender())
	if requester == nil {
		return c.Send("Could not identify user.")
	}
	gameID := gameEntity.GameID(gameIDStr)

	cmd := gameCommand.AssignRolesCommand{
		Requester: *requester,
		GameID:    gameID,
	}

	// Assignments are returned but ignored for now as private sending is commented out.
	_, err := h.assignRolesHandler.Handle(context.Background(), cmd)
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
