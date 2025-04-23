# 4. Repository Port Specification

**Goal:** Define the exact Go interface contracts for data persistence, separating the domain layer from specific storage implementations.

## 4.1. Room Repository (`internal/domain/room/port/room_repository.go`)

```go
package port

import (
	roomEntity "telemafia/internal/domain/room/entity"
	sharedEntity "telemafia/internal/shared/entity"
)

// RoomReader defines the interface for reading room data
type RoomReader interface {
	GetRoomByID(id roomEntity.RoomID) (*roomEntity.Room, error)
	GetRooms() ([]*roomEntity.Room, error)
	GetPlayerRooms(playerID sharedEntity.UserID) ([]*roomEntity.Room, error)
	GetPlayersInRoom(roomID roomEntity.RoomID) ([]*sharedEntity.User, error)

	// CheckChangeFlag checks the current state of the change flag (specific to in-memory refresh logic)
	CheckChangeFlag() bool
}

// RoomWriter defines the interface for writing room data
type RoomWriter interface {
	CreateRoom(room *roomEntity.Room) error
	UpdateRoom(room *roomEntity.Room) error
	AddPlayerToRoom(roomID roomEntity.RoomID, player *sharedEntity.User) error
	RemovePlayerFromRoom(roomID roomEntity.RoomID, playerID sharedEntity.UserID) error
	DeleteRoom(roomID roomEntity.RoomID) error

	// ConsumeChangeFlag checks and resets the change flag (specific to in-memory refresh logic)
	ConsumeChangeFlag() bool
	// RaiseChangeFlag sets the change flag to true (specific to in-memory refresh logic)
	RaiseChangeFlag()
}

// RoomRepository defines the combined interface for room persistence
type RoomRepository interface {
	RoomReader
	RoomWriter
}
```

**Notes:**

*   The `CheckChangeFlag`, `ConsumeChangeFlag`, and `RaiseChangeFlag` methods are specific to the current in-memory implementation's refresh mechanism. They might be removed or refactored if persistence changes.
*   `AssignScenarioToRoom` and `GetRoomScenario` methods have been removed. Room state changes (like setting `ScenarioName` or adding `Description`) should be done on the `Room` entity itself, and then persisted using the `UpdateRoom` method.
*   The repository should handle potential errors like `ErrRoomNotFound`, `ErrRoomAlreadyExists`, `ErrPlayerNotInRoom` as defined in the Room entity package.

## 4.2. Scenario Repository (`internal/domain/scenario/port/scenario_repository.go`)

```go
package port

import (
	scenarioEntity "telemafia/internal/domain/scenario/entity"
)

// ScenarioReader defines the interface for reading scenario data
type ScenarioReader interface {
	GetScenarioByID(id string) (*scenarioEntity.Scenario, error)
	GetAllScenarios() ([]*scenarioEntity.Scenario, error)
}

// ScenarioWriter defines the interface for writing scenario data
type ScenarioWriter interface {
	CreateScenario(scenario *scenarioEntity.Scenario) error
	DeleteScenario(id string) error
	AddRoleToScenario(scenarioID string, role scenarioEntity.Role) error
	RemoveRoleFromScenario(scenarioID string, roleName string) error
}

// ScenarioRepository defines the interface for scenario persistence
type ScenarioRepository interface {
	ScenarioReader
	ScenarioWriter
}

```

**Notes:**

*   Assumes scenario IDs are strings.
*   `RemoveRoleFromScenario` operates based on `roleName`.
*   Should handle errors like scenario not found, scenario already exists, role not found.

## 4.3. Game Repository (`internal/domain/game/port/game_repository.go`)

```go
package port

import (
	gameEntity "telemafia/internal/domain/game/entity"
	roomEntity "telemafia/internal/domain/room/entity" // Needed for RoomID
)

// GameReader defines the interface for reading game data
type GameReader interface {
	GetGameByID(id gameEntity.GameID) (*gameEntity.Game, error)
	GetGameByRoomID(roomID roomEntity.RoomID) (*gameEntity.Game, error)
	GetAllGames() ([]*gameEntity.Game, error)
}

// GameWriter defines the interface for writing game data
type GameWriter interface {
	CreateGame(game *gameEntity.Game) error
	UpdateGame(game *gameEntity.Game) error // Used for saving state changes (like assignments, state transitions)
	DeleteGame(id gameEntity.GameID) error
}

// GameRepository defines the interface for game persistence
type GameRepository interface {
	GameReader
	GameWriter
}
```

**Notes:**

*   `UpdateGame` is crucial for persisting changes made to the Game entity after creation (e.g., adding assignments, changing state).
*   Should handle errors like game not found, game already exists.

## 4.4. Room Client Port (`internal/domain/game/port/room_client.go`)

```go
package port

import (
	roomEntity "telemafia/internal/domain/room/entity"
)

// RoomClient defines an interface for the Game domain to fetch Room data.
// Implementations could be local (monolith) or remote (microservice).
type RoomClient interface {
	FetchRoom(id roomEntity.RoomID) (*roomEntity.Room, error)
}
```

**Notes:**

*   This interface abstracts the fetching of Room data needed by the Game domain.

## 4.5. Scenario Client Port (`internal/domain/game/port/scenario_client.go`)

```go
package port

import (
	scenarioEntity "telemafia/internal/domain/scenario/entity"
)

// ScenarioClient defines an interface for the Game domain to fetch Scenario data.
// Implementations could be local (monolith) or remote (microservice).
type ScenarioClient interface {
	FetchScenario(id string) (*scenarioEntity.Scenario, error)
}
```

**Notes:**

*   This interface abstracts the fetching of Scenario data needed by the Game domain.

// End of file 