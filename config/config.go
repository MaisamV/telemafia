package config

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"gopkg.in/telebot.v3"
	"log"
	"os"
	"strings"
	"telemafia/delivery/common"
	"telemafia/delivery/telegram"
	roomMemory "telemafia/internal/infrastructure/room/memory"
	roomCommand "telemafia/internal/room/usecase/command"
	roomQuery "telemafia/internal/room/usecase/query"
	"telemafia/pkg/event"
	"time"
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
			common.SetAdminUsers(cfg.AdminUsernames)
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
	common.SetAdminUsers(GetGlobalConfig().AdminUsernames)
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

	// Initialize event publishers
	eventPublisher := &EventPublisher{}

	// Initialize command handlers
	createRoomHandler := roomCommand.NewCreateRoomHandler(roomRepo, eventPublisher)
	joinRoomHandler := roomCommand.NewJoinRoomHandler(roomRepo, eventPublisher)
	leaveRoomHandler := roomCommand.NewLeaveRoomHandler(roomRepo, eventPublisher)
	kickUserHandler := roomCommand.NewKickUserHandler(roomRepo, eventPublisher)

	// Initialize query handlers
	getRoomsHandler := roomQuery.NewGetRoomsHandler(roomRepo)
	getPlayerRoomsHandler := roomQuery.NewGetPlayerRoomsHandler(roomRepo)
	getPlayersInRoomsHandler := roomQuery.NewGetPlayersInRoomHandler(roomRepo)

	// Initialize bot handler
	botHandler := telegram.NewBotHandler(
		telegramBot,
		cfg.AdminUsernames,
		createRoomHandler,
		joinRoomHandler,
		leaveRoomHandler,
		kickUserHandler,
		getRoomsHandler,
		getPlayerRoomsHandler,
		getPlayersInRoomsHandler,
	)

	return botHandler, nil
}
