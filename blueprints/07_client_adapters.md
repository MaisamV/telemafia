# Blueprint 07: Client Adapters (Local)

**Source:** `internal/adapters/api/`

**Purpose:** Details the adapters that implement client interfaces defined in domain ports, allowing domain modules (like Game) to fetch data from other modules (like Room, Scenario) without direct dependency.

## 1. Overview

These adapters implement the `RoomClient` and `ScenarioClient` interfaces defined in `internal/domain/game/port/`. They act as intermediaries for cross-domain data fetching.

In the current **monolithic structure**, these adapters are implemented locally. They simply wrap the repository *reader* interfaces of the target domain and call the appropriate repository methods.

If the application were split into **microservices**, these local implementations would be replaced by adapters containing actual network clients (e.g., HTTP or gRPC clients) to communicate with the respective Room or Scenario services.

## 2. `room_client.go` (`LocalRoomClient`)

*   Implements `gamePort.RoomClient`.
*   **Dependency:** Takes a `roomPort.RoomReader` in its constructor.
*   **`FetchRoom(id roomEntity.RoomID) (*roomEntity.Room, error)`:**
    *   Implementation: Calls `c.roomRepo.GetRoomByID(id)`.
    *   Purpose: Allows the Game domain (specifically `CreateGameHandler`) to fetch `Room` details using the `RoomClient` abstraction.

## 3. `scenario_client.go` (`LocalScenarioClient`)

*   Implements `gamePort.ScenarioClient`.
*   **Dependency:** Takes a `scenarioPort.ScenarioReader` in its constructor.
*   **`FetchScenario(id string) (*scenarioEntity.Scenario, error)`:**
    *   Implementation: Calls `c.scenarioRepo.GetScenarioByID(id)`.
    *   Purpose: Allows the Game domain (specifically `CreateGameHandler` and `AssignRolesHandler`) to fetch `Scenario` details using the `ScenarioClient` abstraction. 