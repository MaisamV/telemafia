# 10. Shared Components Specification

**Goal:** Define reusable components like shared entities, common utilities, and the event system structure.

## 10.1. Shared Entity: User (`internal/shared/entity/user.go`)

*   **Purpose:** Represents a unique user across the application.
*   **Definition:**
    ```go
    package entity

    type UserID int64

    type User struct {
        ID         UserID
        TelegramID int64 // Telegram's unique ID
        FirstName  string
        LastName   string
        Username   string // Telegram username (case-sensitive)
        Admin      bool   // Determined at runtime based on config
    }
    ```
*   **Notes:** The `Admin` field is populated by the presentation layer based on configuration, not stored persistently.

## 10.2. Shared Common Utilities (`internal/shared/common/utils.go`)

*   **Purpose:** Provide simple, generic helper functions.
*   **Required Functions:**
    *   `Contains(slice []string, str string) bool`: Checks if a string exists in a slice of strings (case-sensitive).
    *   `StringToInt64(s string) (int64, error)`: Safely converts a string to `int64`.
    *   *(Other general-purpose helpers can be added here as needed)*
*   **Note:** This package might be less used now that Telegram-specific utilities are in `tgutil`.

## 10.3. Shared Event System (`internal/shared/event/`)

*   **Purpose:** Define a simple mechanism for domain events and publishing.
*   **`event.go`:**
    *   `Event` interface:
        ```go
        type Event interface {
            EventName() string
            OccurredAt() time.Time
        }
        ```
    *   `Publisher` interface:
        ```go
        type Publisher interface {
            Publish(event Event) error
        }
        ```
*   **`events.go`:**
    *   Define concrete event structs for significant domain occurrences. Each struct **MUST** implement the `Event` interface.
    *   **Required Events:**
        *   `RoomCreatedEvent { RoomID roomEntity.RoomID; Name string; CreatedAt time.Time; ScenarioName string }` (EventName: "room.created")
        *   `PlayerJoinedEvent { RoomID roomEntity.RoomID; PlayerID sharedEntity.UserID; RoomName string; JoinedAt time.Time }` (EventName: "room.player_joined")
        *   `PlayerLeftEvent { RoomID roomEntity.RoomID; PlayerID sharedEntity.UserID; LeftAt time.Time }` (EventName: "room.player_left")
        *   `PlayerKickedEvent { RoomID roomEntity.RoomID; PlayerID sharedEntity.UserID; KickedAt time.Time }` (EventName: "room.player_kicked")
        *   *(Add other events like `ScenarioCreated`, `GameCreated`, `RolesAssigned` if needed for future integrations or logging)*
*   **Implementation:**
    *   A simple implementation of the `Publisher` interface should be provided in `cmd/telemafia/main.go`.
    *   For this version, the `Publish` method can simply log the event details using the standard `log` package (e.g., `log.Printf("Event published: Type=%s, Data=%+v\n", e.EventName(), e)`).
    *   This implementation should be injected into command handlers that need to publish events.

## 10.4. Shared Telegram Utilities (`internal/shared/tgutil/`)

*   **Purpose:** Consolidate Telegram-specific constants and helper functions used by the presentation layer.
*   **`const.go`:**
    *   Defines `Unique...` constants for callback query identifiers.
*   **`util.go`:**
    *   Defines helper functions:
        *   `SetAdminUsers(usernames []string)`: Stores admin usernames (package-level variable).
        *   `IsAdmin(username string) bool`: Checks if a username is in the admin list.
        *   `ToUser(sender *telebot.User) *sharedEntity.User`: Converts `telebot.User` to `sharedEntity.User`, checking admin status.
        *   `SplitCallbackData(data string) (unique string, payload string)`: Parses callback data.
    *   **`refresh_state.go`:**
        *   Defines `RefreshState` struct with mutex, boolean flag (`needsRefresh`), and map (`activeMessages`) to manage dynamic message updates.
        *   Provides methods like `NewRefreshState`, `RaiseRefreshNeeded`, `ConsumeRefreshNeeded`, `AddActiveMessage`, `RemoveActiveMessage`, `GetAllActiveMessages`. 