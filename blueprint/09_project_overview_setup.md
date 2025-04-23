# 9. Project Overview & Setup Specification

**Goal:** Provide high-level context about the project, basic setup, dependencies, and running instructions.

## 9.1. Project Description

*   **Name:** TeleMafia Bot
*   **Purpose:** A Telegram bot written in Go to facilitate playing the party game Mafia. It allows users (especially designated admins) to manage game rooms, define game scenarios (sets of roles), create games, and assign roles to players.
*   **Core Technology:** Go (Golang)
*   **Key Architectural Feature:** Clean Architecture (Ports & Adapters), Modular Monolith, CQRS.
*   **Persistence:** In-Memory (data is lost on restart).

## 9.2. Setup Requirements

*   **Go Installation:** Go version 1.18 or higher must be installed on the system.
*   **Dependencies:** Project dependencies are managed using Go Modules.

## 9.3. Dependencies (`go.mod`)

The primary external dependency is:

*   `gopkg.in/telebot.v3`: Framework for interacting with the Telegram Bot API.

*(The `go.mod` file will also include indirect dependencies pulled in by `telebot`)*

Example `go.mod` content:

```go
module telemafia

go 1.22 // Specify target Go version (e.g., 1.22 or higher)

require gopkg.in/telebot.v3 v3.3.8 // Or latest compatible version

// Indirect dependencies will be added automatically by Go tooling
```

## 9.4. Running the Bot

1.  **Navigate:** Open a terminal in the project's root directory.
2.  **Build & Run (Recommended):**
    ```bash
    go build -o telemafia_bot ./cmd/telemafia/
    ./telemafia_bot -token "YOUR_TOKEN" -admins "admin1,admin2"
    ```
3.  **Run Directly:**
    ```bash
    go run ./cmd/telemafia/main.go -token "YOUR_TOKEN" -admins "admin1,admin2"
    ```
4.  **Using `config.json`:**
    *   Create a `config.json` file in the root directory (see Configuration Specification for structure).
    *   Run without flags:
        ```bash
        # Using built binary
        ./telemafia_bot
        # Or running directly
        go run ./cmd/telemafia/main.go
        ```

*Replace `YOUR_TOKEN` and `admin1,admin2` with actual values.*

## 9.5. Basic `README.md` Content

A basic `README.md` should be generated including:

*   Project title and brief description.
*   Mention of Clean Architecture, CQRS, and In-Memory persistence.
*   Installation steps (Go required).
*   Running instructions (both flag and `config.json` methods).
*   A summary of the main bot commands (from Feature Specification).
*   A note about the non-commercial license (if applicable). 