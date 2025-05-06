# Blueprint 08: Telegram Presentation Layer

**Source:** `internal/presentation/telegram/`

**Purpose:** Details the components responsible for handling Telegram interactions and translating them into calls to domain use cases.

## 1. `handler/bot_handler.go` (`BotHandler`)

*   **Role:** Acts as the central orchestrator for the Telegram presentation layer.
*   **Dependencies:** Receives the `telebot.Bot` instance, admin usernames, loaded `messages.Messages`, and *all* domain use case handlers (commands and queries) via constructor injection (`NewBotHandler`).
*   **Refresh State:** Holds instances of `tgutil.RefreshingMessageBook` (e.g., `roomListRefreshMessage`, `roomDetailRefreshMessage`) to manage dynamic message updates.
*   **`RegisterHandlers()`:** Maps Telegram commands (e.g., `/create_room`, `/create_game`) and events (e.g., `telebot.OnCallback`, `telebot.OnDocument`) to specific *dispatcher methods* on the `BotHandler` struct (e.g., `h.handleCreateRoom`, `h.handleCreateGame`).
*   **Dispatcher Methods (e.g., `handleCreateRoom`, `handleCallback`, `handleDocument`, `handleCreateGame`):** These methods primarily act as routers. They parse necessary information from the `telebot.Context` (like payload, sender) and call the corresponding **exported handler functions** located in sub-packages (e.g., `room.HandleCreateRoom`, `game.HandleCreateGame`) or other handler files (`callbacks.go`, `document_handler.go`), passing the required dependencies (use case handlers, messages) from the `BotHandler` instance (`h`). Note that `RefreshingMessageBook` instances are generally *not* passed directly to these handlers anymore.
*   **Helper Methods (`Get/GetOrCreate/Delete...`):** (UPDATED) Provides methods for managing the maps holding `RefreshingMessageBook` instances (e.g., `GetOrCreateAdminAssignmentTracker`). The `GetOrCreate...` methods now instantiate the book using `tgutil.NewRefreshState`, passing the appropriate **message generation function** specific to that book's purpose.
*   **`Start()`:** Initializes background tasks (like the refresh timer) and starts the `telebot` polling loop.

## 2. `handler/callbacks.go` (`BotHandler.handleCallback`)

*   **Purpose:** Central dispatcher for *all* inline button callback queries.
*   **Logic:**
    1.  Retrieves the callback data.
    2.  Uses `tgutil.SplitCallbackData` to separate the `unique` identifier and the `payload`.
    3.  Uses a `switch` statement on the `unique` identifier.
    4.  Each `case` corresponds to a specific button type (defined in `tgutil/const.go`, e.g., `tgutil.UniqueCreateGameSelectRoom`, `tgutil.UniqueKickUserSelect`, `tgutil.UniqueChangeModeratorSelect`).
    5.  Calls the relevant **exported handler function** from the appropriate sub-package (e.g., `room.HandleKickUserSelectCallback`, `room.HandleChangeModeratorSelectCallback`), passing the `telebot.Context`, parsed `payload`, and necessary dependencies.
*   **`handleCallback(c telebot.Context) error`:** Acts as the central dispatcher for *all* inline button callback queries.
    *   Reads the callback `Unique` identifier and `Data`.
    *   Uses a `switch` statement on `Unique` to route the callback to the appropriate specialized handler function within the relevant sub-package (`room` or `game`).
    *   Handles `UniqueCancel` directly by deleting the message.
    *   **New Callbacks Routed:**
        *   `tgutil.UniqueChooseCardStart`: Routes to `game.HandleChooseCardStart`.
        *   `tgutil.UniquePlayerSelectsCard`: Routes to `game.HandlePlayerSelectsCard`.
        *   `tgutil.UniqueCancelGame`: Now routes to `game.HandleCancelCreateGame`, passing the `BotHandlerInterface` for cleanup.

## 3. `handler/refresh.go`

*   **Purpose:** Handles the background task for updating dynamic messages.
*   **`BotHandler.StartRefreshTimer()`:** (UPDATED) Runs a goroutine with a `time.Ticker`.
    *   Periodically iterates through all managed `RefreshingMessageBook` instances (global ones like `roomListRefreshMessage`, and those in maps like `adminAssignmentTrackers`).
    *   Checks `book.ConsumeRefreshNeeded()`.
    *   If true, calls `h.RefreshMessages(book)`.
*   **`BotHandler.RefreshMessages(book)`:** (NEW/REFACTORED from `updateMessages`)
    1.  Gets all active messages tracked by the specific `RefreshingMessageBook` using `book.GetAllActiveMessages()`.
    2.  For each message (`chatID`, `payload`), calls the **message generation function stored within the book** itself: `book.GetMessage(chatID, payload.Data)`.
    3.  Uses `h.bot.Edit()` to update the message in Telegram with the content returned by `book.GetMessage`.
    4.  Handles errors (e.g., message not found, user blocked bot) and removes the message from tracking using `book.RemoveActiveMessage(chatID)` if necessary.
*   **Message Preparation Functions (e.g., `room.PrepareRoomListMessage`, `game.PrepareAdminAssignmentMessage`):** Exported functions responsible for fetching current data and formatting message text/options. These are now passed into `tgutil.NewRefreshState` when a book is created (either globally in `NewBotHandler` or dynamically in `GetOrCreate...` methods).

## 4. `handler/<module>/` (e.g., `handler/room/`, `handler/game/`)

*   Contains the **exported handler functions** that implement the actual logic for specific commands and callbacks.
*   **Naming Convention:**
    *   Commands: `HandleCommandName` (e.g., `room.HandleCreateRoom`).
    *   Callbacks: `HandleCallbackNameCallback` (e.g., `room.HandleJoinRoomCallback`) or descriptive names for multi-step flows (e.g., `game.HandleSelectRoomForCreateGame`).
*   **Function Signature:** Typically receive necessary use case handlers, `telebot.Context`, relevant data (payload, callback data), and `*messages.Messages` as arguments. They *no longer* receive `RefreshingMessageBook` instances directly.
*   **Logic:**
    1.  Parse input (payload, callback data).
    2.  Convert sender to `*sharedEntity.User` using `tgutil.ToUser`.
    3.  Perform presentation-level validation/permission checks (though domain checks are preferred).
    4.  Create domain Command/Query structs.
    5.  Call the appropriate injected domain use case handler (`handler.Handle(...)`).
    6.  Process results/errors.
    7.  Prepare response content by calling appropriate **message preparation functions** (e.g., `RoomDetailMessage`), passing necessary context like user admin status if required by the preparation function.
    8.  Send responses/edit messages using `c.Send()`, `c.Edit()`, `c.Respond()`.
    9.  If state relevant to a dynamic message was changed:
        *   Obtain the relevant `RefreshingMessageBook` using the appropriate `h.Get...` or `h.GetOrCreate...` method.
        *   Call `book.RaiseRefreshNeeded()`. **Note:** May need to trigger refresh on multiple books if the user's view changes (e.g., moving from list to detail).
    10. If a *new* dynamic message is sent (one that needs refreshing):
        *   Obtain the relevant `RefreshingMessageBook` using `h.GetOrCreate...`.
        *   Create the `tgutil.RefreshingMessage` struct.
        *   Add it to the book using `book.AddActiveMessage(chatID, refreshMsg)`. **Note:** Ensure the message reference is removed from any *previous* book the user was viewing (e.g., remove from `roomList` when adding to `roomDetail`). Handlers might also need to explicitly delete the previous Telegram message.
    11. If a dynamic message becomes invalid or finalized:
        *   Obtain the relevant `RefreshingMessageBook` using `h.Get...`.
        *   Remove it using `book.RemoveActiveMessage(chatID)` or trigger book deletion via `h.Delete...`.
*   **Specific Handlers:**
    *   `room.HandleCreateRoom`: Handles the `/create_room` command. Requires admin privileges. Parses name, converts sender to `User`, calls `CreateRoomHandler` use case (passing the User), triggers refresh, and sends success message.
    *   `room.HandleJoinRoomCallback`: Handles the join button press. Calls `JoinRoom` use case. **Updates refresh state by removing the user's message from the `roomList` book and adding it to the `roomDetail` book.** Triggers refresh for both books. Calls `room.RoomDetailMessage`, and edits the message using `ModeMarkdownV2`.
    *   `room.HandleLeaveRoomSelectCallback`: Handles the leave button press. Calls `LeaveRoom` use case. **Updates refresh state by removing the user's message from the `roomDetail` book and adding it back to the `roomList` book.** Triggers refresh for both books. Edits message back to the room list view.
    *   `room.HandleKickUserSelectCallback`: Callback handler for the admin "Kick User" button. Fetches players in the room, displays a message (`msgs.Room.KickUserSelectPrompt`) with each player as a button (including the admin). Button payload includes roomID and userIDToKick, unique is `UniqueKickUserConfirm`. Includes a cancel button.
    *   `room.HandleKickUserConfirmCallback`: Callback handler when an admin selects a user to kick. Parses payload, calls `KickUser` use case, triggers refreshes for room list and detail, responds with success (`msgs.Room.KickUserCallbackSuccess`), and edits the message back to the standard room detail view.
    *   `room.HandleChangeModeratorSelectCallback`: Callback handler for the admin "Change Moderator" button. Fetches players, displays message (`msgs.Room.ChangeModeratorSelectPrompt`) with each player as a button (including the current moderator). Button payload includes roomID and userIDToMakeMod, unique is `UniqueChangeModeratorConfirm`. Includes a cancel button.
    *   `room.HandleChangeModeratorConfirmCallback`: Callback handler when an admin selects a new moderator. Parses payload, fetches target user details, calls `ChangeModerator` use case (which updates the moderator and manages player lists), triggers refreshes, responds with success (`msgs.Room.ChangeModeratorCallbackSuccess`), and edits message back to room detail view.
    *   `game.HandleCreateGame`: Initiates the interactive game creation flow. Fetches rooms and filters them based on permissions (global admins see all, room moderators see only their rooms). Sends a message prompting for room selection.
    *   `game.HandleSelectRoomForCreateGame`: Callback handler for room selection. Sends a message prompting for scenario selection.
    *   `game.HandleSelectScenarioForCreateGame`: Callback handler for scenario selection. Fetches data, presents confirmation message (`msgs.Game.CreateGameConfirmPrompt` with role list), with "Start" and "Cancel" buttons. Uses `ModeMarkdownV2`.
    *   `game.HandleStartGameCallback`: Callback handler for the "Start" button. Calls `AssignRoles`, sends role info via PM (using `msgs.Game.AssignRolesSuccessPrivate` with role name and side, using `ModeMarkdownV2`), then edits the original message (`msgs.Game.CreateGameStartedSuccess` with user profile links and roles). Uses `ModeMarkdownV2`.
    *   `game.HandleCancelGameCreationCallback`: Callback handler for the "Cancel" button. Edits message to show cancellation.
    *   `game.HandleAssignRoles`: Handles the `/assign_roles` command (likely becoming obsolete due to the interactive flow). Iterates over the returned assignment map (`map[User]Role`) and sends private messages.
    *   `game.HandleChooseCardStart`: (UPDATED) Gets/Creates admin and player books. Calls `PrepareAdminAssignmentMessage` *before* sending/editing to get initial content. Adds sent messages to the respective books using `AddActiveMessage`.
    *   `game.HandlePlayerSelectsCard`: (UPDATED) Gets/Creates admin and player books. Calls `RaiseRefreshNeeded()` on both. Removes selecting player's message using `RemoveActiveMessage`. If `allSelected`, calls `h.RefreshMessages(adminRefresher)` directly before cleaning up books with `Delete...` methods.

## 5. `handler/document_handler.go`

*   **`handleDocument` (Dispatcher on `BotHandler`):** Routes `telebot.OnDocument` events here.
*   **`HandleDocument` (Exported Function):**
    *   Checks if the document is JSON.
    *   Performs admin check.
    *   Downloads the file content using `c.Bot().File()`.
    *   Calls the `AddScenarioJSONHandler` use case.
    *   Sends success or error messages back to the user.

## 6. `messages/`

*   **`messages.json` (Root directory):** Contains user-facing strings. Includes keys for the kick flow and change moderator flow (`ChangeModeratorButton`, `ChangeModeratorSelectPrompt`, `ChangeModeratorCallbackSuccess`, `ChangeModeratorCallbackError`, `ChangeModeratorNoCandidates`). Text updated for various flows.
*   **`messages.go`:** Defines the Go struct mirroring `messages.json`.
*   **`loader.go`:** Loads messages from JSON.
*   **Usage:** Injected `*Messages` struct used throughout handlers.

### `shared/tgutil/`

*   **`refresh_state.go` (`RefreshingMessageBook`):** (UPDATED)
    *   Now stores the `GetMessage func(...)` responsible for generating its specific content update.
    *   Passed into the `NewRefreshState` constructor.
*   **`callback_data.go`:** Defines constants for callback `Unique` identifiers (`UniqueChooseCardStart`, `UniquePlayerSelectsCard`).
*   **`state.go` (`InteractiveSelectionState`):** (UPDATED) Holds the state for the interactive role selection process for a specific game.
    *   `ShuffledRoles`: The randomized list of roles.
    *   `Selections`: Map of `UserID` to `PlayerSelection` struct (which contains `ChosenIndex` and the `sharedEntity.User`).
    *   `TakenIndices`: Map of card index to boolean (true if taken).
    *   `Mutex`: Ensures thread-safe access during selection.
*   **`state.go` (`PlayerSelection`):** (NEW) Struct holding the `ChosenIndex` and the `sharedEntity.User` who made the selection.

### `messages.json`

*   **New Keys Added under `Game`:**
    *   `ChooseCardButton`: "🃏 Choose Card"
    *   `AssignmentTrackingMessageAdmin`: "Role Selection Progress:\n%s\nWaiting for players..."
    *   `RoleSelectionPromptPlayer`: "Choose your role card:"
    *   `RoleTakenMarker`: "X"
    *   `PlayerHasRoleError`: "You have already selected a role."
    *   `RoleAlreadyTakenError`: "Card %d has already been taken!"
    *   `AssignRolesSuccessPrivate`: "Your role: *%s* \(%s\)"
    *   `AllRolesSelectedAdmin`: "All roles selected!\n%s"

## 7. Refreshing Message Implementation Rule (NEW)

*   **Rule:** Any message in the Telegram presentation layer that needs to be dynamically updated based on state changes **MUST** use the `tgutil.RefreshingMessageBook` pattern.
*   **Implementation Steps:** (UPDATED)
    1.  **Identify Scope:** Determine the scope (e.g., per-room, per-game). Create a map in `BotHandler` (e.g., `adminAssignmentTrackers map[gameEntity.GameID]*tgutil.RefreshingMessageBook`) to hold `RefreshingMessageBook` instances, keyed by the scope identifier. Protect this map with a `sync.RWMutex`.
    2.  **Define Message Generation Function:** Create an **exported message preparation function** (e.g., `game.PrepareAdminAssignmentMessage`) that fetches necessary data and returns `(string, []interface{}, error)`.
    3.  **Manage Books:** Implement `Get/GetOrCreate/Delete` helper methods on `BotHandler` for the map. The `GetOrCreate...` method **MUST** call `tgutil.NewRefreshState`, passing the specific message generation function defined in step 2.
    4.  **Store Messages:** When sending/editing the initial refreshing message, get the book (`GetOrCreate...`), create a `tgutil.RefreshingMessage{ChatID, MessageID, Data}`, and store it using `book.AddActiveMessage(chatID, refreshMsg)`.
    5.  **Trigger Refresh:** Handlers/callbacks modifying state **MUST** get the relevant book (`Get...` or `GetOrCreate...`) and call `book.RaiseRefreshNeeded()`.
    6.  **Implement Refresh Logic (`refresh.go`):**
        *   In `StartRefreshTimer()`, add a loop over the map in `BotHandler` (using mutex).
        *   Check `book.ConsumeRefreshNeeded()`.
        *   If true, call the *single* generic refresh execution method: `h.RefreshMessages(book)`.
    7.  **Implement Refresh Execution (`refresh.go`):**
        *   The `RefreshMessages(book)` method iterates through `book.GetAllActiveMessages()`.
        *   For each message, it calls `book.GetMessage(chatID, payload.Data)` to get the new content/opts.
        *   It uses `h.bot.Edit()` to apply the update.
    8.  **Cleanup:** (Same as previous rule: Use `Delete...` for book scope, `RemoveActiveMessage` for specific finalized messages). 