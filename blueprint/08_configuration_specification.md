# 8. Configuration Specification

**Goal:** Define how the application is configured, including parameters, loading mechanism, and file structure.

## 8.1. Configuration Parameters

The following parameters are required:

1.  **`telegram_bot_token`**
    *   **Type:** `string`
    *   **Description:** The authentication token provided by Telegram's BotFather for your bot.
    *   **Required:** Yes

2.  **`admin_usernames`**
    *   **Type:** `[]string` (Slice of strings)
    *   **Description:** A list of Telegram usernames (case-sensitive) that are granted administrative privileges within the bot.
    *   **Required:** Yes (can be an empty list if no admins are needed initially, but the key should exist in `config.json` if using it).

## 8.2. Loading Mechanism (`internal/config/config.go`)

*   **Priority:** Command-line arguments take precedence over the configuration file.
*   **Command-Line Flags:**
    *   `-token <your_token>`: Specifies the Telegram bot token.
    *   `-admins "admin1,admin2,admin3"`: Specifies a comma-separated list of admin usernames.
*   **Configuration File (`config.json`):**
    *   If `-token` and `-admins` flags are *not* provided, the application attempts to load configuration from a file named `config.json` located in the application's root directory.
    *   The file **MUST** be valid JSON.
*   **Error Handling:** If neither command-line flags nor a valid `config.json` file with the required parameters are found, the application **MUST** fail to start and log an informative error message.

## 8.3. `config.json` Structure

```json
{
  "telegram_bot_token": "YOUR_TELEGRAM_BOT_TOKEN",
  "admin_usernames": ["admin_user_one", "another_admin"]
}
```

## 8.4. Go Implementation (`internal/config/config.go`)

*   Define a `Config` struct to hold the loaded values:
    ```go
    type Config struct {
        TelegramBotToken string   `json:"telegram_bot_token"`
        AdminUsernames   []string `json:"admin_usernames"`
    }
    ```
*   Implement a `LoadConfig(filename string) (*Config, error)` function:
    1.  Define flags using the `flag` package (`flag.String`).
    2.  Call `flag.Parse()`.
    3.  Check if flag values (`*t`, `*admins`) are non-empty. If so, populate the `Config` struct from flags (splitting the admin string by comma) and return.
    4.  If flags are empty, attempt to `os.Open(filename)`.
    5.  If file opens successfully, use `json.NewDecoder(file).Decode(&cfg)` to parse the JSON into the `Config` struct.
    6.  Validate that the token is not empty after decoding.
    7.  Return the populated `Config` struct from the file.
    8.  If file opening or decoding fails (and flags were not used), return `nil` and an appropriate error indicating missing configuration. 