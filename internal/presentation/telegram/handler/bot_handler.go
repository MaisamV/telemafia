package telegram

import (
	"log"
	"telemafia/internal/shared/tgutil"

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

	game "telemafia/internal/presentation/telegram/handler/game"
	room "telemafia/internal/presentation/telegram/handler/room"
	scenario "telemafia/internal/presentation/telegram/handler/scenario"
	messages "telemafia/internal/presentation/telegram/messages" // Import messages package

	"gopkg.in/telebot.v3"
)

// BotHandler holds dependencies and handles Telegram bot setup
type BotHandler struct {
	bot            *telebot.Bot
	msgs           *messages.Messages // Add messages field
	adminUsernames []string

	// Refresh state management (delegated)
	roomListRefreshMessage   *tgutil.RefreshingMessageBook
	roomDetailRefreshMessage *tgutil.RefreshingMessageBook

	// // Refresh state (moved from repository) - REMOVED
	// refreshMutex            sync.RWMutex
	// needsRefresh            bool
	// activeRefreshMessages   map[int64]*telebot.Message // Map ChatID to the message being refreshed

	// Use Case Handlers
	roomRepo                roomPort.RoomWriter                    // Use roomPort
	createRoomHandler       *roomCommand.CreateRoomHandler         // Use roomCommand
	joinRoomHandler         *roomCommand.JoinRoomHandler           // Use roomCommand
	leaveRoomHandler        *roomCommand.LeaveRoomHandler          // Use roomCommand
	kickUserHandler         *roomCommand.KickUserHandler           // Use roomCommand
	deleteRoomHandler       *roomCommand.DeleteRoomHandler         // Use roomCommand
	getRoomsHandler         *roomQuery.GetRoomsHandler             // Use roomQuery
	getPlayerRoomsHandler   *roomQuery.GetPlayerRoomsHandler       // Use roomQuery
	getPlayersInRoomHandler *roomQuery.GetPlayersInRoomHandler     // Use roomQuery
	getRoomHandler          *roomQuery.GetRoomHandler              // Use roomQuery
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
	msgs *messages.Messages, // Add messages parameter
	roomRepo roomPort.RoomWriter, // Use roomPort
	createRoomHandler *roomCommand.CreateRoomHandler, // Use roomCommand
	joinRoomHandler *roomCommand.JoinRoomHandler, // Use roomCommand
	leaveRoomHandler *roomCommand.LeaveRoomHandler, // Use roomCommand
	kickUserHandler *roomCommand.KickUserHandler, // Use roomCommand
	deleteRoomHandler *roomCommand.DeleteRoomHandler, // Use roomCommand
	getRoomsHandler *roomQuery.GetRoomsHandler, // Use roomQuery
	getPlayerRoomsHandler *roomQuery.GetPlayerRoomsHandler, // Use roomQuery
	getPlayersInRoomHandler *roomQuery.GetPlayersInRoomHandler, // Use roomQuery
	getRoomHandler *roomQuery.GetRoomHandler, // Use roomQuery
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
	// Set admin users for util package (now moved)
	tgutil.SetAdminUsers(adminUsernames) // Use tgutil

	h := &BotHandler{
		bot:                      bot,
		msgs:                     msgs, // Assign messages
		adminUsernames:           adminUsernames,
		roomListRefreshMessage:   tgutil.NewRefreshState(), // Initialize RefreshingMessageBook
		roomDetailRefreshMessage: tgutil.NewRefreshState(), // Initialize RefreshingMessageBook
		roomRepo:                 roomRepo,
		createRoomHandler:        createRoomHandler,
		joinRoomHandler:          joinRoomHandler,
		leaveRoomHandler:         leaveRoomHandler,
		kickUserHandler:          kickUserHandler,
		deleteRoomHandler:        deleteRoomHandler,
		getRoomsHandler:          getRoomsHandler,
		getPlayerRoomsHandler:    getPlayerRoomsHandler,
		getPlayersInRoomHandler:  getPlayersInRoomHandler,
		getRoomHandler:           getRoomHandler,
		addDescriptionHandler:    addDescriptionHandler, // Assign handler
		createScenarioHandler:    createScenarioHandler,
		deleteScenarioHandler:    deleteScenarioHandler,
		manageRolesHandler:       manageRolesHandler,
		getScenarioByIDHandler:   getScenarioByIDHandler,
		getAllScenariosHandler:   getAllScenariosHandler,
		assignRolesHandler:       assignRolesHandler,
		createGameHandler:        createGameHandler,
		getGamesHandler:          getGamesHandler,
		getGameByIDHandler:       getGameByIDHandler,
	}
	return h
}

// Start initializes background tasks and starts the bot polling
func (h *BotHandler) Start() {
	// Start the background refresh goroutine if needed (logic might be in handlers.go)
	// go h.StartRefreshTimer() // Assuming this is handled elsewhere or removed
	// Start the bot's main loop (blocking)
	log.Println("Starting bot polling...")
	go h.StartRefreshTimer()
	h.bot.Start()
}

// RegisterHandlers registers all bot command handlers
func (h *BotHandler) RegisterHandlers() {
	// Common Handlers
	//h.bot.Handle(telebot.OnText, h.handleStart)
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
	return HandleStart(h, c, h.msgs)
}

func (h *BotHandler) handleHelp(c telebot.Context) error {
	return HandleHelp(h, c, h.msgs)
}

// --- Room ---
func (h *BotHandler) handleCreateRoom(c telebot.Context) error {
	return room.HandleCreateRoom(h.createRoomHandler, h.roomListRefreshMessage, c, h.msgs)
}

func (h *BotHandler) handleJoinRoom(c telebot.Context) error {
	return room.HandleJoinRoom(h.joinRoomHandler, h.roomListRefreshMessage, h.roomDetailRefreshMessage, c, h.msgs)
}

func (h *BotHandler) handleLeaveRoom(c telebot.Context) error {
	return room.HandleLeaveRoom(h.leaveRoomHandler, h.roomListRefreshMessage, c, h.msgs)
}

func (h *BotHandler) handleListRooms(c telebot.Context) error {
	return room.HandleListRooms(h.getRoomsHandler, h.getPlayersInRoomHandler, h.bot, h.roomListRefreshMessage, c, h.msgs)
}

func (h *BotHandler) handleMyRooms(c telebot.Context) error {
	return room.HandleMyRooms(h.getPlayerRoomsHandler, c, h.msgs)
}

func (h *BotHandler) handleKickUser(c telebot.Context) error {
	return room.HandleKickUser(h.kickUserHandler, h.roomListRefreshMessage, c, h.msgs)
}

func (h *BotHandler) handleDeleteRoom(c telebot.Context) error {
	// Showing the list doesn't need the notifier, but the confirm callback will.
	return room.HandleDeleteRoom(h.getRoomsHandler, c, h.msgs)
}

// --- Scenario ---
func (h *BotHandler) handleCreateScenario(c telebot.Context) error {
	return scenario.HandleCreateScenario(h.createScenarioHandler, c, h.msgs)
}

func (h *BotHandler) handleDeleteScenario(c telebot.Context) error {
	return scenario.HandleDeleteScenario(h.deleteScenarioHandler, c, h.msgs)
}

func (h *BotHandler) handleAddRole(c telebot.Context) error {
	return scenario.HandleAddRole(h.manageRolesHandler, c, h.msgs)
}

func (h *BotHandler) handleRemoveRole(c telebot.Context) error {
	return scenario.HandleRemoveRole(h.manageRolesHandler, c, h.msgs)
}

// --- Game ---
func (h *BotHandler) handleAssignScenario(c telebot.Context) error {
	return game.HandleAssignScenario(h.getRoomHandler, h.getScenarioByIDHandler, h.addDescriptionHandler, h.createGameHandler, c, h.msgs)
}

func (h *BotHandler) handleAssignRoles(c telebot.Context) error {
	return game.HandleAssignRoles(h.assignRolesHandler, h.bot, c, h.msgs)
}

func (h *BotHandler) handleGamesList(c telebot.Context) error {
	return game.HandleGamesList(h.getGamesHandler, c, h.msgs)
}

// --- Callbacks ---
// Removed handleCallback dispatcher method - implementation is in callbacks.go
// func (h *BotHandler) handleCallback(c telebot.Context) error {
// 	return HandleCallback(h, c) // Assuming HandleCallback is now a function in callbacks.go
// }

// --- Internal Helper Handlers (originally part of BotHandler) ---

// HandleStart handles the /start command
func (h *BotHandler) HandleStart(c telebot.Context) error {
	return HandleStart(h, c, h.msgs)
}
