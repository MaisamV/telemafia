# 2. Feature Specification

**Goal:** Define the functional requirements of the bot from the user's perspective, detailing commands, inputs, outputs, permissions, and expected state changes.

## 2.1. General Commands

1.  **Command:** `/start`
    *   **Description:** Start interaction with the bot, show welcome message and available rooms.
    *   **User Role:** Everyone
    *   **Input:** None
    *   **Output/Behavior:** Sends a welcome message (`"Welcome, <username>!"`). Sends/Updates a dynamic message listing available rooms with join buttons.
    *   **State Changes:** None directly, may register the user's chat for dynamic updates.

2.  **Command:** `/help`
    *   **Description:** Show help message listing available commands.
    *   **User Role:** Everyone
    *   **Input:** None
    *   **Output/Behavior:** Sends a formatted text message listing all commands and their usage.
    *   **State Changes:** None

## 2.2. Room Management Commands

1.  **Command:** `/create_room <room_name>`
    *   **Description:** Create a new game room.
    *   **User Role:** Admin Only
    *   **Input:** `<room_name>` (string, length 3-50 characters).
    *   **Output/Behavior:**
        *   Success: `"Room '<room_name>' created successfully! ID: <room_id>"`
        *   Error (No Name): `"Please provide a room name: /create_room [name]"`
        *   Error (Not Admin): `"Only admins can create rooms"`
        *   Error (Exists): `"Error: Room name '<room_name>' already exists."` (Check by ID in implementation)
        *   Error (Invalid Name): `"Error: Invalid room name '<room_name>'. Must be 3-50 characters."`
        *   Other Error: `"Error creating room: <error_details>"`
    *   **State Changes:** Creates a new `Room` entity in the repository. Publishes `RoomCreatedEvent`.

2.  **Command:** `/join_room <room_id>`
    *   **Description:** Join a specific room.
    *   **User Role:** Everyone
    *   **Input:** `<room_id>` (string, valid Room ID).
    *   **Output/Behavior:**
        *   Success: `"Successfully joined room <room_id>"` (with inline button to leave).
        *   Error (No ID): `"Please provide a room ID: /join_room <room_id>"`
        *   Error (Not Found/Other): `"Error joining room '<room_id>': <error_details>"`
    *   **State Changes:** Adds the user (`shared.User`) to the `Players` list of the specified `Room` entity in the repository. Publishes `PlayerJoinedEvent`.

3.  **Command:** `/leave_room <room_id>`
    *   **Description:** Leave a specific room.
    *   **User Role:** Everyone
    *   **Input:** `<room_id>` (string, valid Room ID the user is in).
    *   **Output/Behavior:**
        *   Success: `"Successfully left room <room_id>!"`
        *   Error (No ID): `"Please provide a room ID: /leave_room [room_id]"`
        *   Error (Not Found/Not In Room/Other): `"Error leaving room '<room_id>': <error_details>"`
    *   **State Changes:** Removes the user from the `Players` list of the specified `Room` entity. Publishes `PlayerLeftEvent`.

4.  **Command:** `/list_rooms`
    *   **Description:** List all available rooms.
    *   **User Role:** Everyone
    *   **Input:** None
    *   **Output/Behavior:**
        *   Sends a message listing rooms: `"Available Rooms:\n- <Name> (<ID>) [<PlayerCount>/<MaxPlayers> players]\n..."` (Potentially dynamic/refreshing with Join buttons).
        *   No Rooms: `"No rooms available."`
        *   Error: `"Error getting rooms: <error_details>"`
    *   **State Changes:** None.

5.  **Command:** `/my_rooms`
    *   **Description:** List rooms the user has joined.
    *   **User Role:** Everyone
    *   **Input:** None
    *   **Output/Behavior:**
        *   Sends a message listing joined rooms: `"Rooms you are in:\n- <Name> (<ID>)\n..."`
        *   Not In Rooms: `"You are not in any rooms."`
        *   Error: `"Error getting your rooms: <error_details>"`
    *   **State Changes:** None.

6.  **Command:** `/kick_user <room_id> <user_id>`
    *   **Description:** Kick a user from a room.
    *   **User Role:** Admin Only
    *   **Input:** `<room_id>` (string, valid Room ID), `<user_id>` (int64, Telegram User ID).
    *   **Output/Behavior:**
        *   Success: `"User <user_id> kicked from room <room_id>"`
        *   Error (Usage): `"Usage: /kick_user <room_id> <user_id>"`
        *   Error (Invalid User ID): `"Invalid user ID format."`
        *   Error (Not Admin): `"Only admins can kick users"` (or similar permission error)
        *   Error (Not Found/Not In Room/Other): `"Error kicking user <user_id> from room <room_id>: <error_details>"`
    *   **State Changes:** Removes the specified user from the `Players` list of the `Room` entity. Publishes `PlayerKickedEvent`.

7.  **Command:** `/delete_room`
    *   **Description:** Initiates process to delete a specific room (shows selection).
    *   **User Role:** Admin Only
    *   **Input:** None initially.
    *   **Output/Behavior:**
        *   Success (Rooms Exist): Sends message `"Select a room to delete:"` with inline buttons for each room.
        *   Success (No Rooms): `"No rooms exist to delete."`
        *   Error (Not Admin): `"You are not authorized to use this command."`
        *   Error (Fetch): `"Failed to fetch rooms list."`
    *   **State Changes:** None initially. Subsequent callback performs deletion.

8.  **Command:** `/assign_scenario <room_id> <scenario_id>`
    *   **Description:** Assign a game scenario to a room. This action **also creates** the corresponding Game entity.
    *   **User Role:** Admin Only
    *   **Input:** `<room_id>` (string, valid Room ID), `<scenario_id>` (string, valid Scenario ID).
    *   **Output/Behavior:**
        *   Success: `"Successfully assigned scenario '<ScenarioName>' (ID: <ScenarioID>) to room '<RoomName>' (ID: <RoomID>) and created game '<GameID>'"`
        *   Error (Usage): `"Usage: /assign_scenario <room_id> <scenario_id>"`
        *   Error (Not Admin): `"You are not authorized to use this command."`
        *   Error (Scenario Not Found): `"Error finding scenario '<scenario_id>': <error_details>"`
        *   Error (Room Not Found): `"Error finding room '<room_id>': <error_details>"`
        *   Error (Assignment): `"Error assigning scenario ID '<scenario_id>' to room '<room_id>': <error_details>"`
        *   Error (Game Creation): `"Scenario assigned, but failed to create game: <error_details>"`
    *   **State Changes:**
        *   Updates the `ScenarioName` (or similar reference) field in the `Room` entity.
        *   Creates a new `Game` entity associated with the `Room` and `Scenario`. Sets initial state to `WaitingForPlayers`. Stores the new `Game` in the repository.

## 2.3. Scenario Management Commands

1.  **Command:** `/create_scenario <scenario_name>`
    *   **Description:** Create a new game scenario (ruleset).
    *   **User Role:** Admin Only
    *   **Input:** `<scenario_name>` (string).
    *   **Output/Behavior:**
        *   Success: `"Scenario '<scenario_name>' created successfully! ID: <scenario_id>\nUse /add_role <scenario_id> <role_name> to add roles."`
        *   Error (No Name): `"Please provide a scenario name: /create_scenario [name]"`
        *   Error (Not Admin): `"You are not authorized to use this command."`
        *   Error (Exists/Other): `"Error creating scenario: <error_details>"`
    *   **State Changes:** Creates a new `Scenario` entity with an empty `Roles` list in the repository.

2.  **Command:** `/delete_scenario <scenario_id>`
    *   **Description:** Delete a game scenario.
    *   **User Role:** Admin Only
    *   **Input:** `<scenario_id>` (string, valid Scenario ID).
    *   **Output/Behavior:**
        *   Success: `"Scenario <scenario_id> deleted successfully!"`
        *   Error (No ID): `"Please provide a scenario ID: /delete_scenario <id>"`
        *   Error (Not Admin): `"You are not authorized to use this command."`
        *   Error (Not Found/Other): `"Error deleting scenario '<scenario_id>': <error_details>"`
    *   **State Changes:** Deletes the `Scenario` entity from the repository.

3.  **Command:** `/add_role <scenario_id> <role_name>`
    *   **Description:** Add a role to a scenario.
    *   **User Role:** Admin Only
    *   **Input:** `<scenario_id>` (string), `<role_name>` (string).
    *   **Output/Behavior:**
        *   Success: `"Role '<role_name>' added to scenario <scenario_id> successfully!"`
        *   Error (Usage): `"Usage: /add_role <scenario_id> <role_name>"`
        *   Error (Not Admin): `"You are not authorized to use this command."`
        *   Error (Not Found/Other): `"Error adding role '<role_name>' to scenario '<scenario_id>': <error_details>"`
    *   **State Changes:** Adds the specified `Role` to the `Roles` list of the `Scenario` entity in the repository.

4.  **Command:** `/remove_role <scenario_id> <role_name>`
    *   **Description:** Remove a role from a scenario.
    *   **User Role:** Admin Only
    *   **Input:** `<scenario_id>` (string), `<role_name>` (string).
    *   **Output/Behavior:**
        *   Success: `"Role '<role_name>' removed from scenario <scenario_id> successfully!"`
        *   Error (Usage): `"Usage: /remove_role <scenario_id> <role_name>"`
        *   Error (Not Admin): `"You are not authorized to use this command."`
        *   Error (Not Found/Other): `"Error removing role '<role_name>' from scenario '<scenario_id>': <error_details>"`
    *   **State Changes:** Removes the specified `Role` (first match by name) from the `Roles` list of the `Scenario` entity.

5.  **Command:** `/list_scenarios`
    *   **Description:** List available scenarios (Currently TODO in existing code).
    *   **User Role:** Admin Only
    *   **Input:** None
    *   **Output/Behavior:**
        *   Success: List scenarios with IDs, names, and roles.
        *   No Scenarios: Message indicating none exist.
        *   Error: Error message.
    *   **State Changes:** None.

## 2.4. Game Management Commands

1.  **Command:** `/games`
    *   **Description:** List active games.
    *   **User Role:** Admin Only
    *   **Input:** None
    *   **Output/Behavior:**
        *   Success: Lists active games with `GameID`, `RoomID`, `ScenarioID`, `State`, and `Player/Assignment Count`. May also show assignments if present.
        *   No Games: `"No active games found."`
        *   Error (Not Admin): `"You are not authorized to use this command."`
        *   Error (Fetch): `"Error fetching games list: <error_details>"`
    *   **State Changes:** None.

2.  **Command:** `/assign_roles <game_id>`
    *   **Description:** Assign roles from the game's scenario to players currently in the game's room.
    *   **User Role:** Admin Only
    *   **Input:** `<game_id>` (string, valid Game ID).
    *   **Output/Behavior:**
        *   Success: `"Roles assigned for game <game_id>:\n<User1>: <Role1>\n<User2>: <Role2>..."` (Response shows assignments, actual roles sent privately via callback/separate step).
        *   Error (No ID): `"Please provide a game ID: /assign_roles <game_id>"`
        *   Error (Not Admin): `"You are not authorized to use this command."`
        *   Error (Game Not Found): `"game '<game_id>' not found: <error_details>"`
        *   Error (No Scenario): `"game has no scenario assigned"`
        *   Error (No Room): `"game has no room assigned"`
        *   Error (Player/Role Mismatch): `"role count (X) does not match player count (Y) for game '<game_id>'"`
        *   Error (Fetch/Update/Other): `"Error assigning roles for game '<game_id>': <error_details>"`
    *   **State Changes:**
        *   Populates the `Assignments` map within the `Game` entity (mapping `UserID` to `Role`).
        *   Updates the `Game` entity's `State` to `RolesAssigned`.
        *   Persists the updated `Game` entity.

## 2.5. Callback Interactions (Examples)

*   **Delete Room Confirmation:**
    *   Trigger: User clicks inline button from `/delete_room`.
    *   Behavior: Bot edits message to `"Are you sure you want to delete room <room_id>?"` with "Yes, delete it!" and "Cancel" buttons.
    *   Confirm Click: Deletes room, responds `"Room <room_id> deleted."`, edits message to confirm.
    *   Cancel Click: Deletes confirmation message, responds `"Operation cancelled."`. 
*   **Leave Room Confirmation:**
    *   Trigger: User clicks inline button from `/join_room` success message.
    *   Behavior: Edits message to `"Are you sure you want to leave room <room_id>?"` with "Yes, leave" and "Cancel".
    *   Confirm Click: Removes user, responds `"You left room <room_id>."`, edits message.
    *   Cancel Click: Deletes message, responds `"Operation cancelled."`. 
*   **(Optional) Assign Roles Confirmation/Send:**
    *   Trigger: Could be callback button after `/assign_roles` success message.
    *   Behavior: Fetches the game, iterates through assignments, sends each player their role via private message. Responds `"Roles sent to X players!"`, edits original message.

*(Note: This specification is based on the commands observed in the `README.md` and `handlers.go` of the provided codebase. Features marked TODO or implied logic might require further clarification.)* 