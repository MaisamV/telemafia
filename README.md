# TeleMafia Bot (Golang)

This is a Telegram bot for managing chat rooms, built with Golang. It follows SOLID principles and supports both JSON-based and Command-Line configurations.

---

## 🛠 Installation

### 1️⃣ Install Go
Make sure you have Go installed (at least Go 1.18). You can download it from:
- Official Go Website: https://go.dev/dl/

---

## 🚀 Running the Bot

### Option 1: Using config.json
1. Create a config.json file in the project root with the following structure:
   ```json
   {
   "telegram_bot_token": "YOUR_TELEGRAM_BOT_TOKEN",
   "admin_usernames": ["admin1", "admin2"]
   }

2. Run the bot:

```sh
go run main.go
```

---

### Option 2: Using Command-Line Arguments
If you don't have a config.json, you can pass parameters directly when starting the bot:

```sh
go run main.go -token "YOUR_TELEGRAM_BOT_TOKEN" -admins "admin1,admin2"
```

- Replace YOUR_TELEGRAM_BOT_TOKEN with your actual bot token.
- Replace admin1,admin2 with a comma-separated list of admin usernames.

Example:

```sh
go run main.go -token "123456789:ABCDEF-TOKEN" -admins "john_doe,admin_user"
```

---

## 📝 Bot Commands
The bot supports the following commands:

| Command                   | Description                                 | Who Can Use?  |
|---------------------------|---------------------------------------------|--------------|
| /start                   | Start the bot                               | Everyone    |
| /help                    | Show this help message                      | Everyone    |
| /create_room <room_name> | Create a new room                           | Admins Only |
| /join_room <room_id>     | Join a specific room                        | Everyone    |
| /leave_room <room_id>    | Leave a specific room                       | Everyone    |
| /list_rooms              | List all available rooms                    | Everyone    |
| /my_rooms                | List rooms you have joined                  | Everyone    |
| /kick_user               | Kick a user from a room                     | Admins Only |

---

## 🔧 How It Works
1. The bot reads configuration settings from config.json.
2. If config.json is missing or invalid, it will switch to command-line arguments (-token & -admins).
3. The bot then connects to Telegram and starts processing user commands.
4. Admins are identified by their username instead of ID for security reasons.

---

## 🛠 Troubleshooting
🔹 Problem: "Missing bot token" error  
✔️ Solution: Ensure you provide a token via config.json or CLI (-token argument).

🔹 Problem: "Admins not detected"  
✔️ Solution: Ensure the admin usernames are correctly formatted and match the usernames on Telegram.

---

## 📜 License
This project is released under a **non-commercial license**.

- You **may use, modify, and distribute** this project **for personal or educational purposes**.
- You **cannot use this project** or its derivatives **for commercial purposes**, including selling it, integrating it into paid products, or using it in any way that generates revenue.
- Any modifications or forks must retain this license and include attribution to the original author.

If you wish to use this project for commercial purposes, please contact the author for licensing arrangements. 🚀  
