package main

import (
	"fmt"
	"log"
	apiAdapter "telemafia/internal/adapters/api"
	memrepo "telemafia/internal/adapters/repository/memory"
	"telemafia/internal/config"
	gameCommand "telemafia/internal/domain/game/usecase/command"
	gameQuery "telemafia/internal/domain/game/usecase/query"
	roomCommand "telemafia/internal/domain/room/usecase/command"
	roomQuery "telemafia/internal/domain/room/usecase/query"
	scenarioCommand "telemafia/internal/domain/scenario/usecase/command"
	scenarioQuery "telemafia/internal/domain/scenario/usecase/query"
	telegramHandler "telemafia/internal/presentation/telegram/handler"
	"telemafia/internal/shared/event"
	"time"

	"gopkg.in/telebot.v3"
)

// EventPublisher implements event.Publisher
type EventPublisher struct{}

func (p *EventPublisher) Publish(e event.Event) error {
	// For now, just log the events
	// TODO: Implement a proper event bus (e.g., using channels or a library)
	log.Printf("Event published: Type=%s, Data=%+v\n", e.EventName(), e)
	return nil
}

func main() {
	// Load Configuration
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("Configuration error: %v", err)
	}

	// Initialize Dependencies (Composition Root)
	botHandler, err := initializeDependencies(cfg)
	if err != nil {
		log.Fatalf("Initialization error: %v", err)
	}

	// Register bot handlers
	botHandler.RegisterHandlers()

	log.Println("Bot is running...")
	// Start the bot (blocking call)
	botHandler.Start()
}

// initializeDependencies sets up and wires all components
func initializeDependencies(cfg *config.Config) (*telegramHandler.BotHandler, error) {
	// Initialize Telegram Bot
	botSettings := telebot.Settings{
		Token:  cfg.TelegramBotToken,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	telegramBot, err := telebot.NewBot(botSettings)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Telegram bot: %w", err)
	}

	// Initialize repositories (Adapters)
	roomRepo := memrepo.NewInMemoryRoomRepository()
	scenarioRepo := memrepo.NewInMemoryScenarioRepository()
	gameRepo := memrepo.NewInMemoryGameRepository()

	// Initialize API Client Adapters (using local repos for now)
	roomClient := apiAdapter.NewLocalRoomClient(roomRepo)
	scenarioClient := apiAdapter.NewLocalScenarioClient(scenarioRepo)

	// Initialize event publisher
	eventPublisher := &EventPublisher{}

	// Initialize use case handlers (Interactors) - using new packages and constructor names
	// Room Use Cases
	createRoomHandler := roomCommand.NewCreateRoomHandler(roomRepo, eventPublisher)
	joinRoomHandler := roomCommand.NewJoinRoomHandler(roomRepo, eventPublisher)
	leaveRoomHandler := roomCommand.NewLeaveRoomHandler(roomRepo, eventPublisher)
	kickUserHandler := roomCommand.NewKickUserHandler(roomRepo, eventPublisher)
	deleteRoomHandler := roomCommand.NewDeleteRoomHandler(roomRepo)
	resetRefreshHandler := roomCommand.NewResetChangeFlagHandler(roomRepo)
	raiseChangeFlagHandler := roomCommand.NewRaiseChangeFlagHandler(roomRepo)
	getRoomHandler := roomQuery.NewGetRoomHandler(roomRepo)
	getRoomsHandler := roomQuery.NewGetRoomsHandler(roomRepo)
	getPlayerRoomsHandler := roomQuery.NewGetPlayerRoomsHandler(roomRepo)
	getPlayersInRoomsHandler := roomQuery.NewGetPlayersInRoomHandler(roomRepo)
	checkRefreshHandler := roomQuery.NewCheckChangeFlagHandler(roomRepo)
	addDescriptionHandler := roomCommand.NewAddDescriptionHandler(roomRepo)

	// Scenario Use Cases
	createScenarioHandler := scenarioCommand.NewCreateScenarioHandler(scenarioRepo)
	deleteScenarioHandler := scenarioCommand.NewDeleteScenarioHandler(scenarioRepo)
	manageRolesHandler := scenarioCommand.NewManageRolesHandler(scenarioRepo)
	getScenarioByIDHandler := scenarioQuery.NewGetScenarioByIDHandler(scenarioRepo)
	getAllScenariosHandler := scenarioQuery.NewGetAllScenariosHandler(scenarioRepo)

	// Game Use Cases
	createGameHandler := gameCommand.NewCreateGameHandler(gameRepo, roomClient, scenarioClient)
	assignRolesHandler := gameCommand.NewAssignRolesHandler(gameRepo, scenarioRepo, roomRepo)
	getGamesHandler := gameQuery.NewGetGamesHandler(gameRepo)
	getGameByIDHandler := gameQuery.NewGetGameByIDHandler(gameRepo)

	// Initialize Telegram Bot Handler (Delivery Mechanism)
	// Pass the correctly typed repository (roomRepo satisfies the interface needed by BotHandler)
	botHandler := telegramHandler.NewBotHandler(
		telegramBot,
		cfg.AdminUsernames,
		roomRepo, // Pass the room repository (as RoomWriter)
		createRoomHandler,
		joinRoomHandler,
		leaveRoomHandler,
		kickUserHandler,
		deleteRoomHandler,
		resetRefreshHandler,    // Corrected name
		raiseChangeFlagHandler, // Corrected name
		getRoomsHandler,
		getPlayerRoomsHandler,
		getPlayersInRoomsHandler,
		getRoomHandler,
		checkRefreshHandler,   // Corrected name
		addDescriptionHandler, // ADDED in correct position
		createScenarioHandler,
		deleteScenarioHandler,
		manageRolesHandler,
		getScenarioByIDHandler,
		getAllScenariosHandler,
		assignRolesHandler,
		createGameHandler,
		getGamesHandler,
		getGameByIDHandler,
	)

	return botHandler, nil
}
