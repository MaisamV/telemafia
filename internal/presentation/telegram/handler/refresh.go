package telegram

import (
	"context"
	"log"
	"strings" // Import messages
	gameEntity "telemafia/internal/domain/game/entity"
	gameQuery "telemafia/internal/domain/game/usecase/query"
	roomQuery "telemafia/internal/domain/room/usecase/query"
	"telemafia/internal/shared/entity"
	"telemafia/internal/shared/tgutil"
	"time"

	game "telemafia/internal/presentation/telegram/handler/game"
	roomHandler "telemafia/internal/presentation/telegram/handler/room"

	"gopkg.in/telebot.v3"
)

func (h *BotHandler) updateMessages(book *tgutil.RefreshingMessageBook, getMessage func(user int64, data string) (string, []interface{}, error)) {
	messagesToUpdate := book.GetAllActiveMessages()
	if len(messagesToUpdate) == 0 {
		return
	}

	log.Printf("Refreshing %d messages...", len(messagesToUpdate))
	for chatID, payload := range messagesToUpdate {
		// Prepare the updated message content using the refactored function
		// Pass the necessary handlers from the BotHandler (h)
		newContent, newMarkup, err := getMessage(chatID, payload.Data)
		if err != nil {
			log.Printf(h.msgs.Refresh.ErrorPrepare, chatID, err) // Use msg
			continue
		}

		_, editErr := h.bot.Edit(&telebot.Message{ID: payload.MessageID, Chat: &telebot.Chat{ID: payload.ChatID}}, newContent, newMarkup...)
		if editErr != nil {
			if strings.Contains(editErr.Error(), "message to edit not found") ||
				strings.Contains(editErr.Error(), "bot was blocked by the user") {
				log.Printf(h.msgs.Refresh.ErrorEditRemoving, chatID, editErr) // Use msg
				book.RemoveActiveMessage(chatID)
			} else {
				log.Printf(h.msgs.Refresh.ErrorEdit, chatID, editErr) // Use msg
			}
		}
	}
	log.Println("Finished refreshing messages.")
}

// StartRefreshTimer handles updating dynamic messages.
// It periodically checks the RefreshingMessageBook and updates all active messages if changes occurred.
func (h *BotHandler) StartRefreshTimer() {
	ticker := time.NewTicker(2 * time.Second) // Reduced interval for faster feedback during testing
	defer ticker.Stop()
	for range ticker.C {
		// --- Room List Refresh ---
		if h.roomListRefreshMessage.ConsumeRefreshNeeded() {
			h.updateMessages(h.roomListRefreshMessage, func(user int64, data string) (string, []interface{}, error) {
				message, markup, err := roomHandler.PrepareRoomListMessage(
					h.getRoomsHandler,
					h.getPlayersInRoomHandler,
					h.msgs,
				)
				opts := []interface{}{
					markup,
					telebot.NoPreview,
				}
				return message, opts, err
			})
		}

		// --- Room Detail Refresh ---
		if h.roomDetailRefreshMessage.ConsumeRefreshNeeded() {
			h.updateMessages(h.roomDetailRefreshMessage, func(user int64, data string) (string, []interface{}, error) {
				message, markup, err := roomHandler.RoomDetailMessage(
					h.getRoomsHandler,
					h.getPlayersInRoomHandler,
					h.msgs,
					entity.UserID(user),
					data,
				)
				opts := []interface{}{
					markup,
					telebot.ModeMarkdownV2,
					telebot.NoPreview,
				}
				return message, opts, err
			})
		}

		// --- Game Admin Assignment Tracker Refresh ---
		h.adminRefreshMutex.RLock() // Lock for reading the map
		for gameID, book := range h.adminAssignmentTrackers {
			if book.ConsumeRefreshNeeded() {
				log.Printf("Refresh needed for Admin Tracker Game ID: %s", gameID)
				h.updateMessages(book, func(user int64, data string) (string, []interface{}, error) {
					gameIDFromData := gameEntity.GameID(data)
					// Fetch necessary data (Game, State, Players) - May need error handling
					gameData, _ := h.GetGameByIDHandler().Handle(context.Background(), gameQuery.GetGameByIDQuery{ID: gameIDFromData})
					state, stateExists := h.GetInteractiveSelectionState(gameIDFromData)
					players, _ := h.GetPlayersInRoomHandler().Handle(context.Background(), roomQuery.GetPlayersInRoomQuery{RoomID: gameData.Room.ID})

					if gameData == nil || !stateExists {
						log.Printf("Refresh Admin: Game %s or State not found.", gameIDFromData)
						return "Error: Game data unavailable.", []interface{}{}, nil // Return error state message
					}
					// Prepare message content using the helper from callbacks_game
					message, markup, err := game.PrepareAdminAssignmentMessage(gameData, state, players, h.msgs)
					opts := []interface{}{
						markup,
						telebot.NoPreview, // Keep it plain text for now
					}
					return message, opts, err
				})
			}
		}
		h.adminRefreshMutex.RUnlock()

		// --- Game Player Role Choice Refresh ---
		h.playerRefreshMutex.RLock() // Lock for reading the map
		for gameID, book := range h.playerRoleChoiceRefreshers {
			if book.ConsumeRefreshNeeded() {
				log.Printf("Refresh needed for Player Choice Game ID: %s", gameID)
				h.updateMessages(book, func(user int64, data string) (string, []interface{}, error) {
					gameIDFromData := gameEntity.GameID(data)
					// Fetch necessary data (State)
					state, stateExists := h.GetInteractiveSelectionState(gameIDFromData)
					if !stateExists {
						log.Printf("Refresh Player: State for game %s not found.", gameIDFromData)
						return "Role selection is no longer active.", []interface{}{}, nil // Return inactive state message
					}

					// Prepare markup using the helper from callbacks_game
					markup, err := game.PreparePlayerRoleSelectionMarkup(gameIDFromData, len(state.ShuffledRoles), state.TakenIndices, h.msgs)
					message := h.msgs.Game.RoleSelectionPromptPlayer // Keep the prompt same, just update buttons
					opts := []interface{}{
						markup,
					}
					return message, opts, err
				})
			}
		}
		h.playerRefreshMutex.RUnlock()
	}
}
