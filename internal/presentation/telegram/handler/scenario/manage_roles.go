package telegram

import (
	"context"
	"fmt"
	"strings"

	scenarioEntity "telemafia/internal/domain/scenario/entity"
	scenarioCommand "telemafia/internal/domain/scenario/usecase/command"
	tgutil "telemafia/internal/shared/tgutil"

	"gopkg.in/telebot.v3"
)

// HandleAddRole handles adding a role via /add_role (now a function)
func HandleAddRole(manageRolesHandler *scenarioCommand.ManageRolesHandler, c telebot.Context) error {
	parts := strings.Fields(c.Message().Payload)
	if len(parts) < 2 {
		return c.Send("Usage: /add_role <scenario_id> <role_name>")
	}

	scenarioID := parts[0]
	roleName := parts[1]

	requester := tgutil.ToUser(c.Sender())
	if requester == nil {
		return c.Send("Could not identify requester.")
	}

	cmd := scenarioCommand.AddRoleCommand{
		Requester:  *requester,
		ScenarioID: scenarioID,
		Role:       scenarioEntity.Role{Name: roleName},
	}
	if err := manageRolesHandler.HandleAddRole(context.Background(), cmd); err != nil {
		return c.Send(fmt.Sprintf("Error adding role '%s' to scenario '%s': %v", roleName, scenarioID, err))
	}

	return c.Send(fmt.Sprintf("Role '%s' added to scenario %s successfully!", roleName, scenarioID))
}

// HandleRemoveRole handles removing a role via /remove_role (now a function)
func HandleRemoveRole(manageRolesHandler *scenarioCommand.ManageRolesHandler, c telebot.Context) error {
	parts := strings.Fields(c.Message().Payload)
	if len(parts) != 2 {
		return c.Send("Usage: /remove_role <scenario_id> <role_name>")
	}

	scenarioID := parts[0]
	roleName := parts[1]

	requester := tgutil.ToUser(c.Sender())
	if requester == nil {
		return c.Send("Could not identify requester.")
	}

	cmd := scenarioCommand.RemoveRoleCommand{
		Requester:  *requester,
		ScenarioID: scenarioID,
		RoleName:   roleName,
	}
	if err := manageRolesHandler.HandleRemoveRole(context.Background(), cmd); err != nil {
		return c.Send(fmt.Sprintf("Error removing role '%s' from scenario '%s': %v", roleName, scenarioID, err))
	}

	return c.Send(fmt.Sprintf("Role '%s' removed from scenario %s successfully!", roleName, scenarioID))
}
