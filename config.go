package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"telemafia/config"
	"telemafia/delivery/telegram"
	roomMemory "telemafia/internal/infrastructure/room/memory"
	roomCommand "telemafia/internal/room/usecase/command"
	roomQuery "telemafia/internal/room/usecase/query"
	"time"

	"gopkg.in/telebot.v3"
)

// LoadConfig reads the bot token and admin usernames from CLI arguments first, then falls back to a JSON file if needed.
func LoadConfig(filename string) (*config.Config, error) {
	// Parse command-line arguments
	token := flag.String("token", "", "Telegram bot token")
	admins := flag.String("admins", "", "Comma-separated list of admin usernames")
	flag.Parse()

	// If CLI arguments are provided, use them directly
	if *token != "" && *admins != "" {
		fmt.Println("✅ Loaded configuration from command-line arguments")
		config.GlobalConfig = &config.Config{
			TelegramBotToken: *token,
			AdminUsernames:   strings.Split(*admins, ","),
		}
		return config.GlobalConfig, nil
	}

	// If CLI arguments are missing, try to load from config.json
	file, err := os.Open(filename)
	if err == nil {
		defer file.Close()
		cfg := &config.Config{}
		decoder := json.NewDecoder(file)
		if err := decoder.Decode(cfg); err == nil && cfg.TelegramBotToken != "" {
			fmt.Println("✅ Loaded configuration from", filename)
			config.GlobalConfig = cfg
			return config.GlobalConfig, nil
		}
	}

	// If both CLI arguments and config.json fail, return an error
	return nil, errors.New("❌ Error: Bot token and admin usernames must be provided either via command-line arguments or config.json")
}

// InitializeDependencies initializes and returns all necessary dependencies
func InitializeDependencies(cfg *config.Config) (*telegram.BotHandler, error) {
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
	roomEventPublisher := &RoomEventPublisher{publisher: eventPublisher}

	// Initialize command handlers
	createRoomHandler := roomCommand.NewCreateRoomHandler(roomRepo, roomEventPublisher)
	joinRoomHandler := roomCommand.NewJoinRoomHandler(roomRepo, roomEventPublisher)
	leaveRoomHandler := roomCommand.NewLeaveRoomHandler(roomRepo, roomEventPublisher)
	kickUserHandler := roomCommand.NewKickUserHandler(roomRepo, roomEventPublisher)

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
