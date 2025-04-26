package telegram

import (
	"context"
	"fmt"
	"strings"

	scenarioCommand "telemafia/internal/domain/scenario/usecase/command"
	tgutil "telemafia/internal/shared/tgutil"

	"gopkg.in/telebot.v3"
)

// HandleDeleteScenario handles the /delete_scenario command (now a function)
func HandleDeleteScenario(deleteScenarioHandler *scenarioCommand.DeleteScenarioHandler, c telebot.Context) error {
	scenarioID := strings.TrimSpace(c.Message().Payload)
	if scenarioID == "" {
		return c.Send("Please provide a scenario ID: /delete_scenario <id>")
	}

	requester := tgutil.ToUser(c.Sender())

	cmd := scenarioCommand.DeleteScenarioCommand{
		Requester: *requester,
		ID:        scenarioID,
	}
	if err := deleteScenarioHandler.Handle(context.Background(), cmd); err != nil {
		return c.Send(fmt.Sprintf("Error deleting scenario '%s': %v", scenarioID, err))
	}

	return c.Send(fmt.Sprintf("Scenario %s deleted successfully!", scenarioID))
}
