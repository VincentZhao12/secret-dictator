package game

import (
	"sync"
)

type Manager struct {
	Games map[string]*Game
	mu   sync.RWMutex
}

func NewManager() *Manager {
	return &Manager{
		Games: make(map[string]*Game),
	}
}

func (m *Manager) GetGame(id string) (*Game, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	game, exists := m.Games[id]
	return game, exists
}

func (m *Manager) AddGame(game *Game) string {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Games[game.ID] = game
	return game.ID
}

func (m *Manager) RemoveGame(id string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.Games, id)
}