package storage

import (
	"errors"
	"sync"
)

var _ GameStateStorage = (*inMemStorage)(nil)
var inMemBackend = "inMemory"

type inMemStorage struct {
	syncMap sync.Map
}

func NewInMemStorage() GameStateStorage {
	return &inMemStorage{}
}

func (s *inMemStorage) GetGameState(sessionID string) (*GameState, error) {
	gs, ok := s.syncMap.Load(sessionID)
	if !ok {
		return nil, nil
	}

	state, ok := gs.(GameState)
	if !ok {
		return nil, errors.New("invalid state object")
	}

	return &state, nil
}

func (s *inMemStorage) UpsertGameState(sessionID string, state *GameState) {
	s.syncMap.Store(sessionID, *state)
}

func (s *inMemStorage) DeleteGameState(sessionID string) {
	s.syncMap.Delete(sessionID)
}
