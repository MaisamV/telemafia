package telegram

import (
	"context"
	"fmt"
	"log"

	roomEntity "telemafia/internal/domain/room/entity"
	roomCommand "telemafia/internal/domain/room/usecase/command"
	roomQuery "telemafia/internal/domain/room/usecase/query"
	messages "telemafia/internal/presentation/telegram/messages"
	tgutil "telemafia/internal/shared/tgutil"

	"gopkg.in/telebot.v3"
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
	return c.Edit(message, markup, telebot.ModeMarkdown, telebot.NoPreview)
}
