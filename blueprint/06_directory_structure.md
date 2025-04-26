# 6. Directory Structure Specification

**Goal:** Define the standard directory layout for the project to ensure consistency and adherence to the Clean Architecture layering.

```text
.
├── .git/                 # Git repository data
├── .gitignore            # Files/patterns for Git to ignore
├── blueprint/            # Contains these specification documents
├── cmd/
│   └── telemafia/        # Main application package
│       └── main.go       # Application entry point, composition root
├── config.json           # Example configuration file (fallback)
├── go.mod                # Go module definition
├── go.sum                # Go module checksums
├── internal/             # Private application code (not externally importable)
│   ├── adapters/         # Infrastructure Layer Implementations
│   │   ├── repository/   # Data Repository Adapters
│   │   │   └── memory/   # Concrete implementation for in-memory storage
│   │   │       ├── room_repository.go     # Implements domain.room.port.RoomRepository
│   │   │       ├── scenario_repository.go # Implements domain.scenario.port.ScenarioRepository
│   │   │       └── game_repository.go     # Implements domain.game.port.GameRepository
│   │   └── api/          # NEW: Adapters simulating external API clients (bridge domain ports)
│   │       ├── room_client.go       # NEW: Implements game domain's RoomClient port
│   │       └── scenario_client.go   # NEW: Implements game domain's ScenarioClient port
│   ├── config/           # Application Configuration Loading
│   │   └── config.go     # Config struct definition and loading logic
│   ├── domain/           # Core Domain Layer (Business Logic)
│   │   ├── room/         # Room Domain Module
│   │   │   ├── entity/   # Room specific entities and value objects
│   │   │   │   └── room.go
│   │   │   ├── port/     # Interfaces (Ports) for Room module
│   │   │   │   └── room_repository.go
│   │   │   └── usecase/  # Application Business Rules (Use Cases)
│   │   │       ├── command/ # Commands (state-changing operations)
│   │   │       │   ├── create_room.go
│   │   │       │   ├── join_room.go
│   │   │       │   ├── leave_room.go
│   │   │       │   ├── kick_user.go
│   │   │       │   ├── delete_room.go
│   │   │       │   ├── raise_change_flag.go  # (In-Memory Specific)
│   │   │       │   └── reset_change_flag.go  # (In-Memory Specific)
│   │   │       └── query/   # Queries (data retrieval operations)
│   │   │           ├── get_room.go
│   │   │           ├── get_rooms.go
│   │   │           ├── get_player_rooms.go
│   │   │           ├── get_players_in_room.go
│   │   │           └── flag_query.go         # (In-Memory Specific)
│   │   ├── scenario/     # Scenario Domain Module
│   │   │   ├── entity/
│   │   │   │   └── scenario.go # Defines Scenario and Role
│   │   │   ├── port/
│   │   │   │   └── scenario_repository.go
│   │   │   └── usecase/
│   │   │       ├── command/
│   │   │       │   ├── create_scenario.go
│   │   │       │   ├── delete_scenario.go
│   │   │       │   └── manage_roles.go    # Handles Add/Remove Role commands
│   │   │       └── query/
│   │   │           ├── get_scenario_by_id.go
│   │   │           └── get_all_scenarios.go
│   │   └── game/         # Game Domain Module
│   │       ├── entity/
│   │       │   └── game.go # Defines Game, GameID, GameState
│   │       ├── port/
│   │       │   ├── game_repository.go
│   │       │   ├── room_client.go           # NEW: Interface for fetching room data
│   │       │   └── scenario_client.go       # NEW: Interface for fetching scenario data
│   │       └── usecase/
│   │           ├── command/
│   │           │   ├── create_game.go
│   │           │   └── assign_roles.go
│   │           └── query/
│   │               ├── get_game_by_id.go
│   │               └── get_games.go
│   ├── presentation/     # Presentation Layer (Driving Adapters)
│   │   └── telegram/     # Telegram Bot Adapter
│   │       └── handler/    # Handlers mapping Telegram commands/callbacks to domain use cases
│   │           ├── bot_handler.go       # Main handler struct, DI, registration
│   │           ├── callbacks.go         # Main callback router (HandleCallback)
│   │           ├── refresh.go           # Dynamic message refresh logic (e.g., RefreshRoomsList)
│   │           ├── common_handlers.go   # NEW: Handlers for /start, /help
│   │           ├── room/                # EXPORTED Handlers related to the Room domain
│   │           │   ├── create_room.go
│   │           │   ├── join_room.go
│   │           │   ├── leave_room.go
│   │           │   ├── kick_user.go
│   │           │   ├── delete_room.go
│   │           │   ├── list_rooms.go
│   │           │   ├── my_rooms.go
│   │           │   └── callbacks_room.go
│   │           ├── scenario/            # Handlers related to the Scenario domain
│   │           │   ├── create_scenario.go
│   │           │   ├── delete_scenario.go
│   │           │   ├── manage_roles.go
│   │           │   └── callbacks_scenario.go # (Empty for now, add if needed)
│   │           └── game/                # Handlers related to the Game domain
│   │               ├── assign_scenario.go
│   │               ├── assign_roles.go
│   │               ├── list_games.go
│   │               └── callbacks_game.go
│   └── shared/           # Shared components across layers/domains
│       ├── common/       # General utility functions (DEPRECATED? -> see tgutil)
│       │   └── utils.go
│       ├── entity/       # Shared domain entities
│       │   └── user.go
│       ├── event/        # Domain event definitions and publisher interface
│       │   ├── event.go  # Base Event interface and Publisher interface
│       │   └── events.go # Concrete event structs (RoomCreated, PlayerJoined, etc.)
│       ├── logger/       # (Optional) Shared logger implementation
│       │   └── logger.go
│       └── tgutil/       # NEW: Shared Telegram utility functions and constants
│           ├── const.go
│           └── util.go
├── README.md             # Project overview, setup, and usage instructions
├── structure.txt         # (Optional) Text version of this structure spec
└── tests/                # Unit and integration tests (structure TBD)
    └── ...
```

**Key Points:**

*   All application core logic resides within the `internal` directory.
*   The `domain` layer is strictly separated by business capability (`room`, `scenario`, `game`).
*   Each domain module contains its `entity`, `port`, and `usecase` (split into `command` and `query`).
*   `adapters` contains concrete implementations for `domain/ports`.
*   `presentation` contains concrete implementations for driving the application (like the Telegram bot handler).
*   `cmd/telemafia/main.go` acts as the central point for dependency injection and application startup.
*   `shared` holds code reusable across different parts of the application, avoiding direct dependencies between domain modules where possible.
*   Telegram command and callback handlers are organized:
    *   Domain-specific handlers (`/create_room`, `/assign_scenario`, etc.) reside in EXPORTED functions within subdirectories (`room/`, `scenario/`, `game/`).
    *   Common handlers (`/start`, `/help`) reside in an EXPORTED function file (e.g., `common_handlers.go`) in the parent `handler` directory.
    *   The main `BotHandler` struct, dispatcher methods, and the main callback router (`handleCallback`) remain in the parent `handler` directory (`bot_handler.go`, `callbacks.go`).
    *   Dispatcher methods and the callback router call the exported handler functions from the corresponding files/sub-packages.
*   Shared Telegram-specific utilities (`ToUser`, `IsAdmin`, etc.) and constants (`Unique...`) reside in `internal/shared/tgutil/`.
*   Dynamic message refreshing logic (like `RefreshRoomsList`) resides in `refresh.go` and is initiated in `BotHandler.Start()`. 