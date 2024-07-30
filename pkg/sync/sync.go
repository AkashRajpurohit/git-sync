package sync

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/AkashRajpurohit/git-sync/pkg/config"
	"github.com/AkashRajpurohit/git-sync/pkg/logger"
)

func CloneOrUpdateRepo(repoOwner, repoName string, config config.Config) {
	repoFullName := fmt.Sprintf("%s/%s", repoOwner, repoName)
	repoURL := fmt.Sprintf("https://%s:%s@github.com/%s.git", config.Username, config.Token, repoFullName)
	repoPath := filepath.Join(config.BackupDir, repoOwner, repoName+".git")

	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		logger.Info("Cloning repo: ", repoFullName)

		cmd := exec.Command("git", "clone", "--bare", repoURL, repoPath)
		if err := cmd.Run(); err != nil {
			logger.Fatalf("Error cloning repo %s: %v\n", repoFullName, err)
		} else {
			logger.Info("Cloned repo: ", repoFullName)
		}
	} else {
		logger.Info("Updating repo: ", repoFullName)

		cmd := exec.Command("git", "--git-dir", repoPath, "fetch", "--prune", "origin", "+*:*")
		if err := cmd.Run(); err != nil {
			logger.Fatalf("Error updating repo %s: %v\n", repoFullName, err)
		} else {
			logger.Info("Updated repo: ", repoFullName)
		}
	}
}
