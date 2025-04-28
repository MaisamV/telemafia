# 1. Architecture Guidelines

**Goal:** Maintain a clean, modular, and testable codebase.

## 1.1. Core Style: Clean Architecture (Ports & Adapters / Hexagonal)

*   The project strictly follows the Clean Architecture pattern.
*   Business logic (domain) is independent of frameworks, UI, and infrastructure.
*   Dependencies **MUST** point inwards only.

## 1.2. Layers

Code resides within the `internal/` directory and is organized into these primary layers:

1.  **`domain`:**
    *   Location: `internal/domain/<module>/...`
    *   Contains core business logic: Entities, Value Objects, Use Cases (Commands/Queries), and Ports (interfaces).
    *   **MUST HAVE NO DEPENDENCIES** on outer layers (adapters, presentation, config, shared infrastructure like specific DBs or APIs).
    *   Divided into modules based on business capability (e.g., `room`, `scenario`, `game`). Inter-module dependencies should be minimized and ideally handled via interfaces or events if necessary.
2.  **`adapters`:**
    *   Location: `internal/adapters/...`
    *   Contains implementations of domain ports (interfaces). Examples: Database repositories (`repository/memory`), external service clients (`api`).
    *   Depends **ONLY** on `internal/domain` ports (interfaces).
3.  **`presentation`:**
    *   Location: `internal/presentation/...`
    *   Contains adapters that drive the application. Example: Telegram bot handlers (`telegram`).
    *   Depends **ONLY** on `internal/domain` use cases (command/query handlers) and ports.
4.  **`shared`:**
    *   Location: `internal/shared/...`
    *   Contains code shared across layers/modules, but **SHOULD NOT** contain core business logic. Examples: Shared entities (`User`), event definitions, common Telegram utilities (`tgutil`). Avoid introducing dependencies *back* into specific domain modules from here.
5.  **`config`:**
    *   Location: `internal/config/`
    *   Handles loading application configuration.
6.  **`cmd`:**
    *   Location: `cmd/<app_name>/`
    *   Application entry point (`main.go`). Acts as the **Composition Root**, responsible for instantiating concrete types (repositories, handlers) and wiring them together via Dependency Injection.

## 1.3. The Dependency Rule

*   **STRICTLY ENFORCED:** Source code dependencies can only point inwards.
*   `presentation` -> `domain` <- `adapters`
*   Nothing in an inner layer can know anything about an outer layer. Specifically, `domain` code **MUST NOT** import anything from `adapters`, `presentation`, `config`, or `cmd`.

## 1.4. Modularity

*   The application is structured as a **Modular Monolith**.
*   The `internal/domain` layer is divided into modules (`room`, `scenario`, `game`). Each module encapsulates its specific entities, ports, and use cases.
*   Aim for high cohesion within modules and low coupling between them.

## 1.5. Core Principles

*   Adhere to **SOLID** principles:
    *   **S**ingle Responsibility: Achieved via CQRS, modular design, focused handlers/use cases.
    *   **O**pen/Closed: Domain core should be extensible (new adapters) without modification.
    *   **L**iskov Substitution: Interface implementations must be substitutable.
    *   **I**nterface Segregation: Define granular interfaces (ports).
    *   **D**ependency Inversion: Depend on abstractions (ports), not concretions. 