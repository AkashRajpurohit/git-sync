package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/AkashRajpurohit/git-sync/pkg/config"
	"github.com/AkashRajpurohit/git-sync/pkg/github"
	"github.com/AkashRajpurohit/git-sync/pkg/sync"
	gh "github.com/google/go-github/v62/github"
	"github.com/spf13/cobra"
)

var (
	cfgFile   string
	backupDir string
)

var rootCmd = &cobra.Command{
	Use:   "git-sync",
	Short: "A tool to backup and sync your git repositories",
	Run: func(cmd *cobra.Command, args []string) {
		var cfg config.Config
		cfg, err := config.LoadConfig(cfgFile)

		if err != nil {
			log.Default().Println("Config file not found, creating a new one...")

			cfg = config.Config{
				Username:        "",
				Token:           "",
				Repos:           []string{},
				IncludeAllRepos: true,
				IncludeForks:    false,
				BackupDir:       config.GetBackupDir(backupDir),
			}

			err = config.SaveConfig(cfg, cfgFile)
		}

		if err != nil {
			log.Fatal("Error in saving/loading config file: ", err)
		}

		log.Default().Println("Config loaded from:", config.GetConfigFile(cfgFile))

		if cfg.Username == "" {
			log.Fatal("No username found in config file, please add one.")
		}

		if cfg.Token == "" {
			log.Fatal("No token found in config file, please add one. See here: https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens#about-personal-access-tokens")
		}

		if cfg.BackupDir == "" {
			log.Fatal("No backup directory found in config file, please add one.")
		}

		ghClient := github.NewClient(cfg.Username, cfg.Token)
		var repos []*gh.Repository

		if cfg.IncludeAllRepos {
			r, err := ghClient.FetchAllRepos(cfg)
			if err != nil {
				log.Fatal("Error in fetching repositories: ", err)
			}

			log.Default().Println("Total repositories found:", len(r))

			repos = r
		} else {
			r, err := ghClient.FetchRepos(cfg)
			if err != nil {
				log.Fatal("Error in fetching repositories: ", err)
			}

			log.Default().Println("Total repositories found:", len(r))

			repos = r
		}

		sync.SyncRepos(cfg, repos)
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
}
