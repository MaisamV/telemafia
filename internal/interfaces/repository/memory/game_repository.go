package memory

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"telemafia/internal/entity"
	"telemafia/internal/usecase"
)

// Ensure InMemoryGameRepository implements the usecase.GameRepository interface.
var _ usecase.GameRepository = (*InMemoryGameRepository)(nil)

// InMemoryGameRepository provides an in-memory implementation of the game repository
type InMemoryGameRepository struct {
	games      map[entity.GameID]*entity.Game
	roomToGame map[entity.RoomID]entity.GameID
	mutex      sync.RWMutex
}

// NewInMemoryGameRepository creates a new in-memory game repository
func NewInMemoryGameRepository() usecase.GameRepository {
	return &InMemoryGameRepository{
		games:      make(map[entity.GameID]*entity.Game),
		roomToGame: make(map[entity.RoomID]entity.GameID),
	}
}

// GetGameByID gets a game by its ID
func (r *InMemoryGameRepository) GetGameByID(id entity.GameID) (*entity.Game, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	game, exists := r.games[id]
	if !exists {
		return nil, errors.New("game not found")
	}
	return game, nil
}

// GetGameByRoomID gets a game by room ID
func (r *InMemoryGameRepository) GetGameByRoomID(roomID entity.RoomID) (*entity.Game, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	log.Printf("Looking for game for room '%s'", roomID)

	// Debug: Log all room-to-game mappings
	log.Printf("Available room-to-game mappings: %d", len(r.roomToGame))
	for rID, gID := range r.roomToGame {
		log.Printf("Room '%s' is mapped to game '%s'", rID, gID)
	}

	gameID, exists := r.roomToGame[roomID]
	if !exists {
		return nil, fmt.Errorf("no game found for room '%s'", roomID)
	}

	log.Printf("Found game ID '%s' for room '%s'", gameID, roomID)

	game, exists := r.games[gameID]
	if !exists {
		return nil, fmt.Errorf("found game ID '%s' in mapping but game not found in storage", gameID)
	}

	log.Printf("Successfully retrieved game for room '%s'", roomID)
	return game, nil
}

// GetAllGames returns all games
func (r *InMemoryGameRepository) GetAllGames() ([]*entity.Game, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	games := make([]*entity.Game, 0, len(r.games))
	for _, game := range r.games {
		games = append(games, game)
	}

	log.Printf("Retrieved all games from repository. Found %d games", len(games))
	return games, nil
}

// CreateGame creates a new game
func (r *InMemoryGameRepository) CreateGame(game *entity.Game) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	log.Printf("Creating game '%s' for room '%s'", game.ID, game.Room.ID)

	r.games[game.ID] = game
	r.roomToGame[game.Room.ID] = game.ID

	log.Printf("Successfully created game '%s' for room '%s'", game.ID, game.Room.ID)
	return nil
}

// UpdateGame updates an existing game
func (r *InMemoryGameRepository) UpdateGame(game *entity.Game) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	log.Printf("Updating game '%s' for room '%s'", game.ID, game.Room.ID)

	if _, exists := r.games[game.ID]; !exists {
		return errors.New("game not found")
	}

	r.games[game.ID] = game

	log.Printf("Successfully updated game '%s'", game.ID)
	return nil
}

// DeleteGame deletes a game by ID
func (r *InMemoryGameRepository) DeleteGame(id entity.GameID) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	game, exists := r.games[id]
	if !exists {
		return errors.New("game not found")
	}

	delete(r.roomToGame, game.Room.ID)
	delete(r.games, id)
	return nil
}
