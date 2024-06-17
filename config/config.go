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
}

func GetConfigPath() string {
	if os.Getenv("GIT_SYNC_CONFIG_PATH") != "" {
		return os.Getenv("GIT_SYNC_CONFIG_PATH")
	}

	return filepath.Join(os.Getenv("HOME"), ".config", "git-sync")
}

func LoadConfig() (Config, error) {
	var config Config
	configPath := GetConfigPath()

	viper.AddConfigPath(configPath)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return config, err
	}

	err := viper.Unmarshal(&config)
	return config, err
}

func SaveConfig(config Config) error {
	configPath := GetConfigPath()

	os.MkdirAll(configPath, os.ModePerm)
	viper.SetConfigFile(filepath.Join(configPath, "config.yaml"))

	viper.Set("username", config.Username)
	viper.Set("token", config.Token)
	viper.Set("repos", config.Repos)
	viper.Set("include_all_repos", config.IncludeAllRepos)

	return viper.WriteConfig()
}
