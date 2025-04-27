package telegram

import (
	"context"
	"fmt"
	"log"
	"strings"

	gameEntity "telemafia/internal/domain/game/entity"
	gameCommand "telemafia/internal/domain/game/usecase/command"
	messages "telemafia/internal/presentation/telegram/messages"
	tgutil "telemafia/internal/shared/tgutil"

	"gopkg.in/telebot.v3"
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

	for userID, role := range assignments {
		targetUser := &telebot.User{ID: int64(userID)}
		privateMsg := fmt.Sprintf(msgs.Game.AssignRolesSuccessPrivate, gameID, "<Room Name>", role.Name)
		_, err := bot.Send(targetUser, privateMsg)
		if err != nil {
			log.Printf(msgs.Game.AssignRolesErrorSendingPrivate, userID, err)
		}
	}

	return c.Send(fmt.Sprintf(msgs.Game.AssignRolesSuccessPublic, gameID))
}
