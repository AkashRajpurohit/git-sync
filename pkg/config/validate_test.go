package config

import (
	"testing"
)

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name    string
		cfg     Config
		wantErr bool
	}{
		{
			name: "Empty Username",
			cfg: Config{
				Username: "",
			},
			wantErr: true,
		},
		{
			name: "Empty Token",
			cfg: Config{
				Username: "test",
				Token:    "",
			},
			wantErr: true,
		},
		{
			name: "Empty BackupDir",
			cfg: Config{
				Username:  "test",
				Token:     "test",
				BackupDir: "",
			},
			wantErr: true,
		},
		{
			name: "Empty Server Domain",
			cfg: Config{
				Username:  "test",
				Token:     "test",
				BackupDir: "test",
				Server: Server{
					Domain: "",
				},
			},
			wantErr: true,
		},
		{
			name: "Invalid Server Protocol",
			cfg: Config{
				Username:  "test",
				Token:     "test",
				BackupDir: "test",
				Server: Server{
					Domain:   "test",
					Protocol: "ftp",
				},
			},
			wantErr: true,
		},
		{
			name: "Empty Workspace for Bitbucket",
			cfg: Config{
				Username:  "test",
				Token:     "test",
				BackupDir: "test",
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
			name: "Invalid Cron Expression",
			cfg: Config{
				Username:  "test",
				Token:     "test",
				BackupDir: "test",
				Server: Server{
					Domain:   "test",
					Protocol: "https",
				},
				Cron: "invalid",
			},
			wantErr: true,
		},
		{
			name: "Invalid Clone Type",
			cfg: Config{
				Username:  "test",
				Token:     "test",
				BackupDir: "test",
				Server: Server{
					Domain:   "test",
					Protocol: "https",
				},
				CloneType: "invalid",
			},
			wantErr: true,
		},
		{
			name: "Valid Config",
			cfg: Config{
				Username:  "test",
				Token:     "test",
				BackupDir: "test",
				Server: Server{
					Domain:   "test",
					Protocol: "https",
				},
				CloneType: "bare",
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
