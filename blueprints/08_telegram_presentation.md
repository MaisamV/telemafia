# Blueprint 08: Telegram Presentation Layer

**Source:** `internal/presentation/telegram/`

**Purpose:** Details the components responsible for handling Telegram interactions and translating them into calls to domain use cases.

## 1. `handler/bot_handler.go` (`BotHandler`)

*   **Role:** Acts as the central orchestrator for the Telegram presentation layer.
*   **Dependencies:** Receives the `telebot.Bot` instance, admin usernames, loaded `messages.Messages`, and *all* domain use case handlers (commands and queries) via constructor injection (`NewBotHandler`).
*   **Refresh State:** Holds instances of `tgutil.RefreshingMessageBook` (e.g., `roomListRefreshMessage`, `roomDetailRefreshMessage`) to manage dynamic message updates.
*   **`RegisterHandlers()`:** Maps Telegram commands (e.g., `/create_room`) and events (e.g., `telebot.OnCallback`, `telebot.OnDocument`) to specific *dispatcher methods* on the `BotHandler` struct (e.g., `h.handleCreateRoom`).
*   **Dispatcher Methods (e.g., `handleCreateRoom`, `handleCallback`, `handleDocument`):** These methods primarily act as routers. They parse necessary information from the `telebot.Context` (like payload, sender) and call the corresponding **exported handler functions** located in sub-packages (e.g., `room.HandleCreateRoom`, `scenario.HandleAddScenarioJSON`) or other handler files (`callbacks.go`, `document_handler.go`), passing the required dependencies (use case handlers, notifiers, messages) from the `BotHandler` instance (`h`).
*   **`Start()`:** Initializes background tasks (like the refresh timer) and starts the `telebot` polling loop.

## 2. `handler/callbacks.go` (`BotHandler.handleCallback`)

*   **Purpose:** Central dispatcher for *all* inline button callback queries.
*   **Logic:**
    1.  Retrieves the callback data.
    2.  Uses `tgutil.SplitCallbackData` to separate the `unique` identifier and the `payload`.
    3.  Uses a `switch` statement on the `unique` identifier.
    4.  Each `case` corresponds to a specific button type (defined in `tgutil/const.go`).
    5.  Calls the relevant **exported handler function** from the appropriate sub-package (e.g., `room.HandleJoinRoomCallback`, `game.HandleSelectRoomForCreateGame`), passing the `telebot.Context`, parsed `payload`, and necessary dependencies (use case handlers, notifiers, messages) obtained from the `BotHandler` (`h`).

## 3. `handler/refresh.go`

*   **Purpose:** Handles the background task for updating dynamic messages.
*   **`BotHandler.StartRefreshTimer()`:** Runs a goroutine with a `time.Ticker` (e.g., every 5 seconds).
*   **Timer Loop:**
    1.  Checks if a refresh is needed for each managed `RefreshingMessageBook` instance (e.g., `h.roomListRefreshMessage`) using `ConsumeRefreshNeeded()`.
    2.  If needed, calls `h.updateMessages()` for that book.
*   **`BotHandler.updateMessages()`:**
    1.  Gets all active messages tracked by the specific `RefreshingMessageBook`.
    2.  For each message, calls a **message preparation function** (e.g., `room.PrepareRoomListMessage`) to get the latest content and markup.
    3.  Uses `h.bot.Edit()` to update the message in Telegram.
    4.  Handles errors (e.g., message not found, user blocked bot) and removes the message from tracking if necessary.
*   **`BotHandler.SendOrUpdateRefreshingMessage()`:** Utility used by handlers (like `HandleListRooms`) to either send a new dynamic message and track it, or edit an existing tracked message.
*   **Message Preparation Functions (e.g., `room.PrepareRoomListMessage`, `room.RoomDetailMessage`):** Exported functions (located in handler sub-packages like `room/`) responsible for fetching current data (using injected query handlers) and formatting the message text and `telebot.ReplyMarkup`.

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
    7.  Send responses/edit messages using `c.Send()`, `c.Edit()`, `c.Respond()`, utilizing the injected `msgs` struct for all user-facing text.
    8.  If state relevant to a dynamic message was changed, call `notifier.RaiseRefreshNeeded()`.
    9.  If a new dynamic message is sent, track it using `notifier.AddActiveMessage()`.
    10. If a dynamic message becomes invalid, untrack it using `notifier.RemoveActiveMessage()`.

## 5. `handler/document_handler.go`

*   **`handleDocument` (Dispatcher on `BotHandler`):** Routes `telebot.OnDocument` events here.
*   **`HandleDocument` (Exported Function):**
    *   Checks if the document is JSON.
    *   Performs admin check.
    *   Downloads the file content using `c.Bot().File()`.
    *   Calls the `AddScenarioJSONHandler` use case.
    *   Sends success or error messages back to the user.

## 6. `messages/`

*   **`messages.json` (Root directory):** Contains all user-facing strings organized by category (common, room, game, etc.) using nested JSON objects.
*   **`messages.go`:** Defines the Go struct (`Messages` and nested structs like `RoomMessages`, `GameMessages`, etc.) that mirrors the structure of `messages.json`. Uses `json` tags for unmarshalling.
*   **`loader.go`:**
    *   **`LoadMessages(filename string) (*Messages, error)`:** Reads the specified JSON file, unmarshals the data into the `Messages` struct, performs optional basic validation, and returns the populated struct.
    *   Called once during startup in `main.go`.
*   **Usage:** The loaded `*Messages` struct is injected into the `BotHandler` and passed down to the exported handler functions, which access strings via struct fields (e.g., `msgs.Room.CreateSuccess`, `fmt.Sprintf(msgs.Common.ErrorGeneric, err)`). 