package token

import (
	"sync"
	"testing"
)

func TestNewManager(t *testing.T) {
	tokens := []string{"token1", "token2", "token3"}
	manager := NewManager(tokens)

	if len(manager.tokens) != len(tokens) {
		t.Errorf("Expected %d tokens, got %d", len(tokens), len(manager.tokens))
	}

	for i, token := range tokens {
		if manager.tokens[i] != token {
			t.Errorf("Expected token %s at position %d, got %s", token, i, manager.tokens[i])
		}
	}
}

func TestGetNextToken(t *testing.T) {
	tests := []struct {
		name   string
		tokens []string
		calls  int
		want   string
	}{
		{
			name:   "Single token",
			tokens: []string{"token1"},
			calls:  1,
			want:   "token1",
		},
		{
			name:   "Multiple tokens - first round",
			tokens: []string{"token1", "token2", "token3"},
			calls:  2,
			want:   "token2",
		},
		{
			name:   "Multiple tokens - wrap around",
			tokens: []string{"token1", "token2"},
			calls:  3,
			want:   "token1",
		},
		{
			name:   "Empty tokens",
			tokens: []string{},
			calls:  1,
			want:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			manager := NewManager(tt.tokens)
			var got string
			for i := 0; i < tt.calls; i++ {
				got = manager.GetNextToken()
			}
			if got != tt.want {
				t.Errorf("GetNextToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetNextTokenConcurrent(t *testing.T) {
	tokens := []string{"token1", "token2", "token3"}
	manager := NewManager(tokens)
	iterations := 100
	goroutines := 10

	var wg sync.WaitGroup
	tokenCounts := make(map[string]int)
	var mu sync.Mutex

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				token := manager.GetNextToken()
				mu.Lock()
				tokenCounts[token]++
				mu.Unlock()
			}
		}()
	}

	wg.Wait()

	// Verify that tokens were distributed roughly evenly
	totalCalls := goroutines * iterations
	expectedPerToken := totalCalls / len(tokens)
	tolerance := float64(expectedPerToken) * 0.2 // 20% tolerance

	for token, count := range tokenCounts {
		diff := float64(abs(count - expectedPerToken))
		if diff > tolerance {
			t.Errorf("Token %s was used %d times, expected approximately %d (±%.0f)",
				token, count, expectedPerToken, tolerance)
		}
	}
}

func TestGetAllTokens(t *testing.T) {
	tokens := []string{"token1", "token2", "token3"}
	manager := NewManager(tokens)

	got := manager.GetAllTokens()
	if len(got) != len(tokens) {
		t.Errorf("GetAllTokens() returned %d tokens, want %d", len(got), len(tokens))
	}

	for i, token := range tokens {
		if got[i] != token {
			t.Errorf("GetAllTokens()[%d] = %v, want %v", i, got[i], token)
		}
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
