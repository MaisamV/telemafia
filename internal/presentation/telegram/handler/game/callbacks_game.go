package telegram

import (
	"context"
	"fmt"
	"log"
	"strings"
	"telemafia/internal/shared/tgutil"

	"gopkg.in/telebot.v3"
	gameEntity "telemafia/internal/domain/game/entity"
	gameCommand "telemafia/internal/domain/game/usecase/command"
	gameQuery "telemafia/internal/domain/game/usecase/query"
	roomEntity "telemafia/internal/domain/room/entity"
	roomQuery "telemafia/internal/domain/room/usecase/query"
	scenarioEntity "telemafia/internal/domain/scenario/entity"
	scenarioQuery "telemafia/internal/domain/scenario/usecase/query"
	messages "telemafia/internal/presentation/telegram/messages"
)

// HandleSelectRoomForCreateGame processes the room selection during game creation.
func HandleSelectRoomForCreateGame(
	getAllScenariosHandler *scenarioQuery.GetAllScenariosHandler,
	c telebot.Context,
	roomID string,
	msgs *messages.Messages,
) error {
	// Fetch scenarios
	scenarios, err := getAllScenariosHandler.Handle(context.Background(), scenarioQuery.GetAllScenariosQuery{})
	if err != nil {
		errMsg := fmt.Sprintf(msgs.Game.CreateGameErrorFetchScenarios, err)
		log.Printf("Callback creategame_room: Error fetching scenarios: %v", err)
		_ = c.Respond(&telebot.CallbackResponse{Text: errMsg, ShowAlert: true})
		return c.Edit(msgs.Common.CallbackFailedEdit)
	}
	if len(scenarios) == 0 {
		_ = c.Respond(&telebot.CallbackResponse{Text: "No scenarios available.", ShowAlert: true})
		return c.Edit(msgs.Common.CallbackFailedEdit)
	}

	// Build scenario keyboard
	markup := &telebot.ReplyMarkup{}
	var rows []telebot.Row
	for _, scenario := range scenarios {

		log.Printf("Scenario name: %s", scenario.Name)
		btn := markup.Data(
			scenario.Name,
			tgutil.UniqueCreateGameSelectScenario,
			fmt.Sprintf("%s|%s", roomID, scenario.ID),
		)
		rows = append(rows, markup.Row(btn))
	}
	rows = append(rows, markup.Row(markup.Data(msgs.Game.CreateGameCancelButton, tgutil.UniqueCancelGame)))
	markup.Inline(rows...)

	// Edit message to ask for scenario
	promptMsg := fmt.Sprintf(msgs.Game.CreateGameSelectScenarioPrompt, roomID) // Assuming room name might be better? Fetch room details?
	return c.Edit(promptMsg, markup)
}

// HandleSelectScenarioForCreateGame processes the scenario selection, creates the game, and shows confirmation.
func HandleSelectScenarioForCreateGame(
	createGameHandler *gameCommand.CreateGameHandler,
	getPlayersInRoomHandler *roomQuery.GetPlayersInRoomHandler,
	getScenarioByIDHandler *scenarioQuery.GetScenarioByIDHandler,
	c telebot.Context,
	roomID string,
	scenarioID string,
	msgs *messages.Messages,
) error {
	requester := tgutil.ToUser(c.Sender())

	// 1. Create the Game entity
	cmd := gameCommand.CreateGameCommand{
		RoomID:     roomEntity.RoomID(roomID),
		ScenarioID: scenarioID,
		Requester:  *requester,
	}
	game, err := createGameHandler.Handle(context.Background(), cmd)
	if err != nil {
		errMsg := fmt.Sprintf(msgs.Game.CreateGameErrorCreatingGame, err)
		log.Printf("Callback creategame_scen: Error creating game: %v", err)
		return c.Respond(&telebot.CallbackResponse{Text: errMsg, ShowAlert: true})
	}

	// 2. Fetch players for display
	players, err := getPlayersInRoomHandler.Handle(context.Background(), roomQuery.GetPlayersInRoomQuery{RoomID: roomEntity.RoomID(roomID)})
	if err != nil {
		errMsg := fmt.Sprintf(msgs.Game.CreateGameErrorFetchPlayers, err)
		log.Printf("Callback creategame_scen: Error fetching players for room %s: %v", roomID, err)
		_ = c.Respond(&telebot.CallbackResponse{Text: errMsg, ShowAlert: true})
		return c.Edit(msgs.Common.CallbackFailedEdit)
	}
	var playerNames []string
	for _, p := range players {
		if p != nil {
			playerNames = append(playerNames, p.Username)
		}
	}

	// 3. Fetch scenario details for display (and role count check)
	scenario, err := getScenarioByIDHandler.Handle(context.Background(), scenarioQuery.GetScenarioByIDQuery{ID: scenarioID})
	if err != nil {
		errMsg := fmt.Sprintf(msgs.Game.CreateGameErrorFetchScenarioDetails, err)
		log.Printf("Callback creategame_scen: Error fetching scenario %s: %v", scenarioID, err)
		_ = c.Respond(&telebot.CallbackResponse{Text: errMsg, ShowAlert: true})
		return c.Edit(msgs.Common.CallbackFailedEdit)
	}

	// Flatten roles for display and count
	flatRoles := make([]scenarioEntity.Role, 0)
	var roleNames []string
	for _, side := range scenario.Sides {
		for _, roleName := range side.Roles {
			flatRoles = append(flatRoles, scenarioEntity.Role{Name: roleName, Side: side.Name})
			roleNames = append(roleNames, roleName)
		}
	}

	// Optional: Early check if player count matches role count
	if len(players) != len(flatRoles) {
		errMsg := fmt.Sprintf(msgs.Game.AssignRolesErrorPlayerMismatch, len(flatRoles), len(players), game.ID)
		log.Printf("Callback creategame_scen: Mismatch players(%d) roles(%d) for game %s", len(players), len(flatRoles), game.ID)
		_ = c.Respond(&telebot.CallbackResponse{Text: errMsg, ShowAlert: true})
		// Maybe delete the created game record here?
		return c.Edit(msgs.Common.CallbackFailedEdit)
	}

	// 4. Build confirmation message and keyboard
	markup := &telebot.ReplyMarkup{}
	markup.Inline(markup.Row(
		markup.Data(msgs.Game.CreateGameStartButton, fmt.Sprintf("%s|%s", tgutil.UniqueStartGame, game.ID)),
		markup.Data(msgs.Game.CreateGameCancelButton, fmt.Sprintf("%s|%s", tgutil.UniqueCancelGame, game.ID)), // Pass gameID to cancel
	))

	confirmMsg := fmt.Sprintf(msgs.Game.CreateGameConfirmPrompt,
		roomID, // Use room/scenario names?
		scenarioID,
		strings.Join(playerNames, "\n- "),
		strings.Join(roleNames, "\n- "),
	)
	return c.Edit(confirmMsg, markup)
}

// HandleStartCreatedGame processes the start button press, assigning roles.
func HandleStartCreatedGame(
	assignRolesHandler *gameCommand.AssignRolesHandler,
	bot *telebot.Bot, // Need bot to send private messages
	c telebot.Context,
	gameID string,
	msgs *messages.Messages,
) error {
	requester := tgutil.ToUser(c.Sender())
	if requester == nil {
		_ = c.Respond(&telebot.CallbackResponse{Text: msgs.Common.ErrorIdentifyRequester, ShowAlert: true})
		return c.Edit(msgs.Common.CallbackFailedEdit)
	}

	cmd := gameCommand.AssignRolesCommand{
		Requester: *requester,
		GameID:    gameEntity.GameID(gameID),
	}

	assignments, err := assignRolesHandler.Handle(context.Background(), cmd)
	if err != nil {
		errMsg := fmt.Sprintf(msgs.Game.CreateGameErrorAssigningRoles, err)
		log.Printf("Callback creategame_start: Error assigning roles for game %s: %v", gameID, err)
		_ = c.Respond(&telebot.CallbackResponse{Text: errMsg, ShowAlert: true})
		// Don't edit the message here, let admin retry or cancel?
		return fmt.Errorf("failed to assign roles: %w", err) // Return error for potential upstream handling
	}

	// Send private messages
	var assignResults []string
	for userID, role := range assignments {
		targetUser := &telebot.User{ID: int64(userID)}
		// Need Room Name - should fetch game details first?
		privateMsg := fmt.Sprintf(msgs.Game.AssignRolesSuccessPrivate, gameID, "<Room Name Placeholder>", role.Name)
		_, pmErr := bot.Send(targetUser, privateMsg)
		if pmErr != nil {
			log.Printf(msgs.Game.AssignRolesErrorSendingPrivate, userID, pmErr)
			// Collect errors?
		}
		// Get username for public message (Requires fetching users?)
		assignResults = append(assignResults, fmt.Sprintf("%d -> %s (%s)", userID, role.Name, role.Side))
	}

	// Edit original message to show success
	finalMsg := fmt.Sprintf(msgs.Game.CreateGameStartedSuccess, strings.Join(assignResults, "\n"))
	return c.Edit(finalMsg)
}

// HandleCancelCreateGame handles cancellation at any step.
func HandleCancelCreateGame(c telebot.Context, msgs *messages.Messages, gameIDMaybe ...string) error {
	// Optional: If gameID is passed (e.g., cancel_creategame_game123), maybe delete the game record?
	if len(gameIDMaybe) > 0 {
		gameID := gameIDMaybe[0]
		log.Printf("Game creation cancelled for potential game ID: %s", gameID)
		//TODO: Add logic to delete game record if necessary
	}

	_ = c.Respond(&telebot.CallbackResponse{Text: msgs.Common.CallbackCancelled})
	return c.Delete()
}

// HandleConfirmAssignments sends a public confirmation after roles are assigned (placeholder/example)
func HandleConfirmAssignments(
	getGameByIDHandler *gameQuery.GetGameByIDHandler,
	c telebot.Context,
	data string,
	msgs *messages.Messages,
) error {
	gameID := data // Assuming data is the game ID

	// Optional: Fetch game details if needed for the message
	_, err := getGameByIDHandler.Handle(context.Background(), gameQuery.GetGameByIDQuery{ID: gameEntity.GameID(gameID)})
	if err != nil {
		log.Printf("Callback ConfirmAssignments: Error fetching game %s: %v", gameID, err)
		// Use msg
		_ = c.Respond(&telebot.CallbackResponse{Text: fmt.Sprintf(msgs.Common.CallbackErrorGeneric, err), ShowAlert: true})
		return c.Edit(msgs.Common.CallbackFailedEdit) // Use msg
	}

	// Acknowledge the callback
	_ = c.Respond()

	// Edit the original message or send a new one
	// Use msg
	return c.Edit(fmt.Sprintf(msgs.Game.AssignmentsConfirmedResponse, gameID))
}

// handleShowMyRoleCallback could potentially show a user their role again
// func handleShowMyRoleCallback(h *BotHandler, c telebot.Context, data string) error {
// 	// Requires fetching game, finding user's assignment, and sending PM
// 	// ... implementation needed ...
// }
