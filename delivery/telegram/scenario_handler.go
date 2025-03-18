package telegram

import (
	"context"
	"fmt"
	"strings"
	"telemafia/internal/scenario/entity"
	scenarioCommand "telemafia/internal/scenario/usecase/command"
	"time"

	"gopkg.in/telebot.v3"
)

// HandleCreateScenario handles the /create_scenario command
func (h *BotHandler) HandleCreateScenario(c telebot.Context) error {
	args := strings.TrimSpace(c.Message().Payload)
	if args == "" {
		return c.Send("Please provide a scenario name: /create_scenario [name]")
	}

	cmd := scenarioCommand.CreateScenarioCommand{
		ID:   fmt.Sprintf("scenario_%d", time.Now().UnixNano()),
		Name: args,
	}
	if err := h.createScenarioHandler.Handle(context.Background(), cmd); err != nil {
		return c.Send(fmt.Sprintf("Error creating scenario: %v", err))
	}

	return c.Send(fmt.Sprintf("/add_role %s", cmd.ID))
}

// HandleDeleteScenario handles the /delete_scenario command
func (h *BotHandler) HandleDeleteScenario(c telebot.Context) error {
	args := strings.TrimSpace(c.Message().Payload)
	if args == "" {
		return c.Send("Please provide a scenario ID: /delete_scenario [id]")
	}

	cmd := scenarioCommand.DeleteScenarioCommand{
		ID: args,
	}
	if err := h.deleteScenarioHandler.Handle(context.Background(), cmd); err != nil {
		return c.Send(fmt.Sprintf("Error deleting scenario: %v", err))
	}

	return c.Send("Scenario deleted successfully!")
}

// HandleAddRole handles the /add_role command
func (h *BotHandler) HandleAddRole(c telebot.Context) error {
	args := strings.Split(strings.TrimSpace(c.Message().Payload), " ")
	if len(args) != 2 {
		return c.Send("Please provide a scenario ID and role name: /add_role [scenario_id] [role_name]")
	}

	cmd := scenarioCommand.AddRoleCommand{
		ScenarioID: args[0],
		Role:       entity.Role{Name: args[1]},
	}
	if err := h.manageRolesHandler.HandleAddRole(context.Background(), cmd); err != nil {
		return c.Send(fmt.Sprintf("Error adding role: %v", err))
	}

	return c.Send("Role added successfully!")
}

// HandleRemoveRole handles the /remove_role command
func (h *BotHandler) HandleRemoveRole(c telebot.Context) error {
	args := strings.Split(strings.TrimSpace(c.Message().Payload), " ")
	if len(args) != 2 {
		return c.Send("Please provide a scenario ID and role name: /remove_role [scenario_id] [role_name]")
	}

	cmd := scenarioCommand.RemoveRoleCommand{
		ScenarioID: args[0],
		RoleName:   args[1],
	}
	if err := h.manageRolesHandler.HandleRemoveRole(context.Background(), cmd); err != nil {
		return c.Send(fmt.Sprintf("Error removing role: %v", err))
	}

	return c.Send("Role removed successfully!")
}
