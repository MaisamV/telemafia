package telegram

import (
	"context"
	"fmt"
	"log"

	roomEntity "telemafia/internal/domain/room/entity"
	roomCommand "telemafia/internal/domain/room/usecase/command"
	roomQuery "telemafia/internal/domain/room/usecase/query"
	tgutil "telemafia/internal/shared/tgutil"

	"gopkg.in/telebot.v3"
)

// HandleLeaveRoomSelectCallback shows confirmation for leaving a room
func HandleLeaveRoomSelectCallback(getPlayerRoomsHandler *roomQuery.GetPlayerRoomsHandler, c telebot.Context, data string) error {
	// The data here could be the roomID if the button includes it
	// Or, if it's a generic /leave_room, query the rooms the user is in
	roomIDStr := data
	requester := tgutil.ToUser(c.Sender())
	if requester == nil {
		return c.Respond(&telebot.CallbackResponse{Text: "Could not identify requester.", ShowAlert: true})
	}

	var targetRoomID roomEntity.RoomID
	if roomIDStr != "" {
		targetRoomID = roomEntity.RoomID(roomIDStr)
		// Optional: Fetch room name for display
	} else {
		// If no room ID in data, query user's rooms (requires getPlayerRoomsHandler)
		query := roomQuery.GetPlayerRoomsQuery{PlayerID: requester.ID}
		rooms, err := getPlayerRoomsHandler.Handle(context.Background(), query)
		if err != nil || len(rooms) == 0 {
			_ = c.Respond(&telebot.CallbackResponse{Text: "You are not in any rooms or failed to fetch them."})
			return c.Edit("No rooms to leave.")
		}
		if len(rooms) == 1 {
			targetRoomID = rooms[0].ID
		} else {
			// TODO: Handle multiple rooms - show selection
			_ = c.Respond(&telebot.CallbackResponse{Text: "Please specify which room to leave via /leave_room <id>"})
			return c.Edit("Multiple rooms found.")
		}
	}

	markup := &telebot.ReplyMarkup{}
	btnConfirm := markup.Data("Yes, leave", tgutil.UniqueLeaveRoomConfirm, string(targetRoomID))
	btnCancel := markup.Data("Cancel", tgutil.UniqueCancel, "")
	markup.Inline(markup.Row(btnConfirm, btnCancel))

	// TODO: Get room name if needed for the prompt
	return c.Edit(fmt.Sprintf("Are you sure you want to leave room %s?", targetRoomID), markup)
}

// HandleLeaveRoomConfirmCallback performs the actual leaving action
func HandleLeaveRoomConfirmCallback(leaveRoomHandler *roomCommand.LeaveRoomHandler, c telebot.Context, data string) error {
	roomID := roomEntity.RoomID(data)
	requester := tgutil.ToUser(c.Sender())
	if requester == nil {
		return c.Respond(&telebot.CallbackResponse{Text: "Could not identify requester.", ShowAlert: true})
	}

	cmd := roomCommand.LeaveRoomCommand{
		Requester: *requester,
		RoomID:    roomID,
	}

	if err := leaveRoomHandler.Handle(context.Background(), cmd); err != nil {
		log.Printf("Error leaving room '%s': %v", roomID, err)
		_ = c.Respond(&telebot.CallbackResponse{Text: fmt.Sprintf("Error leaving: %v", err), ShowAlert: true})
		return c.Edit("Failed to leave room.")
	}

	_ = c.Respond(&telebot.CallbackResponse{Text: "You have left the room."}) // Respond to callback
	return c.Edit(fmt.Sprintf("You left room %s.", roomID))                   // Edit original message
}

// HandleDeleteRoomSelectCallback shows confirmation for deleting a room
func HandleDeleteRoomSelectCallback(getRoomHandler *roomQuery.GetRoomHandler, c telebot.Context, data string) error {
	roomID := roomEntity.RoomID(data)

	// Fetch room to display name (optional but good UX)
	room, err := getRoomHandler.Handle(context.Background(), roomQuery.GetRoomQuery{RoomID: roomID})
	roomName := string(roomID)
	if err == nil && room != nil {
		roomName = room.Name
	}

	markup := &telebot.ReplyMarkup{}
	btnConfirm := markup.Data("Yes, delete it!", tgutil.UniqueDeleteRoomConfirm, data)
	btnCancel := markup.Data("Cancel", tgutil.UniqueCancel, "")
	markup.Inline(markup.Row(btnConfirm, btnCancel))

	return c.Edit(fmt.Sprintf("Are you sure you want to delete room '%s' (%s)?", roomName, roomID), markup)
}

// HandleDeleteRoomConfirmCallback performs the actual room deletion
func HandleDeleteRoomConfirmCallback(deleteRoomHandler *roomCommand.DeleteRoomHandler, c telebot.Context, data string) error {
	roomID := roomEntity.RoomID(data)
	requester := tgutil.ToUser(c.Sender()) // Need requester for admin check within handler
	if requester == nil {
		return c.Respond(&telebot.CallbackResponse{Text: "Could not identify requester.", ShowAlert: true})
	}

	cmd := roomCommand.DeleteRoomCommand{
		Requester: *requester,
		RoomID:    roomID,
	}

	if err := deleteRoomHandler.Handle(context.Background(), cmd); err != nil {
		log.Printf("Error deleting room '%s': %v", roomID, err)
		_ = c.Respond(&telebot.CallbackResponse{Text: fmt.Sprintf("Error deleting room: %v", err), ShowAlert: true})
		return c.Edit("Failed to delete room.") // Edit original message
	}

	_ = c.Respond(&telebot.CallbackResponse{Text: "Room deleted!"})
	return c.Edit(fmt.Sprintf("Room %s deleted successfully!", roomID))
}

// HandleJoinRoomCallback handles the inline join button click
func HandleJoinRoomCallback(joinRoomHandler *roomCommand.JoinRoomHandler, c telebot.Context, data string) error {
	roomID := roomEntity.RoomID(data)
	user := tgutil.ToUser(c.Sender())
	if user == nil {
		return c.Respond(&telebot.CallbackResponse{Text: "Could not identify user.", ShowAlert: true})
	}

	cmd := roomCommand.JoinRoomCommand{
		Requester: *user,
		RoomID:    roomID,
	}

	if err := joinRoomHandler.Handle(context.Background(), cmd); err != nil {
		log.Printf("Error handling join room callback for room '%s': %v", roomID, err)
		return c.Respond(&telebot.CallbackResponse{Text: fmt.Sprintf("Error joining room: %v", err), ShowAlert: true})
	}

	// Respond to callback
	_ = c.Respond(&telebot.CallbackResponse{Text: "Joined successfully!"})

	// Edit the original message (if possible) to remove the join button or update state
	// Optional: Send a confirmation message
	// _ = c.Send(fmt.Sprintf("You successfully joined room %s", roomID))

	return nil // Callback handled
}
