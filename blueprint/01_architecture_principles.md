# 1. Architecture & Principles Specification

**Goal:** Define the non-functional requirements, architectural style, patterns, and constraints for the TeleMafia Bot project.

## 1.1. Core Architecture

*   **Style:** Clean Architecture implemented using the **Ports & Adapters (Hexagonal)** pattern.
*   **Layers:**
    *   `internal/domain`: Core business logic. Contains Entities, Use Cases (Commands/Queries), and Ports (interfaces). Must have **NO** dependencies on outer layers.
    *   `internal/adapters`: Infrastructure implementations (e.g., database repositories, external service clients). Depends **only** on `internal/domain` ports.
    *   `internal/presentation`: Adapters for driving the application (e.g., Telegram bot handlers, API controllers). Depends **only** on `internal/domain` use cases/ports.
    *   `internal/config`: Application configuration loading.
    *   `internal/shared`: Common utilities, entities, or interfaces used across multiple domains/layers.
    *   `cmd`: Main application entry point (composition root), responsible for wiring dependencies.
*   **Dependency Rule:** Dependencies **MUST** point inwards only: `Presentation` -> `Domain` <- `Adapters`.

## 1.2. Architectural Style

*   **Type:** Modular Monolith.
*   **Modularity:** The `internal/domain` layer **MUST** be divided into modules based on core business capabilities: `room`, `scenario`, `game`. Each module should encapsulate its specific entities, use cases, and ports.

## 1.3. Key Patterns

The following patterns **MUST** be used:

*   **Command Query Responsibility Segregation (CQRS):**
    *   Within each domain module (`internal/domain/<module>`), separate use cases into `usecase/command` (for state-changing operations) and `usecase/query` (for data retrieval operations).
    *   Commands should typically return minimal data (e.g., ID, error) or nothing, while Queries return data transfer objects (DTOs) or entities.
*   **Repository Pattern:**
    *   Define data access interfaces (`ports`) within the corresponding domain module (`internal/domain/<module>/port`). Examples: `RoomRepository`, `ScenarioRepository`, `GameRepository`.
    *   Use Go interface embedding to separate `Reader` and `Writer` interfaces within the main repository port where appropriate (e.g., `RoomRepository interface { RoomReader; RoomWriter }`).
    *   Place concrete implementations of these interfaces in the `internal/adapters/repository/<type>` directory (e.g., `internal/adapters/repository/memory`).
*   **Dependency Injection (DI):**
    *   Dependencies (e.g., Repositories, other Use Case Handlers, Event Publishers) **MUST** be injected into handlers and services, primarily via constructors.
    *   The application's composition root (`cmd/telemafia/main.go`) is responsible for instantiating concrete types and wiring them together.

## 1.4. Core Principles

*   Adherence to **SOLID** principles is a primary goal, facilitated by the Clean Architecture structure.
    *   **S**ingle Responsibility: Enforced by CQRS, modular design, and focused handlers.
    *   **O**pen/Closed: Domain core should be open for extension (new adapters) but closed for modification.
    *   **L**iskov Substitution: Interface implementations must be substitutable.
    *   **I**nterface Segregation: Define specific, granular interfaces (`ports`).
    *   **D**ependency Inversion: High-level modules depend on abstractions (ports), not concretions.

## 1.5. Constraints & Language

*   **Persistence:** **In-Memory Persistence Only**. All application state (rooms, scenarios, games, users) is lost upon application restart. Repository implementations should reflect this (e.g., using maps and mutexes).
*   **Language:** Go (Version 1.18 or higher).
*   **Concurrency:** Use `sync.RWMutex` appropriately within in-memory repositories to handle potential concurrent access safely. 