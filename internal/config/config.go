package config

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
)

// Config holds the application configuration
type Config struct {
	TelegramBotToken string   `json:"telegram_bot_token"`
	AdminUsernames   []string `json:"admin_usernames"`
}

// LoadConfig reads the bot token and admin usernames from CLI arguments first, then falls back to a JSON file if needed.
func LoadConfig(filename string) (*Config, error) {
	// Define flags locally, don't rely on global state if possible
	t := flag.String("token", "", "Telegram bot token")
	admins := flag.String("admins", "", "Comma-separated list of admin usernames")
	// Consider making filename a flag too: configFile := flag.String("config", "config.json", "Path to JSON config file")
	flag.Parse()

	cfg := &Config{}

	// If CLI arguments are provided, use them directly
	if *t != "" && *admins != "" {
		fmt.Println("✅ Loaded configuration from command-line arguments")
		cfg.TelegramBotToken = *t
		cfg.AdminUsernames = strings.Split(*admins, ",")
		return cfg, nil
	}

	// If CLI arguments are missing, try to load from config.json
	// Use the filename argument passed to the function
	file, err := os.Open(filename)
	if err == nil {
		defer file.Close()
		decoder := json.NewDecoder(file)
		if err := decoder.Decode(cfg); err == nil && cfg.TelegramBotToken != "" {
			fmt.Println("✅ Loaded configuration from", filename)
			// Validate admin usernames? Ensure not empty?
			if cfg.AdminUsernames == nil {
				cfg.AdminUsernames = []string{}
			}
			return cfg, nil
		}
		// Log error if decoding fails but file exists?
		// log.Printf("Warn: Could not decode config file '%s': %v", filename, err)
	} else if !os.IsNotExist(err) {
		// Log error if opening file failed for reasons other than not existing
		// log.Printf("Warn: Could not open config file '%s': %v", filename, err)
	}

	// If both CLI arguments and config.json fail, return an error
	return nil, errors.New("❌ Error: Bot token and admin usernames must be provided either via command-line arguments (-token, -admins) or a valid config.json file")
}
