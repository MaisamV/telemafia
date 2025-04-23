package telegram

import (
	"fmt"
	"log"

	// gameUsecase "telemafia/internal/game/usecase"
	gameCommand "telemafia/internal/domain/game/usecase/command"
	gameQuery "telemafia/internal/domain/game/usecase/query"

	// roomPort "telemafia/internal/room/port" // Import room port for RoomWriter interface
	roomPort "telemafia/internal/domain/room/port"
	// roomUsecase "telemafia/internal/room/usecase"
	roomCommand "telemafia/internal/domain/room/usecase/command"
	roomQuery "telemafia/internal/domain/room/usecase/query"

	// scenarioUsecase "telemafia/internal/scenario/usecase"
	scenarioCommand "telemafia/internal/domain/scenario/usecase/command"
	scenarioQuery "telemafia/internal/domain/scenario/usecase/query"

	"gopkg.in/telebot.v3"
)

// BotHandler holds dependencies and handles Telegram bot setup
type BotHandler struct {
	bot                     *telebot.Bot
	adminUsernames          []string
	roomRepo                roomPort.RoomWriter                    // Use roomPort
	createRoomHandler       *roomCommand.CreateRoomHandler         // Use roomCommand
	joinRoomHandler         *roomCommand.JoinRoomHandler           // Use roomCommand
	leaveRoomHandler        *roomCommand.LeaveRoomHandler          // Use roomCommand
	kickUserHandler         *roomCommand.KickUserHandler           // Use roomCommand
	deleteRoomHandler       *roomCommand.DeleteRoomHandler         // Use roomCommand
	resetRefreshHandler     *roomCommand.ResetChangeFlagHandler    // Use roomCommand
	raiseChangeFlagHandler  *roomCommand.RaiseChangeFlagHandler    // Use roomCommand
	getRoomsHandler         *roomQuery.GetRoomsHandler             // Use roomQuery
	getPlayerRoomsHandler   *roomQuery.GetPlayerRoomsHandler       // Use roomQuery
	getPlayersInRoomHandler *roomQuery.GetPlayersInRoomHandler     // Use roomQuery
	getRoomHandler          *roomQuery.GetRoomHandler              // Use roomQuery
	checkRefreshHandler     *roomQuery.CheckChangeFlagHandler      // Use roomQuery
	addDescriptionHandler   *roomCommand.AddDescriptionHandler     // Add handler field
	createScenarioHandler   *scenarioCommand.CreateScenarioHandler // Use scenarioCommand
	deleteScenarioHandler   *scenarioCommand.DeleteScenarioHandler // Use scenarioCommand
	manageRolesHandler      *scenarioCommand.ManageRolesHandler    // Use scenarioCommand
	getScenarioByIDHandler  *scenarioQuery.GetScenarioByIDHandler  // Use scenarioQuery
	getAllScenariosHandler  *scenarioQuery.GetAllScenariosHandler  // Use scenarioQuery
	assignRolesHandler      *gameCommand.AssignRolesHandler        // Use gameCommand
	createGameHandler       *gameCommand.CreateGameHandler         // Use gameCommand
	getGamesHandler         *gameQuery.GetGamesHandler             // Use gameQuery
	getGameByIDHandler      *gameQuery.GetGameByIDHandler          // Use gameQuery
}

// NewBotHandler creates a new BotHandler with all dependencies
func NewBotHandler(
	bot *telebot.Bot,
	adminUsernames []string,
	roomRepo roomPort.RoomWriter, // Use roomPort
	createRoomHandler *roomCommand.CreateRoomHandler, // Use roomCommand
	joinRoomHandler *roomCommand.JoinRoomHandler, // Use roomCommand
	leaveRoomHandler *roomCommand.LeaveRoomHandler, // Use roomCommand
	kickUserHandler *roomCommand.KickUserHandler, // Use roomCommand
	deleteRoomHandler *roomCommand.DeleteRoomHandler, // Use roomCommand
	resetRefreshHandler *roomCommand.ResetChangeFlagHandler, // Use roomCommand
	raiseChangeFlagHandler *roomCommand.RaiseChangeFlagHandler, // Use roomCommand
	getRoomsHandler *roomQuery.GetRoomsHandler, // Use roomQuery
	getPlayerRoomsHandler *roomQuery.GetPlayerRoomsHandler, // Use roomQuery
	getPlayersInRoomHandler *roomQuery.GetPlayersInRoomHandler, // Use roomQuery
	getRoomHandler *roomQuery.GetRoomHandler, // Use roomQuery
	checkRefreshHandler *roomQuery.CheckChangeFlagHandler, // Use roomQuery
	addDescriptionHandler *roomCommand.AddDescriptionHandler, // Add handler param
	createScenarioHandler *scenarioCommand.CreateScenarioHandler, // Use scenarioCommand
	deleteScenarioHandler *scenarioCommand.DeleteScenarioHandler, // Use scenarioCommand
	manageRolesHandler *scenarioCommand.ManageRolesHandler, // Use scenarioCommand
	getScenarioByIDHandler *scenarioQuery.GetScenarioByIDHandler, // Use scenarioQuery
	getAllScenariosHandler *scenarioQuery.GetAllScenariosHandler, // Use scenarioQuery
	assignRolesHandler *gameCommand.AssignRolesHandler, // Use gameCommand
	createGameHandler *gameCommand.CreateGameHandler, // Use gameCommand
	getGamesHandler *gameQuery.GetGamesHandler, // Use gameQuery
	getGameByIDHandler *gameQuery.GetGameByIDHandler, // Use gameQuery
) *BotHandler {
	// Set admin users for util package (now local)
	SetAdminUsers(adminUsernames)

	return &BotHandler{
		bot:                     bot,
		adminUsernames:          adminUsernames,
		roomRepo:                roomRepo,
		createRoomHandler:       createRoomHandler,
		joinRoomHandler:         joinRoomHandler,
		leaveRoomHandler:        leaveRoomHandler,
		kickUserHandler:         kickUserHandler,
		deleteRoomHandler:       deleteRoomHandler,
		resetRefreshHandler:     resetRefreshHandler,
		raiseChangeFlagHandler:  raiseChangeFlagHandler,
		getRoomsHandler:         getRoomsHandler,
		getPlayerRoomsHandler:   getPlayerRoomsHandler,
		getPlayersInRoomHandler: getPlayersInRoomHandler,
		getRoomHandler:          getRoomHandler,
		checkRefreshHandler:     checkRefreshHandler,
		addDescriptionHandler:   addDescriptionHandler, // Assign handler
		createScenarioHandler:   createScenarioHandler,
		deleteScenarioHandler:   deleteScenarioHandler,
		manageRolesHandler:      manageRolesHandler,
		getScenarioByIDHandler:  getScenarioByIDHandler,
		getAllScenariosHandler:  getAllScenariosHandler,
		assignRolesHandler:      assignRolesHandler,
		createGameHandler:       createGameHandler,
		getGamesHandler:         getGamesHandler,
		getGameByIDHandler:      getGameByIDHandler,
	}
}

// Start initializes background tasks and starts the bot polling
func (h *BotHandler) Start() {
	// Start the background refresh goroutine if needed (logic might be in handlers.go)
	// go h.RefreshRoomsList() // Assuming this is handled elsewhere or removed
	// Start the bot's main loop (blocking)
	log.Println("Starting bot polling...")
	h.bot.Start()
}

// RegisterHandlers registers all bot command handlers
func (h *BotHandler) RegisterHandlers() {
	// Common Handlers
	h.bot.Handle("/start", h.handleStart)
	h.bot.Handle("/help", h.handleHelp)

	// Room Handlers
	h.bot.Handle("/create_room", h.handleCreateRoom)
	h.bot.Handle("/join_room", h.handleJoinRoom)
	h.bot.Handle("/leave_room", h.handleLeaveRoom)
	h.bot.Handle("/list_rooms", h.handleListRooms)
	h.bot.Handle("/my_rooms", h.handleMyRooms)
	h.bot.Handle("/kick_user", h.handleKickUser)
	h.bot.Handle("/delete_room", h.handleDeleteRoom)
	// AddDescription is not a direct command

	// Scenario Handlers
	h.bot.Handle("/create_scenario", h.handleCreateScenario)
	h.bot.Handle("/delete_scenario", h.handleDeleteScenario)
	h.bot.Handle("/add_role", h.handleAddRole)
	h.bot.Handle("/remove_role", h.handleRemoveRole)
	// TODO: Add /list_scenarios handler

	// Game Handlers
	h.bot.Handle("/assign_scenario", h.handleAssignScenario) // Assigns scenario AND creates game
	h.bot.Handle("/assign_roles", h.handleAssignRoles)
	h.bot.Handle("/games", h.handleGamesList)

	// Register handler for callback queries
	h.bot.Handle(telebot.OnCallback, h.handleCallback)

	log.Println("Registered command and callback handlers.")
}

// --- Dispatcher Methods ---

// --- Common ---
func (h *BotHandler) handleStart(c telebot.Context) error {
	return HandleStart(h, c)
}

func (h *BotHandler) handleHelp(c telebot.Context) error {
	return HandleHelp(h, c)
}

// --- Room ---
func (h *BotHandler) handleCreateRoom(c telebot.Context) error {
	return HandleCreateRoom(h, c)
}

func (h *BotHandler) handleJoinRoom(c telebot.Context) error {
	return HandleJoinRoom(h, c)
}

func (h *BotHandler) handleLeaveRoom(c telebot.Context) error {
	return HandleLeaveRoom(h, c)
}

func (h *BotHandler) handleListRooms(c telebot.Context) error {
	return HandleListRooms(h, c)
}

func (h *BotHandler) handleMyRooms(c telebot.Context) error {
	return HandleMyRooms(h, c)
}

func (h *BotHandler) handleKickUser(c telebot.Context) error {
	return HandleKickUser(h, c)
}

func (h *BotHandler) handleDeleteRoom(c telebot.Context) error {
	return HandleDeleteRoom(h, c)
}

// --- Scenario ---
func (h *BotHandler) handleCreateScenario(c telebot.Context) error {
	return HandleCreateScenario(h, c)
}

func (h *BotHandler) handleDeleteScenario(c telebot.Context) error {
	return HandleDeleteScenario(h, c)
}

func (h *BotHandler) handleAddRole(c telebot.Context) error {
	return HandleAddRole(h, c)
}

func (h *BotHandler) handleRemoveRole(c telebot.Context) error {
	return HandleRemoveRole(h, c)
}

// --- Game ---
func (h *BotHandler) handleAssignScenario(c telebot.Context) error {
	return HandleAssignScenario(h, c)
}

func (h *BotHandler) handleAssignRoles(c telebot.Context) error {
	return HandleAssignRoles(h, c)
}

func (h *BotHandler) handleGamesList(c telebot.Context) error {
	return HandleGamesList(h, c)
}

// --- Callbacks ---
// Removed handleCallback dispatcher method - implementation is in callbacks.go
// func (h *BotHandler) handleCallback(c telebot.Context) error {
// 	return HandleCallback(h, c) // Assuming HandleCallback is now a function in callbacks.go
// }

// --- Internal Helper Handlers (originally part of BotHandler) ---

// HandleHelp provides a simple help message.
func (h *BotHandler) HandleHelp(c telebot.Context) error {
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
	_ = c.Send(fmt.Sprintf("Welcome, %s!", c.Sender().Username))
	return h.SendOrUpdateRefreshingMessage(c.Sender().ID, ListRooms, "")
}
