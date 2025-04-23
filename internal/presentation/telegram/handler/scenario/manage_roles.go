package telegram

import (
	"context"
	"fmt"
	"strings"

	scenarioEntity "telemafia/internal/domain/scenario/entity"
	scenarioCommand "telemafia/internal/domain/scenario/usecase/command"

	"gopkg.in/telebot.v3"
)

// HandleAddRole handles the /add_role command (now a function)
func HandleAddRole(h *BotHandler, c telebot.Context) error {
	args := strings.Fields(strings.TrimSpace(c.Message().Payload))
	if len(args) != 2 {
		return c.Send("Usage: /add_role <scenario_id> <role_name>")
	}

	requester := ToUser(c.Sender())
	if requester == nil {
		return c.Send("Could not identify user.")
	}

	cmd := scenarioCommand.AddRoleCommand{
		Requester:  *requester,
		ScenarioID: args[0],
		Role:       scenarioEntity.Role{Name: args[1]},
	}
	if err := h.manageRolesHandler.HandleAddRole(context.Background(), cmd); err != nil {
		return c.Send(fmt.Sprintf("Error adding role '%s' to scenario '%s': %v", args[1], args[0], err))
	}

	return c.Send(fmt.Sprintf("Role '%s' added to scenario %s successfully!", args[1], args[0]))
}

// HandleRemoveRole handles the /remove_role command (now a function)
func HandleRemoveRole(h *BotHandler, c telebot.Context) error {
	args := strings.Fields(strings.TrimSpace(c.Message().Payload))
	if len(args) != 2 {
		return c.Send("Usage: /remove_role <scenario_id> <role_name>")
	}

	requester := ToUser(c.Sender())
	if requester == nil {
		return c.Send("Could not identify user.")
	}

	cmd := scenarioCommand.RemoveRoleCommand{
		Requester:  *requester,
		ScenarioID: args[0],
		RoleName:   args[1],
	}
	if err := h.manageRolesHandler.HandleRemoveRole(context.Background(), cmd); err != nil {
		return c.Send(fmt.Sprintf("Error removing role '%s' from scenario '%s': %v", args[1], args[0], err))
	}

	return c.Send(fmt.Sprintf("Role '%s' removed from scenario %s successfully!", args[1], args[0]))
}
