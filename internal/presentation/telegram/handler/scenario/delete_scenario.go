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

// HandleDeleteScenario handles the /delete_scenario command (now a function)
func HandleDeleteScenario(
	deleteScenarioHandler *scenarioCommand.DeleteScenarioHandler,
	c telebot.Context,
	msgs *messages.Messages,
) error {
	scenarioID := strings.TrimSpace(c.Message().Payload)
	if scenarioID == "" {
		return c.Send(msgs.Scenario.DeletePrompt)
	}

	requester := tgutil.ToUser(c.Sender())
	if requester == nil || !requester.Admin {
		return c.Send(msgs.Common.ErrorPermissionDenied)
	}

	cmd := scenarioCommand.DeleteScenarioCommand{
		Requester: *requester,
		ID:        scenarioID,
	}
	if err := deleteScenarioHandler.Handle(context.Background(), cmd); err != nil {
		return c.Send(fmt.Sprintf(msgs.Scenario.DeleteError, scenarioID, err))
	}

	return c.Send(fmt.Sprintf(msgs.Scenario.DeleteSuccess, scenarioID))
}
