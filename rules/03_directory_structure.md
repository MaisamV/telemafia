# 3. Directory Structure Guide

**Goal:** Ensure code is placed in the correct layer and module for consistency and adherence to Clean Architecture.

Refer to `blueprint/06_directory_structure.md` for the full tree. This guide summarizes where to put common components:

*   **New Business Capability/Domain:**
    *   Create a new directory under `internal/domain/<new_module_name>/`.
    *   Inside, create `entity/`, `port/`, and `usecase/` (with `command/` and `query/` subdirectories).
*   **Domain Entities / Value Objects:**
    *   Specific to a module: `internal/domain/<module_name>/entity/`
    *   Shared across modules (use sparingly): `internal/shared/entity/` (e.g., `User`)
*   **Repository Interfaces (Ports):**
    *   `internal/domain/<module_name>/port/` (e.g., `room_repository.go`)
*   **Use Cases (Commands/Queries):**
    *   Command Handlers: `internal/domain/<module_name>/usecase/command/`
    *   Query Handlers: `internal/domain/<module_name>/usecase/query/`
*   **Repository Implementations (Adapters):**
    *   `internal/adapters/repository/<implementation_type>/` (e.g., `memory/`)
*   **External Service Client Adapters (e.g., other APIs):**
    *   Interfaces (if needed by domain): `internal/domain/<module_name>/port/`
    *   Implementations: `internal/adapters/<client_type>/` (e.g., `api/`)
*   **Telegram Command Handlers (Presentation):**
    *   Logic: Exported functions in `internal/presentation/telegram/handler/<module_name>/` (e.g., `room/create_room.go`)
    *   Common (non-module specific): `internal/presentation/telegram/handler/common_handlers.go`
    *   Dispatcher methods (calling the exported functions): Add method to `BotHandler` in `internal/presentation/telegram/handler/bot_handler.go` and register in `RegisterHandlers`.
*   **Telegram Callback Handlers (Presentation):**
    *   Logic: Exported functions in `internal/presentation/telegram/handler/<module_name>/` (e.g., `room/callbacks_room.go`)
    *   Routing: Add `case` to `switch` statement in `handleCallback` method in `internal/presentation/telegram/handler/callbacks.go`.
*   **Shared Utilities:**
    *   Telegram-specific: `internal/shared/tgutil/`
    *   General Go (rarely needed): `internal/shared/common/`
*   **Domain Events:**
    *   Interface/Publisher: `internal/shared/event/event.go`
    *   Concrete Structs: `internal/shared/event/events.go`
*   **Configuration Loading:**
    *   `internal/config/config.go`
*   **User-Facing Text (Messages):**
    *   JSON definitions: `messages.json` (root)
    *   Loading code: `internal/presentation/telegram/messages/loader.go`
    *   Go struct definitions: `internal/presentation/telegram/messages/messages.go`
*   **Main Application Entry / DI Wiring:**
    *   `cmd/telemafia/main.go`
*   **Tests:**
    *   Unit Tests: `tests/unit/` (for tests focusing on individual components or functions in isolation, e.g., `shuffle_distribution_test.go`).
    *   Integration Tests: `tests/integration/` (for tests verifying interactions between multiple components, e.g., simulating a full game flow via API calls if applicable).
    *   *Note:* While a central `tests` directory is provided, unit tests can also be co-located with the package they test (e.g., `internal/domain/<module_name>/entity/<entity>_test.go`), following standard Go conventions. The chosen approach for new tests should be consistent with existing patterns for similar tests. 