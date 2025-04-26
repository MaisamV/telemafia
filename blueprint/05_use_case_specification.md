# 5. Use Case Specification

**Goal:** Detail the application-specific business logic for each command and query, specifying inputs, outputs, dependencies, and processing steps.

**General Notes on Command Handlers:**

*   **Authorization:** For commands requiring administrative privileges (e.g., KickUser, DeleteRoom, CreateScenario, DeleteScenario, ManageRoles, CreateGame, AssignRoles, AddDescription), the first step within the `Handle` method **MUST** be to check the `Admin` flag on the `Requester` field of the input command struct. If the flag is false, the handler must return an appropriate authorization error immediately without performing further actions.

## 5.1. Room Use Cases (`internal/domain/room/usecase/...`)

**Commands (`command/`)**

1.  **Create Room**
    *   **Command Struct:** `CreateRoomCommand { ID roomEntity.RoomID; Name string; CreatorID sharedEntity.UserID }`
    *   **Handler:** `CreateRoomHandler`
    *   **Dependencies:** `roomPort.RoomWriter`, `sharedEvent.Publisher`
    *   **Steps:**
        1.  Call `roomEntity.NewRoom(cmd.ID, cmd.Name)` to create and validate the entity.
        2.  Call `roomRepo.CreateRoom(room)` to persist.
        3.  Publish `sharedEvent.RoomCreatedEvent`.
        4.  Return the created `*roomEntity.Room` or error.

2.  **Join Room**
    *   **Command Struct:** `JoinRoomCommand { Requester sharedEntity.User; RoomID roomEntity.RoomID }`
    *   **Handler:** `JoinRoomHandler`
    *   **Dependencies:** `roomPort.RoomRepository` (Reader for check, Writer for add), `sharedEvent.Publisher`
    *   **Steps:**
        1.  Call `roomRepo.GetRoomByID(cmd.RoomID)` to ensure room exists.
        2.  Call `roomRepo.AddPlayerToRoom(cmd.RoomID, &cmd.Requester)` to add player.
        3.  Publish `sharedEvent.PlayerJoinedEvent`.
        4.  Return `nil` on success, or error.

3.  **Leave Room**
    *   **Command Struct:** `LeaveRoomCommand { Requester sharedEntity.User; RoomID roomEntity.RoomID }`
    *   **Handler:** `LeaveRoomHandler`
    *   **Dependencies:** `roomPort.RoomRepository` (Reader for check, Writer for remove), `sharedEvent.Publisher`
    *   **Steps:**
        1.  (Optional: `roomRepo.GetRoomByID` check)
        2.  Call `roomRepo.RemovePlayerFromRoom(cmd.RoomID, cmd.Requester.ID)`.
        3.  Publish `sharedEvent.PlayerLeftEvent`.
        4.  Return `nil` on success, or error (`ErrPlayerNotInRoom` etc.).

4.  **Kick User**
    *   **Command Struct:** `KickUserCommand { Requester sharedEntity.User; RoomID roomEntity.RoomID; PlayerID sharedEntity.UserID }`
    *   **Handler:** `KickUserHandler`
    *   **Dependencies:** `roomPort.RoomRepository`, `sharedEvent.Publisher`
    *   **Steps:**
        1.  (Optional: `roomRepo.GetRoomByID` check)
        2.  Call `roomRepo.RemovePlayerFromRoom(cmd.RoomID, cmd.PlayerID)`.
        3.  Publish `sharedEvent.PlayerKickedEvent`.
        4.  Return `nil` on success, or error.

5.  **Delete Room**
    *   **Command Struct:** `DeleteRoomCommand { Requester sharedEntity.User; RoomID roomEntity.RoomID }`
    *   **Handler:** `DeleteRoomHandler`
    *   **Dependencies:** `roomPort.RoomWriter`
    *   **Steps:**
        1.  Call `roomRepo.DeleteRoom(cmd.RoomID)`.
        2.  Return error from repository.

6.  **Add Description**
    *   **Command Struct:** `AddDescriptionCommand { Requester sharedEntity.User; Room *roomEntity.Room; DescriptionName string; Text string }`
    *   **Handler:** `AddDescriptionHandler`
    *   **Dependencies:** `roomPort.RoomRepository`
    *   **Steps:**
        1.  Validate `cmd.Room` is not nil.
        2.  Call `cmd.Room.SetDescription(cmd.DescriptionName, cmd.Text)`.
        3.  Call `roomRepo.UpdateRoom(cmd.Room)` to persist.
        4.  Return error, if any.

**Queries (`query/`)**

1.  **Get Room By ID**
    *   **Query Struct:** `GetRoomQuery { RoomID roomEntity.RoomID }`
    *   **Handler:** `GetRoomHandler`
    *   **Dependencies:** `roomPort.RoomReader`
    *   **Steps:** Call `roomRepo.GetRoomByID(query.RoomID)`, return result.

2.  **Get All Rooms**
    *   **Query Struct:** `GetRoomsQuery {}`
    *   **Handler:** `GetRoomsHandler`
    *   **Dependencies:** `roomPort.RoomReader`
    *   **Steps:** Call `roomRepo.GetRooms()`, return result.

3.  **Get Player Rooms**
    *   **Query Struct:** `GetPlayerRoomsQuery { PlayerID sharedEntity.UserID }`
    *   **Handler:** `GetPlayerRoomsHandler`
    *   **Dependencies:** `roomPort.RoomReader`
    *   **Steps:** Call `roomRepo.GetPlayerRooms(query.PlayerID)`, return result.

4.  **Get Players In Room**
    *   **Query Struct:** `GetPlayersInRoomQuery { RoomID roomEntity.RoomID }`
    *   **Handler:** `GetPlayersInRoomHandler`
    *   **Dependencies:** `roomPort.RoomReader`
    *   **Steps:** Call `roomRepo.GetPlayersInRoom(query.RoomID)`, return result.

## 5.2. Scenario Use Cases (`internal/domain/scenario/usecase/...`)

**Commands (`command/`)**

1.  **Create Scenario**
    *   **Command Struct:** `CreateScenarioCommand { Requester sharedEntity.User; ID string; Name string }`
    *   **Handler:** `CreateScenarioHandler`
    *   **Dependencies:** `scenarioPort.ScenarioWriter`
    *   **Steps:**
        1.  Create `scenarioEntity.Scenario` instance with empty `Roles`.
        2.  Call `scenarioRepo.CreateScenario(scenario)`.
        3.  Return error from repository.

2.  **Delete Scenario**
    *   **Command Struct:** `DeleteScenarioCommand { Requester sharedEntity.User; ID string }`
    *   **Handler:** `DeleteScenarioHandler`
    *   **Dependencies:** `scenarioPort.ScenarioWriter`
    *   **Steps:**
        1.  Call `scenarioRepo.DeleteScenario(cmd.ID)`, return error.

3.  **Manage Roles (Add/Remove)**
    *   **Command Structs:**
        *   `AddRoleCommand { Requester sharedEntity.User; ScenarioID string; Role scenarioEntity.Role }`
        *   `RemoveRoleCommand { Requester sharedEntity.User; ScenarioID string; RoleName string }`
    *   **Handler:** `ManageRolesHandler`
    *   **Dependencies:** `scenarioPort.ScenarioRepository`
    *   **Steps (Add):**
        1.  Call `scenarioRepo.AddRoleToScenario(cmd.ScenarioID, cmd.Role)`, return error.
    *   **Steps (Remove):**
        1.  Call `scenarioRepo.RemoveRoleFromScenario(cmd.ScenarioID, cmd.RoleName)`, return error.

**Queries (`query/`)**

1.  **Get Scenario By ID**
    *   **Query Struct:** `GetScenarioByIDQuery { ID string }`
    *   **Handler:** `GetScenarioByIDHandler`
    *   **Dependencies:** `scenarioPort.ScenarioReader`
    *   **Steps:** Call `scenarioRepo.GetScenarioByID(query.ID)`, return result.

2.  **Get All Scenarios**
    *   **Query Struct:** `GetAllScenariosQuery {}`
    *   **Handler:** `GetAllScenariosHandler`
    *   **Dependencies:** `scenarioPort.ScenarioReader`
    *   **Steps:** Call `scenarioRepo.GetAllScenarios()`, return result.

## 5.3. Game Use Cases (`internal/domain/game/usecase/...`)

**Commands (`command/`)**

1.  **Create Game**
    *   **Command Struct:** `CreateGameCommand { Requester sharedEntity.User; RoomID roomEntity.RoomID; ScenarioID string }`
    *   **Handler:** `CreateGameHandler`
    *   **Dependencies:** `gamePort.GameRepository`, `gamePort.RoomClient`, `gamePort.ScenarioClient`
    *   **Steps:**
        1.  Fetch `Room` using `roomClient.FetchRoom(cmd.RoomID)`.
        2.  Fetch `Scenario` using `scenarioClient.FetchScenario(cmd.ScenarioID)`.
        3.  Create `gameEntity.Game` instance, assigning fetched Room/Scenario pointers, generate `GameID`, set `State` to `WaitingForPlayers`, initialize empty `Assignments` map.
        4.  Call `gameRepo.CreateGame(game)`.
        5.  Return created `*gameEntity.Game` or error.

2.  **Assign Roles**
    *   **Command Struct:** `AssignRolesCommand { Requester sharedEntity.User; GameID gameEntity.GameID }`
    *   **Handler:** `AssignRolesHandler`
    *   **Dependencies:** `gamePort.GameRepository`, `scenarioPort.ScenarioReader`, `roomPort.RoomReader`
    *   **Steps:**
        1.  Fetch `Game` using `gameRepo.GetGameByID(cmd.GameID)`.
        2.  Fetch `Scenario` using `scenarioRepo.GetScenarioByID(game.Scenario.ID)`.
        3.  Fetch `Players` (Users) from `roomRepo.GetPlayersInRoom(game.Room.ID)`.
        4.  Validate player count matches scenario role count. Return error if mismatch.
        5.  Shuffle scenario roles randomly.
        6.  Iterate through players, assign shuffled roles, populate `game.Assignments` map.
        7.  Update `game.State` to `RolesAssigned`.
        8.  Call `gameRepo.UpdateGame(game)` to persist changes.
        9.  Return the `assignments` map or error.

**Queries (`query/`)**

1.  **Get Game By ID**
    *   **Query Struct:** `GetGameByIDQuery { GameID gameEntity.GameID }`
    *   **Handler:** `GetGameByIDHandler`
    *   **Dependencies:** `gamePort.GameReader`
    *   **Steps:** Call `gameRepo.GetGameByID(query.GameID)`, return result.

2.  **Get All Games**
    *   **Query Struct:** `GetAllGamesQuery {}`
    *   **Handler:** `GetAllGamesHandler`
    *   **Dependencies:** `gamePort.GameReader`
    *   **Steps:** Call `gameRepo.GetAllGames()`, return result.

3.  **Get Player Games**
    *   **Query Struct:** `GetPlayerGamesQuery { PlayerID sharedEntity.UserID }`
    *   **Handler:** `GetPlayerGamesHandler`
    *   **Dependencies:** `gamePort.GameReader`
    *   **Steps:** Call `gameRepo.GetPlayerGames(query.PlayerID)`, return result.

4.  **Get Games In Room**
    *   **Query Struct:** `GetGamesInRoomQuery { RoomID roomEntity.RoomID }`
    *   **Handler:** `GetGamesInRoomHandler`
    *   **Dependencies:** `gamePort.GameReader`
    *   **Steps:** Call `gameRepo.GetGamesInRoom(query.RoomID)`, return result.