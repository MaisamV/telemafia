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
    *   Fields: `Name` (string), `DefaultRole` (string, optional), `Roles` ([]string - names of roles).
*   **`Role` struct:**
    *   Represents a single role *after* extraction for assignment (used by Game module).
    *   Fields: `Name` (string), `Side` (string).

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
    *   `AddScenarioJSONHandler`: Depends on `ScenarioWriter`. Handles admin check, unmarshals JSON into `Scenario` entity, performs validation (names, non-empty roles/sides), generates internal ID, calls `ScenarioWriter.CreateScenario`.
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