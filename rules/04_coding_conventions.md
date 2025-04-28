# 4. Coding Conventions & Rules

**Goal:** Ensure code consistency, safety, and adherence to project specifics.

## 4.1. Language & Dependencies

*   **Language:** Go (Version 1.18+)
*   **Dependencies:** Managed via Go Modules (`go.mod`). Primary dependency is `gopkg.in/telebot.v3`.

## 4.2. Error Handling

*   **Domain Errors:** Define custom error variables within domain entity packages (e.g., `roomEntity.ErrRoomNotFound`). Repositories and Use Cases should return these specific errors when appropriate.
*   **Other Errors:** Use standard Go error handling (e.g., `fmt.Errorf`, `errors.New`).
*   **Checking Errors:** Always check errors returned from function calls.
*   **Logging:** Use the standard `log` package for informative logging, especially for non-fatal errors in background tasks or callbacks.

## 4.3. Concurrency (In-Memory Persistence)

*   **Requirement:** The current persistence mechanism is in-memory maps.
*   **Safety:** All access (read and write) to shared in-memory data stores (e.g., maps within repository implementations) **MUST** be protected using `sync.RWMutex`.
    *   Use `mutex.Lock()` / `defer mutex.Unlock()` for write operations.
    *   Use `mutex.RLock()` / `defer mutex.RUnlock()` for read operations.
*   **Data Copying:** When returning slices or maps from repositories, consider returning copies to prevent external modification of internal state.

## 4.4. Configuration (`internal/config`)

*   Load configuration using `config.LoadConfig`.
*   Priority: Command-line flags (`-token`, `-admins`) > `config.json` file.
*   Application **MUST** fail on startup if required configuration (token) is missing.

## 4.5. User-Facing Messages (`messages.json`)

*   **Externalized:** All text shown to the user (prompts, errors, button labels, etc.) **MUST** be defined in `messages.json`.
*   **Access:** Handlers receive the loaded `*messages.Messages` struct via dependency injection.
*   **Usage:** Access strings via the struct (e.g., `msgs.Room.CreatePrompt`, `fmt.Sprintf(msgs.Common.ErrorGeneric, err)`).
*   **DO NOT** hardcode user-facing strings in Go code.

## 4.6. Naming Conventions

*   Follow standard Go naming conventions (CamelCase for exported identifiers, camelCase for unexported).
*   Repository interfaces: `TypeReader`, `TypeWriter`, `TypeRepository` (e.g., `RoomReader`).
*   Use Case Handlers: `ActionNounHandler` (e.g., `CreateRoomHandler`, `GetRoomsHandler`).
*   Command/Query Structs: `ActionNounCommand`/`ActionNounQuery` (e.g., `CreateRoomCommand`).
*   Telegram Handlers (Exported Functions): `HandleActionNoun` (e.g., `HandleCreateRoom`).

## 4.7. Authorization

*   Admin checks **MUST** be performed in:
    *   Use Case Command Handlers (for domain-level authorization) by checking the `Admin` flag on the `Requester` field of the command struct.
    *   Presentation Layer Handlers (for presentation-level checks, if needed, though domain checks are preferred) using `tgutil.IsAdmin` or checking the `Admin` flag on the converted `*sharedEntity.User`.
*   Return a standard permission error message (from `messages.json`) if authorization fails. 