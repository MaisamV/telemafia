package telegram

import (
	"context"
	"log"
	"strings"
	"sync"
	"telemafia/internal/shared/entity"
	"telemafia/internal/shared/tgutil"

	// gameUsecase "telemafia/internal/game/usecase"
	gameEntity "telemafia/internal/domain/game/entity"
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

// InteractiveSelectionState moved to tgutil package

// BotHandler holds dependencies and handles Telegram bot setup
type BotHandler struct {
	bot            *telebot.Bot
	msgs           *messages.Messages // Add messages field
	adminUsernames []string

	// Refresh state management (delegated)
	roomListRefreshMessage   *tgutil.RefreshingMessageBook
	roomDetailRefreshMessage *tgutil.RefreshingMessageBook

	// State for "Choose Card" interactive role selection
	interactiveSelectionsMutex sync.RWMutex                                            // Mutex for the outer map
	interactiveSelections      map[gameEntity.GameID]*tgutil.InteractiveSelectionState // Use tgutil type

	// Refresh books for interactive role selection
	playerRefreshMutex         sync.RWMutex // Mutex for player refreshers map
	playerRoleChoiceRefreshers map[gameEntity.GameID]*tgutil.RefreshingMessageBook
	adminRefreshMutex          sync.RWMutex // Mutex for admin refreshers map
	adminAssignmentTrackers    map[gameEntity.GameID]*tgutil.RefreshingMessageBook

	// // Refresh state (moved from repository) - REMOVED
	// refreshMutex            sync.RWMutex
	// needsRefresh            bool
	// activeRefreshMessages   map[int64]*telebot.Message // Map ChatID to the message being refreshed

	// Use Case Handlers
	roomRepo                roomPort.RoomWriter                     // Use roomPort
	createRoomHandler       *roomCommand.CreateRoomHandler          // Use roomCommand
	joinRoomHandler         *roomCommand.JoinRoomHandler            // Use roomCommand
	leaveRoomHandler        *roomCommand.LeaveRoomHandler           // Use roomCommand
	kickUserHandler         *roomCommand.KickUserHandler            // Use roomCommand
	deleteRoomHandler       *roomCommand.DeleteRoomHandler          // Use roomCommand
	getRoomsHandler         *roomQuery.GetRoomsHandler              // Use roomQuery
	getPlayerRoomsHandler   *roomQuery.GetPlayerRoomsHandler        // Use roomQuery
	getPlayersInRoomHandler *roomQuery.GetPlayersInRoomHandler      // Use roomQuery
	getRoomHandler          *roomQuery.GetRoomHandler               // Use roomQuery
	addDescriptionHandler   *roomCommand.AddDescriptionHandler      // Add handler field
	changeModeratorHandler  *roomCommand.ChangeModeratorHandler     // Add ChangeModeratorHandler field
	createScenarioHandler   *scenarioCommand.CreateScenarioHandler  // Use scenarioCommand
	deleteScenarioHandler   *scenarioCommand.DeleteScenarioHandler  // Use scenarioCommand
	getScenarioByIDHandler  *scenarioQuery.GetScenarioByIDHandler   // Use scenarioQuery
	getAllScenariosHandler  *scenarioQuery.GetAllScenariosHandler   // Use scenarioQuery
	addScenarioJSONHandler  *scenarioCommand.AddScenarioJSONHandler // NEW: Inject AddScenarioJSONHandler
	assignRolesHandler      *gameCommand.AssignRolesHandler         // Use gameCommand
	createGameHandler       *gameCommand.CreateGameHandler          // Use gameCommand
	updateGameHandler       *gameCommand.UpdateGameHandler          // ADDED: Update Game Handler
	getGamesHandler         *gameQuery.GetGamesHandler              // Use gameQuery
	getGameByIDHandler      *gameQuery.GetGameByIDHandler           // Use gameQuery
}

// --- Methods implementing BotHandlerInterface --- (NEW)

func (h *BotHandler) GetGameByIDHandler() *gameQuery.GetGameByIDHandler {
	return h.getGameByIDHandler
}

func (h *BotHandler) GetPlayersInRoomHandler() *roomQuery.GetPlayersInRoomHandler {
	return h.getPlayersInRoomHandler
}

func (h *BotHandler) GetRoomsHandler() *roomQuery.GetRoomsHandler {
	return h.getRoomsHandler
}

func (h *BotHandler) GetScenarioByIDHandler() *scenarioQuery.GetScenarioByIDHandler {
	return h.getScenarioByIDHandler
}

func (h *BotHandler) AssignRolesHandler() *gameCommand.AssignRolesHandler {
	return h.assignRolesHandler
}

func (h *BotHandler) Bot() *telebot.Bot {
	return h.bot
}

func (h *BotHandler) UpdateGameHandler() *gameCommand.UpdateGameHandler {
	return h.updateGameHandler
}

// --- End Interface Methods ---

// --- Refresh Book Management for Game Role Selection ---

// --- End Refresh Book Management ---

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
	changeModeratorHandler *roomCommand.ChangeModeratorHandler, // Add ChangeModeratorHandler param
	createScenarioHandler *scenarioCommand.CreateScenarioHandler, // Use scenarioCommand
	deleteScenarioHandler *scenarioCommand.DeleteScenarioHandler, // Use scenarioCommand
	getScenarioByIDHandler *scenarioQuery.GetScenarioByIDHandler, // Use scenarioQuery
	getAllScenariosHandler *scenarioQuery.GetAllScenariosHandler, // Use scenarioQuery
	addScenarioJSONHandler *scenarioCommand.AddScenarioJSONHandler, // NEW: Inject AddScenarioJSONHandler
	assignRolesHandler *gameCommand.AssignRolesHandler, // Use gameCommand
	createGameHandler *gameCommand.CreateGameHandler, // Use gameCommand
	updateGameHandler *gameCommand.UpdateGameHandler, // ADDED Parameter
	getGamesHandler *gameQuery.GetGamesHandler, // Use gameQuery
	getGameByIDHandler *gameQuery.GetGameByIDHandler, // Use gameQuery
) *BotHandler {
	// Set admin users for util package (now moved)
	tgutil.SetAdminUsers(adminUsernames)

	h := &BotHandler{
		bot:            bot,
		msgs:           msgs,
		adminUsernames: adminUsernames,
		roomListRefreshMessage: tgutil.NewRefreshState(func(user int64, data string) (string, []interface{}, error) {
			message, markup, err := room.PrepareRoomListMessage(
				getRoomsHandler,
				getPlayersInRoomHandler,
				msgs,
			)
			opts := []interface{}{
				markup,
				telebot.NoPreview,
			}
			return message, opts, err
		}),
		roomDetailRefreshMessage: tgutil.NewRefreshState(func(user int64, data string) (string, []interface{}, error) {
			return room.RoomDetailMessage(
				getRoomsHandler,
				getPlayersInRoomHandler,
				msgs,
				entity.UserID(user),
				data,
			)
		}),
		interactiveSelections:      make(map[gameEntity.GameID]*tgutil.InteractiveSelectionState), // Use tgutil type
		playerRoleChoiceRefreshers: make(map[gameEntity.GameID]*tgutil.RefreshingMessageBook),
		adminAssignmentTrackers:    make(map[gameEntity.GameID]*tgutil.RefreshingMessageBook),
		roomRepo:                   roomRepo,
		createRoomHandler:          createRoomHandler,
		joinRoomHandler:            joinRoomHandler,
		leaveRoomHandler:           leaveRoomHandler,
		kickUserHandler:            kickUserHandler,
		deleteRoomHandler:          deleteRoomHandler,
		getRoomsHandler:            getRoomsHandler,
		getPlayerRoomsHandler:      getPlayerRoomsHandler,
		getPlayersInRoomHandler:    getPlayersInRoomHandler,
		getRoomHandler:             getRoomHandler,
		addDescriptionHandler:      addDescriptionHandler,
		changeModeratorHandler:     changeModeratorHandler,
		createScenarioHandler:      createScenarioHandler,
		deleteScenarioHandler:      deleteScenarioHandler,
		getScenarioByIDHandler:     getScenarioByIDHandler,
		getAllScenariosHandler:     getAllScenariosHandler,
		addScenarioJSONHandler:     addScenarioJSONHandler,
		assignRolesHandler:         assignRolesHandler,
		createGameHandler:          createGameHandler,
		updateGameHandler:          updateGameHandler, // ADDED Assignment
		getGamesHandler:            getGamesHandler,
		getGameByIDHandler:         getGameByIDHandler,
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
	h.bot.Handle("/add_scenario_json", h.handleAddScenarioJSON) // NEW: Register command
	// TODO: Add /list_scenarios handler

	// Game Handlers
	h.bot.Handle("/create_game", h.handleCreateGame) // Renamed from /assign_scenario
	h.bot.Handle("/assign_roles", h.handleAssignRoles)
	h.bot.Handle("/games", h.handleGamesList)

	// Register handler for callback queries
	h.bot.Handle(telebot.OnCallback, h.handleCallback)
	h.bot.Handle(telebot.OnDocument, h.handleDocument)

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
	roomIDStr := strings.TrimSpace(c.Message().Payload)
	return room.HandleJoinRoom(h.joinRoomHandler, h.getRoomsHandler, h.getPlayersInRoomHandler, h.roomListRefreshMessage, h.roomDetailRefreshMessage, c, roomIDStr, h.msgs)
}

func (h *BotHandler) handleLeaveRoom(c telebot.Context) error {
	return room.HandleLeaveRoom(h.leaveRoomHandler, h.roomListRefreshMessage, c, h.msgs)
}

func (h *BotHandler) handleListRooms(c telebot.Context) error {
	return room.HandleListRooms(h.getRoomsHandler, h.getPlayersInRoomHandler, h.bot, h.roomListRefreshMessage, h.roomDetailRefreshMessage, c, h.msgs)
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

// NEW: Dispatcher method for Add Scenario JSON
func (h *BotHandler) handleAddScenarioJSON(c telebot.Context) error {
	return scenario.HandleAddScenarioJSON(h.addScenarioJSONHandler, c, h.msgs)
}

// --- Game ---
func (h *BotHandler) handleCreateGame(c telebot.Context) error { // Renamed from handleAssignScenario
	return game.HandleCreateGame(h.getRoomsHandler, c, h.msgs)
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

func (h *BotHandler) handleDocument(c telebot.Context) error {
	return HandleDocument(h.addScenarioJSONHandler, c, h.msgs)
}

// Helper methods to manage interactive state safely (NEW)
func (h *BotHandler) GetInteractiveSelectionState(gameID gameEntity.GameID) (*tgutil.InteractiveSelectionState, bool) { // Use tgutil type
	h.interactiveSelectionsMutex.RLock()
	defer h.interactiveSelectionsMutex.RUnlock()
	state, exists := h.interactiveSelections[gameID]
	return state, exists
}

func (h *BotHandler) SetInteractiveSelectionState(gameID gameEntity.GameID, state *tgutil.InteractiveSelectionState) { // Use tgutil type
	h.interactiveSelectionsMutex.Lock()
	defer h.interactiveSelectionsMutex.Unlock()
	h.interactiveSelections[gameID] = state
}

func (h *BotHandler) DeleteInteractiveSelectionState(gameID gameEntity.GameID) {
	h.interactiveSelectionsMutex.Lock()
	defer h.interactiveSelectionsMutex.Unlock()
	delete(h.interactiveSelections, gameID)
}

// GetPlayerRoleRefresher Helper methods to manage player role choice refreshers safely (NEW)
func (h *BotHandler) GetPlayerRoleRefresher(gameID gameEntity.GameID) (*tgutil.RefreshingMessageBook, bool) {
	h.playerRefreshMutex.RLock()
	defer h.playerRefreshMutex.RUnlock()
	book, exists := h.playerRoleChoiceRefreshers[gameID]
	return book, exists
}

// RemovePlayerRoleActiveMessage Helper methods to remove player role active message from game choice refreshers safely
func (h *BotHandler) RemovePlayerRoleActiveMessage(gameID gameEntity.GameID, userID int64) {
	h.playerRefreshMutex.RLock()
	defer h.playerRefreshMutex.RUnlock()
	book, exists := h.playerRoleChoiceRefreshers[gameID]
	if exists {
		book.RemoveActiveMessage(userID)
	}
}

func (h *BotHandler) GetOrCreatePlayerRoleRefresher(gameID gameEntity.GameID) *tgutil.RefreshingMessageBook {
	h.playerRefreshMutex.Lock()
	defer h.playerRefreshMutex.Unlock()
	book, exists := h.playerRoleChoiceRefreshers[gameID]
	if !exists {
		book = tgutil.NewRefreshState(func(user int64, data string) (string, []interface{}, error) {
			gameIDFromData := gameEntity.GameID(data)
			// Fetch necessary data (State)
			state, stateExists := h.GetInteractiveSelectionState(gameIDFromData)
			if !stateExists {
				log.Printf("Refresh Player: State for game %s not found.", gameIDFromData)
				return "Role selection is no longer active.", []interface{}{}, nil // Return inactive state message
			}

			// Prepare markup using the helper from callbacks_game
			markup, err := game.PreparePlayerRoleSelectionMarkup(gameIDFromData, len(state.ShuffledRoles), state.TakenIndices, h.msgs)
			message := h.msgs.Game.RoleSelectionPromptPlayer // Keep the prompt same, just update buttons
			opts := []interface{}{
				markup,
			}
			return message, opts, err
		})
		h.playerRoleChoiceRefreshers[gameID] = book
		log.Printf("Created new Player Role Refresher book for game %s", gameID)
	}
	return book
}

func (h *BotHandler) DeletePlayerRoleRefresher(gameID gameEntity.GameID) {
	h.playerRefreshMutex.Lock()
	defer h.playerRefreshMutex.Unlock()
	delete(h.playerRoleChoiceRefreshers, gameID)
	log.Printf("Deleted Player Role Refresher book for game %s", gameID)
}

// Helper methods to manage admin assignment trackers safely (NEW)
func (h *BotHandler) GetAdminAssignmentTracker(gameID gameEntity.GameID) (*tgutil.RefreshingMessageBook, bool) {
	h.adminRefreshMutex.RLock()
	defer h.adminRefreshMutex.RUnlock()
	book, exists := h.adminAssignmentTrackers[gameID]
	return book, exists
}

func (h *BotHandler) GetOrCreateAdminAssignmentTracker(gameID gameEntity.GameID) *tgutil.RefreshingMessageBook {
	h.adminRefreshMutex.Lock()
	defer h.adminRefreshMutex.Unlock()
	book, exists := h.adminAssignmentTrackers[gameID]
	if !exists {
		book = tgutil.NewRefreshState(func(user int64, data string) (string, []interface{}, error) {
			gameIDFromData := gameEntity.GameID(data)
			// Fetch necessary data (Game, State, Players) - May need error handling
			gameData, _ := h.GetGameByIDHandler().Handle(context.Background(), gameQuery.GetGameByIDQuery{ID: gameIDFromData})
			state, stateExists := h.GetInteractiveSelectionState(gameIDFromData)

			if gameData == nil || !stateExists {
				log.Printf("Refresh Admin: Game %s or State not found.", gameIDFromData)
				return "Error: Game data unavailable.", []interface{}{}, nil // Return error state message
			}
			// Prepare message content using the helper from callbacks_game
			return game.PrepareAdminAssignmentMessage(gameData, state, h.msgs)
		})
		h.adminAssignmentTrackers[gameID] = book
		log.Printf("Created new Admin Assignment Tracker book for game %s", gameID)
	}
	return book
}

func (h *BotHandler) DeleteAdminAssignmentTracker(gameID gameEntity.GameID) {
	h.adminRefreshMutex.Lock()
	defer h.adminRefreshMutex.Unlock()
	delete(h.adminAssignmentTrackers, gameID)
	log.Printf("Deleted Admin Assignment Tracker book for game %s", gameID)
}
