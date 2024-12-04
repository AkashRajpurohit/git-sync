package config

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/robfig/cron/v3"
)

// validateGitURL validates if the provided URL is a valid git repository URL
func validateGitURL(rawURL string) error {
	// Handle SSH URLs (git@github.com:user/repo.git)
	if strings.HasPrefix(rawURL, "git@") {
		parts := strings.Split(rawURL, ":")
		if len(parts) != 2 || !strings.HasSuffix(parts[1], ".git") {
			return fmt.Errorf("invalid SSH git URL format: %s", rawURL)
		}
		return nil
	}

	// Handle HTTPS URLs
	u, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("invalid git URL: %s", rawURL)
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("git URL must use http or https protocol: %s", rawURL)
	}

	if !strings.HasSuffix(u.Path, ".git") {
		return fmt.Errorf("git URL must end with .git: %s", rawURL)
	}

	return nil
}

func ValidateConfig(cfg Config) error {
	// Validate backup directory (required for all cases)
	if cfg.BackupDir == "" {
		return fmt.Errorf("backup directory cannot be empty")
	}

	// Validate clone type (required for all cases)
	if cfg.CloneType != "bare" && cfg.CloneType != "full" && cfg.CloneType != "mirror" && cfg.CloneType != "shallow" {
		return fmt.Errorf("clone_type can only be `bare`, `full`, `mirror` or `shallow`")
	}

	// Validate concurrency
	if cfg.Concurrency < 1 || cfg.Concurrency > 20 {
		return fmt.Errorf("concurrency must be between 1 and 20")
	}

	// Validate cron if provided
	if cfg.Cron != "" {
		_, err := cron.ParseStandard(cfg.Cron)
		if err != nil {
			return fmt.Errorf("invalid cron expression %s", cfg.Cron)
		}
	}

	// Validate raw git URLs if provided
	for _, url := range cfg.RawGitURLs {
		if err := validateGitURL(url); err != nil {
			return err
		}
	}

	// If there are no raw git URLs, validate platform-specific configuration
	if len(cfg.RawGitURLs) == 0 {
		// Username is required for platform-specific sync
		if cfg.Username == "" {
			return fmt.Errorf("username cannot be empty when no raw git URLs are provided")
		}

		// Token is required for platform-specific sync
		if cfg.Token == "" {
			return fmt.Errorf("token cannot be empty when no raw git URLs are provided. See here: https://github.com/AkashRajpurohit/git-sync/wiki/Configuration")
		}

		// Server configuration is required for platform-specific sync
		if cfg.Server.Domain == "" {
			return fmt.Errorf("server domain cannot be empty when no raw git URLs are provided")
		}

		if cfg.Server.Protocol != "https" && cfg.Server.Protocol != "http" {
			return fmt.Errorf("server protocol can only be http or https")
		}

		// Workspace is required only for Bitbucket
		if cfg.Platform == "bitbucket" && cfg.Workspace == "" {
			return fmt.Errorf("workspace cannot be empty for bitbucket")
		}
	}

	return nil
}
