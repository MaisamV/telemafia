# Blueprint 05: Game Module Details

**Source:** `internal/domain/game/`

**Purpose:** Details the components within the Game domain module.

## 1. `entity/game.go`

*   **`GameID` (type `string`):** Unique identifier for a game instance.
*   **`GameState` (type `string`):** Represents the game's current state (constants: `GameStateWaitingForPlayers`, `GameStateRolesAssigned`, `GameStateInProgress`, `GameStateFinished`).
*   **`Game` struct:**
    *   Fields: `ID` (GameID), `State` (GameState), `Room` (*roomEntity.Room), `Scenario` (*scenarioEntity.Scenario), `Assignments` (map[sharedEntity.UserID]scenarioEntity.Role).
    *   `Room`, `Scenario`: Pointers to the specific Room and Scenario entities associated with this game instance.
    *   `Assignments`: Maps `UserID` to the `scenarioEntity.Role` assigned to that user.
*   **Methods:** `AssignRole()`, `SetRolesAssigned()`, `StartGame()`, `FinishGame()` modify the `Assignments` map or `State` field.

## 2. `port/` (Ports)

Defines interfaces required by the Game domain.

*   **`game_repository.go`:**
    *   `GameReader` interface: `GetGameByID(id GameID)`, `GetGameByRoomID(roomID roomEntity.RoomID)`, `GetAllGames()`.
    *   `GameWriter` interface: `CreateGame(game *Game)`, `UpdateGame(game *Game)`, `DeleteGame(id GameID)`.
    *   `GameRepository` interface: Embeds `GameReader` and `GameWriter`.
*   **`room_client.go`:**
    *   `RoomClient` interface: `FetchRoom(id roomEntity.RoomID) (*roomEntity.Room, error)`.
    *   Abstracts fetching room data. Allows the Game domain to get room details without depending directly on the Room repository implementation. Used by `CreateGameHandler` and potentially `AssignRolesHandler`.
*   **`scenario_client.go`:**
    *   `ScenarioClient` interface: `FetchScenario(id string) (*scenarioEntity.Scenario, error)`.
    *   Abstracts fetching scenario data. Used by `CreateGameHandler` and `AssignRolesHandler`.

## 3. `usecase/command/` (Commands - State Changing)

*   **`assign_roles.go`:**
    *   `AssignRolesCommand`: Contains `Requester`, `GameID`.
    *   `AssignRolesHandler`: Depends on `GameRepository`, `ScenarioReader` (or `ScenarioClient`), `RoomReader` (or `RoomClient`).
        *   Handles admin check.
        *   Fetches `Game`, `Scenario`, and Players (via `RoomReader` using `Game.Room.ID`).
        *   Validates player count matches role count from the scenario.
        *   Flattens roles from the `Scenario` struct.
        *   Uses `common.Shuffle` to randomize the order of the flattened roles.
        *   Assigns the shuffled roles to sorted players (by UserID).
        *   Updates the `Game.Assignments` map and sets `Game.State` to `GameStateRolesAssigned`.
        *   Calls `GameRepository.UpdateGame`.
        *   Returns the assignments map `map[sharedEntity.User]scenarioEntity.Role` (Note: key is the `User` struct).
*   **`create_game.go`:**
    *   `CreateGameCommand`: Contains `Requester`, `RoomID`, `ScenarioID`.
    *   `CreateGameHandler`: Depends on `GameRepository`, `RoomClient`, `ScenarioClient`.
        *   Handles admin check.
        *   Uses `RoomClient` and `ScenarioClient` to fetch the actual `Room` and `Scenario` entities.
        *   Creates a new `Game` entity, associating the fetched Room and Scenario, setting initial state to `GameStateWaitingForPlayers`.
        *   Generates a unique `GameID`.
        *   Calls `GameRepository.CreateGame`.
        *   Returns the created `*Game` entity.

## 4. `usecase/query/` (Queries - Data Retrieval)

*   **`get_game_by_id.go`:**
    *   `GetGameByIDQuery`: Contains `ID` (GameID).
    *   `GetGameByIDHandler`: Depends on `GameReader`. Calls `GameReader.GetGameByID`.
*   **`get_games.go`:**
    *   `GetGamesQuery`: (Empty).
    *   `GetGamesHandler`: Depends on `GameReader`. Calls `GameReader.GetAllGames`. 