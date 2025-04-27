# 8. Configuration Specification

**Goal:** Define how the application is configured, including parameters, loading mechanism, and file structure.

## 8.1. Configuration Parameters

The following parameters **MUST** be configurable:

1.  **`telegram_bot_token`**
    *   **Type:** `string`
    *   **Description:** The authentication token provided by Telegram's BotFather for the bot.
    *   **Required:** Yes

2.  **`admin_usernames`**
    *   **Type:** `list of strings`
    *   **Description:** A list of Telegram usernames (case-sensitive) that are granted administrative privileges within the bot.
    *   **Required:** Yes (can be an empty list if no admins are needed initially, but the key should exist if using a configuration file).

## 8.2. Loading Mechanism (`internal/config/config.go`)

*   Configuration **MUST** be loadable via command-line arguments AND a JSON configuration file.
*   **Priority:** Command-line arguments **MUST** take precedence over the configuration file.
*   **Command-Line Flags:**
    *   `-token <your_token>`: Specifies the Telegram bot token.
    *   `-admins "admin1,admin2,admin3"`: Specifies a comma-separated list of admin usernames.
*   **Configuration File (`config.json`):**
    *   If configuration flags are *not* provided via command-line, the application **MUST** attempt to load configuration from a file named `config.json` located in the application's root directory.
    *   The file **MUST** be valid JSON containing keys matching the parameter names (e.g., `"telegram_bot_token"`, `"admin_usernames"`).
*   **Error Handling:** If neither command-line flags nor a valid `config.json` file provide the required parameters (`telegram_bot_token`), the application **MUST** fail to start and log an informative error message.
*   **(NEW)** **Messages File:** The application also requires a `messages.json` file (typically in the same directory as the executable or project root during development) containing user-facing strings. Loading this file is handled separately by `internal/presentation/telegram/messages.LoadMessages`. If this file is missing or invalid, the application **MUST** also fail to start.

## 8.3. `config.json` Structure Example

```json
// Example Structure (Keys MUST match parameter names)
{
  "telegram_bot_token": "YOUR_TELEGRAM_BOT_TOKEN",
  "admin_usernames": ["admin_user_one", "another_admin"]
}
```

## 8.4. Go Implementation (`internal/config/config.go`)

*   A `Config` struct **MUST** be defined to hold the loaded values:
    *   Field `TelegramBotToken` (`string`) corresponding to `telegram_bot_token`.
    *   Field `AdminUsernames` (`[]string`) corresponding to `admin_usernames`.
    *   Appropriate JSON tags **MUST** be used if field names differ from JSON keys.
*   A `LoadConfig(filename string) (*Config, error)` function **MUST** be implemented:
    1.  Define flags using the Go `flag` package (`flag.String`) for `-token` and `-admins`.
    2.  Parse command-line flags.
    3.  If flag values are non-empty, populate the `Config` struct from flags (splitting the admin string by comma) and return.
    4.  If flags are empty, attempt to open and read the specified `filename`.
    5.  If the file is read successfully, use `encoding/json` to decode the JSON data into the `Config` struct.
    6.  Validate that the required fields (token) are not empty after loading from either source.
    7.  Return the populated `Config` struct and a `nil` error on success.
    8.  If configuration cannot be loaded successfully from either flags or file, return `nil` and an appropriate error indicating missing or invalid configuration. 