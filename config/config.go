package config

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
