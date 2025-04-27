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
	messages "telemafia/internal/presentation/telegram/messages"
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
	msgs *messages.Messages,
) error {
	parts := strings.Fields(c.Message().Payload)
	if len(parts) != 2 {
		return c.Send(msgs.Game.AssignScenarioPrompt)
	}
	roomIDStr := parts[0]
	scenarioID := parts[1]

	requester := tgutil.ToUser(c.Sender())
	if requester == nil || !requester.Admin {
		return c.Send(msgs.Common.ErrorPermissionDenied)
	}

	roomQuery := roomQuery.GetRoomQuery{RoomID: roomEntity.RoomID(roomIDStr)}
	room, err := getRoomHandler.Handle(context.Background(), roomQuery)
	if err != nil {
		return c.Send(fmt.Sprintf(msgs.Game.AssignScenarioErrorRoomFind, roomIDStr, err))
	}
	if room == nil {
		return c.Send(fmt.Sprintf(msgs.Game.AssignScenarioErrorRoomNotfound, roomIDStr))
	}

	scenarioQuery := scenarioQuery.GetScenarioByIDQuery{ID: scenarioID}
	scenario, err := getScenarioByIDHandler.Handle(context.Background(), scenarioQuery)
	if err != nil {
		return c.Send(fmt.Sprintf(msgs.Game.AssignScenarioErrorScenarioFind, scenarioID, err))
	}
	if scenario == nil {
		return c.Send(fmt.Sprintf(msgs.Game.AssignScenarioErrorScenarioNotfound, scenarioID))
	}

	roomDescriptionCmd := roomCommand.AddDescriptionCommand{
		Requester:       *requester,
		Room:            room,
		DescriptionName: "scenario_info",
		Text:            fmt.Sprintf("Scenario: %s (ID: %s)", scenario.Name, scenario.ID),
	}
	if err := addDescriptionHandler.Handle(context.Background(), roomDescriptionCmd); err != nil {
		log.Printf("Error updating room '%s' with scenario info: %v", room.Name, err)
		return c.Send(fmt.Sprintf(msgs.Game.AssignScenarioErrorUpdateRoom, room.Name, err))
	}

	createGameCmd := gameCommand.CreateGameCommand{
		Requester:  *requester,
		RoomID:     room.ID,
		ScenarioID: scenario.ID,
	}
	createdGame, err := createGameHandler.Handle(context.Background(), createGameCmd)
	if err != nil {
		return c.Send(fmt.Sprintf(msgs.Game.AssignScenarioErrorGameCreate, err))
	}

	return c.Send(fmt.Sprintf(msgs.Game.AssignScenarioSuccess,
		scenario.Name, scenario.ID, room.Name, room.ID, createdGame.ID))
}
