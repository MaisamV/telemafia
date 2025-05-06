# Blueprint 05: Game Module Details

**Source:** `internal/domain/game/`

**Purpose:** Details the components within the Game domain module.

## 1. `entity/game.go`

*   **`GameID` (type `string`):** Unique identifier for a game.
*   **`GameState` (type `string`):** Represents the current state of the game (e.g., `WaitingForPlayers`, `RolesAssigned`, `InProgress`, `Finished`). Defined constants for states.
*   **`Game` struct:**
    *   Fields: `ID` (GameID), `Room` (*roomEntity.Room), `Scenario` (*scenarioEntity.Scenario), `State` (GameState), `Assignments` (map[sharedEntity.UserID]scenarioEntity.Role).
    *   `Assignments`: Maps player UserIDs to their assigned Role (which includes Name and Side).
*   **`AssignRole(playerID UserID, role Role)`:** Adds/updates an entry in the `Assignments` map.
*   **`SetRolesAssigned()`:** Sets the `State` to `GameStateRolesAssigned`.
*   **(Other state transition methods as needed...)**

## 2. `port/game_repository.go`

Defines the interfaces required by the Game domain to interact with persistence and other domain services (via clients).

*   **`GameReader` interface:**
    *   `GetGameByID(id GameID) (*Game, error)`
    *   `GetGamesByState(state GameState) ([]*Game, error)`
    *   `GetAllGames() ([]*Game, error)`
    *   `GetGameByRoomID(roomID roomEntity.RoomID) (*Game, error)`
*   **`GameWriter` interface:**
    *   `CreateGame(game *Game) error`
    *   `UpdateGame(game *Game) error`
    *   `DeleteGame(id GameID) error`
*   **`GameRepository` interface:** Embeds `GameReader` and `GameWriter`.
*   **`RoomClient` interface:** (Client for interacting with Room domain - potentially external service)
    *   `FetchRoom(id roomEntity.RoomID) (*roomEntity.Room, error)`
*   **`ScenarioClient` interface:** (Client for interacting with Scenario domain - potentially external service)
    *   `FetchScenario(id string) (*scenarioEntity.Scenario, error)`

## 3. `usecase/command/` (Commands - State Changing)

*   **`assign_roles.go`:**
    *   `AssignRolesCommand`: Contains `Requester` (User), `GameID`.
    *   `AssignRolesHandler`: Depends on `GameRepository`, `ScenarioReader`, `RoomReader`. Fetches game, performs permission check (global admin OR moderator of the game's room), fetches scenario and players, uses `Scenario.FlatRoles()` to get a flat list of roles, checks player/role count match, sorts roles by name hash and users by ID, shuffles roles (using a copy), updates game `Assignments` map, sets game state, updates game in repository. Returns the assignments map.
*   **`create_game.go`:**
    *   `CreateGameCommand`: Contains `Requester` (User), `RoomID`, `ScenarioID`.
    *   `CreateGameHandler`: Depends on `GameRepository`, `RoomClient`, `ScenarioClient`. Fetches room and scenario via clients, performs permission check (global admin OR moderator of the fetched room), creates new `Game` entity, saves game via repository. Returns the created game.

## 4. `usecase/query/` (Queries - Data Retrieval)

*   **`get_game.go`:**
    *   `GetGameByIDQuery`: Contains `GameID`.
    *   `GetGameByIDHandler`: Depends on `GameReader`. Calls `GameReader.GetGameByID`.
*   **`get_games.go`:**
    *   `GetGamesQuery`: Contains `State` (optional filter).
    *   `GetGamesHandler`: Depends on `GameReader`. Calls `GetGamesByState` or `GetAllGames` based on query. 