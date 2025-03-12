package telegram

import (
	"context"
	"gopkg.in/telebot.v3"
	"telemafia/internal/room/usecase/query"
	userEntity "telemafia/internal/user/entity"
)

// HandleStart handles the /start command
func (h *BotHandler) HandleStart(c telebot.Context) error {
	roomsQuery := query.GetPlayerRoomsQuery{PlayerID: userEntity.UserID(c.Sender().ID)}
	_, err := h.getPlayerRoomsHandler.Handle(context.Background(), roomsQuery)
	if err != nil {

	}
	//if handle != nil && len(handle) > 0 {
	//	return h.
	//} else {
	return h.HandleListRooms(c)
	//}
}
