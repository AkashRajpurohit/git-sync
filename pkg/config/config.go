package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	Username        string   `mapstructure:"username"`
	Token           string   `mapstructure:"token"`
	Repos           []string `mapstructure:"repos"`
	IncludeAllRepos bool     `mapstructure:"include_all_repos"`
	IncludeForks    bool     `mapstructure:"include_forks"`
	BackupDir       string   `mapstructure:"backup_dir"`
}

func GetConfigFile(cfgFile string) string {
	if cfgFile != "" {
		return cfgFile
	}

	if os.Getenv("GIT_SYNC_CONFIG_FILE") != "" {
		return os.Getenv("GIT_SYNC_CONFIG_FILE")
	}

	return filepath.Join(os.Getenv("HOME"), ".config", "git-sync", "config.yaml")
}

func GetBackupDir(backupDir string) string {
	if backupDir != "" {
		return backupDir
	}

	if os.Getenv("GIT_SYNC_BACKUP_DIR") != "" {
		return os.Getenv("GIT_SYNC_BACKUP_DIR")
	}

	return filepath.Join(os.Getenv("HOME"), "git-backups")
}

func LoadConfig(cfgFile string) (Config, error) {
	var config Config
	configFile := GetConfigFile(cfgFile)

	viper.AddConfigPath(filepath.Dir(configFile))
	viper.SetConfigName("config")
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
	viper.Set("repos", config.Repos)
	viper.Set("include_all_repos", config.IncludeAllRepos)
	viper.Set("include_forks", config.IncludeForks)
	viper.Set("backup_dir", config.BackupDir)

	return viper.WriteConfig()
}
