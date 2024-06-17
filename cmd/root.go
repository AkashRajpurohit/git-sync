package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/AkashRajpurohit/git-sync/config"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "git-sync",
	Short: "A tool to backup and sync your git repositories",
	Run: func(cmd *cobra.Command, args []string) {
		var cfg config.Config
		_, err := config.LoadConfig()

		if err != nil {
			fmt.Println("Config file not found, creating a new one...")

			cfg = config.Config{
				Username:        "username",
				Token:           "",
				Repos:           []string{},
				IncludeAllRepos: true,
			}

			err = config.SaveConfig(cfg)
		}

		if err != nil {
			log.Fatal("Error in saving/loading config file: ", err)
		}

		fmt.Println("Config loaded from: ", config.GetConfigPath())
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
