package telegram

import (
	"context"
	"fmt"
	"log"
	"strings"

	gameCommand "telemafia/internal/domain/game/usecase/command"
	roomEntity "telemafia/internal/domain/room/entity"
	roomCommand "telemafia/internal/domain/room/usecase/command"
	roomQuery "telemafia/internal/domain/room/usecase/query"
	scenarioQuery "telemafia/internal/domain/scenario/usecase/query"
	tgutil "telemafia/internal/shared/tgutil"

	"gopkg.in/telebot.v3"
)

// HandleAssignScenario handles the /assign_scenario command (now a function)
func HandleAssignScenario(
	getRoomHandler *roomQuery.GetRoomHandler,
	getScenarioByIDHandler *scenarioQuery.GetScenarioByIDHandler,
	addDescriptionHandler *roomCommand.AddDescriptionHandler,
	createGameHandler *gameCommand.CreateGameHandler,
	c telebot.Context,
) error {
	parts := strings.Fields(c.Message().Payload)
	if len(parts) != 2 {
		return c.Send("Usage: /assign_scenario <room_id> <scenario_id>")
	}
	roomIDStr := parts[0]
	scenarioID := parts[1]

	requester := tgutil.ToUser(c.Sender())
	if requester == nil {
		return c.Send("Could not identify requester.")
	}

	roomQuery := roomQuery.GetRoomQuery{RoomID: roomEntity.RoomID(roomIDStr)}
	room, err := getRoomHandler.Handle(context.Background(), roomQuery)
	if err != nil {
		return c.Send(fmt.Sprintf("Error finding room '%s': %v", roomIDStr, err))
	}
	if room == nil {
		return c.Send(fmt.Sprintf("Room '%s' not found.", roomIDStr))
	}

	scenarioQuery := scenarioQuery.GetScenarioByIDQuery{ID: scenarioID}
	scenario, err := getScenarioByIDHandler.Handle(context.Background(), scenarioQuery)
	if err != nil {
		return c.Send(fmt.Sprintf("Error finding scenario '%s': %v", scenarioID, err))
	}
	if scenario == nil {
		return c.Send(fmt.Sprintf("Scenario '%s' not found.", scenarioID))
	}

	roomDescriptionCmd := roomCommand.AddDescriptionCommand{
		Requester:       *requester,
		Room:            room,
		DescriptionName: "scenario_info",
		Text:            fmt.Sprintf("Scenario: %s (ID: %s)", scenario.Name, scenario.ID),
	}
	if err := addDescriptionHandler.Handle(context.Background(), roomDescriptionCmd); err != nil {
		log.Printf("Error updating room '%s' with scenario info: %v", room.Name, err)
		return c.Send(fmt.Sprintf("Error updating room '%s' with scenario info: %v", room.Name, err))
	}

	createGameCmd := gameCommand.CreateGameCommand{
		Requester:  *requester,
		RoomID:     room.ID,
		ScenarioID: scenario.ID,
	}
	createdGame, err := createGameHandler.Handle(context.Background(), createGameCmd)
	if err != nil {
		return c.Send(fmt.Sprintf("Scenario assigned, but failed to create game: %v", err))
	}

	return c.Send(fmt.Sprintf("Successfully assigned scenario '%s' (ID: %s) to room '%s' (ID: %s) and created game '%s'",
		scenario.Name, scenario.ID, room.Name, room.ID, createdGame.ID))
}
