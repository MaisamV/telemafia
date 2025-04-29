package telegram

import (
	"context"
	"fmt"
	"strings"

	scenarioCommand "telemafia/internal/domain/scenario/usecase/command"
	messages "telemafia/internal/presentation/telegram/messages"
	tgutil "telemafia/internal/shared/tgutil"

	"gopkg.in/telebot.v3"
)

// HandleAddScenarioJSON handles the /add_scenario_json command.
func HandleAddScenarioJSON(
	addScenarioJSONHandler *scenarioCommand.AddScenarioJSONHandler,
	c telebot.Context,
	msgs *messages.Messages,
) error {
	jsonData := strings.TrimSpace(c.Message().Payload)
	if jsonData == "" {
		return c.Send(msgs.Scenario.AddScenarioJSONPrompt) // Need new message key
	}

	requester := tgutil.ToUser(c.Sender())
	if requester == nil || !requester.Admin {
		return c.Send(msgs.Common.ErrorPermissionDenied)
	}

	cmd := scenarioCommand.AddScenarioJSONCommand{
		Requester: *requester,
		JSONData:  jsonData,
	}

	createdScenario, err := addScenarioJSONHandler.Handle(context.Background(), cmd)
	if err != nil {
		// Provide specific feedback for JSON errors vs other errors
		if strings.Contains(err.Error(), "invalid JSON format") {
			return c.Send(fmt.Sprintf(msgs.Scenario.AddScenarioJSONInvalidJSON, err))
		} else if strings.Contains(err.Error(), "cannot be empty") {
			return c.Send(fmt.Sprintf(msgs.Scenario.AddScenarioJSONValidationError, err))
		} else {
			return c.Send(fmt.Sprintf(msgs.Scenario.AddScenarioJSONErrorGeneric, err))
		}
	}

	return c.Send(fmt.Sprintf(msgs.Scenario.AddScenarioJSONSuccess, createdScenario.Name, createdScenario.ID))
}
