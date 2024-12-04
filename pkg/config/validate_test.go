package config

import (
	"testing"
)

func TestValidateGitURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{
			name:    "Valid HTTPS URL",
			url:     "https://github.com/user/repo.git",
			wantErr: false,
		},
		{
			name:    "Valid SSH URL",
			url:     "git@github.com:user/repo.git",
			wantErr: false,
		},
		{
			name:    "Invalid HTTPS URL - No .git suffix",
			url:     "https://github.com/user/repo",
			wantErr: true,
		},
		{
			name:    "Invalid SSH URL - Wrong format",
			url:     "git@github.com/user/repo.git",
			wantErr: true,
		},
		{
			name:    "Invalid Protocol",
			url:     "ftp://github.com/user/repo.git",
			wantErr: true,
		},
		{
			name:    "Invalid URL format",
			url:     "not-a-url",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateGitURL(tt.url); (err != nil) != tt.wantErr {
				t.Errorf("validateGitURL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name    string
		cfg     Config
		wantErr bool
	}{
		{
			name: "Empty BackupDir",
			cfg: Config{
				BackupDir: "",
			},
			wantErr: true,
		},
		{
			name: "Invalid Clone Type",
			cfg: Config{
				BackupDir: "test",
				CloneType: "invalid",
			},
			wantErr: true,
		},
		{
			name: "Invalid Cron Expression",
			cfg: Config{
				BackupDir: "test",
				CloneType: "bare",
				Cron:      "invalid",
			},
			wantErr: true,
		},
		{
			name: "Valid Raw Git URLs Only",
			cfg: Config{
				BackupDir: "test",
				CloneType: "bare",
				RawGitURLs: []string{
					"https://github.com/user/repo1.git",
					"git@github.com:user/repo2.git",
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid Raw Git URL",
			cfg: Config{
				BackupDir: "test",
				CloneType: "bare",
				RawGitURLs: []string{
					"https://github.com/user/repo1", // Missing .git
					"git@github.com:user/repo2.git",
				},
			},
			wantErr: true,
		},
		{
			name: "Empty Username with No Raw URLs",
			cfg: Config{
				BackupDir: "test",
				CloneType: "bare",
				Token:     "test",
				Server: Server{
					Domain:   "test",
					Protocol: "https",
				},
			},
			wantErr: true,
		},
		{
			name: "Empty Token with No Raw URLs",
			cfg: Config{
				Username:  "test",
				BackupDir: "test",
				CloneType: "bare",
				Server: Server{
					Domain:   "test",
					Protocol: "https",
				},
			},
			wantErr: true,
		},
		{
			name: "Empty Server Domain with No Raw URLs",
			cfg: Config{
				Username:  "test",
				Token:     "test",
				BackupDir: "test",
				CloneType: "bare",
				Server: Server{
					Protocol: "https",
				},
			},
			wantErr: true,
		},
		{
			name: "Invalid Server Protocol with No Raw URLs",
			cfg: Config{
				Username:  "test",
				Token:     "test",
				BackupDir: "test",
				CloneType: "bare",
				Server: Server{
					Domain:   "test",
					Protocol: "ftp",
				},
			},
			wantErr: true,
		},
		{
			name: "Empty Workspace for Bitbucket with No Raw URLs",
			cfg: Config{
				Username:  "test",
				Token:     "test",
				BackupDir: "test",
				CloneType: "bare",
				Server: Server{
					Domain:   "test",
					Protocol: "https",
				},
				Platform:  "bitbucket",
				Workspace: "",
			},
			wantErr: true,
		},
		{
			name: "Valid Platform Config with No Raw URLs",
			cfg: Config{
				Username:  "test",
				Token:     "test",
				BackupDir: "test",
				CloneType: "bare",
				Server: Server{
					Domain:   "test",
					Protocol: "https",
				},
			},
			wantErr: false,
		},
		{
			name: "Valid Mixed Config (Platform + Raw URLs)",
			cfg: Config{
				Username:  "test",
				Token:     "test",
				BackupDir: "test",
				CloneType: "bare",
				Server: Server{
					Domain:   "test",
					Protocol: "https",
				},
				RawGitURLs: []string{
					"https://github.com/user/repo1.git",
					"git@github.com:user/repo2.git",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateConfig(tt.cfg); (err != nil) != tt.wantErr {
				t.Errorf("ValidateConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
