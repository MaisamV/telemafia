# 5. Telegram Presentation Layer Rules

**Goal:** Ensure consistency and correct implementation within the Telegram adapter.

## 5.1. Handler Structure

*   **Main Struct:** `internal/presentation/telegram/handler/bot_handler.go` defines `BotHandler` which holds dependencies (Bot instance, Use Case Handlers, Messages, Refresh Notifiers).
*   **Dispatchers:** Methods on `BotHandler` (e.g., `handleCreateRoom`) map Telegram commands/callbacks to handler logic. They are registered in `RegisterHandlers`.
*   **Logic:** Actual handling logic resides in **exported functions** within sub-packages (`room/`, `game/`, `scenario/`) or `common_handlers.go`.
*   **Dependencies:** Handler functions receive dependencies (use case handlers, messages, notifiers) directly as arguments.

## 5.2. Command Handling

1.  **Define Exported Function:** Create `HandleCommandName` in the appropriate handler package (e.g., `room.HandleMyRooms`).
2.  **Add Dispatcher Method:** Add `handleCommandName` method to `BotHandler` that calls the exported function, passing necessary dependencies from `h`.
3.  **Register:** Add `h.bot.Handle("/command_name", h.handleCommandName)` in `BotHandler.RegisterHandlers`.
4.  **Implementation:**
    *   Parse payload (`c.Message().Payload`).
    *   Convert sender (`c.Sender()`) to `*sharedEntity.User` using `tgutil.ToUser()`.
    *   Perform presentation-level validation/checks (e.g., admin status if required here, payload format).
    *   Create the domain Command/Query struct.
    *   Call the appropriate domain Use Case handler.
    *   Handle results/errors, sending responses using `msgs` struct (`c.Send(...)`).
    *   If state was changed that affects a dynamic message, call `notifier.RaiseRefreshNeeded()`.

## 5.3. Callback Handling

1.  **Define Unique Constant:** Add `UniqueCallbackName` in `internal/shared/tgutil/const.go`.
2.  **Define Exported Function:** Create `HandleCallbackNameCallback` in the appropriate handler package (e.g., `room.HandleJoinRoomCallback`).
3.  **Add Routing Case:** Add `case tgutil.UniqueCallbackName:` to the `switch` in `handleCallback` method (`callbacks.go`), calling the exported function.
    *   **Note:** The `handleCallback` method itself acts solely as a dispatcher. It uses the unique identifier from the callback data to route the request to the appropriate exported handler function (defined in step 2) where the actual processing logic resides.
4.  **Implementation:**
    *   Parse `data` passed into the function (originally from `tgutil.SplitCallbackData`).
    *   Convert sender (`c.Sender()`) if needed.
    *   Perform logic, potentially calling Use Case handlers.
    *   Acknowledge the callback using `c.Respond()` (use `msgs` for text).
    *   Update the original message using `c.Edit()` or `c.Delete()` (use `msgs` for text).
    *   If state was changed that affects a dynamic message, call `notifier.RaiseRefreshNeeded()`.
*   **Callback Data Format:** Use `unique|payload` format. Generate data using `fmt.Sprintf("%s|%s", tgutil.UniqueCallbackName, payload)`.

## 5.4. Dynamic Message Refreshing (Rule)

*   **Rule:** Any message in the Telegram presentation layer that needs to be dynamically updated based on state changes **MUST** use the `tgutil.RefreshingMessageBook` pattern.
*   **Mechanism:** The `RefreshingMessageBook` struct now encapsulates the logic needed to generate its own updates. It holds a `GetMessage` function provided during its creation.
*   **Implementation Steps:**
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
        *   For each message, it calls the book's own generator: `book.GetMessage(chatID, payload.Data)` to get the new content/opts.
        *   It uses `h.bot.Edit()` to apply the update.
    8.  **Cleanup:**
        *   Use the `Delete...` book helper method when the message scope is no longer valid (e.g., game ends).
        *   Use `book.RemoveActiveMessage(chatID)` when a *specific* user's message in the book is finalized or invalidated (e.g., player confirms role choice).

## 5.5. Messages

*   **MUST** use the injected `*messages.Messages` struct for all user-facing text.
*   Refer to `messages.json` for available keys. 