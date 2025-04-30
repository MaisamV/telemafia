package telegram

import (
	"context"
	"fmt"
	"io"
	"strings"

	scenarioCommand "telemafia/internal/domain/scenario/usecase/command"
	messages "telemafia/internal/presentation/telegram/messages"
	tgutil "telemafia/internal/shared/tgutil"

	"gopkg.in/telebot.v3"
)

// HandleDocument handles json files to add scenario.
func HandleDocument(
	addScenarioJSONHandler *scenarioCommand.AddScenarioJSONHandler,
	c telebot.Context,
	msgs *messages.Messages,
) error {
	if c.Message().Document != nil && c.Message().Document.MIME != "application/json" {
		return c.Send(msgs.Scenario.AddScenarioJSONPrompt)
	}

	requester := tgutil.ToUser(c.Sender())
	if requester == nil || !requester.Admin {
		return c.Send(msgs.Common.ErrorPermissionDenied)
	}

	fmt.Printf("Doc: %v", c.Message().Document)

	file, err2 := c.Bot().File(&c.Message().Document.File)

	if err2 != nil {
		return c.Send(msgs.Common.ErrorGeneric)
	}
	bytes, err := io.ReadAll(file)
	if err != nil {
		return c.Send(msgs.Common.ErrorGeneric)
	}
	jsonData := string(bytes)
	jsonData = strings.TrimSpace(jsonData)
	fmt.Printf("Json data: %s", jsonData)
	if jsonData == "" {
		return c.Send(msgs.Scenario.AddScenarioJSONPrompt) // Need new message key
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
