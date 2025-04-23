package telegram

import (
	"context"
	"fmt"
	"log"

	roomEntity "telemafia/internal/domain/room/entity"
	roomCommand "telemafia/internal/domain/room/usecase/command"

	"gopkg.in/telebot.v3"
)

// handleLeaveRoomSelectCallback asks for confirmation to leave.
func handleLeaveRoomSelectCallback(h *BotHandler, c telebot.Context, roomIDStr string) error {
	markup := &telebot.ReplyMarkup{}
	confirmData := fmt.Sprintf("%s:%s", UniqueLeaveRoomConfirm, roomIDStr)
	btnConfirm := markup.Data("Yes, leave", confirmData)
	btnCancel := markup.Data("Cancel", UniqueCancel)
	markup.Inline(markup.Row(btnConfirm, btnCancel))

	err := c.Edit(fmt.Sprintf("Are you sure you want to leave room %s?", roomIDStr), markup)
	if err != nil {
		log.Printf("Error editing message for leave confirmation: %v", err)
		_ = c.Respond(&telebot.CallbackResponse{Text: "Error showing confirmation."})
		return err
	}
	_ = c.Respond()
	return nil
}

// handleLeaveRoomConfirmCallback performs leaving the room.
func handleLeaveRoomConfirmCallback(h *BotHandler, c telebot.Context, roomIDStr string) error {
	user := ToUser(c.Sender())
	if user == nil {
		// Maybe just respond to callback, editing original message might not be possible
		_ = c.Respond(&telebot.CallbackResponse{Text: "Could not identify user.", ShowAlert: true})
		return fmt.Errorf("could not identify user for leave confirm callback")
	}
	cmd := roomCommand.LeaveRoomCommand{
		RoomID:    roomEntity.RoomID(roomIDStr),
		Requester: *user,
	}
	if err := h.leaveRoomHandler.Handle(context.Background(), cmd); err != nil {
		_ = c.Respond(&telebot.CallbackResponse{Text: fmt.Sprintf("Error leaving room: %v", err), ShowAlert: true})
		return err
	}

	_ = c.Respond(&telebot.CallbackResponse{Text: fmt.Sprintf("You left room %s.", roomIDStr)})
	_ = c.Edit(fmt.Sprintf("You have left room %s.", roomIDStr), &telebot.ReplyMarkup{})
	return nil
}

// handleDeleteRoomSelectCallback asks for confirmation.
func handleDeleteRoomSelectCallback(h *BotHandler, c telebot.Context, roomIDStr string) error {
	markup := &telebot.ReplyMarkup{}
	confirmData := fmt.Sprintf("%s:%s", UniqueDeleteRoomConfirm, roomIDStr)
	btnConfirm := markup.Data("Yes, delete it!", confirmData)
	btnCancel := markup.Data("Cancel", UniqueCancel)
	markup.Inline(markup.Row(btnConfirm, btnCancel))

	err := c.Edit(fmt.Sprintf("Are you sure you want to delete room %s?", roomIDStr), markup)
	if err != nil {
		log.Printf("Error editing message for delete confirmation: %v", err)
		_ = c.Respond(&telebot.CallbackResponse{Text: "Error showing confirmation."})
		return err
	}
	_ = c.Respond()
	return nil
}

// handleDeleteRoomConfirmCallback is called when the admin confirms deletion
func handleDeleteRoomConfirmCallback(h *BotHandler, c telebot.Context, roomIDStr string) error {
	roomID := roomEntity.RoomID(roomIDStr)
	requester := ToUser(c.Sender())
	if requester == nil {
		_ = c.Respond(&telebot.CallbackResponse{Text: "Could not identify user.", ShowAlert: true})
		return fmt.Errorf("could not identify user for delete confirm callback")
	}

	cmd := roomCommand.DeleteRoomCommand{
		Requester: *requester,
		RoomID:    roomID,
	}

	if err := h.deleteRoomHandler.Handle(context.Background(), cmd); err != nil {
		log.Printf("Error deleting room %s after confirmation: %v", roomID, err)
		// Try to edit original message, fallback to simple response if edit fails
		_ = c.Edit(fmt.Sprintf("Error deleting room %s: %v", roomID, err))
		_ = c.Respond(&telebot.CallbackResponse{Text: fmt.Sprintf("Error deleting room: %v", err), ShowAlert: true})
		return err
	}

	// Try to edit original message, fallback to simple response if edit fails
	_ = c.Edit(fmt.Sprintf("Room %s deleted successfully!", roomID), &telebot.ReplyMarkup{})
	_ = c.Respond(&telebot.CallbackResponse{Text: fmt.Sprintf("Room %s deleted!", roomID)})
	return nil
}

// handleJoinRoomCallback handles the callback for joining a room.
func handleJoinRoomCallback(h *BotHandler, c telebot.Context, roomIDStr string) error {
	roomID := roomEntity.RoomID(roomIDStr)
	user := ToUser(c.Sender())
	if user == nil {
		_ = c.Respond(&telebot.CallbackResponse{Text: "Could not identify user.", ShowAlert: true})
		return fmt.Errorf("could not identify user for join callback")
	}

	cmd := roomCommand.JoinRoomCommand{
		Requester: *user,
		RoomID:    roomID,
	}

	if err := h.joinRoomHandler.Handle(context.Background(), cmd); err != nil {
		_ = c.Respond(&telebot.CallbackResponse{Text: fmt.Sprintf("Error: %v", err), ShowAlert: true})
		return err
	}

	_ = c.Respond(&telebot.CallbackResponse{Text: fmt.Sprintf("Joined room %s!", roomID)})
	return nil
}
