package telegram

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	roomQuery "telemafia/internal/room/usecase/query"

	"gopkg.in/telebot.v3"
)

type MinifiedChat struct {
	chatID int64
}

type UserMessage struct {
	userID    int64
	messageID int
}

func (r MinifiedChat) Recipient() string {
	return strconv.FormatInt(r.chatID, 10)
}

// Map to store message IDs for each user
var (
	messageIDs      = make(map[int64]int)
	messageIDsMutex = &sync.RWMutex{}
)

// HandleListRooms handles the /list_rooms command
func (h *BotHandler) HandleListRooms(c telebot.Context) error {
	// Delete the previous list rooms response if it exists
	messageIDsMutex.RLock()
	prevMsgID, exists := messageIDs[c.Sender().ID]
	delete(messageIDs, c.Sender().ID)
	messageIDsMutex.RUnlock()
	if exists {
		h.bot.Delete(&telebot.Message{ID: prevMsgID, Chat: &telebot.Chat{ID: c.Sender().ID}})
	}
	text, markup, err := h.ListRoomsMessage()
	if err != nil {
		c.Respond(&telebot.CallbackResponse{Text: fmt.Sprintf("Error listing rooms: %v", err)})
		return err
	}
	msg, err := h.bot.Send(&telebot.Chat{ID: c.Sender().ID}, text, markup)
	if err == nil {
		messageIDsMutex.Lock()
		messageIDs[c.Sender().ID] = msg.ID
		messageIDsMutex.Unlock()
	}
	return err
}

func (h *BotHandler) ListRoomsMessage() (string, *telebot.ReplyMarkup, error) {
	rooms, err := h.getRoomsHandler.Handle(context.Background(), roomQuery.GetRoomsQuery{})
	if err != nil {
		return "", nil, err
	}

	var text string
	var markup *telebot.ReplyMarkup = nil
	if len(rooms) == 0 {
		text = "فعلا بازی در حال شروع شدن نیست."
	} else {
		// Create inline keyboard
		var buttons [][]telebot.InlineButton
		for _, room := range rooms {
			buttonText := fmt.Sprintf("%s (بازیکنان: %d)", room.Name, len(room.Players))
			buttons = append(buttons, []telebot.InlineButton{
				{
					Unique: UniqueJoinSelectRoom,
					Text:   buttonText,
					Data:   string(room.ID),
				},
			})
		}
		text = "Available rooms:"
		markup = &telebot.ReplyMarkup{
			InlineKeyboard: buttons,
		}
	}

	return text, markup, err
}
