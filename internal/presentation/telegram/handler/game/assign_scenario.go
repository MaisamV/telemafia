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

	"gopkg.in/telebot.v3"
)

// HandleAssignScenario handles the /assign_scenario command (now a function)
func HandleAssignScenario(h *BotHandler, c telebot.Context) error {
	args := strings.Fields(strings.TrimSpace(c.Message().Payload))
	if len(args) != 2 {
		return c.Send("Usage: /assign_scenario <room_id> <scenario_id>")
	}

	roomID := roomEntity.RoomID(args[0])
	scenarioID := args[1]
	requester := ToUser(c.Sender())
	if requester == nil {
		return c.Send("Could not identify user.")
	}

	// 1. Fetch Room (Needed for AddDescription)
	room, err := h.getRoomHandler.Handle(context.Background(), roomQuery.GetRoomQuery{RoomID: roomID})
	if err != nil {
		return c.Send(fmt.Sprintf("Error fetching room '%s': %v", roomID, err))
	}
	if room == nil {
		return c.Send(fmt.Sprintf("Room '%s' not found.", roomID))
	}

	// 2. Fetch Scenario (Needed for Description Text)
	scen, err := h.getScenarioByIDHandler.Handle(context.Background(), scenarioQuery.GetScenarioByIDQuery{ID: scenarioID})
	if err != nil {
		return c.Send(fmt.Sprintf("Error fetching scenario '%s': %v", scenarioID, err))
	}
	if scen == nil {
		return c.Send(fmt.Sprintf("Scenario '%s' not found.", scenarioID))
	}

	// 3. Add/Update Scenario Description in Room (Admin check happens in AddDescriptionHandler)
	descCmd := roomCommand.AddDescriptionCommand{
		Requester:       *requester,
		Room:            room,
		DescriptionName: "scenario_info",
		Text:            fmt.Sprintf("Scenario: %s (%d roles)", scen.Name, len(scen.Roles)),
	}
	if err := h.addDescriptionHandler.Handle(context.Background(), descCmd); err != nil {
		log.Printf("Error adding scenario description to room '%s': %v", roomID, err)
		// Continue even if description fails, game creation is more critical
	}

	// 4. Create the Game (Admin check happens in CreateGameHandler)
	createGameCmd := gameCommand.CreateGameCommand{
		Requester:  *requester,
		RoomID:     roomID,
		ScenarioID: scenarioID,
	}

	game, err := h.createGameHandler.Handle(context.Background(), createGameCmd)
	if err != nil {
		return c.Send(fmt.Sprintf("Error creating game for room '%s' with scenario '%s': %v", roomID, scenarioID, err))
	}

	return c.Send(fmt.Sprintf("Game created (ID: %s) and scenario '%s' assigned to room '%s'!", game.ID, scenarioID, roomID))
}
