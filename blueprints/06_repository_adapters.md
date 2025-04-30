# Blueprint 06: Repository Adapters (In-Memory)

**Source:** `internal/adapters/repository/memory/`

**Purpose:** Details the concrete in-memory implementations of the domain repository ports.

## 1. Overview

This directory contains adapters that implement the `RoomRepository`, `ScenarioRepository`, and `GameRepository` interfaces defined in the respective domain `port/` directories. Since the current persistence mechanism is in-memory, these adapters use Go maps and mutexes.

**Key Characteristics:**

*   **Data Storage:** Uses standard Go maps (e.g., `map[RoomID]*Room`) to store domain entities.
*   **Concurrency Control:** Employs `sync.RWMutex` to protect map access. Write operations (`Create`, `Update`, `Delete`, `AddPlayer`, `RemovePlayer`) use `mutex.Lock() / Unlock()`, while read operations (`Get...`) use `mutex.RLock() / RUnlock()`.
*   **Interface Implementation:** Each repository struct (e.g., `InMemoryRoomRepository`) explicitly implements the corresponding domain port interface (e.g., `roomPort.RoomRepository`).
*   **Constructors:** Provide `NewInMemory...Repository()` functions that return the *interface type*, hiding the concrete implementation from the rest of the application (as seen in `main.go`).
*   **Error Handling:** Returns standard domain errors where appropriate (e.g., `roomEntity.ErrRoomNotFound`) or formatted errors for implementation-specific issues (e.g., record already exists).

## 2. `room_repository.go` (`InMemoryRoomRepository`)

*   Implements `roomPort.RoomRepository`.
*   Stores data in `rooms map[roomEntity.RoomID]*roomEntity.Room`.
*   Uses a single `sync.RWMutex` for all map operations.
*   `AddPlayerToRoom` / `RemovePlayerFromRoom` modify the `Players` slice within the `Room` struct stored in the map.
*   `UpdateRoom` replaces the existing room pointer in the map with the provided one.
*   Getter methods (`GetRooms`, `GetPlayerRooms`, `GetPlayersInRoom`) iterate over the map or room slices as needed.

## 3. `scenario_repository.go` (`InMemoryScenarioRepository`)

*   Implements `scenarioPort.ScenarioRepository`.
*   Stores data in `data map[string]*scenarioEntity.Scenario` (key is `Scenario.ID`).
*   Uses a single `sync.RWMutex`.
*   Provides basic CRUD operations for scenarios.

## 4. `game_repository.go` (`InMemoryGameRepository`)

*   Implements `gamePort.GameRepository`.
*   Stores data in:
    *   `games map[gameEntity.GameID]*gameEntity.Game`
    *   `roomToGame map[roomEntity.RoomID]gameEntity.GameID` (for efficient lookup by RoomID).
*   Uses a single `sync.RWMutex`.
*   `CreateGame` adds entries to both maps.
*   `UpdateGame` replaces the entry in the `games` map.
*   `DeleteGame` removes entries from both maps.
*   `GetGameByRoomID` uses the `roomToGame` map first, then looks up the `Game` in the `games` map. 