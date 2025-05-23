# Blueprint 01: Dependency Wiring (Composition Root)

**Source:** `cmd/telemafia/main.go`

**Purpose:** This document outlines how application components are instantiated and wired together at the application's entry point, following the Dependency Injection (DI) pattern.

## 1. Overview

The `main` function in `cmd/telemafia/main.go` acts as the **Composition Root**. It performs the following steps:

1.  **Initialize Random Seed:** Calls `common.InitSeed()` to ensure randomness for features like role shuffling.
2.  **Load Configuration:** Reads settings (Bot Token, Admin Usernames) from `config.json` or command-line flags using `config.LoadConfig`.
3.  **Load Messages:** Reads user-facing text from `messages.json` using `messages.LoadMessages`.
4.  **Initialize Dependencies:** Calls the `initializeDependencies` function to create and connect all necessary components.
5.  **Register Handlers:** Calls `botHandler.RegisterHandlers()` to map Telegram commands and callbacks to their respective handler methods.
6.  **Start Bot:** Calls `botHandler.Start()` to begin the bot's polling loop and background tasks (like message refreshing).

## 2. `initializeDependencies` Function

This function is the core of the DI setup:

1.  **Telegram Bot:** Initializes the `telebot.Bot` instance using the token from the configuration.
2.  **Repositories (Adapters):**
    *   Instantiates in-memory repositories for `Room`, `Scenario`, and `Game` using their respective `NewInMemory...Repository()` constructors from `internal/adapters/repository/memory/`.
    *   These constructors return the *port interface types* (e.g., `roomPort.RoomRepository`), decoupling the rest of the application from the specific implementation.
3.  **API Client Adapters:**
    *   Instantiates local client adapters (`LocalRoomClient`, `LocalScenarioClient`) from `internal/adapters/api/`.
    *   These adapters currently wrap the in-memory repositories, simulating communication within the monolith but allowing for future replacement with actual network clients if the modules were split into microservices. They depend on the *reader* interfaces of the repositories.
4.  **Event Publisher:**
    *   Instantiates a simple `EventPublisher` (defined locally in `main.go`) that currently just logs events. It implements the `event.Publisher` interface from `internal/shared/event/`.
5.  **Use Case Handlers (Domain Interactors):**
    *   Instantiates command and query handlers for each domain module (`room`, `scenario`, `game`), including the new `roomCommand.NewChangeModeratorHandler(roomRepo)`.
    *   **Constructor Injection:** Dependencies like repositories (ports), other clients (ports), and the event publisher are passed into the handler constructors (e.g., `roomCommand.NewCreateRoomHandler(roomRepo, eventPublisher)`). Handlers depend on *interface types*.
6.  **Telegram Bot Handler (Presentation):**
    *   Instantiates the main `telegramHandler.BotHandler` from `internal/presentation/telegram/handler/`.
    *   **Constructor Injection:** Injects the `telebot.Bot` instance, admin usernames, loaded messages (`*messages.Messages`), and *all* the previously instantiated use case handlers (including the new `changeModeratorHandler`).
7.  **Return:** Returns the fully configured `BotHandler` instance.

## 3. Key Principles Demonstrated

*   **Dependency Inversion:** Components depend on abstractions (interfaces/ports) defined in inner layers (domain), not on concrete implementations from outer layers (adapters, main).
*   **Composition Root:** A single location (`main.go`) is responsible for composing the application object graph.
*   **Constructor Injection:** Dependencies are provided through constructors, making dependencies explicit and facilitating testing.

## 4. Changes from Original Code

*   Added `ChangeModeratorHandler` instantiation and injection.

## 5. Updated Code

```
	getPlayerRoomsHandler := roomQuery.NewGetPlayerRoomsHandler(roomRepo)
	getPlayersInRoomsHandler := roomQuery.NewGetPlayersInRoomHandler(roomRepo)
	addDescriptionHandler := roomCommand.NewAddDescriptionHandler(roomRepo)
	changeModeratorHandler := roomCommand.NewChangeModeratorHandler(roomRepo)

	// Scenario Use Cases
		getPlayerRoomsHandler,
		getPlayersInRoomsHandler,
		getRoomHandler,
		addDescriptionHandler,
		changeModeratorHandler,
		createScenarioHandler,
		assignRolesHandler,
		createGameHandler,
		updateGameHandler,
		getGamesHandler,
		getGameByIDHandler,
		// Initialize maps for interactive selection state and refresh books
		interactiveSelections:      make(map[gameEntity.GameID]*tgutil.InteractiveSelectionState),
		playerRoleChoiceRefreshers: make(map[gameEntity.GameID]*tgutil.RefreshingMessageBook),
		adminAssignmentTrackers:    make(map[gameEntity.GameID]*tgutil.RefreshingMessageBook),
		// Pass message generation logic when initializing global refresh books
		roomListRefreshMessage: tgutil.NewRefreshState(func(user int64, data string) (string, []interface{}, error) {
			message, markup, err := room.PrepareRoomListMessage(getRoomsHandler, getPlayersInRoomsHandler, msgs)
			opts := []interface{}{markup, telebot.NoPreview}
			return message, opts, err
		}),
		roomDetailRefreshMessage: tgutil.NewRefreshState(func(user int64, data string) (string, []interface{}, error) {
			message, markup, err := room.RoomDetailMessage(getRoomsHandler, getPlayersInRoomsHandler, msgs, entity.UserID(user), data)
			opts := []interface{}{markup, telebot.ModeMarkdownV2, telebot.NoPreview}
			return message, opts, err
		}),
	}

	botHandler := telegramHandler.NewBotHandler(telebot.Bot{
		Token:          config.BotToken,
		AdminUsernames: config.AdminUsernames,
		Messages:       messages,
		Handlers:       handlers,
	})

	botHandler.RegisterHandlers()

	botHandler.Start() 