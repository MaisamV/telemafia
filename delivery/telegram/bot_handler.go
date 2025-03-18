package telegram

import (
	"log"
	"os"
	gameCommand "telemafia/internal/game/usecase/command"
	gameQuery "telemafia/internal/game/usecase/query"
	roomCommand "telemafia/internal/room/usecase/command"
	roomQuery "telemafia/internal/room/usecase/query"
	scenarioCommand "telemafia/internal/scenario/usecase/command"
	scenarioQuery "telemafia/internal/scenario/usecase/query"

	"gopkg.in/telebot.v3"
)

// BotHandler handles Telegram bot commands
type BotHandler struct {
	bot                     *telebot.Bot
	adminUsernames          []string
	createRoomHandler       *roomCommand.CreateRoomHandler
	joinRoomHandler         *roomCommand.JoinRoomHandler
	leaveRoomHandler        *roomCommand.LeaveRoomHandler
	kickUserHandler         *roomCommand.KickUserHandler
	deleteRoomHandler       *roomCommand.DeleteRoomHandler
	resetRefreshHandler     *roomCommand.ResetChangeFlagHandler
	raiseChangeFlagHandler  *roomCommand.RaiseChangeFlagHandler
	getRoomsHandler         *roomQuery.GetRoomsHandler
	getPlayerRoomsHandler   *roomQuery.GetPlayerRoomsHandler
	getPlayersInRoomHandler *roomQuery.GetPlayersInRoomHandler
	getRoomHandler          *roomQuery.GetRoomHandler
	checkRefreshHandler     *roomQuery.CheckChangeFlagHandler
	createScenarioHandler   *scenarioCommand.CreateScenarioHandler
	deleteScenarioHandler   *scenarioCommand.DeleteScenarioHandler
	manageRolesHandler      *scenarioCommand.ManageRolesHandler
	getScenarioByIDHandler  *scenarioQuery.GetScenarioByIDHandler
	getAllScenariosHandler  *scenarioQuery.GetAllScenariosHandler
	assignRolesHandler      *gameCommand.AssignRolesHandler
	assignScenarioHandler   *roomCommand.AssignScenarioHandler
	createGameHandler       *gameCommand.CreateGameHandler
	getGamesHandler         *gameQuery.GetGamesHandler
	getGameByIDHandler      *gameQuery.GetGameByIDHandler
}

// NewBotHandler creates a new BotHandler
func NewBotHandler(
	bot *telebot.Bot,
	adminUsernames []string,
	createRoomHandler *roomCommand.CreateRoomHandler,
	joinRoomHandler *roomCommand.JoinRoomHandler,
	leaveRoomHandler *roomCommand.LeaveRoomHandler,
	kickUserHandler *roomCommand.KickUserHandler,
	deleteRoomHandler *roomCommand.DeleteRoomHandler,
	resetRefreshHandler *roomCommand.ResetChangeFlagHandler,
	raiseChangeFlagHandler *roomCommand.RaiseChangeFlagHandler,
	getRoomsHandler *roomQuery.GetRoomsHandler,
	getPlayerRoomsHandler *roomQuery.GetPlayerRoomsHandler,
	getPlayersInRoomHandler *roomQuery.GetPlayersInRoomHandler,
	getRoomHandler *roomQuery.GetRoomHandler,
	checkRefreshHandler *roomQuery.CheckChangeFlagHandler,
	createScenarioHandler *scenarioCommand.CreateScenarioHandler,
	deleteScenarioHandler *scenarioCommand.DeleteScenarioHandler,
	manageRolesHandler *scenarioCommand.ManageRolesHandler,
	getScenarioByIDHandler *scenarioQuery.GetScenarioByIDHandler,
	getAllScenariosHandler *scenarioQuery.GetAllScenariosHandler,
	assignRolesHandler *gameCommand.AssignRolesHandler,
	assignScenarioHandler *roomCommand.AssignScenarioHandler,
	createGameHandler *gameCommand.CreateGameHandler,
	getGamesHandler *gameQuery.GetGamesHandler,
	getGameByIDHandler *gameQuery.GetGameByIDHandler,
) *BotHandler {
	return &BotHandler{
		bot:                     bot,
		adminUsernames:          adminUsernames,
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
		assignScenarioHandler:   assignScenarioHandler,
		createGameHandler:       createGameHandler,
		getGamesHandler:         getGamesHandler,
		getGameByIDHandler:      getGameByIDHandler,
	}
}

func (h *BotHandler) Start() {
	logger := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
	logger.Println("Test")
	go h.RefreshRoomsList() // Start the goroutine to refresh room lists
	h.bot.Start()
}

// RegisterHandlers registers all bot command handlers
func (h *BotHandler) RegisterHandlers() {
	h.bot.SetCommands([]telebot.Command{
		{Text: "empty", Description: ""},
	})
	h.bot.SetCommands([]telebot.Command{
		{Text: "list_rooms", Description: "Show mafia rooms"},
	})
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
	h.bot.Handle("/assign_roles", h.HandleAssignRoles)
	h.bot.Handle("/assign_scenario", h.HandleAssignScenario)
	h.bot.Handle("/games", h.HandleGamesList)
	h.bot.Handle(telebot.OnCallback, h.HandleCallback)
}

// HandleGamesList handles the /games command
func (h *BotHandler) HandleGamesList(c telebot.Context) error {
	// Create a games list handler
	gamesListHandler := NewGamesListHandler(h.bot, h.getGamesHandler)

	// Forward the call to the games list handler
	return gamesListHandler.HandleGamesList(c)
}
