package telegram

import (
	"context"
	"fmt"
	"log"
	"strings"

	"telemafia/internal/game/entity"
	"telemafia/internal/game/usecase/query"

	"gopkg.in/telebot.v3"
)

// GamesListHandler handles the /games command to display all available games
type GamesListHandler struct {
	bot             *telebot.Bot
	getGamesHandler *query.GetGamesHandler
}

// NewGamesListHandler creates a new GamesListHandler
func NewGamesListHandler(
	bot *telebot.Bot,
	getGamesHandler *query.GetGamesHandler,
) *GamesListHandler {
	return &GamesListHandler{
		bot:             bot,
		getGamesHandler: getGamesHandler,
	}
}

// HandleGamesList handles the /games command
func (h *GamesListHandler) HandleGamesList(c telebot.Context) error {
	log.Printf("Handling /games command from user %s", c.Sender().Username)

	games, err := h.getGamesHandler.Handle(context.Background())
	if err != nil {
		log.Printf("Error getting games: %v", err)
		return c.Send("Error retrieving games. Please try again later.")
	}

	if len(games) == 0 {
		return c.Send("No games found. Create a game with /create_game [scenario] [room_id]")
	}

	var response strings.Builder
	response.WriteString("Available games:\n\n")

	for i, game := range games {
		gameStatus := "Not started"
		if game.Status == entity.GameStatusRolesAssigned ||
			game.Status == entity.GameStatusInProgress ||
			game.Status == entity.GameStatusFinished {
			gameStatus = string(game.Status)
		}

		response.WriteString(fmt.Sprintf("%d. Game ID: %s\n", i+1, game.ID))
		response.WriteString(fmt.Sprintf("   Room ID: %s\n", game.Room.ID))
		response.WriteString(fmt.Sprintf("   Scenario: %s\n", game.Scenario.Name))
		response.WriteString(fmt.Sprintf("   Status: %s\n", gameStatus))

		if len(game.Assignments) > 0 {
			response.WriteString("   Player Assignments:\n")
			for userID, role := range game.Assignments {
				response.WriteString(fmt.Sprintf("   - User %d: %s\n", userID, role.Name))
			}
		} else {
			response.WriteString("   No player assignments yet.\n")
		}

		response.WriteString("\n")
	}

	return c.Send(response.String())
}

// HandleCallback handles callback queries related to games list
func (h *GamesListHandler) HandleCallback(c telebot.Context) error {
	log.Printf("Handling callback query with data: %s", c.Callback().Data)
	return nil
}
