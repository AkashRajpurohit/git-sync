package sync

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/AkashRajpurohit/git-sync/pkg/config"
	"github.com/AkashRajpurohit/git-sync/pkg/logger"
	"github.com/google/go-github/v62/github"
)

func SyncRepos(config config.Config, repos []*github.Repository) {
	backupDir := config.BackupDir
	os.MkdirAll(backupDir, os.ModePerm)

	var wg sync.WaitGroup
	sem := make(chan struct{}, 10) // Concurrency of 10

	for _, repo := range repos {
		wg.Add(1)
		go func(repo *github.Repository) {
			defer wg.Done()
			sem <- struct{}{}
			cloneOrUpdateRepo(repo, backupDir, config)
			<-sem
		}(repo)
	}

	wg.Wait()
}

func cloneOrUpdateRepo(repo *github.Repository, backupDir string, config config.Config) {
	repoURL := fmt.Sprintf("https://%s:%s@github.com/%s.git", config.Username, config.Token, repo.GetFullName())
	repoPath := filepath.Join(backupDir, repo.GetName()+".git")

	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		logger.Info("Cloning repo: ", repo.GetName())

		cmd := exec.Command("git", "clone", "--bare", repoURL, repoPath)
		if err := cmd.Run(); err != nil {
			logger.Fatalf("Error cloning repo %s: %v\n", repo.GetName(), err)
		} else {
			logger.Info("Cloned repo: ", repo.GetName())
		}
	} else {
		logger.Info("Updating repo: ", repo.GetName())

		cmd := exec.Command("git", "--git-dir", repoPath, "fetch", "--prune", "origin", "+*:*")
		if err := cmd.Run(); err != nil {
			logger.Fatalf("Error updating repo %s: %v\n", repo.GetName(), err)
		} else {
			logger.Info("Updated repo: ", repo.GetName())
		}
	}
}
