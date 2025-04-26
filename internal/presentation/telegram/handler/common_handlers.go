package telegram

import (
	"fmt"
	"gopkg.in/telebot.v3"
	room "telemafia/internal/presentation/telegram/handler/room"
)

// HandleHelp provides a simple help message.
func HandleHelp(h *BotHandler, c telebot.Context) error {
	help := `Available commands:
/start - Show welcome message & rooms
/help - Show this help message
/list_rooms - List all available rooms
/my_rooms - List rooms you have joined
/join_room <room_id> - Join a specific room
/leave_room <room_id> - Leave the specified room

Admin Commands:
/create_room <room_name> - Create a new room
/delete_room - Select a room to delete
/kick_user <room_id> <user_id> - Kick a user from a room
/create_scenario <scenario_name> - Create a new game scenario
/delete_scenario <scenario_id> - Delete a scenario
/add_role <scenario_id> <role_name> - Add a role to a scenario
/remove_role <scenario_id> <role_name> - Remove a role from a scenario
/assign_scenario <room_id> <scenario_id> - Assign a scenario to a room (creates a game)
/games - List active games and their status
/assign_roles <game_id> - Assign roles to players in a game`
	return c.Send(help, &telebot.SendOptions{DisableWebPagePreview: true})
}

func HandleStart(h *BotHandler, c telebot.Context) error {
	_ = c.Send(fmt.Sprintf("Welcome, %s!", c.Sender().Username))
	return room.HandleListRooms(h.getRoomsHandler, h.getPlayersInRoomHandler, c)
}
