package telegram

import (
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
	roomRepo                roomPort.RoomWriter                 // Use roomPort
	createRoomHandler       *roomCommand.CreateRoomHandler      // Use roomCommand
	joinRoomHandler         *roomCommand.JoinRoomHandler        // Use roomCommand
	leaveRoomHandler        *roomCommand.LeaveRoomHandler       // Use roomCommand
	kickUserHandler         *roomCommand.KickUserHandler        // Use roomCommand
	deleteRoomHandler       *roomCommand.DeleteRoomHandler      // Use roomCommand
	resetRefreshHandler     *roomCommand.ResetChangeFlagHandler // Use roomCommand
	raiseChangeFlagHandler  *roomCommand.RaiseChangeFlagHandler // Use roomCommand
	getRoomsHandler         *roomQuery.GetRoomsHandler          // Use roomQuery
	getPlayerRoomsHandler   *roomQuery.GetPlayerRoomsHandler    // Use roomQuery
	getPlayersInRoomHandler *roomQuery.GetPlayersInRoomHandler  // Use roomQuery
	getRoomHandler          *roomQuery.GetRoomHandler           // Use roomQuery
	checkRefreshHandler     *roomQuery.CheckChangeFlagHandler   // Use roomQuery
	// addDescriptionHandler   *roomCommand.AddDescriptionHandler   // Use roomCommand (if wiring this)
	createScenarioHandler  *scenarioCommand.CreateScenarioHandler // Use scenarioCommand
	deleteScenarioHandler  *scenarioCommand.DeleteScenarioHandler // Use scenarioCommand
	manageRolesHandler     *scenarioCommand.ManageRolesHandler    // Use scenarioCommand
	getScenarioByIDHandler *scenarioQuery.GetScenarioByIDHandler  // Use scenarioQuery
	getAllScenariosHandler *scenarioQuery.GetAllScenariosHandler  // Use scenarioQuery
	assignRolesHandler     *gameCommand.AssignRolesHandler        // Use gameCommand
	createGameHandler      *gameCommand.CreateGameHandler         // Use gameCommand
	getGamesHandler        *gameQuery.GetGamesHandler             // Use gameQuery
	getGameByIDHandler     *gameQuery.GetGameByIDHandler          // Use gameQuery
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
	// addDescriptionHandler *roomCommand.AddDescriptionHandler, // Use roomCommand (if wiring this)
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
		// addDescriptionHandler:   addDescriptionHandler, // If wiring this
		createScenarioHandler:  createScenarioHandler,
		deleteScenarioHandler:  deleteScenarioHandler,
		manageRolesHandler:     manageRolesHandler,
		getScenarioByIDHandler: getScenarioByIDHandler,
		getAllScenariosHandler: getAllScenariosHandler,
		assignRolesHandler:     assignRolesHandler,
		createGameHandler:      createGameHandler,
		getGamesHandler:        getGamesHandler,
		getGameByIDHandler:     getGameByIDHandler,
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
	// Register handlers for commands (mapping to methods in handlers.go)
	// These methods now use the injected handlers.
	h.bot.Handle("/start", h.HandleStart)
	h.bot.Handle("/help", h.HandleHelp)
	h.bot.Handle("/create_room", h.HandleCreateRoom)
	h.bot.Handle("/join_room", h.HandleJoinRoom)
	h.bot.Handle("/leave_room", h.HandleLeaveRoom)
	h.bot.Handle("/list_rooms", h.HandleListRooms)
	h.bot.Handle("/my_rooms", h.HandleMyRooms)
	h.bot.Handle("/kick_user", h.HandleKickUser)
	h.bot.Handle("/delete_room", h.HandleDeleteRoom)
	h.bot.Handle("/create_scenario", h.HandleCreateScenario)
	h.bot.Handle("/delete_scenario", h.HandleDeleteScenario)
	h.bot.Handle("/add_role", h.HandleAddRole)
	h.bot.Handle("/remove_role", h.HandleRemoveRole)
	h.bot.Handle("/assign_scenario", h.HandleAssignScenario) // This likely needs adjustment
	h.bot.Handle("/assign_roles", h.HandleAssignRoles)
	h.bot.Handle("/games", h.HandleGamesList) // This likely needs adjustment
	// TODO: /list_scenarios command mentioned in README needs implementation/wiring

	// Register handler for callback queries
	h.bot.Handle(telebot.OnCallback, h.HandleCallback)

	log.Println("Registered command and callback handlers.")
}

// REMOVED HandleGamesList method (logic moved to handlers.go)
// func (h *BotHandler) HandleGamesList(c telebot.Context) error {
// 	// Create a games list handler
// 	gamesListHandler := NewGamesListHandler(h.bot, h.getGamesHandler)
//
// 	// Forward the call to the games list handler
// 	return gamesListHandler.HandleGamesList(c)
// }
