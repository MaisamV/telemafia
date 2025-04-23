# 7. Adapter Specification

**Goal:** Define the implementation details for the required adapters (In-Memory Persistence and Telegram Presentation).

## 7.1. In-Memory Repository Adapters (`internal/adapters/repository/memory/`)

*   **Requirement:** Implement the `RoomRepository`, `ScenarioRepository`, and `GameRepository` interfaces defined in the respective `internal/domain/<module>/port` packages.
*   **Storage:** Use standard Go maps as the underlying storage mechanism.
    *   `RoomRepository`: `map[roomEntity.RoomID]*roomEntity.Room`
    *   `ScenarioRepository`: `map[string]*scenarioEntity.Scenario` (string key is Scenario ID)
    *   `GameRepository`: `map[gameEntity.GameID]*gameEntity.Game` and potentially a secondary map like `map[roomEntity.RoomID]gameEntity.GameID` for efficient lookup by RoomID.
*   **Concurrency:** All map access (read and write) **MUST** be protected using a `sync.RWMutex` within each repository instance to ensure thread safety.
    *   Use `mutex.Lock()` for write operations (Create, Update, Delete, Add, Remove).
    *   Use `mutex.RLock()` for read operations (GetByID, GetAll, GetPlayers, etc.).
    *   Ensure `defer mutex.Unlock()` or `defer mutex.RUnlock()` is used immediately after acquiring the lock.
*   **Error Handling:** Implementations should return appropriate errors defined in the domain entities (e.g., `roomEntity.ErrRoomNotFound`) or standard Go errors (e.g., using `fmt.Errorf`) when operations fail (e.g., item not found, item already exists).
*   **`UpdateRoom` Implementation:** Must handle updating the room in the map and setting the `changeFlag`.
*   **Removed Methods:** `AssignScenarioToRoom`, `GetRoomScenario` implementations should be removed.
*   **Constructors:** Provide `NewInMemory<Type>Repository()` functions that initialize the maps and mutex and return the *repository port interface type* (e.g., `func NewInMemoryRoomRepository() roomPort.RoomRepository`).
*   **(Room Repository Specific) Change Flag:** Implement the `CheckChangeFlag`, `ConsumeChangeFlag`, and `RaiseChangeFlag` methods using a boolean flag within the repository struct, also protected by the mutex. This flag should be set to `true` on any write operation and reset by `ConsumeChangeFlag`.
*   **Data Copying (Optional but Recommended):** When returning slices or maps from read operations (e.g., `GetRooms`, `GetPlayersInRoom`), consider returning copies to prevent external modification of the internal repository state.

## 7.2. API Client Adapters (Monolith Implementation) (`internal/adapters/api/`)

*   **Purpose:** To implement the client interfaces defined in domain ports (`gamePort.RoomClient`, `gamePort.ScenarioClient`) by calling other local domain repositories. This simulates the client-side interaction pattern expected in microservices, even within the monolith.
*   **`LocalRoomClient` (`room_client.go`):**
    *   Implements `gamePort.RoomClient`.
    *   Depends on `roomPort.RoomReader` (injected).
    *   `FetchRoom` method delegates to `roomRepo.GetRoomByID`.
*   **`LocalScenarioClient` (`scenario_client.go`):**
    *   Implements `gamePort.ScenarioClient`.
    *   Depends on `scenarioPort.ScenarioReader` (injected).
    *   `FetchScenario` method delegates to `scenarioRepo.GetScenarioByID`.
*   **Notes:** In a true microservice architecture, these implementations would be replaced with ones using HTTP/gRPC clients to call separate Room and Scenario services.

## 7.3. Telegram Presentation Adapter (`internal/presentation/telegram/handler/`)

*   **Framework:** Use `gopkg.in/telebot.v3`.
*   **Main Handler (`BotHandler`):**
    *   This struct acts as the primary container for the Telegram presentation layer.
    *   It **MUST** receive the initialized `*telebot.Bot` instance and all necessary domain *use case handlers* (e.g., `*roomCommand.CreateRoomHandler`, `*roomQuery.GetRoomsHandler`, etc.) via its constructor (`NewBotHandler`).
    *   It should also store the list of admin usernames loaded from configuration.
    *   Must receive `*roomCommand.AddDescriptionHandler` via constructor.
    *   `RegisterHandlers()` method: Maps Telegram command strings (e.g., `/start`, `/create_room`) to specific handler methods within this package (e.g., `h.HandleStart`, `h.HandleCreateRoom`) using `bot.Handle()`.
    *   `Start()` method: Calls `bot.Start()` to begin polling for updates.
*   **Command Handler Methods (e.g., `HandleCreateRoom`, `HandleJoinRoom` in `handlers.go`):**
    *   These methods receive the `telebot.Context`.
    *   **Input Parsing:** Extract command arguments from `c.Message().Payload`. Use `strings.TrimSpace`, `strings.Fields`, etc., for parsing.
    *   **User Conversion:** Convert the `*telebot.User` (from `c.Sender()`) to the internal `*sharedEntity.User` using the `util.ToUser()` helper function.
    *   **Use Case Invocation:** Create the appropriate domain Command or Query struct (e.g., `roomCommand.CreateRoomCommand`) with data parsed from the input and the converted User object (including the `Requester` field for commands needing it).
    *   Call the `Handle` method of the corresponding injected use case handler (e.g., `h.createRoomHandler.Handle(context.Background(), cmd)`).
    *   **Response Handling:** Based on the error or result from the use case handler, send appropriate messages back to the user via `c.Send()`. Format messages clearly. Use inline keyboards (`telebot.ReplyMarkup`) for callbacks/actions where necessary.
*   **`HandleAssignScenario` Method:**
    *   Fetches Room and Scenario.
    *   Sets `room.ScenarioName` (e.g., to `scenario.ID`).
    *   Creates `roomCommand.AddDescriptionCommand` passing the fetched `room`, the `requester`, a key like "scenario_info", and descriptive text (e.g., scenario name).
    *   Calls `h.addDescriptionHandler.Handle(...)` (admin check inside).
    *   Creates `gameCommand.CreateGameCommand` passing the `requester`.
    *   Calls `h.createGameHandler.Handle(...)` (admin check inside).
    *   Handles errors and sends appropriate response.
*   **Callback Handling (`HandleCallback`):**
    *   Handles `telebot.OnCallback` events.
    *   Extract the `unique` identifier and `payload` from `c.Callback().Data` (using a helper like `SplitCallbackData`).
    *   Use a `switch` statement on the `unique` identifier to route to specific callback logic functions (e.g., `handleDeleteRoomConfirmCallback`).
    *   Use `c.Respond()` to acknowledge the callback (dismiss loading indicators).
    *   Use `c.Edit()` to modify the original message (e.g., change text, remove keyboard) or `c.Delete()` to remove it.
*   **Utility Functions (`util.go`):**
    *   `SetAdminUsers(usernames []string)`: Stores admin list locally.
    *   `ToUser(sender *telebot.User) *sharedEntity.User`: Converts Telegram user to internal user entity. Includes setting the `Admin` flag based on the stored admin list.
    *   `SplitCallbackData(data string) (unique string, payload string)`: Parses callback data.
*   **Error Handling:** Handle errors from use case handlers gracefully, sending informative messages to the user.
*   **Context:** Pass `context.Background()` to use case handlers for now, unless specific cancellation/deadline logic is required. 