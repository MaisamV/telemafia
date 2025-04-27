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

*   `gopkg.in/telebot.v3` (or compatible version): Framework for interacting with the Telegram Bot API.

*(The `go.mod` file will also include indirect dependencies pulled in by `telebot`)*

Example `go.mod` content:

*   The Go module definition file (`go.mod`) **MUST** specify a compatible Go version (e.g., 1.18 or higher) and require a compatible version of `gopkg.in/telebot.v3`.

## 9.4. Running the Bot

1.  **Navigate:** Open a terminal in the project's root directory.
2.  **Build:** Compile the application using `go build -o <output_name> ./cmd/telemafia/`.
3.  **Run:** Execute the compiled binary.
    *   **Using Flags:** Provide configuration via flags: `./<output_name> -token "YOUR_TOKEN" -admins "admin1,admin2"`.
    *   **Using `config.json`:** Create `config.json` (see Configuration Spec) and `messages.json` (or ensure it exists) in the same directory as the executable and run `./<output_name>` without flags.
4.  **Run Directly (Development):** Use `go run ./cmd/telemafia/main.go` with the same flag/`config.json` options. Ensure `messages.json` exists in the project root.

*Replace placeholders like `<output_name>`, `YOUR_TOKEN`, and `admin1,admin2` with actual values.*

## 9.5. Basic `README.md` Content

A basic `README.md` should be generated including:

*   Project title and brief description.
*   Mention of Clean Architecture, CQRS, and In-Memory persistence.
*   Installation steps (Go required).
*   Running instructions (both flag and `config.json` methods).
*   A summary of the main bot commands (from Feature Specification).
*   A note about the non-commercial license (if applicable). 