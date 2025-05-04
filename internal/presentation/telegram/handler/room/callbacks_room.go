package telegram

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"gopkg.in/telebot.v3"

	roomEntity "telemafia/internal/domain/room/entity"
	roomCommand "telemafia/internal/domain/room/usecase/command"
	roomQuery "telemafia/internal/domain/room/usecase/query"
	messages "telemafia/internal/presentation/telegram/messages"
	sharedEntity "telemafia/internal/shared/entity"
	tgutil "telemafia/internal/shared/tgutil"
)

// RefreshNotifier defines an interface for triggering a refresh.
// *tgutil.RefreshingMessageBook satisfies this interface.
type RefreshNotifier interface {
	RaiseRefreshNeeded()
	AddActiveMessage(chatID int64, msg *tgutil.RefreshingMessage)
	RemoveActiveMessage(chatID int64)
	GetActiveMessage(chatID int64) (*tgutil.RefreshingMessage, bool)
}

// HandleLeaveRoomSelectCallback shows confirmation for leaving a room
func HandleLeaveRoomSelectCallback(
	leaveRoomHandler *roomCommand.LeaveRoomHandler,
	getRoomsHandler *roomQuery.GetRoomsHandler,
	getPlayersInRoomHandler *roomQuery.GetPlayersInRoomHandler,
	roomList RefreshNotifier,
	roomDetail RefreshNotifier,
	c telebot.Context,
	data string,
	msgs *messages.Messages,
) error {
	roomID := roomEntity.RoomID(data)
	requester := tgutil.ToUser(c.Sender())
	if requester == nil {
		return c.Respond(&telebot.CallbackResponse{Text: msgs.Common.ErrorIdentifyRequester, ShowAlert: true})
	}

	cmd := roomCommand.LeaveRoomCommand{
		Requester: *requester,
		RoomID:    roomID,
	}

	if err := leaveRoomHandler.Handle(context.Background(), cmd); err != nil {
		log.Printf("Error leaving room '%s': %v", roomID, err)
		return c.Respond(&telebot.CallbackResponse{Text: fmt.Sprintf(msgs.Common.CallbackErrorGeneric, err), ShowAlert: true})
	}

	chatID := c.Sender().ID
	roomDetail.RemoveActiveMessage(chatID)
	roomList.AddActiveMessage(chatID, &tgutil.RefreshingMessage{
		MessageID: c.Message().ID,
		ChatID:    c.Message().Chat.ID,
		Data:      string(roomID),
	})
	roomList.RaiseRefreshNeeded()
	roomDetail.RaiseRefreshNeeded()
	message, markup, err := PrepareRoomListMessage(getRoomsHandler, getPlayersInRoomHandler, msgs)
	if err != nil {
		return err
	}
	return c.Edit(message, markup)
}

// HandleLeaveRoomConfirmCallback performs the actual leaving action
func HandleLeaveRoomConfirmCallback(
	leaveRoomHandler *roomCommand.LeaveRoomHandler,
	refreshNotifier RefreshNotifier,
	c telebot.Context,
	data string,
	msgs *messages.Messages,
) error {
	roomID := roomEntity.RoomID(data)
	requester := tgutil.ToUser(c.Sender())
	if requester == nil {
		return c.Respond(&telebot.CallbackResponse{Text: msgs.Common.ErrorIdentifyRequester, ShowAlert: true})
	}

	cmd := roomCommand.LeaveRoomCommand{
		Requester: *requester,
		RoomID:    roomID,
	}

	if err := leaveRoomHandler.Handle(context.Background(), cmd); err != nil {
		log.Printf("Error leaving room '%s': %v", roomID, err)
		_ = c.Respond(&telebot.CallbackResponse{Text: fmt.Sprintf(msgs.Common.CallbackErrorGeneric, err), ShowAlert: true})
		return c.Edit(msgs.Room.LeaveCallbackEditFail)
	}

	refreshNotifier.RaiseRefreshNeeded()
	_ = c.Respond(&telebot.CallbackResponse{Text: msgs.Room.LeaveCallbackSuccess})
	return c.Edit(fmt.Sprintf(msgs.Room.LeaveCallbackEditSuccess, roomID))
}

// HandleDeleteRoomSelectCallback shows confirmation for deleting a room
func HandleDeleteRoomSelectCallback(
	getRoomHandler *roomQuery.GetRoomHandler,
	c telebot.Context,
	data string,
	msgs *messages.Messages,
) error {
	roomID := roomEntity.RoomID(data)

	room, err := getRoomHandler.Handle(context.Background(), roomQuery.GetRoomQuery{RoomID: roomID})
	roomName := string(roomID)
	if err == nil && room != nil {
		roomName = room.Name
	}

	markup := &telebot.ReplyMarkup{}
	btnConfirm := markup.Data(msgs.Room.DeleteConfirmButton, tgutil.UniqueDeleteRoomConfirm, data)
	btnCancel := markup.Data(msgs.Room.DeleteCancelButton, tgutil.UniqueCancel, "")
	markup.Inline(markup.Row(btnConfirm, btnCancel))

	return c.Edit(fmt.Sprintf(msgs.Room.DeletePromptConfirm, roomName, roomID), markup)
}

// HandleDeleteRoomConfirmCallback performs the actual room deletion
func HandleDeleteRoomConfirmCallback(
	deleteRoomHandler *roomCommand.DeleteRoomHandler,
	refreshNotifier RefreshNotifier,
	c telebot.Context,
	data string,
	msgs *messages.Messages,
) error {
	roomID := roomEntity.RoomID(data)
	requester := tgutil.ToUser(c.Sender())
	if requester == nil {
		return c.Respond(&telebot.CallbackResponse{Text: msgs.Common.ErrorIdentifyRequester, ShowAlert: true})
	}

	cmd := roomCommand.DeleteRoomCommand{
		Requester: *requester,
		RoomID:    roomID,
	}

	if err := deleteRoomHandler.Handle(context.Background(), cmd); err != nil {
		log.Printf("Error deleting room '%s': %v", roomID, err)
		_ = c.Respond(&telebot.CallbackResponse{Text: fmt.Sprintf(msgs.Room.DeleteCallbackError, err), ShowAlert: true})
		return c.Edit(msgs.Room.DeleteCallbackEditFail)
	}

	refreshNotifier.RaiseRefreshNeeded()
	_ = c.Respond(&telebot.CallbackResponse{Text: msgs.Room.DeleteCallbackSuccess})
	return c.Edit(fmt.Sprintf(msgs.Room.DeleteCallbackEditSuccess, roomID))
}

// HandleJoinRoomCallback handles the inline join button click
func HandleJoinRoomCallback(
	joinRoomHandler *roomCommand.JoinRoomHandler,
	getRoomsHandler *roomQuery.GetRoomsHandler,
	getPlayersInRoomHandler *roomQuery.GetPlayersInRoomHandler,
	roomList RefreshNotifier,
	roomDetail RefreshNotifier,
	c telebot.Context,
	data string,
	msgs *messages.Messages,
) error {
	roomID := roomEntity.RoomID(data)
	user := tgutil.ToUser(c.Sender())
	if user == nil {
		return c.Respond(&telebot.CallbackResponse{Text: msgs.Common.ErrorIdentifyUser, ShowAlert: true})
	}

	cmd := roomCommand.JoinRoomCommand{
		Requester: *user,
		RoomID:    roomID,
	}

	if err := joinRoomHandler.Handle(context.Background(), cmd); err != nil {
		log.Printf("Error handling join room callback for room '%s': %v", roomID, err)
		return c.Respond(&telebot.CallbackResponse{Text: fmt.Sprintf(msgs.Common.CallbackErrorGeneric, err), ShowAlert: true})
	}
	chatID := c.Sender().ID
	roomList.RemoveActiveMessage(chatID)
	roomDetail.AddActiveMessage(chatID, &tgutil.RefreshingMessage{
		MessageID: c.Message().ID,
		ChatID:    c.Message().Chat.ID,
		Data:      string(roomID),
	})
	roomList.RaiseRefreshNeeded()
	roomDetail.RaiseRefreshNeeded()
	_ = c.Respond(&telebot.CallbackResponse{Text: msgs.Room.JoinSuccess})
	message, markup, err := RoomDetailMessage(getRoomsHandler, getPlayersInRoomHandler, msgs, user.Admin, data)
	if err != nil {
		return err
	}
	return c.Edit(message, markup, telebot.ModeMarkdownV2, telebot.NoPreview)
}

// HandleKickUserSelectCallback shows the list of users to kick from a room.
func HandleKickUserSelectCallback(
	getPlayersHandler *roomQuery.GetPlayersInRoomHandler,
	c telebot.Context,
	roomIDStr string, // Room ID passed as data
	msgs *messages.Messages,
) error {
	roomID := roomEntity.RoomID(roomIDStr)
	requester := tgutil.ToUser(c.Sender()) // Assumes this is called by an admin
	if requester == nil {
		return c.Respond(&telebot.CallbackResponse{Text: msgs.Common.ErrorIdentifyRequester, ShowAlert: true})
	}

	// Fetch players in the room
	players, err := getPlayersHandler.Handle(context.Background(), roomQuery.GetPlayersInRoomQuery{RoomID: roomID})
	if err != nil {
		log.Printf("KickUserSelect: Error fetching players for room '%s': %v", roomID, err)
		return c.Respond(&telebot.CallbackResponse{Text: fmt.Sprintf(msgs.Common.CallbackErrorGeneric, err), ShowAlert: true})
	}

	markup := &telebot.ReplyMarkup{}
	var userRows []telebot.Row
	playersToKickCount := 0

	for _, player := range players {
		// Don't list the admin themself
		if player.ID == requester.ID {
			continue
		}
		playersToKickCount++
		// Create payload: roomID|userIDToKick
		payload := fmt.Sprintf("%s|%d", roomIDStr, player.ID)
		btn := markup.Data(player.FirstName, tgutil.UniqueKickUserConfirm, payload)
		userRows = append(userRows, markup.Row(btn))
	}

	if playersToKickCount == 0 {
		_ = c.Respond(&telebot.CallbackResponse{Text: msgs.Room.KickUserNoPlayers})
		// Optionally edit the message back to the standard room detail?
		// For now, just respond and leave the message as is.
		return nil
	}

	// Add cancel button
	cancelPayload := roomIDStr                                                                  // Cancel goes back to room detail view
	cancelBtn := markup.Data(msgs.Room.LeaveCancelButton, tgutil.UniqueJoinRoom, cancelPayload) // Re-use join unique to show detail
	userRows = append(userRows, markup.Row(cancelBtn))

	markup.Inline(userRows...)

	prompt := fmt.Sprintf(msgs.Room.KickUserSelectPrompt, roomIDStr) // Room name would be better
	return c.Edit(prompt, markup)
}

// HandleKickUserConfirmCallback handles the selection of a user to kick.
func HandleKickUserConfirmCallback(
	kickUserHandler *roomCommand.KickUserHandler,
	getRoomsHandler *roomQuery.GetRoomsHandler, // Need these to reconstruct RoomDetailMessage
	getPlayersHandler *roomQuery.GetPlayersInRoomHandler,
	roomList RefreshNotifier,
	roomDetail RefreshNotifier,
	c telebot.Context,
	data string, // Payload: roomID|userIDToKick
	msgs *messages.Messages,
) error {
	requester := tgutil.ToUser(c.Sender()) // Assumes this is called by an admin
	if requester == nil {
		return c.Respond(&telebot.CallbackResponse{Text: msgs.Common.ErrorIdentifyRequester, ShowAlert: true})
	}

	// Parse payload
	parts := strings.Split(data, "|")
	if len(parts) != 2 {
		log.Printf("KickUserConfirm: Invalid payload format: %s", data)
		return c.Respond(&telebot.CallbackResponse{Text: fmt.Sprintf(msgs.Common.CallbackErrorGeneric, "invalid data"), ShowAlert: true})
	}
	roomIDStr := parts[0]
	userIDToKickStr := parts[1]

	roomID := roomEntity.RoomID(roomIDStr)
	userIDToKick, err := strconv.ParseInt(userIDToKickStr, 10, 64)
	if err != nil {
		log.Printf("KickUserConfirm: Invalid user ID in payload: %s, error: %v", userIDToKickStr, err)
		return c.Respond(&telebot.CallbackResponse{Text: fmt.Sprintf(msgs.Common.CallbackErrorGeneric, "invalid user ID"), ShowAlert: true})
	}

	// Call the use case
	kickCmd := roomCommand.KickUserCommand{
		Requester: *requester,
		RoomID:    roomID,
		PlayerID:  sharedEntity.UserID(userIDToKick),
	}
	if err := kickUserHandler.Handle(context.Background(), kickCmd); err != nil {
		log.Printf("KickUserConfirm: Error kicking user %d from room %s: %v", userIDToKick, roomID, err)
		return c.Respond(&telebot.CallbackResponse{Text: fmt.Sprintf(msgs.Room.KickUserCallbackError, err), ShowAlert: true})
	}

	// Success - trigger refreshes and edit back to room detail
	roomList.RaiseRefreshNeeded()
	roomDetail.RaiseRefreshNeeded()

	// Acknowledge the callback first
	// User's name isn't readily available here without another fetch, use ID for now
	ackMsg := fmt.Sprintf(msgs.Room.KickUserCallbackSuccess, userIDToKickStr, roomIDStr)
	_ = c.Respond(&telebot.CallbackResponse{Text: ackMsg})

	// Prepare and edit the message back to the standard room detail
	// Note: We pass requester.Admin which should be true here
	message, markup, err := RoomDetailMessage(getRoomsHandler, getPlayersHandler, msgs, requester.Admin, roomIDStr)
	if err != nil {
		log.Printf("KickUserConfirm: Error preparing room detail after kick for room '%s': %v", roomID, err)
		// Can't easily recover the message here, just log
		return nil
	}

	return c.Edit(message, markup, telebot.ModeMarkdownV2, telebot.NoPreview)
}

// HandleChangeModeratorSelectCallback shows the list of users who can become the new moderator.
func HandleChangeModeratorSelectCallback(
	getPlayersHandler *roomQuery.GetPlayersInRoomHandler,
	c telebot.Context,
	roomIDStr string, // Room ID passed as data
	msgs *messages.Messages,
) error {
	roomID := roomEntity.RoomID(roomIDStr)
	requester := tgutil.ToUser(c.Sender()) // Assumes this is called by an admin
	if requester == nil {
		return c.Respond(&telebot.CallbackResponse{Text: msgs.Common.ErrorIdentifyRequester, ShowAlert: true})
	}

	// Fetch players in the room
	players, err := getPlayersHandler.Handle(context.Background(), roomQuery.GetPlayersInRoomQuery{RoomID: roomID})
	if err != nil {
		log.Printf("ChangeModSelect: Error fetching players for room '%s': %v", roomID, err)
		return c.Respond(&telebot.CallbackResponse{Text: fmt.Sprintf(msgs.Common.CallbackErrorGeneric, err), ShowAlert: true})
	}

	markup := &telebot.ReplyMarkup{}
	var userRows []telebot.Row
	moderatorCandidatesCount := 0

	for _, player := range players {
		moderatorCandidatesCount++
		// Create payload: roomID|userIDToMakeModerator
		payload := fmt.Sprintf("%s|%d", roomIDStr, player.ID)
		btn := markup.Data(player.FirstName, tgutil.UniqueChangeModeratorConfirm, payload)
		userRows = append(userRows, markup.Row(btn))
	}

	if moderatorCandidatesCount == 0 {
		_ = c.Respond(&telebot.CallbackResponse{Text: msgs.Room.ChangeModeratorNoCandidates})
		return nil // Stay on the same message (room detail)
	}

	// Add cancel button
	cancelPayload := roomIDStr                                                                  // Cancel goes back to room detail view
	cancelBtn := markup.Data(msgs.Room.LeaveCancelButton, tgutil.UniqueJoinRoom, cancelPayload) // Re-use join unique to show detail
	userRows = append(userRows, markup.Row(cancelBtn))

	markup.Inline(userRows...)

	prompt := fmt.Sprintf(msgs.Room.ChangeModeratorSelectPrompt, roomIDStr) // Room name would be better
	return c.Edit(prompt, markup)
}

// HandleChangeModeratorConfirmCallback handles the selection of a user to make moderator.
func HandleChangeModeratorConfirmCallback(
	changeModeratorHandler *roomCommand.ChangeModeratorHandler, // Inject the new use case handler
	getRoomsHandler *roomQuery.GetRoomsHandler,
	getPlayersHandler *roomQuery.GetPlayersInRoomHandler, // Needed to fetch the user details
	roomList RefreshNotifier,
	roomDetail RefreshNotifier,
	c telebot.Context,
	data string, // Payload: roomID|userIDToMakeModerator
	msgs *messages.Messages,
) error {
	requester := tgutil.ToUser(c.Sender()) // Assumes this is called by an admin
	if requester == nil {
		return c.Respond(&telebot.CallbackResponse{Text: msgs.Common.ErrorIdentifyRequester, ShowAlert: true})
	}

	// Parse payload
	parts := strings.Split(data, "|")
	if len(parts) != 2 {
		log.Printf("ChangeModConfirm: Invalid payload format: %s", data)
		return c.Respond(&telebot.CallbackResponse{Text: fmt.Sprintf(msgs.Common.CallbackErrorGeneric, "invalid data"), ShowAlert: true})
	}
	roomIDStr := parts[0]
	newModeratorIDStr := parts[1]

	roomID := roomEntity.RoomID(roomIDStr)
	newModeratorID, err := strconv.ParseInt(newModeratorIDStr, 10, 64)
	if err != nil {
		log.Printf("ChangeModConfirm: Invalid user ID in payload: %s, error: %v", newModeratorIDStr, err)
		return c.Respond(&telebot.CallbackResponse{Text: fmt.Sprintf(msgs.Common.CallbackErrorGeneric, "invalid user ID"), ShowAlert: true})
	}

	// Fetch the target user's details (needed for the use case command)
	// Ideally, there'd be a GetUserByID query, but we can get it from players list for now
	players, err := getPlayersHandler.Handle(context.Background(), roomQuery.GetPlayersInRoomQuery{RoomID: roomID})
	if err != nil {
		log.Printf("ChangeModConfirm: Error fetching players for room '%s': %v", roomID, err)
		return c.Respond(&telebot.CallbackResponse{Text: fmt.Sprintf(msgs.Common.CallbackErrorGeneric, err), ShowAlert: true})
	}
	var newModeratorUser *sharedEntity.User
	for _, p := range players {
		if p.ID == sharedEntity.UserID(newModeratorID) {
			newModeratorUser = p
			break
		}
	}
	if newModeratorUser == nil {
		// This might happen if the user left between selection and confirmation
		log.Printf("ChangeModConfirm: Could not find user %d in room %s players list.", newModeratorID, roomID)
		return c.Respond(&telebot.CallbackResponse{Text: fmt.Sprintf(msgs.Common.CallbackErrorGeneric, "selected user not found"), ShowAlert: true})
	}

	// Call the use case
	changeCmd := roomCommand.ChangeModeratorCommand{
		Requester:    requester,
		RoomID:       roomID,
		NewModerator: newModeratorUser,
	}
	if err := changeModeratorHandler.Handle(context.Background(), changeCmd); err != nil {
		log.Printf("ChangeModConfirm: Error changing moderator for room %s: %v", roomID, err)
		return c.Respond(&telebot.CallbackResponse{Text: fmt.Sprintf(msgs.Room.ChangeModeratorCallbackError, err), ShowAlert: true})
	}

	// Success - trigger refreshes and edit back to room detail
	roomList.RaiseRefreshNeeded()
	roomDetail.RaiseRefreshNeeded()

	// Acknowledge the callback first
	ackMsg := fmt.Sprintf(msgs.Room.ChangeModeratorCallbackSuccess, newModeratorUser.GetProfileLink(), roomIDStr)
	_ = c.Respond(&telebot.CallbackResponse{Text: ackMsg})

	// Prepare and edit the message back to the standard room detail
	message, markup, err := RoomDetailMessage(getRoomsHandler, getPlayersHandler, msgs, requester.Admin, roomIDStr)
	if err != nil {
		log.Printf("ChangeModConfirm: Error preparing room detail after change for room '%s': %v", roomID, err)
		return nil // Can't easily recover message
	}

	return c.Edit(message, markup, telebot.ModeMarkdownV2, telebot.NoPreview)
}
