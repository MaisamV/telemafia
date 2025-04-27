package telegram

import (
	"context"
	"fmt"
	"strings"
	"time"

	scenarioCommand "telemafia/internal/domain/scenario/usecase/command"
	messages "telemafia/internal/presentation/telegram/messages"
	tgutil "telemafia/internal/shared/tgutil"

	"gopkg.in/telebot.v3"
)

// HandleCreateScenario handles the /create_scenario command (now a function)
func HandleCreateScenario(
	createScenarioHandler *scenarioCommand.CreateScenarioHandler,
	c telebot.Context,
	msgs *messages.Messages,
) error {
	args := strings.TrimSpace(c.Message().Payload)
	if args == "" {
		return c.Send(msgs.Scenario.CreatePrompt)
	}

	requester := tgutil.ToUser(c.Sender())
	if requester == nil || !requester.Admin {
		return c.Send(msgs.Common.ErrorPermissionDenied)
	}

	cmd := scenarioCommand.CreateScenarioCommand{
		Requester: *requester,
		ID:        fmt.Sprintf("scen_%d", time.Now().UnixNano()),
		Name:      args,
	}
	if err := createScenarioHandler.Handle(context.Background(), cmd); err != nil {
		return c.Send(fmt.Sprintf(msgs.Scenario.CreateError, err))
	}

	return c.Send(fmt.Sprintf(msgs.Scenario.CreateSuccess, cmd.Name, cmd.ID, cmd.ID))
}
