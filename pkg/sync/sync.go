package sync

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

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

		output, err := exec.Command("git", "clone", "--bare", repoURL, repoPath).CombinedOutput()
		logger.Debugf("Output: %s\n", output)
		if err != nil {
			logger.Fatalf("Error cloning repo %s: %v\n", repoFullName, err)
		} else {
			logger.Info("Cloned repo: ", repoFullName)
		}
	} else {
		logger.Info("Updating repo: ", repoFullName)

		output, err := exec.Command("git", "--git-dir", repoPath, "fetch", "--prune", "origin", "+*:*").CombinedOutput()
		logger.Debugf("Output: %s\n", output)
		if err != nil {
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

		output, err := exec.Command("git", "clone", repoWikiURL, repoWikiPath).CombinedOutput()
		logger.Debugf("Output: %s\n", output)
		if err != nil {
			exitErr, ok := err.(*exec.ExitError)
			if ok && exitErr.ExitCode() == 128 {
				// Check if the output contains "not found" to handle the scenario
				// where wiki is enabled but does not exist
				if strings.Contains(string(output), "not found") {
					logger.Warnf("The wiki for repository %s does not exist. Please check your repository settings and make sure that either wiki is disabled if it is not being used or create a wiki page to start with.", repoFullName)
				} else {
					logger.Fatalf("Error cloning wiki %s: %v\n", repoFullName, err)
				}
			} else {
				logger.Fatalf("Error cloning wiki %s: %v\n", repoFullName, err)
			}
		} else {
			logger.Info("Cloned wiki: ", repoFullName)
		}
	} else {
		logger.Info("Updating wiki: ", repoFullName)

		output, err := exec.Command("git", "-C", repoWikiPath, "pull", "--prune", "origin").CombinedOutput()
		logger.Debugf("Output: %s\n", output)
		if err != nil {
			logger.Fatalf("Error updating wiki %s: %v\n", repoFullName, err)
		} else {
			logger.Info("Updated wiki: ", repoFullName)
		}
	}
}
