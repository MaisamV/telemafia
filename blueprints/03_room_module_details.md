# Blueprint 03: Room Module Details

**Source:** `internal/domain/room/`

**Purpose:** Details the components within the Room domain module.

## 1. `entity/room.go`

*   **`RoomID` (type `string`):** Unique identifier for a room.
*   **`Room` struct:**
    *   Fields: `ID` (RoomID), `Name` (string), `CreatedAt` (time.Time), `Players` ([]*sharedEntity.User), `Description` (map[string]string), `ScenarioName` (string).
    *   `Players`: Slice of pointers to shared User entities currently in the room.
    *   `ScenarioName`: Holds the name of the assigned scenario (if any).
*   **Error Variables:** Defines standard errors like `ErrInvalidRoomName`, `ErrRoomNotFound`, `ErrPlayerNotInRoom`.
*   **`NewRoom(id RoomID, name string) (*Room, error)`:** Constructor, validates name length, initializes fields.
*   **`AddPlayer(player *sharedEntity.User)`:** Appends a player to the `Players` slice.
*   **`RemovePlayer(playerID sharedEntity.UserID)`:** Removes a player by ID from the `Players` slice.
*   **`SetDescription(descriptionName string, text string)`:** Adds/updates an entry in the `Description` map.

## 2. `port/room_repository.go`

Defines the interfaces required by the Room domain to interact with persistence.

*   **`RoomReader` interface:**
    *   `GetRoomByID(id RoomID) (*Room, error)`
    *   `GetRooms() ([]*Room, error)`
    *   `GetPlayerRooms(playerID UserID) ([]*Room, error)`
    *   `GetPlayersInRoom(roomID RoomID) ([]*User, error)`
*   **`RoomWriter` interface:**
    *   `CreateRoom(room *Room) error`
    *   `UpdateRoom(room *Room) error`
    *   `AddPlayerToRoom(roomID RoomID, player *User) error`
    *   `RemovePlayerFromRoom(roomID RoomID, playerID UserID) error`
    *   `DeleteRoom(roomID RoomID) error`
*   **`RoomRepository` interface:** Embeds `RoomReader` and `RoomWriter`.

## 3. `usecase/command/` (Commands - State Changing)

*   **`add_description.go`:**
    *   `AddDescriptionCommand`: Contains `Requester` (User), `Room` (*Room), `DescriptionName`, `Text`.
    *   `AddDescriptionHandler`: Depends on `RoomRepository`. Handles admin check, calls `Room.SetDescription()`, and `RoomRepository.UpdateRoom()`.
*   **`create_room.go`:**
    *   `CreateRoomCommand`: Contains `ID`, `Name`, `CreatorID`.
    *   `CreateRoomHandler`: Depends on `RoomWriter` and `event.Publisher`. Calls `entity.NewRoom`, `RoomWriter.CreateRoom`, and publishes `RoomCreatedEvent`.
*   **`delete_room.go`:**
    *   `DeleteRoomCommand`: Contains `Requester`, `RoomID`.
    *   `DeleteRoomHandler`: Depends on `RoomWriter`. Handles admin check, calls `RoomWriter.DeleteRoom`.
*   **`join_room.go`:**
    *   `JoinRoomCommand`: Contains `Requester`, `RoomID`.
    *   `JoinRoomHandler`: Depends on `RoomRepository` and `event.Publisher`. Calls `RoomRepository.GetRoomByID` (to verify existence), `RoomRepository.AddPlayerToRoom`, and publishes `PlayerJoinedEvent`.
*   **`kick_user.go`:**
    *   `KickUserCommand`: Contains `Requester`, `RoomID`, `PlayerID`.
    *   `KickUserHandler`: Depends on `RoomRepository` and `event.Publisher`. Handles admin check, calls `RoomRepository.RemovePlayerFromRoom`, and publishes `PlayerKickedEvent`.
*   **`leave_room.go`:**
    *   `LeaveRoomCommand`: Contains `Requester`, `RoomID`.
    *   `LeaveRoomHandler`: Depends on `RoomRepository` and `event.Publisher`. Calls `RoomRepository.RemovePlayerFromRoom` and publishes `PlayerLeftEvent`.

## 4. `usecase/query/` (Queries - Data Retrieval)

*   **`get_player_rooms.go`:**
    *   `GetPlayerRoomsQuery`: Contains `PlayerID`.
    *   `GetPlayerRoomsHandler`: Depends on `RoomReader`. Calls `RoomReader.GetPlayerRooms`.
*   **`get_players_in_room.go`:**
    *   `GetPlayersInRoomQuery`: Contains `RoomID`.
    *   `GetPlayersInRoomHandler`: Depends on `RoomReader`. Calls `RoomReader.GetPlayersInRoom`.
*   **`get_room.go`:**
    *   `GetRoomQuery`: Contains `RoomID`.
    *   `GetRoomHandler`: Depends on `RoomReader`. Calls `RoomReader.GetRoomByID`.
*   **`get_rooms.go`:**
    *   `GetRoomsQuery`: (Currently empty).
    *   `GetRoomsHandler`: Depends on `RoomReader`. Calls `RoomReader.GetRooms`. 