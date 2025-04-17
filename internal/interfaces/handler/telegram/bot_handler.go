package telegram

import (
	"log"
	"telemafia/internal/usecase"

	"gopkg.in/telebot.v3"
)

// BotHandler handles Telegram bot commands
type BotHandler struct {
	bot                     *telebot.Bot
	adminUsernames          []string
	roomRepo                usecase.RoomWriter
	createRoomHandler       *usecase.CreateRoomHandler
	joinRoomHandler         *usecase.JoinRoomHandler
	leaveRoomHandler        *usecase.LeaveRoomHandler
	kickUserHandler         *usecase.KickUserHandler
	deleteRoomHandler       *usecase.DeleteRoomHandler
	resetRefreshHandler     *usecase.ResetChangeFlagHandler
	raiseChangeFlagHandler  *usecase.RaiseChangeFlagHandler
	getRoomsHandler         *usecase.GetRoomsHandler
	getPlayerRoomsHandler   *usecase.GetPlayerRoomsHandler
	getPlayersInRoomHandler *usecase.GetPlayersInRoomHandler
	getRoomHandler          *usecase.GetRoomHandler
	checkRefreshHandler     *usecase.CheckChangeFlagHandler
	createScenarioHandler   *usecase.CreateScenarioHandler
	deleteScenarioHandler   *usecase.DeleteScenarioHandler
	manageRolesHandler      *usecase.ManageRolesHandler
	getScenarioByIDHandler  *usecase.GetScenarioByIDHandler
	getAllScenariosHandler  *usecase.GetAllScenariosHandler
	assignRolesHandler      *usecase.AssignRolesHandler
	createGameHandler       *usecase.CreateGameHandler
	getGamesHandler         *usecase.GetGamesHandler
	getGameByIDHandler      *usecase.GetGameByIDHandler
}

// NewBotHandler creates a new BotHandler
func NewBotHandler(
	bot *telebot.Bot,
	adminUsernames []string,
	roomRepo usecase.RoomWriter,
	createRoomHandler *usecase.CreateRoomHandler,
	joinRoomHandler *usecase.JoinRoomHandler,
	leaveRoomHandler *usecase.LeaveRoomHandler,
	kickUserHandler *usecase.KickUserHandler,
	deleteRoomHandler *usecase.DeleteRoomHandler,
	resetRefreshHandler *usecase.ResetChangeFlagHandler,
	raiseChangeFlagHandler *usecase.RaiseChangeFlagHandler,
	getRoomsHandler *usecase.GetRoomsHandler,
	getPlayerRoomsHandler *usecase.GetPlayerRoomsHandler,
	getPlayersInRoomHandler *usecase.GetPlayersInRoomHandler,
	getRoomHandler *usecase.GetRoomHandler,
	checkRefreshHandler *usecase.CheckChangeFlagHandler,
	createScenarioHandler *usecase.CreateScenarioHandler,
	deleteScenarioHandler *usecase.DeleteScenarioHandler,
	manageRolesHandler *usecase.ManageRolesHandler,
	getScenarioByIDHandler *usecase.GetScenarioByIDHandler,
	getAllScenariosHandler *usecase.GetAllScenariosHandler,
	assignRolesHandler *usecase.AssignRolesHandler,
	createGameHandler *usecase.CreateGameHandler,
	getGamesHandler *usecase.GetGamesHandler,
	getGameByIDHandler *usecase.GetGameByIDHandler,
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

func (h *BotHandler) Start() {
	// Start the background refresh goroutine
	go h.RefreshRoomsList()
	// Start the bot's main loop (blocking)
	log.Println("Starting bot polling...")
	h.bot.Start()
}

// RegisterHandlers registers all bot command handlers
func (h *BotHandler) RegisterHandlers() {
	// Register handlers for commands (mapping to methods in handlers.go)
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
	h.bot.Handle("/assign_scenario", h.HandleAssignScenario)
	h.bot.Handle("/assign_roles", h.HandleAssignRoles)
	h.bot.Handle("/games", h.HandleGamesList)

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
