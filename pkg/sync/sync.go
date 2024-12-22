package sync

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/AkashRajpurohit/git-sync/pkg/config"
	"github.com/AkashRajpurohit/git-sync/pkg/logger"
	"github.com/AkashRajpurohit/git-sync/pkg/token"
)

var tokenManager *token.Manager

func InitTokenManager(tokens []string) {
	tokenManager = token.NewManager(tokens)
}

func getBaseDirectoryPath(repoOwner, repoName string, config config.Config) string {
	return filepath.Join(config.BackupDir, repoOwner, repoName)
}

func getGitCloneCommand(CloneType, repoPath, repoURL string) *exec.Cmd {
	switch CloneType {
	case "bare":
		logger.Debugf("Cloning repo with bare clone type: %s", repoURL)
		return exec.Command("git", "clone", "--bare", repoURL, repoPath)
	case "full":
		logger.Debugf("Cloning repo with full clone type: %s", repoURL)
		return exec.Command("git", "clone", repoURL, repoPath)
	case "mirror":
		logger.Debugf("Cloning repo with mirror clone type: %s", repoURL)
		return exec.Command("git", "clone", "--mirror", repoURL, repoPath)
	case "shallow":
		logger.Debugf("Cloning repo with shallow clone type: %s", repoURL)
		return exec.Command("git", "clone", "--depth", "1", repoURL, repoPath)
	default:
		logger.Debugf("[Default] Cloning repo with bare clone type: %s", repoURL)
		return exec.Command("git", "clone", "--bare", repoURL, repoPath)
	}
}

func getGitFetchCommand(CloneType, repoPath, repoURL string) *exec.Cmd {
	switch CloneType {
	case "bare":
		logger.Debugf("Updating repo with bare clone type: %s", repoPath)
		return exec.Command("git", "--git-dir", repoPath, "fetch", "--prune", repoURL, "+*:*")
	case "full":
		logger.Debugf("Updating repo with full clone type: %s", repoPath)
		return exec.Command("git", "-C", repoPath, "pull", "--prune", repoURL)
	case "mirror":
		logger.Debugf("Updating repo with mirror clone type: %s", repoPath)
		return exec.Command("git", "-C", repoPath, "fetch", "--prune", repoURL, "+*:*")
	case "shallow":
		logger.Debugf("Updating repo with shallow clone type: %s", repoPath)
		return exec.Command("git", "-C", repoPath, "pull", "--prune", repoURL)
	default:
		logger.Debugf("[Default] Updating repo with bare clone type: %s", repoPath)
		return exec.Command("git", "--git-dir", repoPath, "fetch", "--prune", repoURL, "+*:*")
	}
}

func CloneOrUpdateRepo(repoOwner, repoName string, config config.Config) {
	if tokenManager == nil {
		InitTokenManager(config.Tokens)
	}

	repoFullName := fmt.Sprintf("%s/%s", repoOwner, repoName)
	repoURL := fmt.Sprintf("%s://%s:%s@%s/%s.git", config.Server.Protocol, config.Username, tokenManager.GetNextToken(), config.Server.Domain, repoFullName)
	repoPath := filepath.Join(getBaseDirectoryPath(repoOwner, repoName, config), repoName+".git")

	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		logger.Info("Cloning repo: ", repoFullName)
		command := getGitCloneCommand(config.CloneType, repoPath, repoURL)

		err := retryOperation(config, func() error {
			output, err := command.CombinedOutput()
			logger.Debugf("Output: %s\n", output)
			return err
		}, fmt.Sprintf("clone %s", repoFullName))

		if err != nil {
			logger.Errorf("Failed to clone repo %s: %v", repoFullName, err)
			recordRepoFailure(repoFullName, err)
			return
		}

		logger.Info("Cloned repo: ", repoFullName)
		recordRepoSuccess()
	} else {
		logger.Info("Updating repo: ", repoFullName)
		command := getGitFetchCommand(config.CloneType, repoPath, repoURL)

		err := retryOperation(config, func() error {
			output, err := command.CombinedOutput()
			logger.Debugf("Output: %s\n", output)
			return err
		}, fmt.Sprintf("update %s", repoFullName))

		if err != nil {
			logger.Errorf("Failed to update repo %s: %v", repoFullName, err)
			recordRepoFailure(repoFullName, err)
			return
		}

		logger.Info("Updated repo: ", repoFullName)
		recordRepoSuccess()
	}
}

func CloneOrUpdateRawRepo(repoOwner, repoName, repoURL string, config config.Config) {
	repoPath := filepath.Join(getBaseDirectoryPath(repoOwner, repoName, config), repoName+".git")

	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		logger.Info("Cloning raw repo: ", repoURL)
		command := getGitCloneCommand(config.CloneType, repoPath, repoURL)

		err := retryOperation(config, func() error {
			output, err := command.CombinedOutput()
			logger.Debugf("Output: %s\n", output)
			return err
		}, fmt.Sprintf("clone %s", repoURL))

		if err != nil {
			logger.Errorf("Failed to clone raw repo %s: %v", repoURL, err)
			recordRepoFailure(repoURL, err)
			return
		}

		logger.Info("Cloned raw repo: ", repoURL)
		recordRepoSuccess()
	} else {
		logger.Info("Updating raw repo: ", repoURL)
		command := getGitFetchCommand(config.CloneType, repoPath, repoURL)

		err := retryOperation(config, func() error {
			output, err := command.CombinedOutput()
			logger.Debugf("Output: %s\n", output)
			return err
		}, fmt.Sprintf("update %s", repoURL))

		if err != nil {
			logger.Errorf("Failed to update raw repo %s: %v", repoURL, err)
			recordRepoFailure(repoURL, err)
			return
		}

		logger.Info("Updated raw repo: ", repoURL)
		recordRepoSuccess()
	}
}

func SyncWiki(repoOwner, repoName string, config config.Config) {
	if tokenManager == nil {
		InitTokenManager(config.Tokens)
	}

	repoFullName := fmt.Sprintf("%s/%s", repoOwner, repoName)
	repoWikiURL := fmt.Sprintf("%s://%s:%s@%s/%s.wiki.git", config.Server.Protocol, config.Username, tokenManager.GetNextToken(), config.Server.Domain, repoFullName)
	repoWikiPath := filepath.Join(getBaseDirectoryPath(repoOwner, repoName, config), repoName+".wiki.git")

	// Special handling for bitbucket since it does not follow the traditional pattern for wiki repos
	// @see here: https://support.atlassian.com/bitbucket-cloud/docs/clone-a-wiki/
	if config.Platform == "bitbucket" {
		repoWikiURL = fmt.Sprintf("%s://%s:%s@%s/%s.git/wiki", config.Server.Protocol, config.Username, tokenManager.GetNextToken(), config.Server.Domain, repoFullName)
	}

	if _, err := os.Stat(repoWikiPath); os.IsNotExist(err) {
		logger.Info("Cloning wiki: ", repoFullName)
		command := exec.Command("git", "clone", repoWikiURL, repoWikiPath)
		wikiNotFound := false

		err := retryOperation(config, func() error {
			output, err := command.CombinedOutput()
			logger.Debugf("Output: %s\n", output)
			if err != nil && strings.Contains(string(output), "not found") {
				wikiNotFound = true
				// Don't retry for non-existent wikis
				return nil
			}
			return err
		}, fmt.Sprintf("clone wiki %s", repoFullName))

		if err != nil && !wikiNotFound {
			logger.Errorf("Failed to clone wiki %s: %v", repoFullName, err)
			recordWikiFailure(repoFullName, err)
			return
		}

		if wikiNotFound {
			logger.Warnf("The wiki for repository %s does not exist. Please check your repository settings and make sure that either wiki is disabled if it is not being used or create a wiki page to start with.", repoFullName)
		} else {
			logger.Info("Cloned wiki: ", repoFullName)
			recordWikiSuccess()
		}
	} else {
		logger.Info("Updating wiki: ", repoFullName)
		command := exec.Command("git", "-C", repoWikiPath, "pull", "--prune", "origin")

		err := retryOperation(config, func() error {
			output, err := command.CombinedOutput()
			logger.Debugf("Output: %s\n", output)
			return err
		}, fmt.Sprintf("update wiki %s", repoFullName))

		if err != nil {
			logger.Errorf("Failed to update wiki %s: %v", repoFullName, err)
			recordWikiFailure(repoFullName, err)
			return
		}

		logger.Info("Updated wiki: ", repoFullName)
		recordWikiSuccess()
	}
}
