package telegram

import (
	"context"
	"strconv"
	"telemafia/internal/room/entity"
	roomQuery "telemafia/internal/room/usecase/query"

	"gopkg.in/telebot.v3"
)

// HandleKickUserCallback handles the kick user callback
func (h *BotHandler) HandleKickUserCallback(c telebot.Context, data string) error {
	roomID := data
	// Fetch players in the room and send them as inline keyboard buttons
	players, err := h.getPlayersInRoomHandler.Handle(context.Background(), roomQuery.GetPlayersInRoomQuery{RoomID: entity.RoomID(roomID)})
	if err != nil {
		return c.Respond(&telebot.CallbackResponse{
			Text: "Failed to fetch players.",
		})
	}

	var buttons [][]telebot.InlineButton
	for _, player := range players {
		button := telebot.InlineButton{
			Unique: UniqueKickPlayerFromRoom,
			Text:   player.FirstName + " " + player.LastName + " (" + player.Username + ")",
			Data:   roomID + ":" + strconv.FormatInt(int64(player.ID), 10),
		}
		buttons = append(buttons, []telebot.InlineButton{button})
	}

	markup := &telebot.ReplyMarkup{InlineKeyboard: buttons}
	return c.Send("Select a player to kick:", markup)
}
