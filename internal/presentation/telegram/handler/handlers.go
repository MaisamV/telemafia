package telegram

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	gameEntity "telemafia/internal/domain/game/entity"
	gameCommand "telemafia/internal/domain/game/usecase/command"
	gameQuery "telemafia/internal/domain/game/usecase/query"
	roomEntity "telemafia/internal/domain/room/entity"
	roomCommand "telemafia/internal/domain/room/usecase/command"
	roomQuery "telemafia/internal/domain/room/usecase/query"
	scenarioEntity "telemafia/internal/domain/scenario/entity"
	scenarioCommand "telemafia/internal/domain/scenario/usecase/command"
	scenarioQuery "telemafia/internal/domain/scenario/usecase/query"
	sharedEntity "telemafia/internal/shared/entity"

	"gopkg.in/telebot.v3"
)

// --- Existing handlers from handlers.go ---
// HandleJoinRoom handles the /join_room command
func (h *BotHandler) HandleJoinRoom(c telebot.Context) error {
	roomIDStr := strings.TrimSpace(c.Message().Payload)
	if roomIDStr == "" {
		return c.Send("Please provide a room ID: /join_room <room_id>")
	}

	roomID := roomEntity.RoomID(roomIDStr)
	user := ToUser(c.Sender())

	cmd := roomCommand.JoinRoomCommand{
		Requester: *user,
		RoomID:    roomID,
	}

	if err := h.joinRoomHandler.Handle(context.Background(), cmd); err != nil {
		return c.Send(fmt.Sprintf("Error joining room '%s': %v", roomID, err))
	}

	markup := &telebot.ReplyMarkup{}
	btnLeave := markup.Data(fmt.Sprintf("Leave Room %s", roomID), UniqueLeaveRoomSelectRoom, string(roomID))
	markup.Inline(markup.Row(btnLeave))

	return c.Send(fmt.Sprintf("Successfully joined room %s", roomID), markup)
}

// HandleListRooms handles the /list_rooms command
func (h *BotHandler) HandleListRooms(c telebot.Context) error {
	rooms, err := h.getRoomsHandler.Handle(context.Background(), roomQuery.GetRoomsQuery{})
	if err != nil {
		return c.Send(fmt.Sprintf("Error getting rooms: %v", err))
	}

	if len(rooms) == 0 {
		return c.Send("No rooms available.")
	}

	var response strings.Builder
	response.WriteString("Available Rooms:\n")
	for _, room := range rooms {
		players, _ := h.getPlayersInRoomHandler.Handle(context.Background(), roomQuery.GetPlayersInRoomQuery{RoomID: room.ID})
		playerCount := len(players)
		maxPlayers := 10
		response.WriteString(fmt.Sprintf("- %s (%s) [%d/%d players]\n", room.Name, room.ID, playerCount, maxPlayers))
	}

	return c.Send(response.String())
}

// HandleMyRooms handles the /my_rooms command
func (h *BotHandler) HandleMyRooms(c telebot.Context) error {
	user := ToUser(c.Sender())
	query := roomQuery.GetPlayerRoomsQuery{PlayerID: user.ID}
	rooms, err := h.getPlayerRoomsHandler.Handle(context.Background(), query)
	if err != nil {
		return c.Send(fmt.Sprintf("Error getting your rooms: %v", err))
	}

	if len(rooms) == 0 {
		return c.Send("You are not in any rooms.")
	}

	var response strings.Builder
	response.WriteString("Rooms you are in:\n")
	for _, room := range rooms {
		response.WriteString(fmt.Sprintf("- %s (%s)\n", room.Name, room.ID))
	}

	return c.Send(response.String())
}

// HandleKickUser handles the /kick_user command
func (h *BotHandler) HandleKickUser(c telebot.Context) error {
	args := strings.Fields(c.Message().Payload)
	if len(args) != 2 {
		return c.Send("Usage: /kick_user <room_id> <user_id>")
	}

	roomID := roomEntity.RoomID(args[0])
	playerIDStr := args[1]

	playerID, err := strconv.ParseInt(playerIDStr, 10, 64)
	if err != nil {
		return c.Send("Invalid user ID format.")
	}

	requester := ToUser(c.Sender())

	cmd := roomCommand.KickUserCommand{
		Requester: *requester,
		RoomID:    roomID,
		PlayerID:  sharedEntity.UserID(playerID),
	}

	if err := h.kickUserHandler.Handle(context.Background(), cmd); err != nil {
		return c.Send(fmt.Sprintf("Error kicking user %d from room %s: %v", playerID, roomID, err))
	}

	return c.Send(fmt.Sprintf("User %d kicked from room %s", playerID, roomID))
}

// --- Added handlers from other files (and fixed imports/types) ---

// HandleLeaveRoom handles the /leave_room command (from leave_room_handler.go)
func (h *BotHandler) HandleLeaveRoom(c telebot.Context) error {
	args := strings.TrimSpace(c.Message().Payload)
	if args == "" {
		return c.Send("Please provide a room ID: /leave_room [room_id]")
	}

	user := ToUser(c.Sender())
	cmd := roomCommand.LeaveRoomCommand{
		RoomID:    roomEntity.RoomID(args),
		Requester: *user,
	}
	if err := h.leaveRoomHandler.Handle(context.Background(), cmd); err != nil {
		return c.Send(fmt.Sprintf("Error leaving room '%s': %v", args, err))
	}

	return c.Send(fmt.Sprintf("Successfully left room %s!", args))
}

// HandleDeleteRoom handles the /delete_room command (from delete_room_handler.go)
func (h *BotHandler) HandleDeleteRoom(c telebot.Context) error {
	if !IsAdmin(c.Sender().Username) {
		return c.Send("You are not authorized to use this command.")
	}

	rooms, err := h.getRoomsHandler.Handle(context.Background(), roomQuery.GetRoomsQuery{})
	if err != nil {
		log.Printf("Error fetching rooms for deletion selection: %v", err)
		return c.Send("Failed to fetch rooms list.")
	}

	if len(rooms) == 0 {
		return c.Send("No rooms exist to delete.")
	}

	var buttons [][]telebot.InlineButton
	for _, room := range rooms {
		button := telebot.InlineButton{
			Unique: UniqueDeleteRoomSelectRoom,
			Text:   fmt.Sprintf("%s (%s)", room.Name, room.ID),
			Data:   string(room.ID),
		}
		buttons = append(buttons, []telebot.InlineButton{button})
	}

	markup := &telebot.ReplyMarkup{InlineKeyboard: buttons}
	return c.Send("Select a room to delete:", markup)
}

// HandleCreateScenario handles the /create_scenario command (from scenario_handler.go)
func (h *BotHandler) HandleCreateScenario(c telebot.Context) error {
	if !IsAdmin(c.Sender().Username) {
		return c.Send("You are not authorized to use this command.")
	}
	args := strings.TrimSpace(c.Message().Payload)
	if args == "" {
		return c.Send("Please provide a scenario name: /create_scenario [name]")
	}

	cmd := scenarioCommand.CreateScenarioCommand{
		ID:   fmt.Sprintf("scen_%d", time.Now().UnixNano()),
		Name: args,
	}
	if err := h.createScenarioHandler.Handle(context.Background(), cmd); err != nil {
		return c.Send(fmt.Sprintf("Error creating scenario: %v", err))
	}

	return c.Send(fmt.Sprintf("Scenario '%s' created successfully! ID: %s\nUse /add_role %s <role_name> to add roles.", cmd.Name, cmd.ID, cmd.ID))
}

// HandleDeleteScenario handles the /delete_scenario command (from scenario_handler.go)
func (h *BotHandler) HandleDeleteScenario(c telebot.Context) error {
	if !IsAdmin(c.Sender().Username) {
		return c.Send("You are not authorized to use this command.")
	}
	args := strings.TrimSpace(c.Message().Payload)
	if args == "" {
		return c.Send("Please provide a scenario ID: /delete_scenario <id>")
	}

	cmd := scenarioCommand.DeleteScenarioCommand{
		ID: args,
	}
	if err := h.deleteScenarioHandler.Handle(context.Background(), cmd); err != nil {
		return c.Send(fmt.Sprintf("Error deleting scenario '%s': %v", args, err))
	}

	return c.Send(fmt.Sprintf("Scenario %s deleted successfully!", args))
}

// HandleAddRole handles the /add_role command (from scenario_handler.go)
func (h *BotHandler) HandleAddRole(c telebot.Context) error {
	if !IsAdmin(c.Sender().Username) {
		return c.Send("You are not authorized to use this command.")
	}
	args := strings.Fields(strings.TrimSpace(c.Message().Payload))
	if len(args) != 2 {
		return c.Send("Usage: /add_role <scenario_id> <role_name>")
	}

	// Use new usecase path and command struct
	cmd := scenarioCommand.AddRoleCommand{
		ScenarioID: args[0],
		Role:       scenarioEntity.Role{Name: args[1]}, // Use new type
	}
	if err := h.manageRolesHandler.HandleAddRole(context.Background(), cmd); err != nil {
		return c.Send(fmt.Sprintf("Error adding role '%s' to scenario '%s': %v", args[1], args[0], err))
	}

	return c.Send(fmt.Sprintf("Role '%s' added to scenario %s successfully!", args[1], args[0]))
}

// HandleRemoveRole handles the /remove_role command (from scenario_handler.go)
func (h *BotHandler) HandleRemoveRole(c telebot.Context) error {
	if !IsAdmin(c.Sender().Username) {
		return c.Send("You are not authorized to use this command.")
	}
	args := strings.Fields(strings.TrimSpace(c.Message().Payload))
	if len(args) != 2 {
		return c.Send("Usage: /remove_role <scenario_id> <role_name>")
	}

	// Use new usecase path and command struct
	cmd := scenarioCommand.RemoveRoleCommand{
		ScenarioID: args[0],
		RoleName:   args[1],
	}
	if err := h.manageRolesHandler.HandleRemoveRole(context.Background(), cmd); err != nil {
		return c.Send(fmt.Sprintf("Error removing role '%s' from scenario '%s': %v", args[1], args[0], err))
	}

	return c.Send(fmt.Sprintf("Role '%s' removed from scenario %s successfully!", args[1], args[0]))
}

// HandleAssignScenario handles the /assign_scenario command
func (h *BotHandler) HandleAssignScenario(c telebot.Context) error {
	if !IsAdmin(c.Sender().Username) {
		return c.Send("You are not authorized to use this command.")
	}
	args := strings.Fields(strings.TrimSpace(c.Message().Payload)) // Use Fields for splitting
	if len(args) != 2 {                                            // Expect exactly 2 args
		return c.Send("Usage: /assign_scenario <room_id> <scenario_id>")
	}
	roomIDStr := args[0]
	scenarioIDStr := args[1]
	roomID := roomEntity.RoomID(roomIDStr) // Use new type

	// Verify the scenario exists first
	scenario, err := h.getScenarioByIDHandler.Handle(context.Background(), scenarioQuery.GetScenarioByIDQuery{ID: scenarioIDStr})
	if err != nil {
		return c.Send(fmt.Sprintf("Error finding scenario '%s': %v", scenarioIDStr, err))
	}

	// Verify the room exists and get its name for the message
	room, err := h.getRoomHandler.Handle(context.Background(), roomQuery.GetRoomQuery{RoomID: roomID})
	if err != nil {
		return c.Send(fmt.Sprintf("Error finding room '%s': %v", roomIDStr, err))
	}

	// Assign the scenario to the room using the injected RoomWriter
	// Pass scenario ID (string) not name, as repo interface expects name
	// NOTE: The RoomRepository interface currently has AssignScenarioToRoom(roomID, scenarioName string)
	// This seems inconsistent. Should probably assign by ID or update the repo interface.
	// For now, using the scenario.ID as the name for the repo call.
	if err := h.roomRepo.AssignScenarioToRoom(roomID, scenario.ID); err != nil {
		return c.Send(fmt.Sprintf("Error assigning scenario ID '%s' to room '%s': %v", scenario.ID, roomIDStr, err))
	}

	// Create a new game associated with this room and scenario
	// Use new usecase path and command struct
	createGameCmd := gameCommand.CreateGameCommand{
		RoomID:     roomID,
		ScenarioID: scenario.ID, // Use the verified scenario ID
	}
	game, err := h.createGameHandler.Handle(context.Background(), createGameCmd)
	if err != nil {
		// TODO: Should we roll back the scenario assignment if game creation fails?
		return c.Send(fmt.Sprintf("Scenario assigned, but failed to create game: %v", err))
	}

	return c.Send(fmt.Sprintf("Successfully assigned scenario '%s' (ID: %s) to room '%s' (%s) and created game '%s'",
		scenario.Name, scenario.ID, room.Name, room.ID, game.ID))
}

// HandleAssignRoles handles the /assign_roles command
func (h *BotHandler) HandleAssignRoles(c telebot.Context) error {
	if !IsAdmin(c.Sender().Username) {
		return c.Send("You are not authorized to use this command.")
	}
	gameIDStr := strings.TrimSpace(c.Message().Payload) // Use TrimSpace
	if gameIDStr == "" {
		// TODO: Maybe list available games?
		return c.Send("Please provide a game ID: /assign_roles <game_id>")
	}

	gameID := gameEntity.GameID(gameIDStr) // Use new type

	// Use new usecase path and command struct
	cmd := gameCommand.AssignRolesCommand{GameID: gameID}

	assignments, err := h.assignRolesHandler.Handle(context.Background(), cmd)
	if err != nil {
		return c.Send(fmt.Sprintf("Error assigning roles for game '%s': %v", gameID, err))
	}

	// Format and send the assignments (maybe privately?)
	responseText := h.formatAssignments(assignments) // Uses new types

	// TODO: The confirmation button logic seems misplaced here. The assignRolesHandler
	// already updates the game state and persists it. Confirmation/sending roles
	// might be a separate step or handled differently.
	// For now, just return the assignment list.
	// confirmButton := telebot.InlineButton{
	// 	Unique: UniqueConfirm, // Constant from callback_handler
	// 	Text:   "Confirm and Send Roles to Players",
	// 	Data:   string(gameID),
	// }
	// markup := &telebot.ReplyMarkup{InlineKeyboard: [][]telebot.InlineButton{{confirmButton}}}

	return c.Send(fmt.Sprintf("Roles assigned for game %s:\n%s", gameID, responseText))
}

// formatAssignments formats the role assignments map for display.
// Input map uses new types: map[sharedEntity.UserID]scenarioEntity.Role
func (h *BotHandler) formatAssignments(assignments map[sharedEntity.UserID]scenarioEntity.Role) string {
	var b strings.Builder
	// Sort by User ID for consistent output?
	ids := make([]int64, 0, len(assignments))
	for uid := range assignments {
		ids = append(ids, int64(uid))
	}
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })

	for _, id := range ids {
		uid := sharedEntity.UserID(id) // Use new type
		role := assignments[uid]
		// TODO: Get username from ID if possible/needed for better display
		// This currently fetches ALL users from ALL rooms, which is inefficient.
		// Should ideally fetch only users relevant to the current game.
		userName := h.getUserDisplayName(uid) // Use helper function

		b.WriteString(fmt.Sprintf("%s: %s\n", userName, role.Name))
	}
	return b.String()
}

// getUserDisplayName is a helper to get a display name (requires BotHandler access)
func (h *BotHandler) getUserDisplayName(userID sharedEntity.UserID) string {
	// Inefficiently search all rooms for the user - Needs optimization
	// Consider adding a GetUserByID method to a shared user repository/service.
	allRooms, _ := h.getRoomsHandler.Handle(context.Background(), roomQuery.GetRoomsQuery{})
	for _, room := range allRooms {
		players, _ := h.getPlayersInRoomHandler.Handle(context.Background(), roomQuery.GetPlayersInRoomQuery{RoomID: room.ID})
		for _, player := range players {
			if player != nil && player.ID == userID {
				if player.Username != "" {
					return "@" + player.Username
				} else if player.FirstName != "" {
					name := player.FirstName
					if player.LastName != "" {
						name += " " + player.LastName
					}
					return name
				}
			}
		}
	}
	return fmt.Sprintf("User %d", userID)
}

// HandleConfirmAssignments handles the callback to send roles privately.
// NOTE: This was originally part of role_assignment_handler.go, moved here.
// It also seems redundant if roles are assigned and saved in HandleAssignRoles.
// Keeping the logic but might need review.
func (h *BotHandler) HandleConfirmAssignments(c telebot.Context, gameIDStr string) error {
	gameID := gameEntity.GameID(gameIDStr) // Use new type
	if gameID == "" {
		return c.Respond(&telebot.CallbackResponse{Text: "Error: Missing game ID", ShowAlert: true})
	}
	log.Printf("Confirming assignments for game: '%s'", gameID)

	// Use new usecase path and query struct
	targetGame, err := h.getGameByIDHandler.Handle(context.Background(), gameQuery.GetGameByIDQuery{ID: gameID})
	if err != nil {
		log.Printf("Error fetching game with ID '%s': %v", gameID, err)
		return c.Respond(&telebot.CallbackResponse{Text: fmt.Sprintf("Game '%s' not found: %v", gameID, err), ShowAlert: true})
	}

	if targetGame.Room == nil {
		log.Printf("Game '%s' has no associated room.", gameID)
		return c.Respond(&telebot.CallbackResponse{Text: "Game has no room.", ShowAlert: true})
	}
	roomID := targetGame.Room.ID
	log.Printf("Found game '%s' for room '%s'", targetGame.ID, roomID)

	assignments := targetGame.Assignments // Get assignments directly from the fetched game entity
	if len(assignments) == 0 {
		return c.Respond(&telebot.CallbackResponse{Text: "No role assignments found for this game", ShowAlert: true})
	}

	log.Printf("Found %d role assignments for game '%s'", len(assignments), gameID)

	successCount := 0
	for userID, role := range assignments { // userID is sharedEntity.UserID, role is scenarioEntity.Role
		userChat := &telebot.Chat{ID: int64(userID)} // Convert UserID to int64 for telebot
		userName := h.getUserDisplayName(userID)

		log.Printf("Sending role %s to %s (ID: %d)", role.Name, userName, userID)
		message := fmt.Sprintf("🎭 *Your Role Assignment* 🎭\n\nYou have been assigned the role: *%s*\n\nKeep your role secret and follow the game master's instructions!", role.Name)
		_, err = h.bot.Send(userChat, message, &telebot.SendOptions{ParseMode: telebot.ModeMarkdown})

		if err != nil {
			log.Printf("Failed to send role to %s (ID: %d): %v", userName, userID, err)
		} else {
			log.Printf("Successfully sent role to %s (ID: %d)", userName, userID)
			successCount++
		}
	}

	// Edit the original message that had the confirm button
	_ = c.Edit(fmt.Sprintf("Roles sent privately to %d players for game %s.", successCount, gameID))

	return c.Respond(&telebot.CallbackResponse{Text: fmt.Sprintf("Roles sent to %d players!", successCount)})
}

// HandleGamesList handles the /games command
func (h *BotHandler) HandleGamesList(c telebot.Context) error {
	if !IsAdmin(c.Sender().Username) {
		return c.Send("You are not authorized to use this command.")
	}

	// Use new usecase path and query struct
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
	for _, game := range games { // game is *gameEntity.Game
		roomInfo := "(No Room Linked)"
		if game.Room != nil {
			roomInfo = fmt.Sprintf("(Room: %s)", game.Room.ID)
		}
		scenarioInfo := "(No Scenario Linked)"
		if game.Scenario != nil {
			// Need Scenario Name, not just ID. Fetch it? Or store Name in Game entity?
			// Assuming Scenario entity has Name field based on previous code
			scenarioName := game.Scenario.Name // Might be empty if only ID was stored
			if scenarioName == "" {
				scenarioName = "(Name Unknown)"
			}
			scenarioInfo = fmt.Sprintf("(Scenario: %s, ID: %s)", scenarioName, game.Scenario.ID)
		}
		// Assignments map key is sharedEntity.UserID
		assignmentCount := len(game.Assignments)

		response.WriteString(fmt.Sprintf("- Game ID: %s %s %s State: %s Players/Assignments: %d\n",
			game.ID, roomInfo, scenarioInfo, game.State, assignmentCount))

		// Display assignments if present
		if assignmentCount > 0 {
			response.WriteString("  Assignments:\n")
			// Use the helper function to format assignments
			response.WriteString(h.formatAssignments(game.Assignments))
		} else if game.State == gameEntity.GameStateWaitingForPlayers {
			// Suggest assigning roles if none are assigned yet and waiting for players
			response.WriteString(fmt.Sprintf("  Roles not assigned. Use /assign_roles %s\n", game.ID))
		}
		response.WriteString("\n") // Add space between games
	}

	return c.Send(response.String())
}

// --- Refresh Logic (potentially move to its own package/file later) ---

// RefreshingMessageType defines the type of content being refreshed.
type RefreshingMessageType int

const (
	ListRooms  RefreshingMessageType = iota
	RoomDetail                       // Assuming this exists or will be added
)

// RefreshingMessage stores info about a message that needs periodic updates.
// Renamed to TrackedMessage for clarity
type TrackedMessage struct {
	Msg         *telebot.Message // Store the message pointer
	MessageType RefreshingMessageType
	Data        string // e.g., room ID for RoomDetail
}

var (
	// Stores the latest known refreshing message info for each user (ChatID)
	refreshingMessages      = make(map[int64]TrackedMessage) // Store TrackedMessage struct
	refreshingMessagesMutex sync.RWMutex
)

// RefreshRoomsList handles updating dynamic messages.
// It periodically checks a flag (set by commands that modify rooms)
// and updates all registered messages if changes occurred.
func (h *BotHandler) RefreshRoomsList() {
	updateMessages := func() {
		refreshingMessagesMutex.RLock()
		messagesToUpdate := make(map[int64]TrackedMessage) // Map stores TrackedMessage now
		for userID, trackedMsg := range refreshingMessages {
			messagesToUpdate[userID] = trackedMsg
		}
		refreshingMessagesMutex.RUnlock()

		if len(messagesToUpdate) == 0 {
			return // No messages to update
		}

		log.Printf("Refreshing %d messages...", len(messagesToUpdate))
		for userID, trackedMsg := range messagesToUpdate {
			// Prepare content based on the tracked message type
			newContent, newMarkup, err := h.prepareMessageContent(trackedMsg.MessageType, trackedMsg.Data)
			if err != nil {
				log.Printf("Error preparing refresh content for user %d (type %v, data %s): %v", userID, trackedMsg.MessageType, trackedMsg.Data, err)
				continue // Skip this message if content prep fails
			}

			// Edit the existing message using the stored pointer
			_, editErr := h.bot.Edit(trackedMsg.Msg, newContent, newMarkup)
			if editErr != nil {
				// If editing fails (e.g., message deleted, chat blocked), remove it from the list
				if strings.Contains(editErr.Error(), "message to edit not found") ||
					strings.Contains(editErr.Error(), "message is not modified") ||
					strings.Contains(editErr.Error(), "bot was blocked by the user") {
					log.Printf("Removing message for user %d from refresh list (edit failed: %v)", userID, editErr)
					RemoveRefreshingChat(userID)
				} else {
					log.Printf("Non-fatal error editing message for user %d: %v", userID, editErr)
				}
			}
		}
		log.Println("Finished refreshing messages.")
	}

	ticker := time.NewTicker(5 * time.Second) // Check every 5 seconds (adjust as needed)
	defer ticker.Stop()
	for range ticker.C {
		// Use new usecase paths and command/query structs
		if h.checkRefreshHandler.Handle(context.Background(), roomQuery.CheckChangeFlagQuery{}) {
			h.resetRefreshHandler.Handle(context.Background(), roomCommand.ResetChangeFlagCommand{}) // Consume flag
			updateMessages()
		}
	}
}

// prepareMessageContent generates content based on message type.
func (h *BotHandler) prepareMessageContent(messageType RefreshingMessageType, data string) (string, *telebot.ReplyMarkup, error) {
	switch messageType {
	case ListRooms:
		return h.prepareListRoomsMessage()
	// case RoomDetail:
	// 	 return h.prepareRoomDetailMessage(data) // Assuming this func exists
	default:
		return "", nil, fmt.Errorf("unsupported refreshing message type: %v", messageType)
	}
}

// prepareListRoomsMessage generates the text and markup for the list rooms view
func (h *BotHandler) prepareListRoomsMessage() (string, *telebot.ReplyMarkup, error) {
	// Use new usecase path and query struct
	rooms, err := h.getRoomsHandler.Handle(context.Background(), roomQuery.GetRoomsQuery{})
	if err != nil {
		return "", nil, fmt.Errorf("error getting rooms: %w", err)
	}

	var response strings.Builder
	markup := &telebot.ReplyMarkup{}
	var rows []telebot.Row

	if len(rooms) == 0 {
		response.WriteString("No rooms available.")
	} else {
		response.WriteString("Available Rooms (refreshed):\n")
		for _, room := range rooms {
			// Use new usecase path and query struct
			players, _ := h.getPlayersInRoomHandler.Handle(context.Background(), roomQuery.GetPlayersInRoomQuery{RoomID: room.ID})
			playerCount := len(players)
			maxPlayers := 10 // TODO: Define elsewhere
			response.WriteString(fmt.Sprintf("- %s (%s) [%d/%d players]\n", room.Name, room.ID, playerCount, maxPlayers))
			// Add join button
			btnJoin := markup.Data(fmt.Sprintf("Join %s", room.Name), UniqueJoinRoom, string(room.ID))
			rows = append(rows, markup.Row(btnJoin))
		}
	}
	markup.Inline(rows...)
	return response.String(), markup, nil
}

// SendOrUpdateRefreshingMessage sends a new message and registers it for refreshing, or updates an existing one.
func (h *BotHandler) SendOrUpdateRefreshingMessage(userID int64, messageType RefreshingMessageType, data string) error {
	// Prepare content based on type
	content, markup, err := h.prepareMessageContent(messageType, data)
	if err != nil {
		log.Printf("Error preparing content for user %d type %v: %v", userID, messageType, err)
		_, sendErr := h.bot.Send(telebot.ChatID(userID), "Error preparing dynamic message content.")
		if sendErr != nil {
			log.Printf("Failed to send error message to user %d: %v", userID, sendErr)
		}
		return err
	}

	refreshingMessagesMutex.Lock()
	defer refreshingMessagesMutex.Unlock()

	if existingTrackedMsg, ok := refreshingMessages[userID]; ok {
		// Attempt to edit the existing message
		updatedMsg, editErr := h.bot.Edit(existingTrackedMsg.Msg, content, markup)
		if editErr == nil {
			// Update stored message pointer and metadata
			refreshingMessages[userID] = TrackedMessage{
				Msg:         updatedMsg,
				MessageType: messageType,
				Data:        data,
			}
			log.Printf("Successfully updated refreshing message for user %d", userID)
			return nil
		}
		log.Printf("Failed to edit refreshing message %d for user %d, sending new: %v", existingTrackedMsg.Msg.ID, userID, editErr)
		delete(refreshingMessages, userID)
	}

	// Send a new message
	msg, sendErr := h.bot.Send(telebot.ChatID(userID), content, markup)
	if sendErr != nil {
		log.Printf("Error sending new refreshing message to user %d: %v", userID, sendErr)
		return sendErr
	}

	// Register the new message
	refreshingMessages[userID] = TrackedMessage{
		Msg:         msg,
		MessageType: messageType,
		Data:        data,
	}
	log.Printf("Sent and registered new refreshing message %d for user %d", msg.ID, userID)
	return nil
}

// ChangeRefreshType updates the type of message being refreshed for a user.
func ChangeRefreshType(userID int64, messageType RefreshingMessageType, data string) {
	refreshingMessagesMutex.Lock()
	defer refreshingMessagesMutex.Unlock()
	if trackedMsg, exists := refreshingMessages[userID]; exists {
		// Update the metadata, keep the same *telebot.Message pointer
		trackedMsg.MessageType = messageType
		trackedMsg.Data = data
		refreshingMessages[userID] = trackedMsg
		log.Printf("Changed refresh type for user %d to %v (data: %s)", userID, messageType, data)
	} else {
		log.Printf("Attempted to change refresh type for user %d, but no existing message found.", userID)
	}
}

// GetRefreshingChats returns a snapshot of the chats being refreshed.
// The returned slice contains copies, safe for concurrent reading.
type RefreshingChat struct { // Define struct locally (consider moving if used elsewhere)
	userID  int64
	message TrackedMessage // Use TrackedMessage
}

func GetRefreshingChats() []RefreshingChat {
	refreshingMessagesMutex.RLock()
	defer refreshingMessagesMutex.RUnlock()
	userMessages := make([]RefreshingChat, 0, len(refreshingMessages))
	for userID, trackedMsg := range refreshingMessages {
		userMessages = append(userMessages, RefreshingChat{userID: userID, message: trackedMsg})
	}
	return userMessages
}

// RemoveRefreshingChat removes a chat from the refreshing list.
func RemoveRefreshingChat(userID int64) {
	refreshingMessagesMutex.Lock()
	defer refreshingMessagesMutex.Unlock()
	delete(refreshingMessages, userID)
	log.Printf("Removed user %d from refreshing messages list.", userID)
}

// SplitCallbackData helper function to parse callback data (unique:payload)
func SplitCallbackData(data string) (unique string, payload string) {
	parts := strings.SplitN(data, ":", 2)
	if len(parts) == 2 {
		return parts[0], parts[1]
	} else if len(parts) == 1 {
		// Handle cases where only unique identifier might be present (e.g., cancel buttons)
		switch parts[0] {
		case UniqueCancel: // Check against known payload-less constants
			return parts[0], ""
		default:
			// If it doesn't match known payload-less uniques, assume it's the unique part
			// and the payload is empty. This might need adjustment based on actual usage.
			log.Printf("Warning: Callback data '%s' has no colon, assuming empty payload.", data)
			return parts[0], ""
		}
	}
	log.Printf("Warning: Could not split callback data '%s' correctly.", data)
	return "", data // Fallback, should ideally not happen with standard format
}

// --- Callback Logic ---

// Unique identifiers for inline buttons
// NOTE: These were scattered, consolidating here.
// Ensure these match the constants used when creating buttons.
const (
	// Join/Leave related
	UniqueJoinRoom            = "join_room"
	UniqueLeaveRoomSelectRoom = "leave_room_select"
	UniqueLeaveRoomConfirm    = "leave_room_confirm"

	// Delete Room related
	UniqueDeleteRoomSelectRoom = "delete_room_select"
	UniqueDeleteRoomConfirm    = "delete_room_confirm"

	// Game/Assignment related
	UniqueConfirmAssignments = "confirm_assignments"

	// Generic Cancel (might need context)
	UniqueCancel = "cancel"
)

// HandleCallback routes all callback queries from inline buttons.
func (h *BotHandler) HandleCallback(c telebot.Context) error {
	callback := c.Callback()
	if callback == nil {
		log.Println("Received update that is not a callback")
		return nil // Ignore non-callback updates
	}

	// Data often contains unique:payload or just unique
	// Use the helper to split consistently.
	unique, data := SplitCallbackData(callback.Data)
	userID := c.Sender().ID
	log.Printf("Callback received: User=%d, Unique=%s, Data=%s", userID, unique, data)

	switch unique {
	case UniqueJoinRoom:
		return h.handleJoinRoomCallback(c, data)
	case UniqueDeleteRoomSelectRoom:
		return h.handleDeleteRoomSelectCallback(c, data)
	case UniqueDeleteRoomConfirm:
		return h.handleDeleteRoomConfirmCallback(c, data)
	case UniqueLeaveRoomSelectRoom:
		return h.handleLeaveRoomSelectCallback(c, data)
	case UniqueLeaveRoomConfirm:
		return h.handleLeaveRoomConfirmCallback(c, data)
	case UniqueConfirmAssignments:
		return h.HandleConfirmAssignments(c, data)
	case UniqueCancel:
		// Handle generic cancel: usually just edit the message back or delete it.
		log.Printf("User %d cancelled operation.", userID)
		_ = c.Respond(&telebot.CallbackResponse{Text: "Operation cancelled."})
		_ = c.Delete() // Delete the message with the cancel button
		return nil
	default:
		log.Printf("Unknown callback unique identifier: %s", unique)
		return c.Respond(&telebot.CallbackResponse{Text: "Unknown action."})
	}
}

// --- Simple Handlers ---

// HandleHelp provides a simple help message.
func (h *BotHandler) HandleHelp(c telebot.Context) error {
	// Update help text based on actual implemented commands & structure
	help := `Available commands:
/start - Show welcome message & rooms
/help - Show this help message
/list_rooms - List all available rooms
/my_rooms - List rooms you have joined
/join_room <room_id> - Join a specific room
/leave_room <room_id> - Leave the specified room

Admin Commands:
/create_room <room_name> - Create a new room
/delete_room - Select a room to delete
/kick_user <room_id> <user_id> - Kick a user from a room
/create_scenario <scenario_name> - Create a new game scenario
/delete_scenario <scenario_id> - Delete a scenario
/add_role <scenario_id> <role_name> - Add a role to a scenario
/remove_role <scenario_id> <role_name> - Remove a role from a scenario
/assign_scenario <room_id> <scenario_id> - Assign a scenario to a room (creates a game)
/games - List active games and their status
/assign_roles <game_id> - Assign roles to players in a game`
	return c.Send(help, &telebot.SendOptions{DisableWebPagePreview: true})
}

// HandleStart handles the /start command
func (h *BotHandler) HandleStart(c telebot.Context) error {
	// Send the dynamic list rooms message initially?
	// Or send a welcome message first?
	_ = c.Send(fmt.Sprintf("Welcome, %s!", c.Sender().Username))
	return h.SendOrUpdateRefreshingMessage(c.Sender().ID, ListRooms, "")
}

// handleJoinRoomCallback handles the callback for joining a room.
func (h *BotHandler) handleJoinRoomCallback(c telebot.Context, roomIDStr string) error {
	roomID := roomEntity.RoomID(roomIDStr) // Use new type
	user := ToUser(c.Sender())             // Assuming *sharedEntity.User

	// Use new usecase path and command struct
	cmd := roomCommand.JoinRoomCommand{
		Requester: *user,
		RoomID:    roomID,
	}

	if err := h.joinRoomHandler.Handle(context.Background(), cmd); err != nil {
		_ = c.Respond(&telebot.CallbackResponse{Text: fmt.Sprintf("Error: %v", err), ShowAlert: true})
		return err // Return error for logging
	}

	_ = c.Respond(&telebot.CallbackResponse{Text: fmt.Sprintf("Joined room %s!", roomID)})
	// Optionally edit the original message
	// _ = c.Edit(fmt.Sprintf("You joined room %s.", roomID))
	return nil
}

// handleDeleteRoomSelectCallback asks for confirmation.
func (h *BotHandler) handleDeleteRoomSelectCallback(c telebot.Context, roomIDStr string) error {
	if !IsAdmin(c.Sender().Username) {
		return c.Respond(&telebot.CallbackResponse{Text: "Unauthorized.", ShowAlert: true})
	}

	markup := &telebot.ReplyMarkup{}
	// Use the constants for unique identifiers
	confirmData := fmt.Sprintf("%s:%s", UniqueDeleteRoomConfirm, roomIDStr)
	btnConfirm := markup.Data("Yes, delete it!", confirmData)
	btnCancel := markup.Data("Cancel", UniqueCancel)
	markup.Inline(markup.Row(btnConfirm, btnCancel))

	// Edit the original message to ask for confirmation
	err := c.Edit(fmt.Sprintf("Are you sure you want to delete room %s?", roomIDStr), markup)
	if err != nil {
		log.Printf("Error editing message for delete confirmation: %v", err)
		// Attempt to respond to callback even if edit fails
		_ = c.Respond(&telebot.CallbackResponse{Text: "Error showing confirmation."})
		return err // Return error for logging
	}
	// Respond to the initial callback to dismiss the loading indicator
	_ = c.Respond()
	return nil
}

// handleDeleteRoomConfirmCallback performs the deletion.
func (h *BotHandler) handleDeleteRoomConfirmCallback(c telebot.Context, roomIDStr string) error {
	if !IsAdmin(c.Sender().Username) {
		return c.Respond(&telebot.CallbackResponse{Text: "Unauthorized.", ShowAlert: true})
	}

	roomID := roomEntity.RoomID(roomIDStr) // Use new type
	// Use new usecase path and command struct
	cmd := roomCommand.DeleteRoomCommand{RoomID: roomID}

	if err := h.deleteRoomHandler.Handle(context.Background(), cmd); err != nil {
		_ = c.Respond(&telebot.CallbackResponse{Text: fmt.Sprintf("Error deleting room: %v", err), ShowAlert: true})
		return err
	}

	_ = c.Respond(&telebot.CallbackResponse{Text: fmt.Sprintf("Room %s deleted.", roomID)})
	// Edit the original message to confirm deletion and remove buttons
	_ = c.Edit(fmt.Sprintf("Room %s has been deleted.", roomID), &telebot.ReplyMarkup{}) // Clear keyboard
	return nil
}

// handleLeaveRoomSelectCallback asks for confirmation to leave.
func (h *BotHandler) handleLeaveRoomSelectCallback(c telebot.Context, roomIDStr string) error {
	markup := &telebot.ReplyMarkup{}
	confirmData := fmt.Sprintf("%s:%s", UniqueLeaveRoomConfirm, roomIDStr)
	btnConfirm := markup.Data("Yes, leave", confirmData)
	btnCancel := markup.Data("Cancel", UniqueCancel)
	markup.Inline(markup.Row(btnConfirm, btnCancel))

	err := c.Edit(fmt.Sprintf("Are you sure you want to leave room %s?", roomIDStr), markup)
	if err != nil {
		log.Printf("Error editing message for leave confirmation: %v", err)
		_ = c.Respond(&telebot.CallbackResponse{Text: "Error showing confirmation."})
		return err
	}
	_ = c.Respond()
	return nil
}

// handleLeaveRoomConfirmCallback performs leaving the room.
func (h *BotHandler) handleLeaveRoomConfirmCallback(c telebot.Context, roomIDStr string) error {
	user := ToUser(c.Sender()) // Assuming *sharedEntity.User
	// Use new usecase path and command struct
	cmd := roomCommand.LeaveRoomCommand{
		RoomID:    roomEntity.RoomID(roomIDStr), // Use new type
		Requester: *user,
	}
	if err := h.leaveRoomHandler.Handle(context.Background(), cmd); err != nil {
		_ = c.Respond(&telebot.CallbackResponse{Text: fmt.Sprintf("Error leaving room: %v", err), ShowAlert: true})
		return err
	}

	_ = c.Respond(&telebot.CallbackResponse{Text: fmt.Sprintf("You left room %s.", roomIDStr)})
	// Edit the original message to confirm and remove buttons
	_ = c.Edit(fmt.Sprintf("You have left room %s.", roomIDStr), &telebot.ReplyMarkup{}) // Clear keyboard
	return nil
}
