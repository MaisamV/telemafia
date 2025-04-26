package telegram

import (
	"context"
	"fmt"
	"log"
	"strings"

	gameQuery "telemafia/internal/domain/game/usecase/query"
	tgutil "telemafia/internal/shared/tgutil"

	"gopkg.in/telebot.v3"
)

// HandleGamesList handles the /games command (now a function)
func HandleGamesList(getGamesHandler *gameQuery.GetGamesHandler, c telebot.Context) error {
	requester := tgutil.ToUser(c.Sender())
	if requester == nil {
		return c.Send("Could not identify requester.")
	}

	// Admin check remains here for now, as it's a query displaying potentially sensitive info
	if !tgutil.IsAdmin(c.Sender().Username) {
		return c.Send("You are not authorized to use this command.")
	}

	query := gameQuery.GetGamesQuery{}
	games, err := getGamesHandler.Handle(context.Background(), query)
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
		roomName := "<Unknown>"
		if game.Room != nil {
			roomName = game.Room.Name
		}
		scenarioName := "<Unknown>"
		if game.Scenario != nil {
			scenarioName = game.Scenario.Name
		}
		response.WriteString(fmt.Sprintf("- GameID: `%s`\n  Room: %s (`%s`)\n  Scenario: %s (`%s`)\n  State: `%s`\n  Players/Assignments: %d\n",
			game.ID, roomName, game.Room.ID, scenarioName, game.Scenario.ID, game.State, len(game.Assignments)))

		// // Optionally, fetch and display assignments (might be too verbose)
		// if len(game.Assignments) > 0 {
		// 	response.WriteString("  Assignments:\n")
		// 	for userID, role := range game.Assignments {
		// 		// Need a way to get username from userID - this requires another query/lookup
		// 		// For now, just showing ID
		// 		response.WriteString(fmt.Sprintf("    %d: %s\n", userID, role.Name))
		// 	}
		// }
		response.WriteString("\n") // Add a newline between games
	}

	return c.Send(response.String(), &telebot.SendOptions{ParseMode: telebot.ModeMarkdown})
}

// formatAssignments removed - moved to util.go as method on BotHandler
