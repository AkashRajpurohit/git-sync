package config

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/AkashRajpurohit/git-sync/pkg/logger"
	"github.com/spf13/viper"
)

type Server struct {
	Domain   string `mapstructure:"domain"`
	Protocol string `mapstructure:"protocol"`
}

type RetryConfig struct {
	Count int `mapstructure:"count"`
	Delay int `mapstructure:"delay"` // in seconds
}

type Config struct {
	Username     string      `mapstructure:"username"`
	Token        string      `mapstructure:"token"`  // Deprecated: Use Tokens instead
	Tokens       []string    `mapstructure:"tokens"` // New field for multiple tokens
	Platform     string      `mapstructure:"platform"`
	Server       Server      `mapstructure:"server"`
	IncludeRepos []string    `mapstructure:"include_repos"`
	ExcludeRepos []string    `mapstructure:"exclude_repos"`
	IncludeOrgs  []string    `mapstructure:"include_orgs"`
	ExcludeOrgs  []string    `mapstructure:"exclude_orgs"`
	IncludeForks bool        `mapstructure:"include_forks"`
	IncludeWiki  bool        `mapstructure:"include_wiki"`
	BackupDir    string      `mapstructure:"backup_dir"`
	Workspace    string      `mapstructure:"workspace"`
	Cron         string      `mapstructure:"cron"`
	CloneType    string      `mapstructure:"clone_type"`
	RawGitURLs   []string    `mapstructure:"raw_git_urls"`
	Concurrency  int         `mapstructure:"concurrency"`
	Retry        RetryConfig `mapstructure:"retry"`
}

// PreprocessConfig handles backward compatibility for token field
func PreprocessConfig(cfg *Config) {
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

func expandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		homeDir, _ := os.UserHomeDir()
		return filepath.Join(homeDir, path[2:])
	}
	return path
}

func GetConfigFile(cfgFile string) string {
	if cfgFile != "" {
		logger.Debug("Using config file: ", cfgFile)
		return expandPath(cfgFile)
	}

	if os.Getenv("GIT_SYNC_CONFIG_FILE") != "" {
		logger.Debug("Using OS env GIT_SYNC_CONFIG_FILE: ", os.Getenv("GIT_SYNC_CONFIG_FILE"))
		return expandPath(os.Getenv("GIT_SYNC_CONFIG_FILE"))
	}

	defaultConfigFilePath := filepath.Join(os.Getenv("HOME"), ".config", "git-sync", "config.yaml")
	logger.Debug("Using default config file: ", defaultConfigFilePath)
	return expandPath(defaultConfigFilePath)
}

func GetBackupDir(backupDir string) string {
	if backupDir != "" {
		logger.Debug("Using backup directory: ", backupDir)
		return expandPath(backupDir)
	}

	if os.Getenv("GIT_SYNC_BACKUP_DIR") != "" {
		logger.Debug("Using OS env GIT_SYNC_BACKUP_DIR: ", os.Getenv("GIT_SYNC_BACKUP_DIR"))
		return expandPath(os.Getenv("GIT_SYNC_BACKUP_DIR"))
	}

	defaultBackupDir := filepath.Join(os.Getenv("HOME"), "git-backups")
	logger.Debug("Using default backup directory: ", defaultBackupDir)
	return expandPath(defaultBackupDir)
}

func LoadConfig(cfgFile string) (Config, error) {
	var config Config
	configFile := GetConfigFile(cfgFile)

	viper.AddConfigPath(filepath.Dir(configFile))
	viper.SetConfigFile(configFile)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return config, err
	}

	err := viper.Unmarshal(&config)
	if err != nil {
		return config, err
	}

	PreprocessConfig(&config)
	return config, nil
}

func SaveConfig(config Config, cfgFile string) error {
	configFile := GetConfigFile(cfgFile)

	os.MkdirAll(filepath.Dir(configFile), os.ModePerm)
	viper.SetConfigFile(configFile)

	viper.Set("username", config.Username)
	viper.Set("tokens", config.Tokens) // Save only tokens array
	viper.Set("include_repos", config.IncludeRepos)
	viper.Set("exclude_repos", config.ExcludeRepos)
	viper.Set("include_orgs", config.IncludeOrgs)
	viper.Set("exclude_orgs", config.ExcludeOrgs)
	viper.Set("include_forks", config.IncludeForks)
	viper.Set("include_wiki", config.IncludeWiki)
	viper.Set("backup_dir", config.BackupDir)
	viper.Set("platform", config.Platform)
	viper.Set("server", config.Server)
	viper.Set("workspace", config.Workspace)
	viper.Set("cron", config.Cron)
	viper.Set("clone_type", config.CloneType)
	viper.Set("raw_git_urls", config.RawGitURLs)
	viper.Set("concurrency", config.Concurrency)
	viper.Set("retry", config.Retry)

	return viper.WriteConfig()
}

func GetInitialConfig() Config {
	return Config{
		Username: "",
		Tokens:   []string{},
		Platform: "github",
		Server: Server{
			Domain:   "github.com",
			Protocol: "https",
		},
		IncludeRepos: []string{},
		ExcludeRepos: []string{},
		IncludeOrgs:  []string{},
		ExcludeOrgs:  []string{},
		IncludeForks: false,
		IncludeWiki:  true,
		Workspace:    "",
		Cron:         "",
		BackupDir:    GetBackupDir(""),
		CloneType:    "bare",
		RawGitURLs:   []string{},
		Concurrency:  5,
		Retry: RetryConfig{
			Count: 3,
			Delay: 5,
		},
	}
}
