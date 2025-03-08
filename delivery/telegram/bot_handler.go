package telegram

import (
	roomCommand "telemafia/internal/room/usecase/command"
	roomQuery "telemafia/internal/room/usecase/query"

	"gopkg.in/telebot.v3"
)

// BotHandler handles Telegram bot commands
type BotHandler struct {
	bot                     *telebot.Bot
	adminUsernames          []string
	createRoomHandler       *roomCommand.CreateRoomHandler
	joinRoomHandler         *roomCommand.JoinRoomHandler
	leaveRoomHandler        *roomCommand.LeaveRoomHandler
	kickUserHandler         *roomCommand.KickUserHandler
	getRoomsHandler         *roomQuery.GetRoomsHandler
	getPlayerRoomsHandler   *roomQuery.GetPlayerRoomsHandler
	getPlayersInRoomHandler *roomQuery.GetPlayersInRoomHandler
}

// NewBotHandler creates a new BotHandler
func NewBotHandler(
	bot *telebot.Bot,
	adminUsernames []string,
	createRoomHandler *roomCommand.CreateRoomHandler,
	joinRoomHandler *roomCommand.JoinRoomHandler,
	leaveRoomHandler *roomCommand.LeaveRoomHandler,
	kickUserHandler *roomCommand.KickUserHandler,
	getRoomsHandler *roomQuery.GetRoomsHandler,
	getPlayerRoomsHandler *roomQuery.GetPlayerRoomsHandler,
	getPlayersInRoomHandler *roomQuery.GetPlayersInRoomHandler,
) *BotHandler {
	return &BotHandler{
		bot:                     bot,
		adminUsernames:          adminUsernames,
		createRoomHandler:       createRoomHandler,
		joinRoomHandler:         joinRoomHandler,
		leaveRoomHandler:        leaveRoomHandler,
		kickUserHandler:         kickUserHandler,
		getRoomsHandler:         getRoomsHandler,
		getPlayerRoomsHandler:   getPlayerRoomsHandler,
		getPlayersInRoomHandler: getPlayersInRoomHandler,
	}
}

func (h *BotHandler) Start() {
	h.bot.Start()
}

// RegisterHandlers registers all bot command handlers
func (h *BotHandler) RegisterHandlers() {
	h.bot.SetCommands([]telebot.Command{
		{Text: "empty", Description: ""},
	})
	h.bot.SetCommands([]telebot.Command{
		{Text: "list_rooms", Description: "Show mafia rooms"},
	})
	h.bot.Handle("/start", h.HandleStart)
	h.bot.Handle("/help", h.HandleHelp)
	h.bot.Handle("/create_room", h.HandleCreateRoom)
	h.bot.Handle("/join_room", h.HandleJoinRoom)
	h.bot.Handle("/leave_room", h.HandleLeaveRoom)
	h.bot.Handle("/list_rooms", h.HandleListRooms)
	h.bot.Handle("/my_rooms", h.HandleMyRooms)
	h.bot.Handle("/kick_user", h.HandleKickUser)
	h.bot.Handle(telebot.OnCallback, h.HandleCallback)
}
