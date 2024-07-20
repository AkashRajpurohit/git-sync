package sync

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/AkashRajpurohit/git-sync/pkg/config"
	"github.com/google/go-github/v62/github"
)

func SyncRepos(config config.Config, repos []*github.Repository) {
	backupDir := config.BackupDir
	os.MkdirAll(backupDir, os.ModePerm)

	var wg sync.WaitGroup
	sem := make(chan struct{}, 10) // Concurrency of 10

	for _, repo := range repos {
		if ShouldSync(repo.GetName(), config.Repos) {
			wg.Add(1)
			go func(repo *github.Repository) {
				defer wg.Done()
				sem <- struct{}{}
				CloneOrUpdateRepo(repo, backupDir, config)
				<-sem
			}(repo)
		}
	}

	wg.Wait()
}

func ShouldSync(repoName string, configuredRepos []string) bool {
	if len(configuredRepos) == 0 {
		return true
	}

	for _, name := range configuredRepos {
		if name == repoName {
			return true
		}
	}

	return false
}

func CloneOrUpdateRepo(repo *github.Repository, backupDir string, config config.Config) {
	repoURL := fmt.Sprintf("https://%s:%s@github.com/%s.git", config.Username, config.Token, repo.GetFullName())
	repoPath := filepath.Join(backupDir, repo.GetName()+".git")

	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		log.Default().Println("Cloning repo:", repo.GetName())

		cmd := exec.Command("git", "clone", "--bare", repoURL, repoPath)
		if err := cmd.Run(); err != nil {
			log.Printf("Error cloning repo %s: %v\n", repo.GetName(), err)
		} else {
			log.Default().Println("Cloned repo:", repo.GetName())
		}
	} else {
		log.Default().Println("Updating repo:", repo.GetName())

		cmd := exec.Command("git", "--git-dir", repoPath, "fetch", "--prune", "origin", "*:*")
		if err := cmd.Run(); err != nil {
			log.Printf("Error updating repo %s: %v\n", repo.GetName(), err)
		} else {
			log.Default().Println("Updated repo:", repo.GetName())
		}
	}
}
