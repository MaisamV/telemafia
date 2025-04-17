# TeleMafia Bot (Golang)

This is a Telegram bot for managing Mafia game rooms, built with Golang. It features a layered architecture, dependency injection, and uses Command Query Responsibility Segregation (CQRS) principles for handling actions. Configuration is supported via both `config.json` and command-line arguments.

**Important Note:** This bot currently uses **in-memory storage**. This means all rooms, scenarios, games, and user data will be **lost** when the bot restarts.

---

## 🛠 Installation

### 1️⃣ Install Go
Make sure you have Go installed (at least Go 1.18). You can download it from:
- Official Go Website: https://go.dev/dl/

---

## 🚀 Running the Bot

### Option 1: Using config.json
1. Create a `config.json` file in the project root with the following structure:
   ```json
   {
     "telegram_bot_token": "YOUR_TELEGRAM_BOT_TOKEN",
     "admin_usernames": ["admin1", "admin2"]
   }
   ```
   * Replace `YOUR_TELEGRAM_BOT_TOKEN` with your actual bot token.
   * Replace `admin1`, `admin2` with a comma-separated list of admin usernames (case-sensitive).

2. Run the bot:
   ```sh
   go run main.go
   ```

---

### Option 2: Using Command-Line Arguments
If you don't have a `config.json`, you can pass parameters directly when starting the bot:

```sh
go run main.go -token "YOUR_TELEGRAM_BOT_TOKEN" -admins "admin1,admin2"
```

* Replace `YOUR_TELEGRAM_BOT_TOKEN` with your actual bot token.
* Replace `admin1,admin2` with a comma-separated list of admin usernames (case-sensitive).

Example:

```sh
go run main.go -token "123456789:ABCDEF-TOKEN" -admins "john_doe,admin_user"
```

---

## 📝 Bot Commands

The bot supports the following commands:

### General Commands
| Command          | Description                      | Who Can Use? |
|------------------|----------------------------------|--------------|
| `/start`         | Start interaction with the bot   | Everyone     |
| `/help`          | Show help message                | Everyone     |

### Room Management
| Command                   | Description                                      | Who Can Use?  |
|---------------------------|--------------------------------------------------|---------------|
| `/create_room <room_name>`| Create a new room                                | Admins Only   |
| `/join_room <room_id>`    | Join a specific room                             | Everyone      |
| `/leave_room <room_id>`   | Leave a specific room                            | Everyone      |
| `/list_rooms`             | List all available rooms                         | Everyone      |
| `/my_rooms`               | List rooms you have joined                       | Everyone      |
| `/kick_user <room_id> <username>` | Kick a user from a room                  | Admins Only   |
| `/delete_room <room_id>`  | Delete a specific room                           | Admins Only   |
| `/assign_scenario <room_id> <scenario_id>` | Assign a game scenario to a room | Admins Only   |

### Scenario Management (Game Rulesets)
*Scenarios define the set of roles available for a game.*
| Command                           | Description                         | Who Can Use?  |
|-----------------------------------|-------------------------------------|---------------|
| `/create_scenario <scenario_name>`| Create a new game scenario          | Admins Only   |
| `/delete_scenario <scenario_id>`  | Delete a game scenario              | Admins Only   |
| `/add_role <scenario_id> <role_name>` | Add a role to a scenario        | Admins Only   |
| `/remove_role <scenario_id> <role_name>`| Remove a role from a scenario | Admins Only   |
| `/list_scenarios` (TODO)          | List available scenarios            | Admins Only   | (*Not yet implemented*) |

### Game Management
| Command                   | Description                                   | Who Can Use?  |
|---------------------------|-----------------------------------------------|---------------|
| `/games`                  | List active games (shows game ID and room ID) | Admins Only   |
| `/assign_roles <game_id>` | Assign roles from the room's scenario to players in the game | Admins Only   |

---

## 🔧 How It Works

1.  **Configuration:** The bot reads settings (`telegram_bot_token`, `admin_usernames`) from command-line arguments (`-token`, `-admins`) first. If not provided, it falls back to `config.json`.
2.  **Initialization:** It sets up the Telegram bot connection, initializes in-memory repositories for rooms, scenarios, and games, and prepares command/query handlers.
3.  **Architecture:** The bot follows a layered architecture:
    *   **Delivery (Telegram):** Handles incoming Telegram commands and callbacks, interacts with users.
    *   **Internal (Core Logic):** Contains use cases (commands/queries), domain entities (Room, Game, Scenario, User), and repository interfaces. Uses CQRS pattern.
    *   **Infrastructure:** Provides concrete implementations (currently in-memory) for repositories.
4.  **Command Handling:** Incoming commands are routed to specific handlers in the delivery layer, which then execute corresponding use cases in the internal layer.
5.  **Data Storage:** All application data (rooms, players, scenarios, games) is stored **in memory** and is **lost** upon bot restart.
6.  **Admin Identification:** Admins are identified by their Telegram **username** (case-sensitive) as provided in the configuration.

---

## 🛠 Troubleshooting

🔹 **Problem:** "Missing bot token" error
✔️ **Solution:** Ensure you provide a valid Telegram bot token via `config.json` or the `-token` command-line argument.

🔹 **Problem:** "Admins not detected" or admin commands don't work
✔️ **Solution:** Ensure the admin usernames in `config.json` or the `-admins` argument exactly match the Telegram usernames (case-sensitive).

🔹 **Problem:** Rooms, scenarios, or game data disappear after restarting the bot.
✔️ **Solution:** This is expected behavior. The bot currently uses **in-memory storage**, which does not persist data.

---

## 📜 License

This project is released under a **non-commercial license**.

- You **may use, modify, and distribute** this project **for personal or educational purposes**.
- You **cannot use this project** or its derivatives **for commercial purposes**, including selling it, integrating it into paid products, or using it in any way that generates revenue.
- Any modifications or forks must retain this license and include attribution to the original author.

If you wish to use this project for commercial purposes, please contact the author for licensing arrangements. 🚀  
