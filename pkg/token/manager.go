package token

import "sync/atomic"

// Manager handles token distribution in a round-robin fashion
type Manager struct {
	tokens []string
	index  atomic.Uint32
}

func NewManager(tokens []string) *Manager {
	return &Manager{
		tokens: tokens,
	}
}

// GetNextToken returns the next token in round-robin fashion
func (m *Manager) GetNextToken() string {
	if len(m.tokens) == 0 {
		return ""
	}

	index := m.index.Add(1) - 1
	return m.tokens[index%uint32(len(m.tokens))]
}

func (m *Manager) GetAllTokens() []string {
	return m.tokens
}
