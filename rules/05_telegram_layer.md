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

## 5.4. Dynamic Message Refreshing

*   **Mechanism:** Uses `RefreshingMessageBook` from `internal/shared/tgutil/refresh_state.go`.
*   **State:** `BotHandler` holds `RefreshingMessageBook` instances.
*   **Triggering:** Command/Callback handlers that modify state relevant to a dynamic message **MUST** call `RaiseRefreshNeeded()` on the appropriate `RefreshingMessageBook` instance (passed via the `RefreshNotifier` interface).
*   **Updating:** The background timer in `refresh.go` periodically calls message preparation functions (e.g., `room.PrepareRoomListMessage`, `room.RoomDetailMessage`) and uses `bot.Edit` to update tracked messages.
*   **Message Preparation:** Functions like `PrepareRoomListMessage` are responsible for fetching current data and formatting the message text and `telebot.ReplyMarkup`.
*   **Tracking:** Handlers that send *new* dynamic messages should add them to the relevant `RefreshingMessageBook` using `AddActiveMessage`.
*   **Untracking:** Handlers or callbacks that invalidate a dynamic message (e.g., leaving a room makes the room detail message obsolete) should remove it using `RemoveActiveMessage`.

## 5.5. Messages

*   **MUST** use the injected `*messages.Messages` struct for all user-facing text.
*   Refer to `messages.json` for available keys. 