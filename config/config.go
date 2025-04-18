package config

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"telemafia/delivery/telegram"
	"telemafia/delivery/util"
	gameMemory "telemafia/internal/infrastructure/game/memory"
	roomMemory "telemafia/internal/infrastructure/room/memory"
	roomCommand "telemafia/internal/room/usecase/command"
	roomQuery "telemafia/internal/room/usecase/query"
	"telemafia/pkg/event"
	"time"

	scenarioMemory "telemafia/internal/infrastructure/scenario/memory"
	scenarioCommand "telemafia/internal/scenario/usecase/command"
	scenarioQuery "telemafia/internal/scenario/usecase/query"

	gameCommand "telemafia/internal/game/usecase/command"
	gameQuery "telemafia/internal/game/usecase/query"

	"gopkg.in/telebot.v3"
)

// Config holds the application configuration
type Config struct {
	TelegramBotToken string   `json:"telegram_bot_token"`
	AdminUsernames   []string `json:"admin_usernames"`
}

// GlobalConfig is a global instance of Config
var GlobalConfig *Config

// GetGlobalConfig returns the global configuration instance
func GetGlobalConfig() *Config {
	return GlobalConfig
}

// LoadConfig reads the bot token and admin usernames from CLI arguments first, then falls back to a JSON file if needed.
func LoadConfig(filename string) (*Config, error) {
	// Parse command-line arguments
	token := flag.String("token", "", "Telegram bot token")
	admins := flag.String("admins", "", "Comma-separated list of admin usernames")
	flag.Parse()

	// If CLI arguments are provided, use them directly
	if *token != "" && *admins != "" {
		fmt.Println("✅ Loaded configuration from command-line arguments")
		GlobalConfig = &Config{
			TelegramBotToken: *token,
			AdminUsernames:   strings.Split(*admins, ","),
		}
		return GlobalConfig, nil
	}

	// If CLI arguments are missing, try to load from config.json
	file, err := os.Open(filename)
	if err == nil {
		defer file.Close()
		cfg := &Config{}
		decoder := json.NewDecoder(file)
		if err := decoder.Decode(cfg); err == nil && cfg.TelegramBotToken != "" {
			fmt.Println("✅ Loaded configuration from", filename)
			GlobalConfig = cfg
			util.SetAdminUsers(cfg.AdminUsernames)
			return GlobalConfig, nil
		}
	}

	// If both CLI arguments and config.json fail, return an error
	return nil, errors.New("❌ Error: Bot token and admin usernames must be provided either via command-line arguments or config.json")
}

// EventPublisher implements both room and user event publishers
type EventPublisher struct{}

func (p *EventPublisher) Publish(event event.Event) error {
	// For now, just log the events
	log.Printf("Event published: %+v\n", event)
	return nil
}

// InitializeDependencies initializes and returns all necessary dependencies
func InitializeDependencies(cfg *Config) (*telegram.BotHandler, error) {
	util.SetAdminUsers(GetGlobalConfig().AdminUsernames)
	// Initialize Telegram Bot
	botSettings := telebot.Settings{
		Token:  cfg.TelegramBotToken,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	telegramBot, err := telebot.NewBot(botSettings)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Telegram bot: %w", err)
	}

	// Initialize repositories from infrastructure layer
	roomRepo := roomMemory.NewInMemoryRepository()

	// Initialize scenario repository
	scenarioRepo := scenarioMemory.NewInMemoryRepository()

	// Initialize game repository
	gameRepo := gameMemory.NewInMemoryGameRepository()

	// Initialize event publishers
	eventPublisher := &EventPublisher{}

	// Initialize command handlers
	createRoomHandler := roomCommand.NewCreateRoomHandler(roomRepo, eventPublisher)
	joinRoomHandler := roomCommand.NewJoinRoomHandler(roomRepo, eventPublisher)
	leaveRoomHandler := roomCommand.NewLeaveRoomHandler(roomRepo, eventPublisher)
	kickUserHandler := roomCommand.NewKickUserHandler(roomRepo, eventPublisher)
	deleteRoomHandler := roomCommand.NewDeleteRoomHandler(roomRepo)
	resetChangeFlagHandler := roomCommand.NewResetChangeFlagCommand(roomRepo)
	raiseChangeFlagHandler := roomCommand.NewRaiseChangeFlagHandler(roomRepo)
	getRoomHandler := roomQuery.NewGetRoomHandler(roomRepo)

	// Initialize scenario command handlers
	createScenarioHandler := scenarioCommand.NewCreateScenarioHandler(scenarioRepo)
	deleteScenarioHandler := scenarioCommand.NewDeleteScenarioHandler(scenarioRepo)
	manageRolesHandler := scenarioCommand.NewManageRolesHandler(scenarioRepo)

	// Initialize scenario query handlers
	getScenarioByIDHandler := scenarioQuery.NewGetScenarioByIDHandler(scenarioRepo)
	getAllScenariosHandler := scenarioQuery.NewGetAllScenariosHandler(scenarioRepo)

	// Initialize query handlers
	getRoomsHandler := roomQuery.NewGetRoomsHandler(roomRepo)
	getPlayerRoomsHandler := roomQuery.NewGetPlayerRoomsHandler(roomRepo)
	getPlayersInRoomsHandler := roomQuery.NewGetPlayersInRoomHandler(roomRepo)
	checkChangeFlagHandler := roomQuery.NewCheckChangeFlagHandler(roomRepo)

	// Initialize game handlers
	createGameHandler := gameCommand.NewCreateGameHandler(gameRepo)
	assignRolesHandler := gameCommand.NewAssignRolesHandler(gameRepo, scenarioRepo)
	getGamesHandler := gameQuery.NewGetGamesHandler(gameRepo)
	getGameByIDHandler := gameQuery.NewGetGameByIDHandler(gameRepo)

	// Initialize room command handlers that use game functionality
	assignScenarioHandler := roomCommand.NewAssignScenarioHandler(roomRepo)

	// Initialize bot handler
	botHandler := telegram.NewBotHandler(
		telegramBot,
		cfg.AdminUsernames,
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
		assignScenarioHandler,
		createGameHandler,
		getGamesHandler,
		getGameByIDHandler,
	)

	return botHandler, nil
}
