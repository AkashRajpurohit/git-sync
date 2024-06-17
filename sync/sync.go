package sync

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/AkashRajpurohit/git-sync/config"
	"github.com/google/go-github/v62/github"
)

func SyncRepos(config config.Config, repos []*github.Repository) {
	backupDir := config.BackupDir
	if backupDir == "" {
		backupDir = filepath.Join(os.Getenv("HOME"), "git-backups")
	}

	os.MkdirAll(backupDir, os.ModePerm)

	var wg sync.WaitGroup
	sem := make(chan struct{}, 10) // Concurrency of 10

	for _, repo := range repos {
		if ShouldSync(repo.GetName(), config.Repos) {
			wg.Add(1)
			go func(repo *github.Repository) {
				defer wg.Done()
				sem <- struct{}{}
				CloneOrUpdateRepo(repo, backupDir)
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

func CloneOrUpdateRepo(repo *github.Repository, backupDir string) {
	repoURL := repo.GetSSHURL()
	repoPath := filepath.Join(backupDir, repo.GetName()+".git")

	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		log.Default().Println("Cloning: ", repo.GetName())

		if _, err := exec.Command("git", "clone", "--bare", repoURL, repoPath).Output(); err != nil {
			log.Printf("Error cloning repo %s: %v", repo.GetName(), err)
		} else {
			log.Printf("Cloned repo: %s\n", repo.GetName())
		}
	} else {
		log.Printf("Updating repo: %s\n", repo.GetName())
		if _, err := exec.Command("git", "-C", repoPath, "fetch", "--all").Output(); err != nil {
			log.Printf("Error updating repo %s: %v", repo.GetName(), err)
		} else {
			log.Printf("Updated repo: %s\n", repo.GetName())
		}
	}
}
