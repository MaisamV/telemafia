# Blueprint 02: Shared Components

**Source:** `internal/shared/`

**Purpose:** This document describes components intended for use across different layers or domain modules, ensuring consistency and avoiding code duplication while adhering to dependency rules.

## 1. `entity/user.go`

*   **`UserID` (type `int64`):** Unique identifier for a user, typically derived from the Telegram User ID.
*   **`User` struct:**
    *   Fields: `ID` (UserID), `TelegramID` (int64), `FirstName`, `LastName`, `Username` (string), `Admin` (bool).
    *   Represents a user interacting with the bot.
    *   The `Admin` flag is determined during user conversion (e.g., in `tgutil.ToUser`) based on the configured admin usernames.
*   **Methods:**
    *   `CanCreateRoom() bool`: Checks if the user is an admin (used for domain-level authorization, though potentially obsolete).
    *   `GetProfileLink() string`: Returns a Markdown-formatted link to the user's Telegram profile (e.g., `[FirstName](tg://user?id=...)` or `[FirstName](https://t.me/username)`).

## 2. `event/`

*   **`event.go`:**
    *   **`Event` interface:** Defines the basic contract for domain events (`EventName() string`, `OccurredAt() time.Time`).
    *   **`Publisher` interface:** Defines the contract for publishing events (`Publish(event Event) error`). The current implementation in `main.go` simply logs events.
*   **`events.go`:**
    *   Defines concrete event structs (e.g., `RoomCreatedEvent`, `PlayerJoinedEvent`, `PlayerLeftEvent`, `PlayerKickedEvent`).
    *   Each struct embeds the necessary data for the event (IDs, names, timestamps).
    *   Each struct implements the `Event` interface.
    *   Used by command handlers to signal significant domain state changes.

## 3. `tgutil/` (Telegram Utilities)

*   **`const.go`:**
    *   Defines string constants for unique callback query identifiers (e.g., `UniqueJoinRoom`, `UniqueCreateGameSelectRoom`, `UniqueKickUserSelect`, `UniqueKickUserConfirm`). Used for creating inline buttons and routing callbacks in `handler/callbacks.go`.
*   **`refresh_state.go`:**
    *   **`RefreshingMessageBook` struct:** Manages the state for dynamic message updates (like the room list).
        *   Tracks active messages per chat ID (`activeMessages map[int64]*RefreshingMessage`).
        *   Uses a `needsRefresh` flag (protected by mutex) to signal when updates are required.
        *   Provides methods: `RaiseRefreshNeeded()`, `ConsumeRefreshNeeded()`, `AddActiveMessage()`, `RemoveActiveMessage()`, `GetAllActiveMessages()`.
    *   **`RefreshNotifier` interface (e.g., in `handler/room/callbacks_room.go`):** Implemented by `RefreshingMessageBook` to allow handlers to trigger refreshes.
*   **`util.go`:**
    *   `SetAdminUsers()`: Stores the list of admin usernames (used by `IsAdmin`).
    *   `IsAdmin(username string) bool`: Checks if a username is in the admin list (case-insensitive).
    *   `ToUser(sender *telebot.User) *sharedEntity.User`: Converts a `telebot.User` to the shared `entity.User`, setting the `Admin` flag based on `IsAdmin`.
    *   `SplitCallbackData(data string) (unique string, payload string)`: Parses callback data strings (format: `unique|payload`).

## 4. `common/utils.go`

*   Contains general utility functions not specific to any layer or domain.
*   **`StringToInt64(s string) (int64, error)`:** Converts a string to int64.
*   **`ContainsString(slice []string, s string) bool`:** Checks if a string slice contains a specific string.
*   **`GenerateHash(s string) string`:** Generates a SHA256 hash of a string.
*   **`InitSeed()`:** Initializes the shared `math/rand` random number generator using the current time. Should be called once at application startup.
*   **`Shuffle(n int, swap func(i, j int))`:** Shuffles a sequence of length `n` using the provided swap function and the initialized random number generator. Used for randomizing role assignments.

## 5. `logger/`

*   Placeholder for a potential shared logging setup (currently standard `log` package is used). 