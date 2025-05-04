# Blueprint 08: Telegram Presentation Layer

**Source:** `internal/presentation/telegram/`

**Purpose:** Details the components responsible for handling Telegram interactions and translating them into calls to domain use cases.

## 1. `handler/bot_handler.go` (`BotHandler`)

*   **Role:** Acts as the central orchestrator for the Telegram presentation layer.
*   **Dependencies:** Receives the `telebot.Bot` instance, admin usernames, loaded `messages.Messages`, and *all* domain use case handlers (commands and queries) via constructor injection (`NewBotHandler`).
*   **Refresh State:** Holds instances of `tgutil.RefreshingMessageBook` (e.g., `roomListRefreshMessage`, `roomDetailRefreshMessage`) to manage dynamic message updates.
*   **`RegisterHandlers()`:** Maps Telegram commands (e.g., `/create_room`, `/create_game`) and events (e.g., `telebot.OnCallback`, `telebot.OnDocument`) to specific *dispatcher methods* on the `BotHandler` struct (e.g., `h.handleCreateRoom`, `h.handleCreateGame`).
*   **Dispatcher Methods (e.g., `handleCreateRoom`, `handleCallback`, `handleDocument`, `handleCreateGame`):** These methods primarily act as routers. They parse necessary information from the `telebot.Context` (like payload, sender) and call the corresponding **exported handler functions** located in sub-packages (e.g., `room.HandleCreateRoom`, `game.HandleCreateGame`) or other handler files (`callbacks.go`, `document_handler.go`), passing the required dependencies (use case handlers, notifiers, messages) from the `BotHandler` instance (`h`).
*   **`Start()`:** Initializes background tasks (like the refresh timer) and starts the `telebot` polling loop.

## 2. `handler/callbacks.go` (`BotHandler.handleCallback`)

*   **Purpose:** Central dispatcher for *all* inline button callback queries.
*   **Logic:**
    1.  Retrieves the callback data.
    2.  Uses `tgutil.SplitCallbackData` to separate the `unique` identifier and the `payload`.
    3.  Uses a `switch` statement on the `unique` identifier.
    4.  Each `case` corresponds to a specific button type (defined in `tgutil/const.go`, e.g., `tgutil.UniqueCreateGameSelectRoom`, `tgutil.UniqueKickUserSelect`, `tgutil.UniqueKickUserConfirm`).
    5.  Calls the relevant **exported handler function** from the appropriate sub-package (e.g., `room.HandleJoinRoomCallback`, `room.HandleKickUserSelectCallback`), passing the `telebot.Context`, parsed `payload`, and necessary dependencies.

## 3. `handler/refresh.go`

*   **Purpose:** Handles the background task for updating dynamic messages.
*   **`BotHandler.StartRefreshTimer()`:** Runs a goroutine with a `time.Ticker` (e.g., every 5 seconds).
*   **Timer Loop:**
    1.  Checks if a refresh is needed for each managed `RefreshingMessageBook` instance (e.g., `h.roomListRefreshMessage`) using `ConsumeRefreshNeeded()`.
    2.  If needed, calls `h.updateMessages()` for that book.
*   **`BotHandler.updateMessages()`:**
    1.  Gets all active messages tracked by the specific `RefreshingMessageBook`.
    2.  For each message, calls a **message preparation function** (e.g., `room.PrepareRoomListMessage`, `room.RoomDetailMessage`) to get the latest content and markup.
        *   **Note:** As message preparation functions may require user context (e.g., admin status) for conditional elements, the refresh mechanism may not display these elements correctly since it lacks the original user context.
    3.  Uses `h.bot.Edit()` to update the message in Telegram (using `ModeMarkdownV2`).
    4.  Handles errors (e.g., message not found, user blocked bot) and removes the message from tracking if necessary.
*   **`BotHandler.SendOrUpdateRefreshingMessage()`:** Utility used by handlers (like `HandleListRooms`) to either send a new dynamic message and track it, or edit an existing tracked message.
*   **Message Preparation Functions (e.g., `room.PrepareRoomListMessage`, `room.RoomDetailMessage`):** Exported functions (located in handler sub-packages like `room/`) responsible for fetching current data (using injected query handlers) and formatting the message text and `telebot.ReplyMarkup`.
    *   These functions may accept additional parameters, such as the requesting user's admin status (as seen in `RoomDetailMessage`), to render conditional UI elements like admin-only buttons.
    *   `RoomDetailMessage`: Now always uses the `msgs.Room.RoomDetail` format string, formats the player list using `user.GetProfileLink()`, and includes admin-only buttons for "Start Game" and "Kick User".

## 4. `handler/<module>/` (e.g., `handler/room/`, `handler/game/`)

*   Contains the **exported handler functions** that implement the actual logic for specific commands and callbacks.
*   **Naming Convention:**
    *   Commands: `HandleCommandName` (e.g., `room.HandleCreateRoom`).
    *   Callbacks: `HandleCallbackNameCallback` (e.g., `room.HandleJoinRoomCallback`) or descriptive names for multi-step flows (e.g., `game.HandleSelectRoomForCreateGame`).
*   **Function Signature:** Typically receive necessary use case handlers, notifiers (like `RefreshNotifier`), `telebot.Context`, relevant data (payload, callback data), and `*messages.Messages` as arguments.
*   **Logic:**
    1.  Parse input (payload, callback data).
    2.  Convert sender to `*sharedEntity.User` using `tgutil.ToUser`.
    3.  Perform presentation-level validation/permission checks (though domain checks are preferred).
    4.  Create domain Command/Query structs.
    5.  Call the appropriate injected domain use case handler (`handler.Handle(...)`).
    6.  Process results/errors.
    7.  Prepare response content by calling appropriate **message preparation functions** (e.g., `RoomDetailMessage`), passing necessary context like user admin status if required by the preparation function.
    8.  Send responses/edit messages using `c.Send()`, `c.Edit()`, `c.Respond()` (typically with `telebot.ModeMarkdownV2`), utilizing the injected `msgs` struct for text and the markup from preparation functions.
    9.  If state relevant to a dynamic message was changed, call `notifier.RaiseRefreshNeeded()`.
    10. If a new dynamic message is sent, track it using `notifier.AddActiveMessage()`.
    11. If a dynamic message becomes invalid, untrack it using `notifier.RemoveActiveMessage()`.
*   **Specific Handlers:**
    *   `room.HandleJoinRoomCallback`: Handles the join button press. Calls `JoinRoom` use case, updates refresh state, calls `room.RoomDetailMessage`, and edits the message using `ModeMarkdownV2`.
    *   `room.HandleKickUserSelectCallback`: Callback handler for the admin "Kick User" button. Fetches players in the room (excluding the admin), displays a message (`msgs.Room.KickUserSelectPrompt`) with each player as a button. Button payload includes roomID and userIDToKick, unique is `UniqueKickUserConfirm`. Includes a cancel button.
    *   `room.HandleKickUserConfirmCallback`: Callback handler when an admin selects a user to kick. Parses payload, calls `KickUser` use case, triggers refreshes for room list and detail, responds with success (`msgs.Room.KickUserCallbackSuccess`), and edits the message back to the standard room detail view.
    *   `game.HandleCreateGame`: Initiates the interactive game creation flow (admin only). Sends a message prompting for room selection.
    *   `game.HandleSelectRoomForCreateGame`: Callback handler for room selection. Sends a message prompting for scenario selection.
    *   `game.HandleSelectScenarioForCreateGame`: Callback handler for scenario selection. Fetches data, presents confirmation message (`msgs.Game.CreateGameConfirmPrompt` with role list), with "Start" and "Cancel" buttons. Uses `ModeMarkdownV2`.
    *   `game.HandleStartGameCallback`: Callback handler for the "Start" button. Calls `AssignRoles`, sends role info via PM (using `msgs.Game.AssignRolesSuccessPrivate` with role name and side, using `ModeMarkdownV2`), then edits the original message (`msgs.Game.CreateGameStartedSuccess` with user profile links and roles). Uses `ModeMarkdownV2`.
    *   `game.HandleCancelGameCreationCallback`: Callback handler for the "Cancel" button. Edits message to show cancellation.
    *   `game.HandleAssignRoles`: Handles the `/assign_roles` command (likely becoming obsolete due to the interactive flow). Iterates over the returned assignment map (`map[User]Role`) and sends private messages.

## 5. `handler/document_handler.go`

*   **`handleDocument` (Dispatcher on `BotHandler`):** Routes `telebot.OnDocument` events here.
*   **`HandleDocument` (Exported Function):**
    *   Checks if the document is JSON.
    *   Performs admin check.
    *   Downloads the file content using `c.Bot().File()`.
    *   Calls the `AddScenarioJSONHandler` use case.
    *   Sends success or error messages back to the user.

## 6. `messages/`

*   **`messages.json` (Root directory):** Contains user-facing strings. Includes keys for the interactive kick flow (`KickUserButton`, `KickUserSelectPrompt`, `KickUserCallbackSuccess`, `KickUserCallbackError`, `KickUserNoPlayers`). Text for game creation flow, role assignment PMs, and room details updated (primarily Farsi translations and formatting, including MarkdownV2 syntax like `||` for spoilers and escaped characters `\\`).
*   **`messages.go`:** Defines the Go struct mirroring `messages.json`.
*   **`loader.go`:** Loads messages from JSON.
*   **Usage:** Injected `*Messages` struct used throughout handlers. 