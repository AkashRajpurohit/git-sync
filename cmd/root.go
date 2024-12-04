package cmd

import (
	"fmt"
	"os"

	"github.com/AkashRajpurohit/git-sync/pkg/bitbucket"
	"github.com/AkashRajpurohit/git-sync/pkg/client"
	"github.com/AkashRajpurohit/git-sync/pkg/config"
	"github.com/AkashRajpurohit/git-sync/pkg/forgejo"
	"github.com/AkashRajpurohit/git-sync/pkg/github"
	"github.com/AkashRajpurohit/git-sync/pkg/gitlab"
	"github.com/AkashRajpurohit/git-sync/pkg/logger"
	"github.com/AkashRajpurohit/git-sync/pkg/raw"
	ch "github.com/robfig/cron/v3"
	"github.com/spf13/cobra"
)

var (
	cfgFile   string
	backupDir string
	logLevel  string = "info"
	cron      string
)

var rootCmd = &cobra.Command{
	Use:   "git-sync",
	Short: "A tool to backup and sync your git repositories",
	Run: func(cmd *cobra.Command, args []string) {
		logger.InitLogger(logLevel)

		var cfg config.Config
		cfg, err := config.LoadConfig(cfgFile)

		if err != nil {
			logger.Debugf("Error in loading config file: ", err)
			logger.Info("Config file not found, creating a new one...")

			cfg = config.GetInitialConfig()

			err = config.SaveConfig(cfg, cfgFile)
			if err != nil {
				logger.Fatal("Error in saving config file: ", err)
			}
		}

		// If backupDir option is passed in the command line, use that instead of the one in the config file
		if backupDir != "" {
			cfg.BackupDir = config.GetBackupDir(backupDir)
		}

		// If cron option is passed in the command line, use that instead of the one in the config file
		if cron != "" {
			cfg.Cron = cron
		}

		// If no clone_type is not set in the config file, set it to bare
		if cfg.CloneType == "" {
			cfg.CloneType = "bare"
		}

		logger.Info("Config loaded from: ", config.GetConfigFile(cfgFile))
		logger.Info("Validating config ⏳")

		err = config.ValidateConfig(cfg)
		if err != nil {
			logger.Fatalf("Error validating config: %s", err)
		}

		// Create backup directory if it doesn't exist
		os.MkdirAll(cfg.BackupDir, os.ModePerm)

		var platformClient client.Client
		var hasRawURLs bool = len(cfg.RawGitURLs) > 0

		// Only initialize platform client if raw URLs are not provided or if both are needed
		if !hasRawURLs || (cfg.Username != "" && cfg.Token != "") {
			switch cfg.Platform {
			case "github":
				platformClient = github.NewGitHubClient(cfg.Token)
			case "gitlab":
				platformClient = gitlab.NewGitlabClient(cfg.Server, cfg.Token)
			case "bitbucket":
				platformClient = bitbucket.NewBitbucketClient(cfg.Username, cfg.Token)
			case "forgejo":
				platformClient = forgejo.NewForgejoClient(cfg.Server, cfg.Token)
			default:
				if !hasRawURLs {
					logger.Fatalf("Platform %s not supported", cfg.Platform)
				}
			}
		}

		logger.Info("Valid config found ✅")
		if platformClient != nil {
			logger.Infof("Using Platform: %s", cfg.Platform)
		}
		if hasRawURLs {
			logger.Infof("Found %d raw git URLs to sync", len(cfg.RawGitURLs))
		}

		if cfg.Cron != "" {
			c := ch.New()
			_, err := c.AddFunc(cfg.Cron, func() {
				// First sync platform repositories if configured
				if platformClient != nil {
					if err := platformClient.Sync(cfg); err != nil {
						logger.Errorf("Error syncing platform repositories: %s", err)
					}
				}

				// Then sync raw git URLs if any
				if hasRawURLs {
					rawClient := raw.NewRawClient()
					if err := rawClient.Sync(cfg); err != nil {
						logger.Errorf("Error syncing raw repositories: %s", err)
					}
				}
			})

			if err != nil {
				logger.Fatalf("Error adding cron job: %s", err)
			}

			c.Start()
			logger.Infof("Cron job scheduled to run at: %s", cfg.Cron)

			// Wait indefinitely
			select {}
		} else {
			// First sync platform repositories if configured
			if platformClient != nil {
				if err := platformClient.Sync(cfg); err != nil {
					logger.Errorf("Error syncing platform repositories: %s", err)
				}
			}

			// Then sync raw git URLs if any
			if hasRawURLs {
				rawClient := raw.NewRawClient()
				if err := rawClient.Sync(cfg); err != nil {
					logger.Errorf("Error syncing raw repositories: %s", err)
				}
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
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/git-sync/config.yaml)")
	rootCmd.PersistentFlags().StringVar(&backupDir, "backup-dir", "", "directory to backup repositories (default is $HOME/git-backups)")
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "info", "log level (debug, info, warn, error, fatal)")
	rootCmd.PersistentFlags().StringVar(&cron, "cron", "", "cron expression to run the sync job periodically")
}
