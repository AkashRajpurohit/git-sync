package sync

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/AkashRajpurohit/git-sync/pkg/config"
	"github.com/AkashRajpurohit/git-sync/pkg/logger"
)

func getBaseDirectoryPath(repoOwner, repoName string, config config.Config) string {
	return filepath.Join(config.BackupDir, repoOwner, repoName)
}

func CloneOrUpdateRepo(repoOwner, repoName string, config config.Config) {
	repoFullName := fmt.Sprintf("%s/%s", repoOwner, repoName)
	repoURL := fmt.Sprintf("%s://%s:%s@%s/%s.git", config.Server.Protocol, config.Username, config.Token, config.Server.Domain, repoFullName)
	repoPath := filepath.Join(getBaseDirectoryPath(repoOwner, repoName, config), repoName+".git")

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

func SyncWiki(repoOwner, repoName string, config config.Config) {
	repoFullName := fmt.Sprintf("%s/%s", repoOwner, repoName)
	repoWikiURL := fmt.Sprintf("%s://%s:%s@%s/%s.wiki.git", config.Server.Protocol, config.Username, config.Token, config.Server.Domain, repoFullName)
	repoWikiPath := filepath.Join(getBaseDirectoryPath(repoOwner, repoName, config), repoName+".wiki.git")

	// Special handling for bitbucket since it does not follow the traditional pattern for wiki repos
	// @see here: https://support.atlassian.com/bitbucket-cloud/docs/clone-a-wiki/
	if config.Platform == "bitbucket" {
		repoWikiURL = fmt.Sprintf("%s://%s:%s@%s/%s.git/wiki", config.Server.Protocol, config.Username, config.Token, config.Server.Domain, repoFullName)
	}

	if _, err := os.Stat(repoWikiPath); os.IsNotExist(err) {
		logger.Info("Cloning wiki: ", repoFullName)

		cmd := exec.Command("git", "clone", repoWikiURL, repoWikiPath)
		if err := cmd.Run(); err != nil {
			logger.Fatalf("Error cloning wiki %s: %v\n", repoFullName, err)
		} else {
			logger.Info("Cloned wiki: ", repoFullName)
		}
	} else {
		logger.Info("Updating wiki: ", repoFullName)

		cmd := exec.Command("git", "-C", repoWikiPath, "pull", "--prune", "origin")
		if err := cmd.Run(); err != nil {
			logger.Fatalf("Error updating wiki %s: %v\n", repoFullName, err)
		} else {
			logger.Info("Updated wiki: ", repoFullName)
		}
	}
}
