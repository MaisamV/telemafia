package telegram

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"telemafia/delivery/util"
	gameEntity "telemafia/internal/game/entity"
	gameCommand "telemafia/internal/game/usecase/command"
	gameQuery "telemafia/internal/game/usecase/query"
	"telemafia/internal/room/entity"
	roomCommand "telemafia/internal/room/usecase/command"
	roomQuery "telemafia/internal/room/usecase/query"
	scenarioEntity "telemafia/internal/scenario/entity"
	userEntity "telemafia/internal/user/entity"

	"gopkg.in/telebot.v3"
)

// HandleAssignRoles handles the /assign_roles command
func (h *BotHandler) HandleAssignRoles(c telebot.Context) error {
	if !util.IsAdmin(c.Sender().Username) {
		return c.Send("You are not authorized to use this command.")
	}

	// Format: /assign_roles [game_id]
	args := strings.Split(c.Message().Payload, " ")
	if len(args) < 1 {
		return c.Send("Please provide a game ID: /assign_roles [game_id]")
	}

	gameID := args[0]
	log.Printf("Assigning roles for game '%s'", gameID)

	// Get the game directly by ID
	targetGame, err := h.getGameByIDHandler.Handle(context.Background(),
		gameQuery.GetGameByIDQuery{ID: gameEntity.GameID(gameID)})
	if err != nil {
		log.Printf("Error fetching game with ID '%s': %v", gameID, err)
		return c.Send(fmt.Sprintf("Game with ID '%s' not found: %v", gameID, err))
	}

	log.Printf("Found game '%s' for room '%s'", targetGame.ID, targetGame.Room.ID)

	// Fetch users for the game's room
	users, err := h.getPlayersInRoomHandler.Handle(context.Background(),
		roomQuery.GetPlayersInRoomQuery{RoomID: targetGame.Room.ID})
	if err != nil {
		log.Printf("Error fetching users for room '%s': %v", targetGame.Room.ID, err)
		return c.Send(fmt.Sprintf("Error fetching users: %v", err))
	}

	log.Printf("Found %d users in room '%s'", len(users), targetGame.Room.ID)

	// Convert users to the correct type
	userList := make([]userEntity.User, len(users))
	for i, u := range users {
		userList[i] = *u
	}

	// Create and handle the assign roles command
	cmd := gameCommand.AssignRolesCommand{
		GameID: gameEntity.GameID(gameID),
		Users:  userList,
		Roles:  []scenarioEntity.Role{}, // Roles will be fetched from the game's scenario
	}
	assignments, err := h.assignRolesHandler.Handle(cmd)
	if err != nil {
		log.Printf("Error assigning roles: %v", err)
		return c.Send(fmt.Sprintf("Error assigning roles: %v", err))
	}
	targetGame.Assignments = assignments

	log.Printf("Successfully assigned %d roles", len(assignments))

	// Signal refresh to update room information
	h.raiseChangeFlagHandler.Handle(context.Background(), roomCommand.RaiseChangeFlagCommand{})

	// Show assignments to the admin with a confirm button
	assignmentText := h.formatAssignments(assignments)
	confirmButton := telebot.InlineButton{
		Unique: UniqueConfirm,
		Text:   "Confirm and Send Roles to Players",
		Data:   gameID,
	}
	markup := &telebot.ReplyMarkup{InlineKeyboard: [][]telebot.InlineButton{{confirmButton}}}
	return c.Send(fmt.Sprintf("Role Assignments for Game '%s':\n\n%s", gameID, assignmentText), markup)
}

// formatAssignments formats the role assignments for display
func (h *BotHandler) formatAssignments(assignments map[userEntity.UserID]scenarioEntity.Role) string {
	var result string

	// Prepare a map of user IDs to users for easy lookup
	usersMap := make(map[userEntity.UserID]*userEntity.User)

	// Try to get all users from all rooms to have a better chance of finding user information
	allRooms, _ := h.getRoomsHandler.Handle(context.Background(), roomQuery.GetRoomsQuery{})
	for _, room := range allRooms {
		users, _ := h.getPlayersInRoomHandler.Handle(context.Background(),
			roomQuery.GetPlayersInRoomQuery{RoomID: room.ID})
		for _, user := range users {
			usersMap[user.ID] = user
		}
	}

	// Format each assignment
	for userID, role := range assignments {
		userIDInt := userID
		uid := userEntity.UserID(userIDInt)

		// Try to find the user in our map
		var userName string
		if user, found := usersMap[uid]; found {
			if user.Username != "" {
				userName = "@" + user.Username
			} else {
				userName = user.FirstName
				if user.LastName != "" {
					userName += " " + user.LastName
				}
			}
		}

		// Fall back to user ID if name not found
		if userName == "" {
			userName = fmt.Sprintf("User %d", userIDInt)
		}

		result += fmt.Sprintf("%s: %s\n", userName, role.Name)
	}
	return result
}

// HandleConfirmAssignments handles the 'confirm_assignments' callback
func (h *BotHandler) HandleConfirmAssignments(c telebot.Context, gameID string) error {
	if gameID == "" {
		return c.Respond(&telebot.CallbackResponse{
			Text:      "Error: Missing game ID",
			ShowAlert: true,
		})
	}

	// Log the gameID for debugging
	log.Printf("Confirming assignments for game: '%s'", gameID)

	// Get the game directly by ID
	targetGame, err := h.getGameByIDHandler.Handle(context.Background(),
		gameQuery.GetGameByIDQuery{ID: gameEntity.GameID(gameID)})
	if err != nil {
		log.Printf("Error fetching game with ID '%s': %v", gameID, err)
		return c.Respond(&telebot.CallbackResponse{
			Text:      fmt.Sprintf("Game with ID '%s' not found: %v", gameID, err),
			ShowAlert: true,
		})
	}

	roomID := string(targetGame.Room.ID)
	log.Printf("Found game '%s' for room '%s'", targetGame.ID, roomID)

	// Get the users in the room first to ensure we have their information
	users, err := h.getPlayersInRoomHandler.Handle(context.Background(),
		roomQuery.GetPlayersInRoomQuery{RoomID: entity.RoomID(roomID)})
	if err != nil {
		log.Printf("Error getting players in room '%s': %v", roomID, err)
	} else {
		log.Printf("Found %d players in room '%s'", len(users), roomID)
	}

	// Get user maps for better display
	usersMap := make(map[int64]*userEntity.User)
	for _, user := range users {
		usersMap[int64(user.ID)] = user
	}

	// Attempt to get assignments
	assignments, err := h.assignRolesHandler.GetAssignments(gameID)
	if err != nil {
		log.Printf("Error retrieving assignments for game '%s': %v", gameID, err)
		return c.Respond(&telebot.CallbackResponse{
			Text:      fmt.Sprintf("Error retrieving assignments: %v", err),
			ShowAlert: true,
		})
	}

	if len(assignments) == 0 {
		return c.Respond(&telebot.CallbackResponse{
			Text:      "No role assignments found for this game",
			ShowAlert: true,
		})
	}

	// Log the assignments
	log.Printf("Found %d role assignments for game '%s'", len(assignments), gameID)
	for userID, role := range assignments {
		log.Printf("Assignment: User %s - Role %s", userID, role.Name)
	}

	// Send role assignments to each user
	successCount := 0
	for userID, role := range assignments {
		userIDInt, err := strconv.ParseInt(userID, 10, 64)
		if err != nil {
			log.Printf("Invalid user ID '%s': %v", userID, err)
			continue // Skip invalid user IDs
		}

		userChat := &telebot.Chat{ID: userIDInt}

		// Get user name for logs
		var userName string
		if user, found := usersMap[userIDInt]; found {
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

		// Format a nicer message for the user
		message := fmt.Sprintf("🎭 *Your Role Assignment* 🎭\n\n"+
			"You have been assigned the role: *%s*\n\n"+
			"Keep your role secret and follow the game master's instructions!",
			role.Name)

		_, err = h.bot.Send(userChat, message, &telebot.SendOptions{
			ParseMode: telebot.ModeMarkdown,
		})

		if err != nil {
			log.Printf("Failed to send role to %s (ID: %d): %v", userName, userIDInt, err)
		} else {
			log.Printf("Successfully sent role to %s (ID: %d)", userName, userIDInt)
			successCount++
		}
	}

	return c.Respond(&telebot.CallbackResponse{
		Text:      fmt.Sprintf("Roles have been sent to %d players!", successCount),
		ShowAlert: true,
	})
}
