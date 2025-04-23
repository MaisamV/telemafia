package telegram

import (
	"context"
	"fmt"
	"log"
	"strings"

	gameEntity "telemafia/internal/domain/game/entity"
	gameQuery "telemafia/internal/domain/game/usecase/query"

	"gopkg.in/telebot.v3"
)

// HandleGamesList handles the /games command (now a function)
func HandleGamesList(h *BotHandler, c telebot.Context) error {
	// Admin check remains here for now, as it's a query displaying potentially sensitive info
	if !IsAdmin(c.Sender().Username) {
		return c.Send("You are not authorized to use this command.")
	}

	games, err := h.getGamesHandler.Handle(context.Background(), gameQuery.GetGamesQuery{})
	if err != nil {
		log.Printf("Error fetching games list: %v", err)
		return c.Send(fmt.Sprintf("Error fetching games list: %v", err))
	}

	if len(games) == 0 {
		return c.Send("No active games found.")
	}

	var response strings.Builder
	response.WriteString("Active Games:\n")
	for _, game := range games {
		roomInfo := "(No Room Linked)"
		if game.Room != nil {
			roomInfo = fmt.Sprintf("(Room: %s)", game.Room.ID)
		}
		scenarioInfo := "(No Scenario Linked)"
		if game.Scenario != nil {
			scenarioName := game.Scenario.Name
			if scenarioName == "" {
				scenarioName = "(Name Unknown)"
			}
			scenarioInfo = fmt.Sprintf("(Scenario: %s, ID: %s)", scenarioName, game.Scenario.ID)
		}
		assignmentCount := len(game.Assignments)

		response.WriteString(fmt.Sprintf("- Game ID: %s %s %s State: %s Players/Assignments: %d\n",
			game.ID, roomInfo, scenarioInfo, game.State, assignmentCount))

		if assignmentCount > 0 {
			response.WriteString("  Assignments:\n")
			response.WriteString(h.formatAssignments(game.Assignments))
		} else if game.State == gameEntity.GameStateWaitingForPlayers {
			response.WriteString(fmt.Sprintf("  Roles not assigned. Use /assign_roles %s\n", game.ID))
		}
		response.WriteString("\n")
	}

	return c.Send(response.String())
}

// formatAssignments removed - moved to util.go as method on BotHandler
