# 2. Feature Specification

**Goal:** Define the functional requirements of the bot from the user's perspective, detailing commands, inputs, outputs, permissions, and expected state changes.

## 2.1. General Commands

1.  **Command:** `/start`
    *   **Description:** Start interaction with the bot, show welcome message and available rooms.
    *   **User Role:** Everyone
    *   **Input:** None
    *   **Output/Behavior:** Sends a welcome message (`"Welcome, <username>!"`). Sends/Updates a dynamic message with the text `"Available Rooms:\n"` and **inline buttons** for each room. Each button text **MUST** be `"<Name> (<PlayerCount>/<MaxPlayers>)"` and its callback data **MUST** be `"UniqueJoinRoom:<RoomID>"`.
    *   **State Changes:** Registers the user's chat for dynamic updates.

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
        *   Sends a message with text `"Available Rooms:\n"` and **inline buttons** for each room. Each button text **MUST** be `"<Name> (<PlayerCount>/<MaxPlayers>)"` and its callback data **MUST** be `"UniqueJoinRoom:<RoomID>"`.
        *   No Rooms: `"No rooms available."` (without buttons)
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
        *   Success: `