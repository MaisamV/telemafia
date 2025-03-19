package telegram

import (
	"context"
	"fmt"
	"strings"
	"telemafia/delivery/util"
	gameCommand "telemafia/internal/game/usecase/command"
	"telemafia/internal/room/entity"
	roomCommand "telemafia/internal/room/usecase/command"
	scenarioQuery "telemafia/internal/scenario/usecase/query"

	"gopkg.in/telebot.v3"
)

// HandleAssignScenario handles the /assign_scenario command
func (h *BotHandler) HandleAssignScenario(c telebot.Context) error {
	if !util.IsAdmin(c.Sender().Username) {
		return c.Send("You are not authorized to use this command.")
	}

	// Format: /assign_scenario [room_id] [scenario_name]
	args := strings.Split(c.Message().Payload, " ")
	if len(args) < 2 {
		return c.Send("Please provide a room ID and scenario name: /assign_scenario [room_id] [scenario_name]")
	}

	roomID := args[0]
	scenarioID := args[1]

	// Verify the scenario exists
	scenario, err := h.getScenarioByIDHandler.Handle(context.Background(), scenarioQuery.GetScenarioByIDQuery{ID: scenarioID})
	if err != nil {
		return c.Send(fmt.Sprintf("Error fetching scenario: %v", err))
	}

	// Assign the scenario to the room
	cmd := roomCommand.AssignScenarioCommand{
		RoomID:          entity.RoomID(roomID),
		DescriptionName: "Scenario",
		Text:            scenario.Name,
	}
	err = h.assignScenarioHandler.Handle(context.Background(), cmd)
	if err != nil {
		return c.Send(fmt.Sprintf("Error assigning scenario: %v", err))
	}

	// Create a new game with this room and scenario
	createGameCmd := gameCommand.CreateGameCommand{
		RoomID:       entity.RoomID(roomID),
		ScenarioID:   scenarioID,
		ScenarioName: scenario.Name,
	}
	game, err := h.createGameHandler.Handle(createGameCmd)
	if err != nil {
		return c.Send(fmt.Sprintf("Error creating game: %v", err))
	}

	// Signal refresh to update room information
	h.raiseChangeFlagHandler.Handle(context.Background(), roomCommand.RaiseChangeFlagCommand{})

	return c.Send(fmt.Sprintf("Successfully assigned scenario '%s' to room '%s' and created game '%s'",
		scenario.Name, roomID, game.ID))
}
