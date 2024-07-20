package cmd

import (
	"fmt"
	"os"

	"github.com/AkashRajpurohit/git-sync/pkg/config"
	"github.com/AkashRajpurohit/git-sync/pkg/github"
	"github.com/AkashRajpurohit/git-sync/pkg/logger"
	"github.com/AkashRajpurohit/git-sync/pkg/sync"
	"github.com/spf13/cobra"
)

var (
	cfgFile   string
	backupDir string
	logLevel  string = "info"
)

var rootCmd = &cobra.Command{
	Use:   "git-sync",
	Short: "A tool to backup and sync your git repositories",
	Run: func(cmd *cobra.Command, args []string) {
		logger.InitLogger(logLevel)

		var cfg config.Config
		cfg, err := config.LoadConfig(cfgFile)

		if err != nil {
			logger.Info("Config file not found, creating a new one...")

			cfg = config.Config{
				Username:     "",
				Token:        "",
				IncludeRepos: []string{},
				ExcludeRepos: []string{},
				IncludeForks: false,
				BackupDir:    config.GetBackupDir(backupDir),
			}

			err = config.SaveConfig(cfg, cfgFile)
		}

		if err != nil {
			logger.Fatal("Error in saving/loading config file: ", err)
		}

		logger.Info("Config loaded from: ", config.GetConfigFile(cfgFile))
		logger.Info("Validating config ⏳")

		if cfg.Username == "" {
			logger.Fatal("No username found in config file, please add one.")
		}

		if cfg.Token == "" {
			logger.Fatal("No token found in config file, please add one. See here: https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens#about-personal-access-tokens")
		}

		if cfg.BackupDir == "" {
			logger.Fatal("No backup directory found in config file, please add one.")
		}

		if len(cfg.IncludeRepos) > 0 && len(cfg.ExcludeRepos) > 0 {
			logger.Warn("Both include and exclude repos are set, ignoring exclude repos")
		}

		logger.Info("Valid config found ✅")

		repos, err := github.GetGitHubRepos(cfg)

		if err != nil {
			logger.Fatalf("Error fetching repositories: ", err)
		}

		logger.Info("Total repositories: ", len(repos))

		sync.SyncRepos(cfg, repos)

		logger.Info("All repositories synced ✅")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/git-sync/config.yaml)")
	rootCmd.PersistentFlags().StringVar(&backupDir, "backup-dir", "", "directory to backup repositories (default is $HOME/git-backups)")
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "info", "log level (debug, info, warn, error, fatal)")
}
