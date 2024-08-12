package config

import (
	"fmt"

	"github.com/robfig/cron/v3"
)

func ValidateConfig(cfg Config) error {
	if cfg.Username == "" {
		return fmt.Errorf("username cannot be empty")
	}

	if cfg.Token == "" {
		return fmt.Errorf("token cannot be empty. See here: https://github.com/AkashRajpurohit/git-sync/wiki/Configuration")
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

	if cfg.Platform == "bitbucket" && cfg.Workspace == "" {
		return fmt.Errorf("workspace cannot be empty for bitbucket")
	}

	if cfg.Cron != "" {
		_, err := cron.ParseStandard(cfg.Cron)
		if err != nil {
			return fmt.Errorf("invalid cron expression %s", cfg.Cron)
		}
	}

	return nil
}
