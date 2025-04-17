package main

import (
	"fmt"
	"log"
	"telemafia/internal/config"
	telegramHandler "telemafia/internal/interfaces/handler/telegram"
	memrepo "telemafia/internal/interfaces/repository/memory"
	"telemafia/internal/usecase"
	"telemafia/pkg/event"
	"time"

	"gopkg.in/telebot.v3"
)

// EventPublisher implements event.Publisher (moved from config.go)
type EventPublisher struct{}

func (p *EventPublisher) Publish(e event.Event) error {
	// For now, just log the events
	// TODO: Implement a proper event bus (e.g., using channels or a library)
	log.Printf("Event published: %+v\n", e)
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

// initializeDependencies sets up and wires all components (moved from config.go)
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
	roomRepo := memrepo.NewInMemoryRepository()
	scenarioRepo := memrepo.NewInMemoryScenarioRepository()
	gameRepo := memrepo.NewInMemoryGameRepository()

	// Initialize event publisher
	eventPublisher := &EventPublisher{}

	// Initialize use case handlers (Interactors)
	createRoomHandler := usecase.NewCreateRoomHandler(roomRepo, eventPublisher)
	joinRoomHandler := usecase.NewJoinRoomHandler(roomRepo, eventPublisher)
	leaveRoomHandler := usecase.NewLeaveRoomHandler(roomRepo, eventPublisher)
	kickUserHandler := usecase.NewKickUserHandler(roomRepo, eventPublisher)
	deleteRoomHandler := usecase.NewDeleteRoomHandler(roomRepo)
	resetChangeFlagHandler := usecase.NewResetChangeFlagCommand(roomRepo)
	raiseChangeFlagHandler := usecase.NewRaiseChangeFlagHandler(roomRepo)
	getRoomHandler := usecase.NewGetRoomHandler(roomRepo)
	getRoomsHandler := usecase.NewGetRoomsHandler(roomRepo)
	getPlayerRoomsHandler := usecase.NewGetPlayerRoomsHandler(roomRepo)
	getPlayersInRoomsHandler := usecase.NewGetPlayersInRoomHandler(roomRepo)
	checkChangeFlagHandler := usecase.NewCheckChangeFlagHandler(roomRepo)
	createScenarioHandler := usecase.NewCreateScenarioHandler(scenarioRepo)
	deleteScenarioHandler := usecase.NewDeleteScenarioHandler(scenarioRepo)
	manageRolesHandler := usecase.NewManageRolesHandler(scenarioRepo)
	getScenarioByIDHandler := usecase.NewGetScenarioByIDHandler(scenarioRepo)
	getAllScenariosHandler := usecase.NewGetAllScenariosHandler(scenarioRepo)

	createGameHandler := usecase.NewCreateGameHandler(gameRepo)
	assignRolesHandler := usecase.NewAssignRolesHandler(gameRepo, scenarioRepo)
	getGamesHandler := usecase.NewGetGamesHandler(gameRepo)
	getGameByIDHandler := usecase.NewGetGameByIDHandler(gameRepo)

	// Initialize Telegram Bot Handler (Delivery Mechanism)
	botHandler := telegramHandler.NewBotHandler(
		telegramBot,
		cfg.AdminUsernames,
		roomRepo,
		createRoomHandler,
		joinRoomHandler,
		leaveRoomHandler,
		kickUserHandler,
		deleteRoomHandler,
		resetChangeFlagHandler,
		raiseChangeFlagHandler,
		getRoomsHandler,
		getPlayerRoomsHandler,
		getPlayersInRoomsHandler,
		getRoomHandler,
		checkChangeFlagHandler,
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
