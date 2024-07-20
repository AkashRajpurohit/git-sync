package cmd

import (
	"fmt"

	"github.com/AkashRajpurohit/git-sync/version"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of git-sync",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("git-sync\nVersion: %s\nBuild: %s", version.Version, version.Build)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
