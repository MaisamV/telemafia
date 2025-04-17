package telegram

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"telemafia/internal/entity"
	"telemafia/internal/usecase"
	"time"

	"gopkg.in/telebot.v3"
)

// --- Existing handlers from handlers.go ---
// HandleJoinRoom handles the /join_room command
func (h *BotHandler) HandleJoinRoom(c telebot.Context) error {
	roomIDStr := strings.TrimSpace(c.Message().Payload)
	if roomIDStr == "" {
		return c.Send("Please provide a room ID: /join_room <room_id>")
	}

	roomID := entity.RoomID(roomIDStr)
	user := ToUser(c.Sender())

	cmd := usecase.JoinRoomCommand{
		Requester: *user,
		RoomID:    roomID,
	}

	if err := h.joinRoomHandler.Handle(context.Background(), cmd); err != nil {
		return c.Send(fmt.Sprintf("Error joining room: %v", err))
	}

	markup := &telebot.ReplyMarkup{}
	btnLeave := markup.Data(fmt.Sprintf("Leave Room %s", roomID), UniqueLeaveRoomSelectRoom, string(roomID))
	markup.Inline(markup.Row(btnLeave))

	return c.Send(fmt.Sprintf("Successfully joined room %s", roomID), markup)
}

// HandleListRooms handles the /list_rooms command
func (h *BotHandler) HandleListRooms(c telebot.Context) error {
	rooms, err := h.getRoomsHandler.Handle(context.Background())
	if err != nil {
		return c.Send(fmt.Sprintf("Error getting rooms: %v", err))
	}

	if len(rooms) == 0 {
		return c.Send("No rooms available.")
	}

	var response strings.Builder
	response.WriteString("Available Rooms:\n")
	for _, room := range rooms {
		players, _ := h.getPlayersInRoomHandler.Handle(context.Background(), usecase.GetPlayersInRoomQuery{RoomID: room.ID})
		response.WriteString(fmt.Sprintf("- %s (%s) [%d/%d players]\n", room.Name, room.ID, len(players), 10))
	}

	return c.Send(response.String())
}

// HandleMyRooms handles the /my_rooms command
func (h *BotHandler) HandleMyRooms(c telebot.Context) error {
	user := ToUser(c.Sender())
	query := usecase.GetPlayerRoomsQuery{PlayerID: user.ID}
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

	roomID := entity.RoomID(args[0])
	playerIDStr := args[1]

	playerID, err := strconv.ParseInt(playerIDStr, 10, 64)
	if err != nil {
		return c.Send("Invalid user ID format.")
	}

	requester := ToUser(c.Sender())

	cmd := usecase.KickUserCommand{
		Requester: *requester,
		RoomID:    roomID,
		PlayerID:  entity.UserID(playerID),
	}

	if err := h.kickUserHandler.Handle(context.Background(), cmd); err != nil {
		return c.Send(fmt.Sprintf("Error kicking user: %v", err))
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
	cmd := usecase.LeaveRoomCommand{
		RoomID:    entity.RoomID(args),
		Requester: *user,
	}
	if err := h.leaveRoomHandler.Handle(context.Background(), cmd); err != nil {
		return c.Send(fmt.Sprintf("Error leaving room: %v", err)) // Simple error formatting
	}

	return c.Send("Successfully left the room!")
}

// HandleDeleteRoom handles the /delete_room command (from delete_room_handler.go)
func (h *BotHandler) HandleDeleteRoom(c telebot.Context) error {
	if !IsAdmin(c.Sender().Username) {
		return c.Send("You are not authorized to use this command.")
	}

	rooms, err := h.getRoomsHandler.Handle(context.Background()) // No query needed for GetRooms
	if err != nil {
		return c.Send("Failed to fetch rooms.")
	}

	var buttons [][]telebot.InlineButton
	for _, room := range rooms {
		button := telebot.InlineButton{
			Unique: UniqueDeleteRoomSelectRoom, // Constant from callback_handler
			Text:   room.Name,
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

	cmd := usecase.CreateScenarioCommand{
		ID:   fmt.Sprintf("scenario_%d", time.Now().UnixNano()),
		Name: args,
	}
	if err := h.createScenarioHandler.Handle(context.Background(), cmd); err != nil {
		return c.Send(fmt.Sprintf("Error creating scenario: %v", err))
	}

	return c.Send(fmt.Sprintf("Scenario '%s' created successfully! ID: %s\nUse /add_role %s [role_name] to add roles.", cmd.Name, cmd.ID, cmd.ID))
}

// HandleDeleteScenario handles the /delete_scenario command (from scenario_handler.go)
func (h *BotHandler) HandleDeleteScenario(c telebot.Context) error {
	if !IsAdmin(c.Sender().Username) {
		return c.Send("You are not authorized to use this command.")
	}
	args := strings.TrimSpace(c.Message().Payload)
	if args == "" {
		return c.Send("Please provide a scenario ID: /delete_scenario [id]")
	}

	cmd := usecase.DeleteScenarioCommand{
		ID: args,
	}
	if err := h.deleteScenarioHandler.Handle(context.Background(), cmd); err != nil {
		return c.Send(fmt.Sprintf("Error deleting scenario: %v", err))
	}

	return c.Send("Scenario deleted successfully!")
}

// HandleAddRole handles the /add_role command (from scenario_handler.go)
func (h *BotHandler) HandleAddRole(c telebot.Context) error {
	if !IsAdmin(c.Sender().Username) {
		return c.Send("You are not authorized to use this command.")
	}
	args := strings.Split(strings.TrimSpace(c.Message().Payload), " ")
	if len(args) != 2 {
		return c.Send("Usage: /add_role [scenario_id] [role_name]")
	}

	cmd := usecase.AddRoleCommand{
		ScenarioID: args[0],
		Role:       entity.Role{Name: args[1]}, // Use entity package
	}
	if err := h.manageRolesHandler.HandleAddRole(context.Background(), cmd); err != nil {
		return c.Send(fmt.Sprintf("Error adding role: %v", err))
	}

	return c.Send("Role added successfully!")
}

// HandleRemoveRole handles the /remove_role command (from scenario_handler.go)
func (h *BotHandler) HandleRemoveRole(c telebot.Context) error {
	if !IsAdmin(c.Sender().Username) {
		return c.Send("You are not authorized to use this command.")
	}
	args := strings.Split(strings.TrimSpace(c.Message().Payload), " ")
	if len(args) != 2 {
		return c.Send("Usage: /remove_role [scenario_id] [role_name]")
	}

	cmd := usecase.RemoveRoleCommand{
		ScenarioID: args[0],
		RoleName:   args[1],
	}
	if err := h.manageRolesHandler.HandleRemoveRole(context.Background(), cmd); err != nil {
		return c.Send(fmt.Sprintf("Error removing role: %v", err))
	}

	return c.Send("Role removed successfully!")
}

// HandleAssignScenario handles the /assign_scenario command
func (h *BotHandler) HandleAssignScenario(c telebot.Context) error {
	if !IsAdmin(c.Sender().Username) {
		return c.Send("You are not authorized to use this command.")
	}

	args := strings.Split(c.Message().Payload, " ")
	if len(args) < 2 {
		return c.Send("Usage: /assign_scenario [room_id] [scenario_id]")
	}

	roomIDStr := args[0]
	scenarioID := args[1]
	roomID := entity.RoomID(roomIDStr)

	// Verify the room exists (optional, AssignScenarioToRoom might handle it)
	_, err := h.getRoomHandler.Handle(context.Background(), usecase.GetRoomQuery{RoomID: roomID})
	if err != nil {
		return c.Send(fmt.Sprintf("Error finding room: %v", err))
	}

	// Verify the scenario exists
	scenario, err := h.getScenarioByIDHandler.Handle(context.Background(), usecase.GetScenarioByIDQuery{ID: scenarioID})
	if err != nil {
		return c.Send(fmt.Sprintf("Error fetching scenario: %v", err))
	}

	// Assign the scenario to the room using the injected RoomWriter
	err = h.roomRepo.AssignScenarioToRoom(roomID, scenario.Name) // Using roomRepo now
	if err != nil {
		return c.Send(fmt.Sprintf("Error assigning scenario to room: %v", err))
	}

	// Create a new game with this room and scenario
	createGameCmd := usecase.CreateGameCommand{
		RoomID:     roomID,
		ScenarioID: scenarioID,
	}
	game, err := h.createGameHandler.Handle(context.Background(), createGameCmd)
	if err != nil {
		return c.Send(fmt.Sprintf("Error creating game: %v", err))
	}

	// Repository should handle raising the change flag internally now

	return c.Send(fmt.Sprintf("Successfully assigned scenario '%s' to room '%s' and created game '%s'",
		scenario.Name, roomID, game.ID))
}

// HandleAssignRoles handles the /assign_roles command (from role_assignment_handler.go)
func (h *BotHandler) HandleAssignRoles(c telebot.Context) error {
	if !IsAdmin(c.Sender().Username) {
		return c.Send("You are not authorized to use this command.")
	}

	args := strings.Split(c.Message().Payload, " ")
	if len(args) < 1 {
		return c.Send("Please provide a game ID: /assign_roles [game_id]")
	}

	gameIDStr := args[0]
	gameID := entity.GameID(gameIDStr)
	log.Printf("Assigning roles for game '%s'", gameID)

	targetGame, err := h.getGameByIDHandler.Handle(context.Background(), usecase.GetGameByIDQuery{ID: gameID})
	if err != nil {
		log.Printf("Error fetching game with ID '%s': %v", gameID, err)
		return c.Send(fmt.Sprintf("Game with ID '%s' not found: %v", gameID, err))
	}

	log.Printf("Found game '%s' for room '%s'", targetGame.ID, targetGame.Room.ID)

	users, err := h.getPlayersInRoomHandler.Handle(context.Background(), usecase.GetPlayersInRoomQuery{RoomID: targetGame.Room.ID})
	if err != nil {
		log.Printf("Error fetching users for room '%s': %v", targetGame.Room.ID, err)
		return c.Send(fmt.Sprintf("Error fetching users: %v", err))
	}

	log.Printf("Found %d users in room '%s'", len(users), targetGame.Room.ID)

	userList := make([]entity.User, len(users))
	for i, u := range users {
		userList[i] = *u
	}

	cmd := usecase.AssignRolesCommand{
		GameID: gameID,
		Users:  userList,
		Roles:  []entity.Role{}, // Roles will be fetched from the game's scenario
	}
	assignments, err := h.assignRolesHandler.Handle(cmd)
	if err != nil {
		log.Printf("Error assigning roles: %v", err)
		return c.Send(fmt.Sprintf("Error assigning roles: %v", err))
	}
	targetGame.Assignments = assignments // Should the handler return the updated game?

	log.Printf("Successfully assigned %d roles", len(assignments))

	assignmentText := h.formatAssignments(assignments)
	confirmButton := telebot.InlineButton{
		Unique: UniqueConfirm, // Constant from callback_handler
		Text:   "Confirm and Send Roles to Players",
		Data:   string(gameID),
	}
	markup := &telebot.ReplyMarkup{InlineKeyboard: [][]telebot.InlineButton{{confirmButton}}}
	return c.Send(fmt.Sprintf("Role Assignments for Game '%s':\n\n%s", gameID, assignmentText), markup)
}

// formatAssignments formats the role assignments for display (from role_assignment_handler.go)
func (h *BotHandler) formatAssignments(assignments map[entity.UserID]entity.Role) string {
	var result strings.Builder
	usersMap := make(map[entity.UserID]*entity.User)

	allRooms, _ := h.getRoomsHandler.Handle(context.Background())
	for _, room := range allRooms {
		users, _ := h.getPlayersInRoomHandler.Handle(context.Background(), usecase.GetPlayersInRoomQuery{RoomID: room.ID})
		for _, user := range users {
			usersMap[user.ID] = user
		}
	}

	for userID, role := range assignments {
		var userName string
		if user, found := usersMap[userID]; found {
			if user.Username != "" {
				userName = "@" + user.Username
			} else {
				userName = user.FirstName
				if user.LastName != "" {
					userName += " " + user.LastName
				}
			}
		} else {
			userName = fmt.Sprintf("User %d", userID)
		}
		result.WriteString(fmt.Sprintf("%s: %s\n", userName, role.Name))
	}
	return result.String()
}

// HandleConfirmAssignments handles the 'confirm_assignments' callback (from role_assignment_handler.go)
func (h *BotHandler) HandleConfirmAssignments(c telebot.Context, gameIDStr string) error {
	gameID := entity.GameID(gameIDStr)
	if gameID == "" {
		return c.Respond(&telebot.CallbackResponse{Text: "Error: Missing game ID", ShowAlert: true})
	}
	log.Printf("Confirming assignments for game: '%s'", gameID)

	targetGame, err := h.getGameByIDHandler.Handle(context.Background(), usecase.GetGameByIDQuery{ID: gameID})
	if err != nil {
		log.Printf("Error fetching game with ID '%s': %v", gameID, err)
		return c.Respond(&telebot.CallbackResponse{Text: fmt.Sprintf("Game with ID '%s' not found: %v", gameID, err), ShowAlert: true})
	}

	roomID := targetGame.Room.ID
	log.Printf("Found game '%s' for room '%s'", targetGame.ID, roomID)

	users, err := h.getPlayersInRoomHandler.Handle(context.Background(), usecase.GetPlayersInRoomQuery{RoomID: roomID})
	if err != nil {
		log.Printf("Error getting players in room '%s': %v", roomID, err)
	} else {
		log.Printf("Found %d players in room '%s'", len(users), roomID)
	}

	usersMap := make(map[entity.UserID]*entity.User)
	for _, user := range users {
		usersMap[user.ID] = user
	}

	assignments, err := h.assignRolesHandler.GetAssignments(string(gameID)) // Assuming GetAssignments returns map[string]entity.Role
	if err != nil {
		log.Printf("Error retrieving assignments for game '%s': %v", gameID, err)
		return c.Respond(&telebot.CallbackResponse{Text: fmt.Sprintf("Error retrieving assignments: %v", err), ShowAlert: true})
	}

	if len(assignments) == 0 {
		return c.Respond(&telebot.CallbackResponse{Text: "No role assignments found for this game", ShowAlert: true})
	}

	log.Printf("Found %d role assignments for game '%s'", len(assignments), gameID)

	successCount := 0
	for userIDStr, role := range assignments {
		userIDInt, err := strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			log.Printf("Invalid user ID '%s': %v", userIDStr, err)
			continue
		}
		userID := entity.UserID(userIDInt)
		userChat := &telebot.Chat{ID: userIDInt}

		var userName string
		if user, found := usersMap[userID]; found {
			if user.Username != "" {
				userName = "@" + user.Username
			} else {
				userName = user.FirstName
				if user.LastName != "" {
					userName += " " + user.LastName
				}
			}
		} else {
			userName = fmt.Sprintf("User %d", userIDInt)
		}

		log.Printf("Sending role %s to %s (ID: %d)", role.Name, userName, userIDInt)
		message := fmt.Sprintf("🎭 *Your Role Assignment* 🎭\n\nYou have been assigned the role: *%s*\n\nKeep your role secret and follow the game master's instructions!", role.Name)
		_, err = h.bot.Send(userChat, message, &telebot.SendOptions{ParseMode: telebot.ModeMarkdown})

		if err != nil {
			log.Printf("Failed to send role to %s (ID: %d): %v", userName, userIDInt, err)
		} else {
			log.Printf("Successfully sent role to %s (ID: %d)", userName, userIDInt)
			successCount++
		}
	}

	return c.Respond(&telebot.CallbackResponse{Text: fmt.Sprintf("Roles have been sent to %d players!", successCount), ShowAlert: true})
}

// HandleGamesList handles the /games command
func (h *BotHandler) HandleGamesList(c telebot.Context) error {
	if !IsAdmin(c.Sender().Username) {
		return c.Send("You are not authorized to use this command.")
	}
	log.Printf("Handling /games command from user %s", c.Sender().Username)

	games, err := h.getGamesHandler.Handle(context.Background())
	if err != nil {
		log.Printf("Error getting games: %v", err)
		return c.Send("Error retrieving games. Please try again later.")
	}

	if len(games) == 0 {
		return c.Send("No active games found.")
	}

	var response strings.Builder
	response.WriteString("Active Games:\n\n")
	for i, game := range games {
		response.WriteString(fmt.Sprintf("%d. Game ID: %s\n", i+1, game.ID))
		response.WriteString(fmt.Sprintf("   Room ID: %s\n", game.Room.ID))
		response.WriteString(fmt.Sprintf("   Scenario: %s (ID: %s)\n", game.Scenario.Name, game.Scenario.ID))
		response.WriteString(fmt.Sprintf("   Status: %s\n", game.State))
		if len(game.Assignments) > 0 {
			response.WriteString("   Player Assignments:\n")
			users, _ := h.getPlayersInRoomHandler.Handle(context.Background(), usecase.GetPlayersInRoomQuery{RoomID: game.Room.ID})
			usersMap := make(map[entity.UserID]*entity.User)
			for _, u := range users {
				usersMap[u.ID] = u
			}

			for userID, role := range game.Assignments {
				userName := fmt.Sprintf("User %d", userID)
				if user, found := usersMap[userID]; found {
					if user.Username != "" {
						userName = "@" + user.Username
					} else {
						userName = user.FirstName
					}
				}
				response.WriteString(fmt.Sprintf("   - %s: %s\n", userName, role.Name))
			}
		} else {
			response.WriteString(fmt.Sprintf("   Roles not assigned yet. Use /assign_roles %s\n", game.ID))
		}
		response.WriteString("\n")
	}
	return c.Send(response.String())
}

// --- Refresh Logic (from refresh_rooms.go) ---

type RefreshingMessageType int

const (
	ListRooms  RefreshingMessageType = iota
	RoomDetail                       // Assuming this exists or will be added
)

type RefreshingMessage struct {
	ID          int
	messageType RefreshingMessageType
	data        string // e.g., room ID for RoomDetail
}

var (
	refreshingChats      = make(map[int64]RefreshingMessage)
	refreshingChatsMutex = &sync.RWMutex{}
)

// RefreshRoomsList handles updating dynamic messages (from refresh_rooms.go)
func (h *BotHandler) RefreshRoomsList() {
	updateMessages := func() {
		refreshingChats := GetRefreshingChats()

		for _, chatInfo := range refreshingChats {
			var text string
			var markup *telebot.ReplyMarkup
			var err error

			switch chatInfo.message.messageType {
			case ListRooms:
				text, markup, err = h.prepareListRoomsMessage()
			// case RoomDetail: // Add logic for RoomDetail if needed
			// 	 text, markup, err = h.prepareRoomDetailMessage(chatInfo.message.data)
			// 	 if err != nil {
			// 		 // If room detail fails (e.g., room deleted), fallback to list view
			// 		 ChangeRefreshType(chatInfo.userID, ListRooms, "")
			// 		 text, markup, err = h.prepareListRoomsMessage()
			// 	 }
			default:
				log.Printf("Unknown refreshing message type: %v", chatInfo.message.messageType)
				continue
			}

			if err != nil {
				log.Printf("Error preparing refresh message for user %d: %v", chatInfo.userID, err)
			} else {
				_, editErr := h.bot.Edit(&telebot.Message{ID: chatInfo.message.ID, Chat: &telebot.Chat{ID: chatInfo.userID}}, text, markup)
				if editErr != nil {
					log.Printf("Error editing message %d for user %d: %v", chatInfo.message.ID, chatInfo.userID, editErr)
					// Maybe remove from refreshing chats if message is gone?
					if strings.Contains(editErr.Error(), "message to edit not found") {
						RemoveRefreshingChat(chatInfo.userID)
					}
				}
			}
		}
	}

	ticker := time.NewTicker(2 * time.Second) // Check every 2 seconds
	defer ticker.Stop()
	for range ticker.C {
		if h.checkRefreshHandler.Handle(context.Background()) { // Check flag using usecase handler
			h.resetRefreshHandler.Handle(context.Background()) // Consume flag using usecase handler
			updateMessages()
		}
	}
}

// prepareListRoomsMessage generates the text and markup for the list rooms view
func (h *BotHandler) prepareListRoomsMessage() (string, *telebot.ReplyMarkup, error) {
	rooms, err := h.getRoomsHandler.Handle(context.Background())
	if err != nil {
		return "", nil, fmt.Errorf("error getting rooms: %w", err)
	}

	if len(rooms) == 0 {
		return "No rooms available.", nil, nil
	}

	var response strings.Builder
	response.WriteString("Available Rooms:\n")
	markup := &telebot.ReplyMarkup{}
	var rows []telebot.Row

	for _, room := range rooms {
		players, _ := h.getPlayersInRoomHandler.Handle(context.Background(), usecase.GetPlayersInRoomQuery{RoomID: room.ID})
		response.WriteString(fmt.Sprintf("- %s (%s) [%d/%d players]\n", room.Name, room.ID, len(players), 10))
		// Add join button
		btnJoin := markup.Data(fmt.Sprintf("Join %s", room.Name), UniqueJoinSelectRoom, string(room.ID))
		rows = append(rows, markup.Row(btnJoin))
	}
	markup.Inline(rows...)
	return response.String(), markup, nil
}

// SendOrUpdateRefreshingMessage sends a new message and registers it for refreshing, or updates an existing one.
func (h *BotHandler) SendOrUpdateRefreshingMessage(userID int64, messageType RefreshingMessageType, data string) error {
	text, markup, err := h.prepareMessageContent(messageType, data)
	if err != nil {
		log.Printf("Error preparing message content for user %d: %v", userID, err)
		// Attempt to send a fallback message
		fallbackText, fallbackMarkup, fallbackErr := h.prepareListRoomsMessage()
		if fallbackErr == nil {
			_, sendErr := h.bot.Send(&telebot.Chat{ID: userID}, fallbackText, fallbackMarkup)
			if sendErr != nil {
				log.Printf("Error sending fallback message to user %d: %v", userID, sendErr)
			}
		}
		return fmt.Errorf("failed to prepare original message: %w", err)
	}

	refreshingChatsMutex.Lock()
	prevMsg, exists := refreshingChats[userID]

	if exists {
		_, editErr := h.bot.Edit(&telebot.Message{ID: prevMsg.ID, Chat: &telebot.Chat{ID: userID}}, text, markup)
		if editErr == nil {
			// Update message type and data if edit succeeds
			refreshingChats[userID] = RefreshingMessage{ID: prevMsg.ID, messageType: messageType, data: data}
			refreshingChatsMutex.Unlock()
			return nil
		}
		log.Printf("Failed to edit message %d for user %d, will send new: %v", prevMsg.ID, userID, editErr)
		// If edit fails (e.g., message deleted), proceed to send a new message
		delete(refreshingChats, userID)
	}

	// Send new message
	msg, sendErr := h.bot.Send(&telebot.Chat{ID: userID}, text, markup)
	if sendErr != nil {
		refreshingChatsMutex.Unlock()
		log.Printf("Error sending new refreshing message to user %d: %v", userID, sendErr)
		return sendErr
	}

	// Register the new message
	refreshingChats[userID] = RefreshingMessage{ID: msg.ID, messageType: messageType, data: data}
	refreshingChatsMutex.Unlock()
	return nil
}

func (h *BotHandler) prepareMessageContent(messageType RefreshingMessageType, data string) (string, *telebot.ReplyMarkup, error) {
	switch messageType {
	case ListRooms:
		return h.prepareListRoomsMessage()
	// case RoomDetail: // Add logic if needed
	// 	 return h.prepareRoomDetailMessage(data)
	default:
		return "", nil, fmt.Errorf("unknown message type: %v", messageType)
	}
}

// ChangeRefreshType updates the type of message being refreshed for a user (from refresh_rooms.go)
func ChangeRefreshType(userID int64, messageType RefreshingMessageType, data string) {
	refreshingChatsMutex.Lock()
	defer refreshingChatsMutex.Unlock()
	if msg, exists := refreshingChats[userID]; exists {
		msg.messageType = messageType
		msg.data = data
		refreshingChats[userID] = msg
	} else {
		// If no message exists, we probably can't change its type.
		// Maybe log this? Or should we register a new one?
		log.Printf("Attempted to change refresh type for user %d, but no existing message found.", userID)
	}
}

// GetRefreshingChats returns a snapshot of the chats being refreshed (from refresh_rooms.go)
type RefreshingChat struct { // Define struct locally
	userID  int64
	message RefreshingMessage
}

func GetRefreshingChats() []RefreshingChat {
	refreshingChatsMutex.RLock()
	defer refreshingChatsMutex.RUnlock()
	userMessages := make([]RefreshingChat, 0, len(refreshingChats))
	for userID, message := range refreshingChats {
		userMessages = append(userMessages, RefreshingChat{userID: userID, message: message})
	}
	return userMessages
}

// RemoveRefreshingChat removes a chat from the refreshing list (e.g., if message deleted)
func RemoveRefreshingChat(userID int64) {
	refreshingChatsMutex.Lock()
	defer refreshingChatsMutex.Unlock()
	delete(refreshingChats, userID)
	log.Printf("Removed user %d from refreshing chats list.", userID)
}

// --- Callback Logic (from callback_handler.go) ---

// Unique identifiers for inline buttons
const (
	UniqueJoinSelectRoom = "join"
	// UniqueKickSelectRoom           = "kick_selectRoom" // Example, adjust as needed
	// UniqueKickFromRoomSelectPlayer = "kickFromRoom_selectPlayer"
	UniqueDeleteRoomSelectRoom = "delete_room"
	UniqueLeaveRoomSelectRoom  = "leave_room"
	UniqueConfirm              = "confirm_assignments"
)

// HandleCallback handles all callbacks from inline buttons
func (h *BotHandler) HandleCallback(c telebot.Context) error {
	callback := c.Callback()
	if callback == nil {
		return nil // Ignore non-callback updates
	}

	unique := callback.Unique
	data := callback.Data // This is the specific data (e.g., room ID, game ID)

	log.Printf("Received callback: unique=%s, data=%s", unique, data)

	switch unique {
	case UniqueConfirm:
		return h.HandleConfirmAssignments(c, data)
	case UniqueJoinSelectRoom:
		// Need a join room callback handler? Or handle directly?
		log.Printf("Join room callback for room: %s", data)
		return c.Respond(&telebot.CallbackResponse{Text: fmt.Sprintf("Joining room %s...", data)})
	// case UniqueKickSelectRoom:
	// 	 return h.HandleKickUserCallback(c, data)
	// case UniqueKickFromRoomSelectPlayer:
	// 	 return h.HandleKickUserFromRoomCallback(c, data)
	case UniqueDeleteRoomSelectRoom:
		return h.HandleDeleteRoomCallback(c, data)
	case UniqueLeaveRoomSelectRoom:
		return h.handleLeaveCallback(c, data) // Use the existing leave callback logic
	default:
		log.Printf("Unknown callback unique ID: %s", unique)
		return c.Respond(&telebot.CallbackResponse{Text: "Unknown callback!"})
	}
}

// HandleDeleteRoomCallback handles the callback to delete a specific room (from delete_room_handler.go)
func (h *BotHandler) HandleDeleteRoomCallback(c telebot.Context, roomIDStr string) error {
	if !IsAdmin(c.Sender().Username) {
		return c.Respond(&telebot.CallbackResponse{Text: "Unauthorized."})
	}
	cmd := usecase.DeleteRoomCommand{
		RoomID: entity.RoomID(roomIDStr),
	}
	if err := h.deleteRoomHandler.Handle(context.Background(), cmd); err != nil {
		return c.Respond(&telebot.CallbackResponse{Text: "Failed to delete room."})
	}
	// Optionally edit the original message to remove buttons or show confirmation
	c.Edit("Room successfully deleted!", &telebot.ReplyMarkup{}) // Remove keyboard
	return c.Respond(&telebot.CallbackResponse{Text: "Room successfully deleted!"})
}

// handleLeaveCallback handles the leave room callback (from join_room_handler.go)
func (h *BotHandler) handleLeaveCallback(c telebot.Context, roomIDStr string) error {
	roomID := entity.RoomID(roomIDStr)
	user := ToUser(c.Sender())

	cmd := usecase.LeaveRoomCommand{
		Requester: *user,
		RoomID:    roomID,
	}

	if err := h.leaveRoomHandler.Handle(context.Background(), cmd); err != nil {
		return c.Respond(&telebot.CallbackResponse{Text: fmt.Sprintf("Error leaving room: %v", err)})
	}

	// Update the message view for the user
	ChangeRefreshType(c.Sender().ID, ListRooms, "")               // Change view back to list
	h.SendOrUpdateRefreshingMessage(c.Sender().ID, ListRooms, "") // Trigger update

	return c.Respond(&telebot.CallbackResponse{Text: "You have left the room."})
}

// --- Simple Handlers (from help_handler.go, start_handler.go) ---

// HandleHelp handles the /help command
func (h *BotHandler) HandleHelp(c telebot.Context) error {
	// Update help text based on actual implemented commands
	help := `Available commands:
/start - Show available rooms
/help - Show this help message
/list_rooms - List all available rooms
/my_rooms - List rooms you have joined
/join_room <room_id> - Join a specific room
/leave_room <room_id> - Leave the specified room

Admin Commands:
/create_room <room_name> - Create a new room
/delete_room <room_id> - Delete a room
/kick_user <room_id> <user_id> - Kick a user from a room
/create_scenario <scenario_name> - Create a new game scenario
/delete_scenario <scenario_id> - Delete a scenario
/add_role <scenario_id> <role_name> - Add a role to a scenario
/remove_role <scenario_id> <role_name> - Remove a role from a scenario
/assign_scenario <room_id> <scenario_id> - Assign a scenario to a room (creates a game)
/games - List active games and their status
/assign_roles <game_id> - Assign roles to players in a game`
	return c.Send(help)
}

// HandleStart handles the /start command
func (h *BotHandler) HandleStart(c telebot.Context) error {
	// Send the dynamic list rooms message
	return h.SendOrUpdateRefreshingMessage(c.Sender().ID, ListRooms, "")
}
