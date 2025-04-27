package telegram

import (
	"context"
	"fmt"
	"strings"

	gameQuery "telemafia/internal/domain/game/usecase/query"
	messages "telemafia/internal/presentation/telegram/messages"
	tgutil "telemafia/internal/shared/tgutil"

	"gopkg.in/telebot.v3"
)

// HandleGamesList handles the /games command
func HandleGamesList(
	getGamesHandler *gameQuery.GetGamesHandler,
	c telebot.Context,
	msgs *messages.Messages,
) error {
	requester := tgutil.ToUser(c.Sender())
	if requester == nil || !requester.Admin {
		return c.Send(msgs.Common.ErrorPermissionDenied)
	}

	query := gameQuery.GetGamesQuery{}
	games, err := getGamesHandler.Handle(context.Background(), query)
	if err != nil {
		return c.Send(fmt.Sprintf(msgs.Game.ListGamesError, err))
	}

	if len(games) == 0 {
		return c.Send(msgs.Game.ListGamesNoGames)
	}

	var response strings.Builder
	response.WriteString(msgs.Game.ListGamesTitle)
	for _, game := range games {
		roomName := "<unknown>"
		roomID := "<unknown>"
		if game.Room != nil {
			roomName = game.Room.Name
			roomID = string(game.Room.ID)
		}
		scenarioName := "<unknown>"
		scenarioID := "<unknown>"
		if game.Scenario != nil {
			scenarioName = game.Scenario.Name
			scenarioID = game.Scenario.ID
		}
		playerCount := len(game.Assignments)
		if playerCount == 0 && game.Room != nil {
			playerCount = len(game.Room.Players)
		}
		response.WriteString(fmt.Sprintf(msgs.Game.ListGamesEntry,
			game.ID, roomName, roomID, scenarioName, scenarioID, game.State, playerCount,
		))
	}

	return c.Send(response.String())
}

// formatAssignments removed - moved to util.go as method on BotHandler
