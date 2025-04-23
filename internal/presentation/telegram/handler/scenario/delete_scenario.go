package telegram

import (
	"context"
	"fmt"
	"strings"

	scenarioCommand "telemafia/internal/domain/scenario/usecase/command"

	"gopkg.in/telebot.v3"
)

// HandleDeleteScenario handles the /delete_scenario command
func HandleDeleteScenario(h *BotHandler, c telebot.Context) error {
	args := strings.TrimSpace(c.Message().Payload)
	if args == "" {
		return c.Send("Please provide a scenario ID: /delete_scenario <id>")
	}

	requester := ToUser(c.Sender())

	cmd := scenarioCommand.DeleteScenarioCommand{
		Requester: *requester,
		ID:        args,
	}
	if err := h.deleteScenarioHandler.Handle(context.Background(), cmd); err != nil {
		return c.Send(fmt.Sprintf("Error deleting scenario '%s': %v", args, err))
	}

	return c.Send(fmt.Sprintf("Scenario %s deleted successfully!", args))
}
