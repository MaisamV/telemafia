# 2. Design Patterns

**Goal:** Ensure consistent application of established design patterns for clarity and maintainability.

The following patterns **MUST** be used:

## 2.1. Command Query Responsibility Segregation (CQRS)

*   **Location:** `internal/domain/<module>/usecase/`
*   Separate application logic into:
    *   **Commands:** State-changing operations (`command/`). Handlers typically receive data via a Command struct, perform actions using domain entities and repositories, and return minimal data (e.g., ID, error) or nothing. Command structs should include the `Requester` user for authorization checks.
    *   **Queries:** Data retrieval operations (`query/`). Handlers typically receive query parameters, fetch data using repositories, and return Data Transfer Objects (DTOs) or domain entities.
*   This promotes the Single Responsibility Principle for use cases.

## 2.2. Repository Pattern

*   **Purpose:** Decouple the domain layer from data persistence concerns.
*   **Interface Definition (`Ports`):**
    *   Location: `internal/domain/<module>/port/`
    *   Define data access interfaces (e.g., `RoomRepository`, `ScenarioRepository`, `GameRepository`).
    *   Use Go interface embedding to separate `Reader` and `Writer` interfaces where appropriate (e.g., `RoomRepository interface { RoomReader; RoomWriter }`).
    *   These interfaces **MUST NOT** leak implementation details (e.g., SQL queries, specific database types, change flags).
*   **Implementation (`Adapters`):**
    *   Location: `internal/adapters/repository/<type>/` (e.g., `internal/adapters/repository/memory/`)
    *   Provide concrete implementations of the repository ports.
    *   Implementations handle the specifics of the chosen storage (e.g., maps with mutexes for in-memory).
    *   Constructors (e.g., `NewInMemoryRoomRepository`) **MUST** return the *port interface type*, not the concrete struct type.

## 2.3. Dependency Injection (DI)

*   **Purpose:** Decouple components and facilitate testing.
*   **Method:** Primarily Constructor Injection.
*   **Implementation:**
    *   Dependencies (e.g., Repositories, other Use Case Handlers, Event Publishers, Message Loaders, Refresh Notifiers) **MUST** be injected into their consumers (e.g., Use Case Handlers, Presentation Handlers) via their constructor functions.
    *   The application's **Composition Root** (`cmd/telemafia/main.go`) is responsible for instantiating concrete types and wiring them together.
    *   Components should depend on abstractions (interfaces/ports) where possible, not concrete implementations (except at the Composition Root). 