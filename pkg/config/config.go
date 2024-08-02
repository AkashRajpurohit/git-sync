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

type Config struct {
	Username     string   `mapstructure:"username"`
	Token        string   `mapstructure:"token"`
	Platform     string   `mapstructure:"platform"`
	Server       Server   `mapstructure:"server"`
	IncludeRepos []string `mapstructure:"include_repos"`
	ExcludeRepos []string `mapstructure:"exclude_repos"`
	IncludeOrgs  []string `mapstructure:"include_orgs"`
	ExcludeOrgs  []string `mapstructure:"exclude_orgs"`
	IncludeForks bool     `mapstructure:"include_forks"`
	BackupDir    string   `mapstructure:"backup_dir"`
	Workspace    string   `mapstructure:"workspace"`
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
	return config, err
}

func SaveConfig(config Config, cfgFile string) error {
	configFile := GetConfigFile(cfgFile)

	os.MkdirAll(filepath.Dir(configFile), os.ModePerm)
	viper.SetConfigFile(configFile)

	viper.Set("username", config.Username)
	viper.Set("token", config.Token)
	viper.Set("include_repos", config.IncludeRepos)
	viper.Set("exclude_repos", config.ExcludeRepos)
	viper.Set("include_orgs", config.IncludeOrgs)
	viper.Set("exclude_orgs", config.ExcludeOrgs)
	viper.Set("include_forks", config.IncludeForks)
	viper.Set("backup_dir", config.BackupDir)
	viper.Set("platform", config.Platform)
	viper.Set("server", config.Server)
	viper.Set("workspace", config.Workspace)

	return viper.WriteConfig()
}
