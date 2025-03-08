package telegram

import (
	"fmt"
	"telemafia/common"
	userEntity "telemafia/internal/user/entity"

	"gopkg.in/telebot.v3"
)

// HandleStart handles the /start command
func (h *BotHandler) HandleStart(c telebot.Context) error {
	sender := c.Sender()
	user := userEntity.User{
		ID:        userEntity.UserID(sender.ID),
		FirstName: sender.FirstName,
		LastName:  sender.LastName,
		Username:  sender.Username,
		Admin:     common.Contains(h.adminUsernames, c.Sender().Username),
	}

	c.Send(fmt.Sprintf("%s، به کافه‌مافیا خوش اومدی، سناریو رو انتخاب کن", user.FirstName))
	return h.HandleListRooms(c)
}
