package telegram

import (
	"context"
	"fmt"
	"strings"

	scenarioEntity "telemafia/internal/domain/scenario/entity"
	scenarioCommand "telemafia/internal/domain/scenario/usecase/command"
	messages "telemafia/internal/presentation/telegram/messages"
	tgutil "telemafia/internal/shared/tgutil"

	"gopkg.in/telebot.v3"
)

// HandleAddRole handles adding a role via /add_role (now a function)
func HandleAddRole(
	manageRolesHandler *scenarioCommand.ManageRolesHandler,
	c telebot.Context,
	msgs *messages.Messages,
) error {
	parts := strings.Fields(c.Message().Payload)
	if len(parts) < 2 {
		return c.Send(msgs.Scenario.AddRolePrompt)
	}

	scenarioID := parts[0]
	roleName := parts[1]

	requester := tgutil.ToUser(c.Sender())
	if requester == nil || !requester.Admin {
		return c.Send(msgs.Common.ErrorPermissionDenied)
	}

	cmd := scenarioCommand.AddRoleCommand{
		Requester:  *requester,
		ScenarioID: scenarioID,
		Role:       scenarioEntity.Role{Name: roleName},
	}
	if err := manageRolesHandler.HandleAddRole(context.Background(), cmd); err != nil {
		return c.Send(fmt.Sprintf(msgs.Scenario.AddRoleError, roleName, scenarioID, err))
	}

	return c.Send(fmt.Sprintf(msgs.Scenario.AddRoleSuccess, roleName, scenarioID))
}

// HandleRemoveRole handles removing a role via /remove_role (now a function)
func HandleRemoveRole(
	manageRolesHandler *scenarioCommand.ManageRolesHandler,
	c telebot.Context,
	msgs *messages.Messages,
) error {
	parts := strings.Fields(c.Message().Payload)
	if len(parts) < 2 {
		return c.Send(msgs.Scenario.RemoveRolePrompt)
	}

	scenarioID := parts[0]
	roleName := parts[1]

	requester := tgutil.ToUser(c.Sender())
	if requester == nil || !requester.Admin {
		return c.Send(msgs.Common.ErrorPermissionDenied)
	}

	cmd := scenarioCommand.RemoveRoleCommand{
		Requester:  *requester,
		ScenarioID: scenarioID,
		RoleName:   roleName,
	}
	if err := manageRolesHandler.HandleRemoveRole(context.Background(), cmd); err != nil {
		return c.Send(fmt.Sprintf(msgs.Scenario.RemoveRoleError, roleName, scenarioID, err))
	}

	return c.Send(fmt.Sprintf(msgs.Scenario.RemoveRoleSuccess, roleName, scenarioID))
}
