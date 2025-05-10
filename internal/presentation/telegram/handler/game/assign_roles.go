package telegram

import (
	"context"
	"fmt"
	"log"
	"strings"
	"telemafia/internal/domain/scenario/entity"
	"telemafia/internal/shared/common"

	gameEntity "telemafia/internal/domain/game/entity"
	gameCommand "telemafia/internal/domain/game/usecase/command"
	messages "telemafia/internal/presentation/telegram/messages"
	tgutil "telemafia/internal/shared/tgutil"

	"gopkg.in/telebot.v4"
)

// HandleAssignRoles handles the /assign_roles command (now a function)
func HandleAssignRoles(
	assignRolesHandler *gameCommand.AssignRolesHandler,
	bot *telebot.Bot,
	c telebot.Context,
	msgs *messages.Messages,
) error {
	gameIDStr := strings.TrimSpace(c.Message().Payload)
	if gameIDStr == "" {
		return c.Send(msgs.Game.AssignRolesPrompt)
	}
	gameID := gameEntity.GameID(gameIDStr)

	requester := tgutil.ToUser(c.Sender())
	if requester == nil || !requester.Admin {
		return c.Send(msgs.Common.ErrorPermissionDenied)
	}

	cmd := gameCommand.AssignRolesCommand{
		Requester: *requester,
		GameID:    gameID,
	}

	assignments, err := assignRolesHandler.Handle(context.Background(), cmd)
	if err != nil {
		return c.Send(fmt.Sprintf(msgs.Game.AssignRolesError, gameID, err))
	}

	for user, role := range assignments {
		targetUser := &telebot.User{ID: int64(user.ID)}
		privateMsg, opts := PrepareAssignRoleMessage(msgs, role)
		_, err := bot.Send(targetUser, privateMsg, opts...)
		if err != nil {
			log.Printf(msgs.Game.AssignRolesErrorSendingPrivate, user.ID, err)
		}
	}

	return c.Send(fmt.Sprintf(msgs.Game.AssignRolesSuccessPublic, gameID))
}

func PrepareAssignRoleMessage(msgs *messages.Messages, role entity.Role) (interface{}, []interface{}) {
	confirmMsgText := fmt.Sprintf(msgs.Game.AssignRolesSuccessPrivate, role.Name, role.Side)
	if role.Description != "" {
		confirmMsgText = fmt.Sprintf("%s\nتوضیحات: ||%s||", confirmMsgText, common.EscapeMarkdownV2(role.Description))
	}
	var what interface{}
	if role.ImageID != "" {
		what = &telebot.Photo{
			File:       telebot.File{FileID: role.ImageID},
			HasSpoiler: true,
			Caption:    confirmMsgText,
		}
	} else {
		what = confirmMsgText
	}
	opts := []interface{}{
		telebot.ModeMarkdownV2,
	}
	return what, opts
}
