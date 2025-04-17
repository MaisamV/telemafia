package memory

import (
	"errors"
	"fmt"
	"log"
	"sync"

	// gameEntity "telemafia/internal/game/entity"
	gameEntity "telemafia/internal/domain/game/entity" // New path
	// gamePort "telemafia/internal/game/port"
	gamePort "telemafia/internal/domain/game/port" // New path
	// roomEntity "telemafia/internal/room/entity" // Needed for RoomID
	roomEntity "telemafia/internal/domain/room/entity" // New path
)

// Ensure InMemoryGameRepository implements the gamePort.GameRepository interface.
var _ gamePort.GameRepository = (*InMemoryGameRepository)(nil)

// InMemoryGameRepository provides an in-memory implementation of the game repository
type InMemoryGameRepository struct {
	games      map[gameEntity.GameID]*gameEntity.Game  // Use imported types
	roomToGame map[roomEntity.RoomID]gameEntity.GameID // Use imported types
	mutex      sync.RWMutex
}

// NewInMemoryGameRepository creates a new in-memory game repository
func NewInMemoryGameRepository() gamePort.GameRepository { // Return the port interface
	return &InMemoryGameRepository{
		games:      make(map[gameEntity.GameID]*gameEntity.Game),
		roomToGame: make(map[roomEntity.RoomID]gameEntity.GameID),
	}
}

// GetGameByID gets a game by its ID
func (r *InMemoryGameRepository) GetGameByID(id gameEntity.GameID) (*gameEntity.Game, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	game, exists := r.games[id]
	if !exists {
		return nil, fmt.Errorf("game with ID %s not found", id)
	}
	return game, nil
}

// GetGameByRoomID gets a game by room ID
func (r *InMemoryGameRepository) GetGameByRoomID(roomID roomEntity.RoomID) (*gameEntity.Game, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	log.Printf("Looking for game for room '%s'", roomID)

	gameID, exists := r.roomToGame[roomID]
	if !exists {
		return nil, fmt.Errorf("no game found for room '%s'", roomID)
	}

	log.Printf("Found game ID '%s' for room '%s'", gameID, roomID)

	game, gameExists := r.games[gameID]
	if !gameExists {
		// This indicates an inconsistency in the in-memory state
		log.Printf("ERROR: Found game ID '%s' in room->game map but game doesn't exist in games map!", gameID)
		return nil, fmt.Errorf("internal inconsistency: game '%s' mapped but not found", gameID)
	}

	log.Printf("Successfully retrieved game '%s' for room '%s'", game.ID, roomID)
	return game, nil
}

// GetAllGames returns all games
func (r *InMemoryGameRepository) GetAllGames() ([]*gameEntity.Game, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	games := make([]*gameEntity.Game, 0, len(r.games))
	for _, game := range r.games {
		games = append(games, game)
	}

	log.Printf("Retrieved %d games from repository", len(games))
	return games, nil
}

// CreateGame creates a new game
func (r *InMemoryGameRepository) CreateGame(game *gameEntity.Game) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if game == nil || game.Room == nil {
		return errors.New("cannot create game: game or game room is nil")
	}

	log.Printf("Attempting to create game '%s' for room '%s'", game.ID, game.Room.ID)

	if _, exists := r.games[game.ID]; exists {
		return fmt.Errorf("game with ID %s already exists", game.ID)
	}
	if _, exists := r.roomToGame[game.Room.ID]; exists {
		// Maybe allow multiple games per room later? For now, enforce 1-1
		return fmt.Errorf("room %s already has an associated game", game.Room.ID)
	}

	r.games[game.ID] = game
	r.roomToGame[game.Room.ID] = game.ID

	log.Printf("Successfully created game '%s' linked to room '%s'", game.ID, game.Room.ID)
	return nil
}

// UpdateGame updates an existing game
func (r *InMemoryGameRepository) UpdateGame(game *gameEntity.Game) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if game == nil {
		return errors.New("cannot update nil game")
	}

	log.Printf("Updating game '%s'", game.ID)

	if _, exists := r.games[game.ID]; !exists {
		return fmt.Errorf("cannot update game: game with ID %s not found", game.ID)
	}

	// Simply replace the existing game object in the map
	r.games[game.ID] = game

	// Note: This assumes the RoomID associated with the GameID doesn't change.
	// If it could, the roomToGame map would also need updating.

	log.Printf("Successfully updated game '%s'", game.ID)
	return nil
}

// DeleteGame deletes a game by ID
func (r *InMemoryGameRepository) DeleteGame(id gameEntity.GameID) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	game, exists := r.games[id]
	if !exists {
		return fmt.Errorf("cannot delete game: game with ID %s not found", id)
	}

	if game.Room != nil {
		delete(r.roomToGame, game.Room.ID)
		log.Printf("Removed room->game mapping for room '%s'", game.Room.ID)
	} else {
		log.Printf("Warning: Deleting game '%s' which has a nil room reference", id)
	}
	delete(r.games, id)
	log.Printf("Successfully deleted game '%s'", id)
	return nil
}
