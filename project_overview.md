# TeleMafia Bot - Project Overview

**Goal:** This document provides a high-level overview of the TeleMafia Bot project. For detailed architectural rules, design patterns, coding conventions, and directory structure guidelines that **MUST** be followed during development, please refer to the files within the `/rules` directory.

## 1. Project Purpose

TeleMafia is a Telegram bot written in Go designed to facilitate playing the party game Mafia. It allows users to join game rooms, and designated administrators to manage these rooms, define game scenarios (sets of roles), assign scenarios to rooms, and manage the game lifecycle (e.g., assigning roles).

## 2. Core Domain Concepts

The bot's logic is centered around these main concepts:

*   **User (`internal/shared/entity/user.go`):** Represents a Telegram user interacting with the bot. Users can have an `Admin` status derived from configuration.
*   **Room (`internal/domain/room/...`):** A virtual space where players gather before a game starts. Rooms have a name, an ID, a list of players, and can have an assigned Scenario.
*   **Scenario (`internal/domain/scenario/...`):** Defines the roles and rules for a specific Mafia game variant (e.g., "Classic 7 Player"). Contains a name, ID, and a list of Roles.
*   **Game (`internal/domain/game/...`):** Represents an active instance of a Mafia game tied to a specific Room and Scenario. It tracks the game's state (e.g., `WaitingForPlayers`, `RolesAssigned`) and the assignment of Roles to Users.

## 3. Key Features & Commands

The bot offers the following core functionalities:

*   **General:**
    *   `/start`: Welcome message and initial view (typically the room list).
    *   `/help`: Displays available commands.
*   **Room Management:**
    *   `/list_rooms`: Shows available rooms with player counts (dynamic message).
    *   `/join_room <id>` / Inline Buttons: Allows users to join a room.
    *   `/leave_room <id>` / Inline Buttons: Allows users to leave a room.
    *   `/my_rooms`: Lists rooms the user is currently in.
*   **Admin - Room Management:**
    *   `/create_room <name>`: Creates a new room.
    *   `/delete_room`: Initiates the process to select and delete a room.
    *   `/kick_user <room_id> <user_id>`: Removes a player from a room.
*   **Admin - Scenario Management:**
    *   `/create_scenario <name>`: Creates a new scenario definition.
    *   `/delete_scenario <id>`: Deletes a scenario definition.
    *   `/add_role <scenario_id> <role_name>`: Adds a role to a scenario.
    *   `/remove_role <scenario_id> <role_name>`: Removes a role from a scenario.
    *   *(Listing/viewing scenarios might be future features)*
*   **Admin - Game Management:**
    *   `/assign_scenario <room_id> <scenario_id>`: Assigns a scenario to a room and creates the corresponding Game entity.
    *   `/games`: Lists currently active game instances.
    *   `/assign_roles <game_id>`: Distributes roles to players in the specified game's room.

## 4. Technical Stack & Setup

*   **Language:** Go (v1.18+)
*   **Primary Dependency:** `gopkg.in/telebot.v3` (Telegram Bot API framework)
*   **Architecture:** Clean Architecture (Ports & Adapters), Modular Monolith, CQRS.
*   **Persistence:** In-Memory (All state is lost on bot restart).
*   **Configuration:** Requires `config.json` (for bot token, admin usernames) OR command-line flags (`-token`, `-admins`). See `rules/04_coding_conventions.md`.
*   **User Messages:** All user-facing text is stored in `messages.json` and loaded at startup. See `rules/04_coding_conventions.md` and `rules/05_telegram_layer.md`.
*   **Running:** Build with `go build ./cmd/telemafia/` and run the executable, or use `go run ./cmd/telemafia/main.go`. Ensure `config.json` and `messages.json` are present as needed.

## 5. Development Guidelines

**Crucially, all development MUST adhere to the principles and rules outlined in the `/rules` directory.** This ensures consistency in architecture, patterns, code style, and structure. 