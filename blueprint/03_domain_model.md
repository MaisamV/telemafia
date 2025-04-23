# 3. Domain Model Specification

**Goal:** Define the core business entities, their attributes, value objects, and key relationships.

## 3.1. Shared Entities & Value Objects (`internal/shared/...`)

*   **`entity.UserID`**
    *   **Type:** `int64`
    *   **Description:** Unique identifier for a user.

*   **`entity.User`**
    *   **Description:** Represents a user interacting with the bot.
    *   **Attributes:**
        *   `ID UserID`: Primary identifier.
        *   `TelegramID int64`: The user's Telegram-specific ID.
        *   `FirstName string`: User's first name from Telegram.
        *   `LastName string`: User's last name from Telegram.
        *   `Username string`: User's Telegram username (case-sensitive).
        *   `Admin bool`: Flag indicating if the user is an administrator (derived from configuration at runtime).

*   **`event.Event`**
    *   **Type:** `interface`
    *   **Description:** Base interface for domain events.
    *   **Methods:**
        *   `EventName() string`: Returns the unique name of the event (e.g., "room.created").
        *   `OccurredAt() time.Time`: Returns the timestamp when the event occurred.

*   **`event.Publisher`**
    *   **Type:** `interface`
    *   **Description:** Interface for publishing domain events.
    *   **Methods:**
        *   `Publish(event Event) error`: Publishes the given event.

## 3.2. Room Domain (`internal/domain/room/...`)

*   **`entity.RoomID`**
    *   **Type:** `string`
    *   **Description:** Unique identifier for a game room (e.g., "room_1678886400").

*   **`entity.Room`**
    *   **Description:** Represents a game room where users gather.
    *   **Attributes:**
        *   `ID RoomID`: Primary identifier.
        *  \
        * `Name string`: Display name of the room (validation: 3-50 chars).
        *   `CreatedAt time.Time`: Timestamp of room creation.
        *   `Players []*sharedEntity.User`: Slice of pointers to users currently in the room.
        *   `ScenarioName string`: ID of the scenario currently assigned to this room. **Note:** This is primarily for informational/display purposes related to the room itself. The definitive link for an active game is within the `Game` entity.
        *   `Description map[string]string`: Map for storing descriptive texts related to the room (e.g., scenario details via key "scenario_info").
    *   **Key Methods (Conceptual):**
        *   `NewRoom(id RoomID, name string) (*Room, error)`: Constructor with name validation.
        *   `AddPlayer(player *sharedEntity.User)`: Adds a user to the Players slice.
        *   `RemovePlayer(playerID sharedEntity.UserID)`: Removes a user from the Players slice by ID.
        *   `SetDescription(...)`: (Potentially unused) Method to add to Description map.
    *   **Validation:** Room name must be between 3 and 50 characters.
    *   **Errors:** `ErrInvalidRoomName`, `ErrRoomNotFound`, `ErrRoomAlreadyExists`, `ErrPlayerNotInRoom`.

## 3.3. Scenario Domain (`internal/domain/scenario/...`)

*   **`entity.Role`**
    *   **Description:** Represents a role within a game scenario.
    *   **Attributes:**
        *   `Name string`: Name of the role (e.g., "Mafia", "Doctor", "Civilian").

*   **`entity.Scenario`**
    *   **Description:** Defines a set of roles for a specific Mafia game variant.
    *   **Attributes:**
        *   `ID string`: Unique identifier for the scenario (e.g., "scen_1678886400").
        *   `Name string`: Display name of the scenario (e.g., "Classic 7 Player").
        *   `Roles []Role`: Slice containing the roles available in this scenario.

## 3.4. Game Domain (`internal/domain/game/...`)

*   **`entity.GameID`**
    *   **Type:** `string`
    *   **Description:** Unique identifier for an active game instance (e.g., "game_1678886400").

*   **`entity.GameState`**
    *   **Type:** `string` (Enum-like consts)
    *   **Description:** Represents the current state of a game.
    *   **Values:**
        *   `GameStateWaitingForPlayers` ("waiting_for_players")
        *   `GameStateRolesAssigned` ("roles_assigned")
        *   `GameStateInProgress` ("in_progress")
        *   `GameStateFinished` ("finished")

*   **`entity.Game`**
    *   **Description:** Represents an instance of a Mafia game being played or prepared.
    *   **Attributes:**
        *   `ID GameID`: Primary identifier.
        *   `State GameState`: The current state of the game.
        *   `Room *roomEntity.Room`: Pointer to the Room entity where the game takes place.
        *   `Scenario *scenarioEntity.Scenario`: Pointer to the Scenario entity defining the roles.
        *   `Assignments map[sharedEntity.UserID]scenarioEntity.Role`: Map associating User IDs with their assigned Roles for this game instance.
    *   **Key Methods (Conceptual):**
        *   `AssignRole(userID sharedEntity.UserID, role scenarioEntity.Role)`: Adds an entry to the Assignments map.
        *   `SetRolesAssigned()`: Changes State to `GameStateRolesAssigned`.
        *   `StartGame()`: Changes State to `GameStateInProgress`.
        *   `FinishGame()`: Changes State to `GameStateFinished`.

## 3.5. Relationships Summary

*   A `Room` contains zero or more `User`s (Players).
*   A `Room` *can have* an assigned `Scenario` indicated (e.g., via `ScenarioName` or `Description`) for informational purposes.
*   A `Game` **is** associated with exactly one `Room`.
*   A `Game` **uses** exactly one `Scenario`.
*   A `Game` contains zero or more role `Assignments` (mapping `User` to `Role`).
*   A `Scenario` contains one or more `Role`s. 