package config

import (
	"testing"

	"github.com/AkashRajpurohit/git-sync/pkg/logger"
)

func TestSetSensibleDefaults(t *testing.T) {
	tests := []struct {
		name     string
		cfg      Config
		expected Config
	}{
		{
			name: "GitHub platform defaults",
			cfg: Config{
				Platform: "github",
			},
			expected: Config{
				Platform: "github",
				Server: Server{
					Domain:   "github.com",
					Protocol: "https",
				},
				CloneType:   "bare",
				Concurrency: 5,
			},
		},
		{
			name: "GitLab platform defaults",
			cfg: Config{
				Platform: "gitlab",
			},
			expected: Config{
				Platform: "gitlab",
				Server: Server{
					Domain:   "gitlab.com",
					Protocol: "https",
				},
				CloneType:   "bare",
				Concurrency: 5,
			},
		},
		{
			name: "Bitbucket platform defaults",
			cfg: Config{
				Platform: "bitbucket",
			},
			expected: Config{
				Platform: "bitbucket",
				Server: Server{
					Domain:   "bitbucket.org",
					Protocol: "https",
				},
				CloneType:   "bare",
				Concurrency: 5,
			},
		},
		{
			name: "Forgejo platform defaults",
			cfg: Config{
				Platform: "forgejo",
			},
			expected: Config{
				Platform: "forgejo",
				Server: Server{
					Domain:   "v9.next.forgejo.org",
					Protocol: "https",
				},
				CloneType:   "bare",
				Concurrency: 5,
			},
		},
		{
			name: "Default concurrency when not set",
			cfg: Config{
				Platform: "github",
				Server: Server{
					Domain:   "github.com",
					Protocol: "https",
				},
			},
			expected: Config{
				Platform: "github",
				Server: Server{
					Domain:   "github.com",
					Protocol: "https",
				},
				CloneType:   "bare",
				Concurrency: 5,
			},
		},
		{
			name: "Default clone type when not set",
			cfg: Config{
				Platform: "github",
				Server: Server{
					Domain:   "github.com",
					Protocol: "https",
				},
				Concurrency: 5,
			},
			expected: Config{
				Platform: "github",
				Server: Server{
					Domain:   "github.com",
					Protocol: "https",
				},
				CloneType:   "bare",
				Concurrency: 5,
			},
		},
		{
			name: "Merge single token with tokens array",
			cfg: Config{
				Token:  "single-token",
				Tokens: []string{"token1", "token2"},
			},
			expected: Config{
				Tokens:      []string{"single-token", "token1", "token2"},
				CloneType:   "bare",
				Concurrency: 5,
			},
		},
		{
			name: "Convert single token to tokens array",
			cfg: Config{
				Token: "single-token",
			},
			expected: Config{
				Tokens:      []string{"single-token"},
				CloneType:   "bare",
				Concurrency: 5,
			},
		},
		{
			name: "Keep existing values if already set",
			cfg: Config{
				Platform:    "github",
				CloneType:   "mirror",
				Concurrency: 10,
				Server: Server{
					Domain:   "custom.github.com",
					Protocol: "ssh",
				},
			},
			expected: Config{
				Platform:    "github",
				CloneType:   "mirror",
				Concurrency: 10,
				Server: Server{
					Domain:   "custom.github.com",
					Protocol: "ssh",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger.InitLogger("fatal")
			cfg := tt.cfg
			SetSensibleDefaults(&cfg)

			if cfg.Server.Domain != tt.expected.Server.Domain {
				t.Errorf("Server Domain = %v, want %v", cfg.Server.Domain, tt.expected.Server.Domain)
			}
			if cfg.Server.Protocol != tt.expected.Server.Protocol {
				t.Errorf("Server Protocol = %v, want %v", cfg.Server.Protocol, tt.expected.Server.Protocol)
			}
			if cfg.CloneType != tt.expected.CloneType {
				t.Errorf("CloneType = %v, want %v", cfg.CloneType, tt.expected.CloneType)
			}
			if cfg.Concurrency != tt.expected.Concurrency {
				t.Errorf("Concurrency = %v, want %v", cfg.Concurrency, tt.expected.Concurrency)
			}
			if cfg.Token != "" {
				t.Error("Token should be cleared after conversion")
			}
			if len(cfg.Tokens) != len(tt.expected.Tokens) {
				t.Errorf("Tokens length = %v, want %v", len(cfg.Tokens), len(tt.expected.Tokens))
			}
			for i, token := range tt.expected.Tokens {
				if i < len(cfg.Tokens) && cfg.Tokens[i] != token {
					t.Errorf("Token at index %d = %v, want %v", i, cfg.Tokens[i], token)
				}
			}
		})
	}
}
