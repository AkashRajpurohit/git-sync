package config

import "github.com/AkashRajpurohit/git-sync/pkg/logger"

func SetSensibleDefaults(cfg *Config) {
	if cfg.Platform != "" && (cfg.Server.Domain == "" || cfg.Server.Protocol == "") {
		if cfg.Platform == "github" {
			cfg.Server.Domain = "github.com"
			cfg.Server.Protocol = "https"
		}

		if cfg.Platform == "gitlab" {
			cfg.Server.Domain = "gitlab.com"
			cfg.Server.Protocol = "https"
		}

		if cfg.Platform == "bitbucket" {
			cfg.Server.Domain = "bitbucket.org"
			cfg.Server.Protocol = "https"
		}

		if cfg.Platform == "forgejo" {
			cfg.Server.Domain = "v9.next.forgejo.org"
			cfg.Server.Protocol = "https"
		}

		if cfg.Platform == "gitea" {
			cfg.Server.Domain = "gitea.com"
			cfg.Server.Protocol = "https"
		}
	}

	// TODO: Remove these before v1.0.0 release
	// If concurrency is not set, set it to 5
	if cfg.Concurrency == 0 {
		logger.Warn("Concurrency is required but not set. Add the 'concurrency' field to the config file as mentioned in the docs: https://github.com/AkashRajpurohit/git-sync/wiki/Configuration. Setting it to 5.")
		cfg.Concurrency = 5
	}

	// If no clone_type is not set in the config file, set it to bare
	if cfg.CloneType == "" {
		logger.Warn("Clone type is required but not set. Add the 'clone_type' field to the config file as mentioned in the docs: https://github.com/AkashRajpurohit/git-sync/wiki/Configuration. Setting it to 'bare'.")
		cfg.CloneType = "bare"
	}

	// If both are set, merge them with single token being first
	if cfg.Token != "" && len(cfg.Tokens) > 0 {
		logger.Warn("Both 'token' and 'tokens' fields are set. 'token' field is deprecated and will be merged with 'tokens'.")
		cfg.Tokens = append([]string{cfg.Token}, cfg.Tokens...)
	}

	// If single token is set and tokens array is empty, convert single token to array
	if cfg.Token != "" && len(cfg.Tokens) == 0 {
		logger.Warn("Using 'token' field is deprecated. Please use 'tokens' array instead.")
		cfg.Tokens = []string{cfg.Token}
	}

	// Clear the deprecated token field
	cfg.Token = ""
}
