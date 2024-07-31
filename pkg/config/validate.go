package config

import "fmt"

func ValidateConfig(cfg Config) error {
	if cfg.Username == "" {
		return fmt.Errorf("username cannot be empty")
	}

	if cfg.Token == "" {
		return fmt.Errorf("token cannot be empty. See here: https://github.com/AkashRajpurohit/git-sync#how-do-i-create-a-github-personal-access-token")
	}

	if cfg.BackupDir == "" {
		return fmt.Errorf("backup directory cannot be empty")
	}

	if cfg.Server.Domain == "" {
		return fmt.Errorf("server domain cannot be empty")
	}

	if cfg.Server.Protocol != "https" && cfg.Server.Protocol != "http" {
		return fmt.Errorf("server protocol can only be http or https")
	}

	return nil
}
