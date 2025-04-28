# 6. Domain Modeling Guidelines

**Goal:** Ensure the domain layer accurately reflects business concepts and remains independent.

## 6.1. Entities & Value Objects

*   **Location:** `internal/domain/<module>/entity/`
*   **Entities:** Represent objects with identity and lifecycle (e.g., `Room`, `Game`, `Scenario`).
    *   Have a unique identifier (e.g., `RoomID`).
    *   Contain attributes representing their state.
    *   May contain methods that enforce business rules (invariants) related to their state (e.g., validation in a `NewRoom` constructor).
*   **Value Objects:** Represent descriptive aspects of the domain, identified by their value, not identity (e.g., `Role`, `UserID`, `GameState`).
    *   Typically immutable.
    *   Used as attributes within entities.
*   **Independence:** Domain objects **MUST NOT** depend on:
    *   Infrastructure details (e.g., database specifics, specific API clients).
    *   Framework details (e.g., `telebot` types).
    *   Other domain modules directly (use interfaces or events for interaction).
*   **Validation:** Business rule validation (e.g., room name length) should ideally occur within entity constructors or methods.
*   **Errors:** Define domain-specific errors (e.g., `ErrRoomNotFound`) within the entity package.

## 6.2. Ports (Interfaces)

*   **Location:** `internal/domain/<module>/port/`
*   Define interfaces required by the domain/use cases to interact with the outside world (typically infrastructure).
*   Examples: `RoomRepository`, `ScenarioRepository`, `GameRepository`, `event.Publisher`, `game.RoomClient`.
*   Define the *contract* based on what the domain needs, not what a specific adapter provides.
*   **MUST NOT** leak implementation details.

## 6.3. Use Cases (Interactors)

*   **Location:** `internal/domain/<module>/usecase/command/` and `internal/domain/<module>/usecase/query/`
*   Implement application-specific business rules.
*   Orchestrate domain entities and repositories (via ports) to fulfill a specific command or query.
*   Receive dependencies (repositories, other handlers, publishers) via constructor injection.
*   **MUST NOT** contain infrastructure or UI logic.
*   Should handle transaction management or unit-of-work if applicable (though not implemented with current in-memory approach).

## 6.4. Domain Events

*   **Purpose:** Decouple modules and notify other parts of the system about significant occurrences in the domain.
*   **Definition:**
    *   Base `Event` interface: `internal/shared/event/event.go`
    *   Concrete event structs: `internal/shared/event/events.go` (e.g., `RoomCreatedEvent`).
*   **Publishing:**
    *   `Publisher` interface: `internal/shared/event/event.go`.
    *   Command handlers responsible for the state change should publish events via the injected `Publisher`.
    *   The current implementation logs events, but this allows for future expansion (e.g., triggering other use cases, sending notifications). 