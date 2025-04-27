# 4. Repository Port Specification

**Goal:** Define the exact Go interface contracts for data persistence, separating the domain layer from specific storage implementations.

## 4.1. Room Repository (`internal/domain/room/port/room_repository.go`)

*   **`RoomReader` Interface**
    *   **Purpose:** Defines methods for reading room data.
    *   **Methods:**
        *   `GetRoomByID(id roomEntity.RoomID) (*roomEntity.Room, error)`: Get room by ID.
        *   `GetRooms() ([]*roomEntity.Room, error)`: Get all rooms.
        *   `GetPlayerRooms(playerID sharedEntity.UserID) ([]*roomEntity.Room, error)`: Get rooms a specific player is in.
        *   `GetPlayersInRoom(roomID roomEntity.RoomID) ([]*sharedEntity.User, error)`: Get players in a specific room.

*   **`RoomWriter` Interface**
    *   **Purpose:** Defines methods for writing/modifying room data.
    *   **Methods:**
        *   `CreateRoom(room *roomEntity.Room) error`: Create a new room.
        *   `UpdateRoom(room *roomEntity.Room) error`: Update an existing room.
        *   `AddPlayerToRoom(roomID roomEntity.RoomID, player *sharedEntity.User) error`: Add a player to a room.
        *   `RemovePlayerFromRoom(roomID roomEntity.RoomID, playerID sharedEntity.UserID) error`: Remove a player from a room.
        *   `DeleteRoom(roomID roomEntity.RoomID) error`: Delete a room by ID.

*   **`RoomRepository` Interface**
    *   **Purpose:** Combined interface embedding `RoomReader` and `RoomWriter` for room persistence.
    *   **Composition:** Embeds `RoomReader`, `RoomWriter`.

**Notes:**

*   Refresh state management (previously handled via change flags in the repository) is now handled within the presentation layer using `tgutil.RefreshState`.
*   The repository port interfaces (`RoomReader`, `RoomWriter`, `RoomRepository`) **MUST NOT** include any methods related to change flags (e.g., `CheckChangeFlag`, `ConsumeChangeFlag`). This state is managed entirely by the presentation layer (`internal/presentation/telegram/handler` and `internal/shared/tgutil`).
*   `AssignScenarioToRoom` and `GetRoomScenario` methods are not part of the repository. Room state changes (like setting `ScenarioName` or adding `Description`) should be done on the `Room` entity itself, and then persisted using the `UpdateRoom` method.
*   The repository implementation should handle potential errors like `roomEntity.ErrRoomNotFound`, `roomEntity.ErrRoomAlreadyExists`, `roomEntity.ErrPlayerNotInRoom`.

## 4.2. Scenario Repository (`internal/domain/scenario/port/scenario_repository.go`)

*   **`ScenarioReader` Interface**
    *   **Purpose:** Defines methods for reading scenario data.
    *   **Methods:**
        *   `GetScenarioByID(id string) (*scenarioEntity.Scenario, error)`: Get scenario by ID.
        *   `GetAllScenarios() ([]*scenarioEntity.Scenario, error)`: Get all scenarios.

*   **`ScenarioWriter` Interface**
    *   **Purpose:** Defines methods for writing/modifying scenario data.
    *   **Methods:**
        *   `CreateScenario(scenario *scenarioEntity.Scenario) error`: Create a new scenario.
        *   `DeleteScenario(id string) error`: Delete a scenario by ID.
        *   `AddRoleToScenario(scenarioID string, role scenarioEntity.Role) error`: Add a role to a specific scenario.
        *   `RemoveRoleFromScenario(scenarioID string, roleName string) error`: Remove a role (by name) from a specific scenario.

*   **`ScenarioRepository` Interface**
    *   **Purpose:** Combined interface embedding `ScenarioReader` and `ScenarioWriter` for scenario persistence.
    *   **Composition:** Embeds `ScenarioReader`, `ScenarioWriter`.

**Notes:**

*   Assumes scenario IDs are strings.
*   `RemoveRoleFromScenario` operates based on `roleName`.
*   Should handle errors like scenario not found, scenario already exists, role not found.

## 4.3. Game Repository (`internal/domain/game/port/game_repository.go`)

*   **`GameReader` Interface**
    *   **Purpose:** Defines methods for reading game data.
    *   **Methods:**
        *   `GetGameByID(id gameEntity.GameID) (*gameEntity.Game, error)`: Get game by ID.
        *   `GetGameByRoomID(roomID roomEntity.RoomID) (*gameEntity.Game, error)`: Get the game associated with a specific room ID.
        *   `GetAllGames() ([]*gameEntity.Game, error)`: Get all active games.

*   **`GameWriter` Interface**
    *   **Purpose:** Defines methods for writing/modifying game data.
    *   **Methods:**
        *   `CreateGame(game *gameEntity.Game) error`: Create a new game instance.
        *   `UpdateGame(game *gameEntity.Game) error`: Update an existing game (used for saving state changes like assignments, state transitions).
        *   `DeleteGame(id gameEntity.GameID) error`: Delete a game by ID.

*   **`GameRepository` Interface**
    *   **Purpose:** Combined interface embedding `GameReader` and `GameWriter` for game persistence.
    *   **Composition:** Embeds `GameReader`, `GameWriter`.

**Notes:**

*   `UpdateGame` is crucial for persisting changes made to the Game entity after creation.
*   Should handle errors like game not found, game already exists.

## 4.4. Room Client Port (`internal/domain/game/port/room_client.go`)

*   **`RoomClient` Interface**
    *   **Purpose:** Defines an interface for the Game domain to fetch Room data (abstraction for potential microservice communication).
    *   **Methods:**
        *   `FetchRoom(id roomEntity.RoomID) (*roomEntity.Room, error)`: Fetches room data by ID.

**Notes:**

*   This interface abstracts the fetching of Room data needed by the Game domain.

## 4.5. Scenario Client Port (`internal/domain/game/port/scenario_client.go`)

*   **`ScenarioClient` Interface**
    *   **Purpose:** Defines an interface for the Game domain to fetch Scenario data (abstraction for potential microservice communication).
    *   **Methods:**
        *   `FetchScenario(id string) (*scenarioEntity.Scenario, error)`: Fetches scenario data by ID.

**Notes:**

*   This interface abstracts the fetching of Scenario data needed by the Game domain.

// End of file 