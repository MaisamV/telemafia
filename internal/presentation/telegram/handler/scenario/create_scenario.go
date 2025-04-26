package telegram

import (
	"context"
	"fmt"
	"strings"
	"time"

	scenarioCommand "telemafia/internal/domain/scenario/usecase/command"
	tgutil "telemafia/internal/shared/tgutil"

	"gopkg.in/telebot.v3"
)

// HandleCreateScenario handles the /create_scenario command (now a function)
func HandleCreateScenario(createScenarioHandler *scenarioCommand.CreateScenarioHandler, c telebot.Context) error {
	args := strings.TrimSpace(c.Message().Payload)
	if args == "" {
		return c.Send("Please provide a scenario name: /create_scenario [name]")
	}

	requester := tgutil.ToUser(c.Sender())

	cmd := scenarioCommand.CreateScenarioCommand{
		Requester: *requester,
		ID:        fmt.Sprintf("scen_%d", time.Now().UnixNano()),
		Name:      args,
	}
	if err := createScenarioHandler.Handle(context.Background(), cmd); err != nil {
		return c.Send(fmt.Sprintf("Error creating scenario: %v", err))
	}

	return c.Send(fmt.Sprintf("Scenario '%s' created successfully! ID: %s\nUse /add_role %s <role_name> to add roles.", cmd.Name, cmd.ID, cmd.ID))
}
