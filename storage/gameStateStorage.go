package storage

import (
	"fmt"
)

type GameState struct {
	Deck       []int
	DealerHand []int
	PlayerHand []int
	Wins       int
	Losses     int
	Ties       int
	Result     int
}

type GameStateStorage interface {
	GetGameState(sessionID string) (*GameState, error)
	UpsertGameState(sessionID string, state *GameState)
	DeleteGameState(sessionID string)
}

func NewGameStateStorage(backend string) (GameStateStorage, error) {
	switch backend {
	case inMemBackend:
		return NewInMemStorage(), nil
	default:
		return nil, fmt.Errorf("unknown game state storage backend '%s'", backend)
	}
}
