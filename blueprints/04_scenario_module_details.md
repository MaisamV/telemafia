# Blueprint 04: Scenario Module Details

**Source:** `internal/domain/scenario/`

**Purpose:** Details the components within the Scenario domain module.

## 1. `entity/scenario.go`

*   **`Scenario` struct:**
    *   Top-level structure representing a game scenario definition.
    *   Fields: `ID` (string, internal), `Name` (string), `Sides` ([]Side).
    *   Intended to be populated from JSON data (e.g., via `add_scenario_json.go` use case).
*   **`Side` struct:**
    *   Represents a group of roles (e.g., Mafia, Civilians).
    *   Fields: `Name` (string), `PopulationRate` (*float32, optional - for dynamic role count calculation), `DefaultRole` (*Role, optional - used if explicit roles are insufficient for player count), `Roles` ([]Role - list of actual Role structs explicitly defined for the side).
*   **`Role` struct:**
    *   Represents a single role with its `Name`, optional `ImageID`, and `Side` (populated during flattening).
    *   Used *after* extracting roles from the `Scenario` structure for assignment (by Game module) or display.
*   **`Scenario.FlatRoles() []Role` (Method):**
    *   Returns a flattened list of all *explicitly defined* `Role` structs from all `Sides` within the scenario.
    *   During flattening, the `Side` field of each `Role` struct is populated with the name of the `Side` it belongs to.
*   **`Scenario.GetRoles(playerNum int) []Role` (Method):**
    *   NEW: Returns a list of roles dynamically calculated for the given `playerNum`.
    *   It starts with explicitly defined roles from `Scenario.Sides[...].Roles`.
    *   If the total number of explicit roles is less than `playerNum`, it can fill the remaining slots using the `DefaultRole` of sides.
    *   If `Side.PopulationRate` is defined, it influences how many roles (including default ones) are taken from that side relative to `playerNum`.
    *   The method sorts roles by a hash of their name before returning.
*   **`Scenario.GetShuffledRoles(playerNum int) []Role` (Method):**
    *   NEW: Calls `GetRoles(playerNum)` and then shuffles the resulting list.
    *   This is the primary method used by the Game domain to get roles for assignment. The shuffling aims for a statistically reasonable distribution of roles, which can be verified (as exemplified by tests like `TestRoleShuffleDistribution`).

## 2. `port/scenario_repository.go`

Defines the interfaces required by the Scenario domain to interact with persistence.

*   **`ScenarioReader` interface:**
    *   `GetScenarioByID(id string) (*Scenario, error)`
    *   `GetAllScenarios() ([]*Scenario, error)`
*   **`ScenarioWriter` interface:**
    *   `CreateScenario(scenario *Scenario) error`
    *   `DeleteScenario(id string) error`
*   **`ScenarioRepository` interface:** Embeds `ScenarioReader` and `ScenarioWriter`.

## 3. `usecase/command/` (Commands - State Changing)

*   **`add_scenario_json.go`:**
    *   `AddScenarioJSONCommand`: Contains `Requester`, `JSONData` (string).
    *   `AddScenarioJSONHandler`: Depends on `ScenarioWriter`. Handles admin check, unmarshals JSON into `Scenario` entity, performs validation (names, non-empty roles/sides unless a DefaultRole is present, non-empty role names, at least one role overall if no default roles are viable), generates internal ID, calls `ScenarioWriter.CreateScenario`.
*   **`create_scenario.go`:** (Note: This seems superseded by `add_scenario_json` for complex scenarios, but might be used for basic name/ID creation initially).
    *   `CreateScenarioCommand`: Contains `Requester`, `ID`, `Name`.
    *   `CreateScenarioHandler`: Depends on `ScenarioWriter`. Handles admin check, creates a basic `Scenario` struct, calls `ScenarioWriter.CreateScenario`.
*   **`delete_scenario.go`:**
    *   `DeleteScenarioCommand`: Contains `Requester`, `ID`.
    *   `DeleteScenarioHandler`: Depends on `ScenarioWriter`. Handles admin check, calls `ScenarioWriter.DeleteScenario`.

## 4. `usecase/query/` (Queries - Data Retrieval)

*   **`get_all_scenarios.go`:**
    *   `GetAllScenariosQuery`: (Empty).
    *   `GetAllScenariosHandler`: Depends on `ScenarioReader`. Calls `ScenarioReader.GetAllScenarios`.
*   **`get_scenario_by_id.go`:**
    *   `GetScenarioByIDQuery`: Contains `ID`.
    *   `GetScenarioByIDHandler`: Depends on `ScenarioReader`. Calls `ScenarioReader.GetScenarioByID`. 