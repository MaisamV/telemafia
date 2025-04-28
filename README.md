# TeleMafia Bot

[![Go Reference](https://pkg.go.dev/badge/github.com/username/telemafia.svg)](https://pkg.go.dev/github.com/username/telemafia)

A Telegram bot written in Go to facilitate playing the party game Mafia.

**Important Note:** This bot currently uses **in-memory storage**. This means all rooms, scenarios, games, and user data will be **lost** when the bot restarts.

---

## 🚀 Getting Started

### Prerequisites

*   Go (Version 1.18 or higher installed)
*   A Telegram Bot Token from BotFather

### Configuration

The bot requires configuration for the Telegram token and admin usernames.

1.  **`config.json` (Recommended):**
    *   Create `config.json` in the project root:
      ```json
      {
        "telegram_bot_token": "YOUR_TELEGRAM_BOT_TOKEN",
        "admin_usernames": ["your_admin_username", "another_admin"]
      }
      ```
    *   Replace placeholders with your actual token and desired admin Telegram usernames (case-sensitive).
2.  **Command-line Flags (Overrides `config.json`):**
    *   `-token "YOUR_TOKEN"`: Specifies the bot token.
    *   `-admins "admin1,admin2"`: Specifies a comma-separated list of admin usernames.

Additionally, the bot requires a `messages.json` file in the project root containing user-facing text. A default version is included.

### Running the Bot

1.  **Navigate** to the project root directory.
2.  **Run directly:**
    ```bash
    # Using config.json & messages.json
    go run ./cmd/telemafia/main.go 
    
    # Using flags (ensure messages.json exists)
    go run ./cmd/telemafia/main.go -token "YOUR_TOKEN" -admins "admin1,admin2"
    ```
3.  **Build and Run:**
    ```bash
    # Build the executable (e.g., named 'telemafia_bot')
    go build -o telemafia_bot ./cmd/telemafia/
    
    # Run the executable (ensure config.json & messages.json are in the same directory)
    ./telemafia_bot
    
    # Or run with flags (ensure messages.json is present)
    ./telemafia_bot -token "YOUR_TOKEN" -admins "admin1,admin2"
    ```

---

## 📖 Documentation & Guidelines

*   **Project Overview:** For a high-level understanding of the bot's purpose, core concepts, and key features, see [project_overview.md](./project_overview.md).
*   **Development Rules:** To ensure consistency in architecture, patterns, and code style, all contributors **MUST** follow the guidelines outlined in the `/rules` directory.
    *   [Architecture Guidelines](./rules/01_architecture_guidelines.md)
    *   [Design Patterns](./rules/02_design_patterns.md)
    *   [Directory Structure Guide](./rules/03_directory_structure.md)
    *   [Coding Conventions & Rules](./rules/04_coding_conventions.md)
    *   [Telegram Presentation Layer Rules](./rules/05_telegram_layer.md)
    *   [Domain Modeling Guidelines](./rules/06_domain_modeling.md)

---

## 📜 License

This project is released under a **non-commercial license**.

- You **may use, modify, and distribute** this project **for personal or educational purposes**.
- You **cannot use this project** or its derivatives **for commercial purposes**.
- Any modifications or forks must retain this license and include attribution.
