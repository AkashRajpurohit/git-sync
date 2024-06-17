package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/AkashRajpurohit/git-sync/config"
	"github.com/AkashRajpurohit/git-sync/github"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "git-sync",
	Short: "A tool to backup and sync your git repositories",
	Run: func(cmd *cobra.Command, args []string) {
		var cfg config.Config
		cfg, err := config.LoadConfig()

		if err != nil {
			log.Default().Println("Config file not found, creating a new one...")

			cfg = config.Config{
				Username:        "",
				Token:           "",
				Repos:           []string{},
				IncludeAllRepos: true,
			}

			err = config.SaveConfig(cfg)
		}

		if err != nil {
			log.Fatal("Error in saving/loading config file: ", err)
		}

		log.Default().Println("Config loaded from: ", config.GetConfigPath())

		if cfg.Username == "" {
			log.Fatal("No username found in config file, please add one.")
		}

		if cfg.Token == "" {
			log.Fatal("No token found in config file, please add one. See here: https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens#about-personal-access-tokens")
		}

		ghClient := github.NewClient(cfg.Username, cfg.Token)

		if cfg.IncludeAllRepos {
			repos, err := ghClient.FetchAllRepos()
			if err != nil {
				log.Fatal("Error in fetching repositories: ", err)
			}

			log.Default().Println("Total repositories found: ", len(repos))

			for _, repo := range repos {
				log.Default().Println(*repo.FullName)
			}

		} else {
			repos, err := ghClient.FetchRepos(cfg.Repos)
			if err != nil {
				log.Fatal("Error in fetching repositories: ", err)
			}

			log.Default().Println("Total repositories found: ", len(repos))

			for _, repo := range repos {
				log.Default().Println(*repo.FullName)
			}
		}

	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringP("config", "c", "", "config file (default is $HOME/.config/git-sync/config.yaml)")
}
