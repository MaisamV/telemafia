package telegram

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"telemafia/internal/shared/tgutil"

	gameEntity "telemafia/internal/domain/game/entity"
	gameCommand "telemafia/internal/domain/game/usecase/command"
	gameQuery "telemafia/internal/domain/game/usecase/query"
	roomEntity "telemafia/internal/domain/room/entity"
	roomQuery "telemafia/internal/domain/room/usecase/query"
	scenarioEntity "telemafia/internal/domain/scenario/entity"
	scenarioQuery "telemafia/internal/domain/scenario/usecase/query"
	messages "telemafia/internal/presentation/telegram/messages"
	"telemafia/internal/shared/common"
	sharedEntity "telemafia/internal/shared/entity"

	"gopkg.in/telebot.v3"
)

// ---- Game Creation Callbacks ----

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
		return c.Respond(&telebot.CallbackResponse{Text: errMsg, ShowAlert: true})
	}
	if len(scenarios) == 0 {
		return c.Respond(&telebot.CallbackResponse{Text: "No scenarios available.", ShowAlert: true})
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

	promptMsg := fmt.Sprintf(msgs.Game.CreateGameSelectScenarioPrompt)
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
	if requester == nil {
		return c.Respond(&telebot.CallbackResponse{Text: msgs.Common.ErrorIdentifyRequester, ShowAlert: true})
	}

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
		return c.Respond(&telebot.CallbackResponse{Text: errMsg, ShowAlert: true})
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
		return c.Respond(&telebot.CallbackResponse{Text: errMsg, ShowAlert: true})
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
		return c.Respond(&telebot.CallbackResponse{Text: errMsg, ShowAlert: true})
	}

	// 4. Build confirmation message and keyboard
	markup := &telebot.ReplyMarkup{}
	gameIDStr := string(game.ID)
	markup.Inline(markup.Row(
		markup.Data(msgs.Game.CreateGameStartButton, tgutil.UniqueStartGame+"|"+gameIDStr),   // Direct role assignment
		markup.Data(msgs.Game.ChooseCardButton, tgutil.UniqueChooseCardStart+"|"+gameIDStr),  // Interactive role selection
		markup.Data(msgs.Game.CreateGameCancelButton, tgutil.UniqueCancelGame+"|"+gameIDStr), // Cancel creation
	))

	confirmMsg := fmt.Sprintf(msgs.Game.CreateGameConfirmPrompt,
		strings.Join(roleNames, "\n- "),
	)
	return c.Edit(confirmMsg, markup, telebot.ModeMarkdownV2)
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
		return c.Respond(&telebot.CallbackResponse{Text: msgs.Common.ErrorIdentifyRequester, ShowAlert: true})
	}

	// Permission check already done in AssignRolesHandler
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
	for user, role := range assignments {
		targetUser := &telebot.User{ID: int64(user.ID)}
		// Need Room Name - should fetch game details first?
		privateMsg := fmt.Sprintf(msgs.Game.AssignRolesSuccessPrivate, role.Name, role.Side)
		_, pmErr := bot.Send(targetUser, privateMsg, telebot.ModeMarkdownV2)
		if pmErr != nil {
			log.Printf(msgs.Game.AssignRolesErrorSendingPrivate, user.ID, pmErr)
			// Collect errors?
		}
		// Get username for public message (Requires fetching users?)
		assignResults = append(assignResults, fmt.Sprintf("%s \\-\\> %s \\(%s\\)", user.GetProfileLink(), role.Name, role.Side))
	}

	// Edit original message to show success
	finalMsg := fmt.Sprintf(msgs.Game.CreateGameStartedSuccess, strings.Join(assignResults, "\n"))
	return c.Edit(finalMsg, telebot.ModeMarkdownV2, telebot.NoPreview)
}

// --- Choose Card Flow Callbacks --- (NEW - Core Logic Implementation)

// BotHandlerInterface defines methods needed from BotHandler by callbacks
type BotHandlerInterface interface {
	GetGameByIDHandler() *gameQuery.GetGameByIDHandler
	GetPlayersInRoomHandler() *roomQuery.GetPlayersInRoomHandler
	GetScenarioByIDHandler() *scenarioQuery.GetScenarioByIDHandler
	AssignRolesHandler() *gameCommand.AssignRolesHandler
	UpdateGameHandler() *gameCommand.UpdateGameHandler
	Bot() *telebot.Bot
	GetInteractiveSelectionState(gameID gameEntity.GameID) (*tgutil.InteractiveSelectionState, bool)
	SetInteractiveSelectionState(gameID gameEntity.GameID, state *tgutil.InteractiveSelectionState)
	DeleteInteractiveSelectionState(gameID gameEntity.GameID)
	GetOrCreatePlayerRoleRefresher(gameID gameEntity.GameID) *tgutil.RefreshingMessageBook
	GetPlayerRoleRefresher(gameID gameEntity.GameID) (*tgutil.RefreshingMessageBook, bool)
	DeletePlayerRoleRefresher(gameID gameEntity.GameID)
	GetOrCreateAdminAssignmentTracker(gameID gameEntity.GameID) *tgutil.RefreshingMessageBook
	GetAdminAssignmentTracker(gameID gameEntity.GameID) (*tgutil.RefreshingMessageBook, bool)
	DeleteAdminAssignmentTracker(gameID gameEntity.GameID)
}

// HandleChooseCardStart initiates the interactive role selection process.
func HandleChooseCardStart(
	h BotHandlerInterface,
	c telebot.Context,
	gameIDStr string,
	msgs *messages.Messages,
) error {
	requester := tgutil.ToUser(c.Sender())
	if requester == nil {
		return c.Respond(&telebot.CallbackResponse{Text: msgs.Common.ErrorIdentifyRequester, ShowAlert: true})
	}
	gameID := gameEntity.GameID(gameIDStr)

	// 1. Fetch Game & Scenario, check permissions
	game, err := h.GetGameByIDHandler().Handle(context.Background(), gameQuery.GetGameByIDQuery{ID: gameID})
	if err != nil || game == nil {
		log.Printf("ChooseCardStart: Error fetching game %s: %v", gameID, err)
		return c.Respond(&telebot.CallbackResponse{Text: "Error finding game.", ShowAlert: true})
	}
	if game.Room == nil || game.Scenario == nil {
		log.Printf("ChooseCardStart: Game %s is missing room or scenario link.", gameID)
		return c.Respond(&telebot.CallbackResponse{Text: "Game data incomplete.", ShowAlert: true})
	}
	if game.State != gameEntity.GameStateWaitingForPlayers {
		return c.Respond(&telebot.CallbackResponse{Text: "Game is not in the correct state to choose cards.", ShowAlert: true})
	}

	isRoomModerator := game.Room.Moderator != nil && game.Room.Moderator.ID == requester.ID
	if !requester.Admin && !isRoomModerator {
		return c.Respond(&telebot.CallbackResponse{Text: msgs.Common.ErrorPermissionDenied, ShowAlert: true})
	}

	scenario, err := h.GetScenarioByIDHandler().Handle(context.Background(), scenarioQuery.GetScenarioByIDQuery{ID: game.Scenario.ID})
	if err != nil {
		log.Printf("ChooseCardStart: Error fetching scenario %s: %v", game.Scenario.ID, err)
		return c.Respond(&telebot.CallbackResponse{Text: "Error fetching scenario details.", ShowAlert: true})
	}

	// 2. Fetch Players
	players, err := h.GetPlayersInRoomHandler().Handle(context.Background(), roomQuery.GetPlayersInRoomQuery{RoomID: game.Room.ID})
	if err != nil {
		log.Printf("ChooseCardStart: Error fetching players for room %s: %v", game.Room.ID, err)
		return c.Respond(&telebot.CallbackResponse{Text: "Error fetching players.", ShowAlert: true})
	}

	// 3. Flatten and Shuffle Roles
	flatRoles := make([]scenarioEntity.Role, 0)
	for _, side := range scenario.Sides {
		for _, roleName := range side.Roles {
			flatRoles = append(flatRoles, scenarioEntity.Role{Name: roleName, Side: side.Name})
		}
	}
	if len(players) != len(flatRoles) {
		errMsg := fmt.Sprintf(msgs.Game.AssignRolesErrorPlayerMismatch, len(flatRoles), len(players), game.ID)
		log.Printf("ChooseCardStart: Mismatch players(%d) roles(%d) for game %s", len(players), len(flatRoles), game.ID)
		return c.Respond(&telebot.CallbackResponse{Text: errMsg, ShowAlert: true})
	}

	sort.Slice(flatRoles, func(i, j int) bool {
		return common.Hash(flatRoles[i].Name) < common.Hash(flatRoles[j].Name)
	})
	shuffledRoles := make([]scenarioEntity.Role, len(flatRoles))
	copy(shuffledRoles, flatRoles)
	common.Shuffle(len(shuffledRoles), func(i, j int) {
		shuffledRoles[i], shuffledRoles[j] = shuffledRoles[j], shuffledRoles[i]
	})
	log.Printf("Shuffled Roles: %v", shuffledRoles)

	// 4. Initialize Interactive State
	initialSelections := make(map[sharedEntity.UserID]int)
	initialTaken := make(map[int]bool)
	newState := &tgutil.InteractiveSelectionState{
		ShuffledRoles: shuffledRoles,
		Selections:    initialSelections,
		TakenIndices:  initialTaken,
	}
	h.SetInteractiveSelectionState(gameID, newState)

	// 5. Update Game State
	game.State = gameEntity.GameStateRoleSelection
	updateCmd := gameCommand.UpdateGameCommand{Game: game}
	if err := h.UpdateGameHandler().Handle(context.Background(), updateCmd); err != nil {
		log.Printf("ChooseCardStart: Failed to update game state for %s: %v", gameID, err)
		// Non-fatal, but log. The interactive state might become stale if server restarts.
	}

	// 6. Prepare & Edit Admin's Tracking Message (Placeholder Text for now)
	adminMsgContent := fmt.Sprintf(msgs.Game.AssignmentTrackingMessageAdmin, "-") // Placeholder
	adminMarkup := &telebot.ReplyMarkup{}
	adminMarkup.Inline(adminMarkup.Row(adminMarkup.Data(msgs.Game.CreateGameCancelButton, tgutil.UniqueCancelGame+"|"+string(gameID)))) // Keep cancel

	// Edit or Send the admin message
	// var adminMsg telebot.Editable // No longer needed
	sentAdminMsg, err := h.Bot().Edit(c.Message(), adminMsgContent, adminMarkup) // Try editing first
	if err != nil {                                                              // If edit fails...
		log.Printf("ChooseCardStart: Failed to EDIT admin message (%v), trying to SEND new one for %s", err, gameID)
		sentAdminMsg, err = h.Bot().Send(c.Sender(), adminMsgContent, adminMarkup) // ...try sending new
		if err != nil {
			log.Printf("ChooseCardStart: Failed to SEND new admin message for %s: %v", gameID, err)
			// If sending also fails, we can't store it. Return error.
			return c.Respond(&telebot.CallbackResponse{Text: "Failed to initiate game setup message.", ShowAlert: true})
		}
	}

	// If edit or send was successful, store the message in the refresh book
	if sentAdminMsg != nil {
		adminRefresher := h.GetOrCreateAdminAssignmentTracker(gameID)
		refreshMsg := &tgutil.RefreshingMessage{
			ChatID:    sentAdminMsg.Chat.ID, // Access field directly
			MessageID: sentAdminMsg.ID,
			Data:      string(gameID), // Store gameID as data
		}
		adminRefresher.AddActiveMessage(sentAdminMsg.Chat.ID, refreshMsg) // Use AddActiveMessage, use Chat.ID field
		log.Printf("ChooseCardStart: Added/Updated admin message in tracker for game %s", gameID)
	} else {
		// This case should theoretically not happen if error handling above is correct
		log.Printf("ChooseCardStart: Error: admin message is nil after edit/send attempt for game %s", gameID)
		return c.Respond(&telebot.CallbackResponse{Text: "Internal error storing game message.", ShowAlert: true})
	}
	// No adding to refresh book in Step 1 -> Now handled by storing message

	// 7. Send Role Selection Message to Each Player
	playerMsgMarkup, _ := PreparePlayerRoleSelectionMarkup(gameID, len(shuffledRoles), newState.TakenIndices, msgs)
	playerRefresher := h.GetOrCreatePlayerRoleRefresher(gameID)
	h.DeletePlayerRoleRefresher(gameID)                        // Use Delete method to ensure a clean slate
	playerRefresher = h.GetOrCreatePlayerRoleRefresher(gameID) // Recreate after delete
	for _, player := range players {
		if player == nil {
			continue
		}
		targetUser := &telebot.User{ID: int64(player.ID)}
		// Send player message and store it
		sentPlayerMsg, err := h.Bot().Send(targetUser, msgs.Game.RoleSelectionPromptPlayer, playerMsgMarkup)
		if err != nil {
			log.Printf("ChooseCardStart: Failed to send role selection to player %d: %v", player.ID, err)
			// Continue trying to send to others
		} else if sentPlayerMsg != nil {
			refreshMsg := &tgutil.RefreshingMessage{
				ChatID:    sentPlayerMsg.Chat.ID, // Use sentPlayerMsg.Chat.ID
				MessageID: sentPlayerMsg.ID,
				Data:      string(gameID), // Store gameID as data
			}
			playerRefresher.AddActiveMessage(sentPlayerMsg.Chat.ID, refreshMsg) // Use AddActiveMessage
			log.Printf("ChooseCardStart: Added/Updated player message in refresher for player %d in game %s", player.ID, gameID)
		}
		// No adding to refresh book in Step 1 -> Now handled by storing message
	}

	return c.Respond() // Acknowledge callback
}

// HandlePlayerSelectsCard processes a player clicking a numbered card button.
func HandlePlayerSelectsCard(
	h BotHandlerInterface,
	c telebot.Context,
	data string, // format: "gameID|index"
	msgs *messages.Messages,
) error {
	player := tgutil.ToUser(c.Sender())
	if player == nil {
		return c.Respond(&telebot.CallbackResponse{Text: msgs.Common.ErrorIdentifyUser, ShowAlert: true})
	}

	parts := strings.Split(data, "|")
	if len(parts) != 2 {
		log.Printf("PlayerSelectsCard: Invalid data format: %s", data)
		return c.Respond(&telebot.CallbackResponse{Text: "Invalid action data.", ShowAlert: true})
	}
	gameIDStr := parts[0]
	indexStr := parts[1]
	gameID := gameEntity.GameID(gameIDStr)

	chosenIndex, err := strconv.Atoi(indexStr)
	if err != nil || chosenIndex < 1 {
		log.Printf("PlayerSelectsCard: Invalid index %s: %v", indexStr, err)
		return c.Respond(&telebot.CallbackResponse{Text: "Invalid card selected.", ShowAlert: true})
	}

	// 1. Get Interactive State
	state, exists := h.GetInteractiveSelectionState(gameID)
	if !exists {
		log.Printf("PlayerSelectsCard: No interactive state found for game %s", gameID)
		_ = c.Edit("Role selection is no longer active.")
		return c.Respond(&telebot.CallbackResponse{Text: "Role selection inactive.", ShowAlert: true})
	}

	// 2. Lock State & Validate Selection
	state.Mutex.Lock()
	defer state.Mutex.Unlock()

	if chosenIndex > len(state.ShuffledRoles) {
		return c.Respond(&telebot.CallbackResponse{Text: "Invalid card number.", ShowAlert: true})
	}

	if _, alreadySelected := state.Selections[player.ID]; alreadySelected {
		return c.Respond(&telebot.CallbackResponse{Text: msgs.Game.PlayerHasRoleError, ShowAlert: true})
	}
	if state.TakenIndices[chosenIndex] {
		msg := fmt.Sprintf(msgs.Game.RoleAlreadyTakenError, chosenIndex)
		return c.Respond(&telebot.CallbackResponse{Text: msg, ShowAlert: true})
	}

	// 3. Process Selection
	state.TakenIndices[chosenIndex] = true
	state.Selections[player.ID] = chosenIndex
	selectedRole := state.ShuffledRoles[chosenIndex-1] // Adjust for 0-based index

	log.Printf("Player %d selected card %d (Role: %s) for game %s", player.ID, chosenIndex, selectedRole.Name, gameID)

	// 4. Update Game Entity Assignments
	// Fetch game again to ensure we have the latest state before updating
	game, err := h.GetGameByIDHandler().Handle(context.Background(), gameQuery.GetGameByIDQuery{ID: gameID})
	if err != nil || game == nil {
		log.Printf("PlayerSelectsCard: Failed to fetch game %s to update assignment: %v", gameID, err)
		// Non-fatal for selection, but log it
	} else {
		game.AssignRole(player.ID, selectedRole) // Assign in entity
		updateCmd := gameCommand.UpdateGameCommand{Game: game}
		if err := h.UpdateGameHandler().Handle(context.Background(), updateCmd); err != nil {
			log.Printf("PlayerSelectsCard: Failed to update game assignment %s: %v", gameID, err)
		}
	}

	// 5. Confirm to Player & Clean Up Player Message -> EDIT instead of delete
	confirmMsgText := fmt.Sprintf(msgs.Game.AssignRolesSuccessPrivate, selectedRole.Name, selectedRole.Side)
	err = c.Edit(confirmMsgText, telebot.ModeMarkdownV2) // Edit the original message
	if err != nil {
		log.Printf("PlayerSelectsCard: Failed to EDIT player confirmation message for user %d: %v", player.ID, err)
		// If editing fails, maybe try sending?
		// We don't remove from the refresher if edit fails, it might get updated later.
	} else {
		// Successfully edited, remove this player's message from the player refresher
		if playerRefresher, exists := h.GetPlayerRoleRefresher(gameID); exists {
			playerRefresher.RemoveActiveMessage(c.Message().Chat.ID)
			log.Printf("PlayerSelectsCard: Removed player %d message from refresher for game %s", player.ID, gameID)
		}
	}

	// 6. Trigger Refreshes (Added Refresh Logic)
	adminRefresher := h.GetOrCreateAdminAssignmentTracker(gameID)
	playerRefresher := h.GetOrCreatePlayerRoleRefresher(gameID)
	adminRefresher.RaiseRefreshNeeded()
	playerRefresher.RaiseRefreshNeeded()
	log.Printf("PlayerSelectsCard: Raised refresh needed flags for game %s", gameID)

	// 7. Check Completion & Update Admin Message
	allSelected := len(state.Selections) == len(state.ShuffledRoles)

	// --- Admin Message Update is now handled by the refresh timer ---
	// We don't need to prepare/edit it here directly anymore.

	if allSelected {
		log.Printf("All roles selected for game %s", gameID)
		// Final update to game state
		if game != nil { // Re-use fetched game if successful
			game.SetRolesAssigned()
			updateCmd := gameCommand.UpdateGameCommand{Game: game}
			if err := h.UpdateGameHandler().Handle(context.Background(), updateCmd); err != nil {
				log.Printf("PlayerSelectsCard: Failed to set final game state for %s: %v", gameID, err)
			}
		}
		// Trigger one last refresh to show the final state
		adminRefresher.RaiseRefreshNeeded()
		playerRefresher.RaiseRefreshNeeded()
		log.Printf("All roles selected: Triggering final refresh and cleanup for game %s", gameID)

		// Clean up state and refreshers
		h.DeleteInteractiveSelectionState(gameID)
		h.DeleteAdminAssignmentTracker(gameID) // Delete admin book
		h.DeletePlayerRoleRefresher(gameID)    // Delete player book
	}

	return c.Respond() // Acknowledge button press
}

// HandleCancelCreateGame processes the cancel button press during game creation.
func HandleCancelCreateGame(h BotHandlerInterface, c telebot.Context, msgs *messages.Messages, gameIDStr string) error {
	// Check if gameIDStr is empty (cancellation before game creation vs during role select)
	gameID := gameEntity.GameID(gameIDStr)
	if gameID != "" {
		log.Printf("Game creation cancelled for game ID: %s", gameID)
		// Clean up interactive state and refreshers if they exist
		h.DeleteInteractiveSelectionState(gameID)
		h.DeleteAdminAssignmentTracker(gameID)
		h.DeletePlayerRoleRefresher(gameID)
		log.Printf("Cancelled game %s: Cleaned up interactive state and refresh books.", gameID)
		// Optionally: Update game state to Cancelled if needed?
	} else {
		log.Println("Game creation cancelled before game ID was assigned.")
	}

	_ = c.Respond(&telebot.CallbackResponse{Text: msgs.Common.CallbackCancelled})
	return c.Delete()
}

// HandleConfirmAssignments processes the confirm button press during game creation.
func HandleConfirmAssignments(
	getGameByIDHandler *gameQuery.GetGameByIDHandler,
	c telebot.Context,
	data string,
	msgs *messages.Messages,
) error {
	gameID := data
	_, err := getGameByIDHandler.Handle(context.Background(), gameQuery.GetGameByIDQuery{ID: gameEntity.GameID(gameID)})
	if err != nil {
		log.Printf("Callback ConfirmAssignments: Error fetching game %s: %v", gameID, err)
		_ = c.Respond(&telebot.CallbackResponse{Text: fmt.Sprintf(msgs.Common.CallbackErrorGeneric, err), ShowAlert: true})
		return c.Edit(msgs.Common.CallbackFailedEdit)
	}

	_ = c.Respond()
	return c.Edit(fmt.Sprintf(msgs.Game.AssignmentsConfirmedResponse, gameID))
}

// --- Helper Functions --- (NEW)

// PreparePlayerRoleSelectionMarkup creates the numbered button grid for players.
func PreparePlayerRoleSelectionMarkup(gameID gameEntity.GameID, roleCount int, takenIndices map[int]bool, msgs *messages.Messages) (*telebot.ReplyMarkup, error) {
	markup := &telebot.ReplyMarkup{}
	var rows []telebot.Row
	var currentRow []telebot.Btn

	for i := 1; i <= roleCount; i++ {
		text := strconv.Itoa(i)
		// Payload includes game ID for context
		payload := fmt.Sprintf("%s|%d", string(gameID), i)
		unique := tgutil.UniquePlayerSelectsCard

		if takenIndices[i] {
			text = msgs.Game.RoleTakenMarker // Show taken marker
			// Rely on handler check for already taken
		}

		btn := markup.Data(text, unique, payload)
		currentRow = append(currentRow, btn)

		if len(currentRow) == 3 || i == roleCount { // 3 buttons per row or last button
			rows = append(rows, markup.Row(currentRow...))
			currentRow = []telebot.Btn{}
		}
	}
	markup.Inline(rows...)
	return markup, nil
}

// PrepareAdminAssignmentMessage creates the text for the admin's tracking message.
func PrepareAdminAssignmentMessage(game *gameEntity.Game, state *tgutil.InteractiveSelectionState, players []*sharedEntity.User, msgs *messages.Messages) (string, *telebot.ReplyMarkup, error) {
	var assignmentLines []string

	// Create a map for quick lookup of player ID -> player link (Use plain name for now)
	playerNameMap := make(map[sharedEntity.UserID]string)
	if players != nil {
		for _, p := range players {
			if p != nil {
				playerNameMap[p.ID] = p.FirstName // Use FirstName instead of GetProfileLink
			}
		}
	}

	// Create a reverse map for selection: index -> playerID
	selectionMap := make(map[int]sharedEntity.UserID)
	for userID, index := range state.Selections {
		selectionMap[index] = userID
	}

	// Iterate through shuffled roles (cards)
	for i, role := range state.ShuffledRoles {
		cardIndex := i + 1
		// Use plain text formatting
		line := fmt.Sprintf("%d. %s -> ", cardIndex, role.Name)

		if state.TakenIndices[cardIndex] {
			if playerID, ok := selectionMap[cardIndex]; ok {
				if playerName, nameOk := playerNameMap[playerID]; nameOk {
					line += playerName // Append player name
				} else {
					line += fmt.Sprintf("User %d", playerID) // Fallback to ID
				}
			} else {
				line += "(Error: Taken but no player found)"
			}
		} else {
			line += "-"
		}
		assignmentLines = append(assignmentLines, line)
	}

	allSelected := len(state.Selections) == len(state.ShuffledRoles)
	var messageText string
	markup := &telebot.ReplyMarkup{}

	if allSelected {
		// Use simplified message key without Markdown
		messageText = fmt.Sprintf("All roles selected!\n%s", strings.Join(assignmentLines, "\n"))
	} else {
		// Use simplified message key without Markdown
		messageText = fmt.Sprintf("Role Selection Progress:\n%s\nWaiting for players...", strings.Join(assignmentLines, "\n"))
		// Add cancel button
		if game != nil {
			markup.Inline(markup.Row(markup.Data(msgs.Game.CreateGameCancelButton, tgutil.UniqueCancelGame+"|"+string(game.ID))))
		}
	}

	return messageText, markup, nil
}
