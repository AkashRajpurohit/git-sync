package sync

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/AkashRajpurohit/git-sync/config"
	"github.com/google/go-github/v62/github"
)

func SyncRepos(config config.Config, repos []*github.Repository) {
	backupDir := config.BackupDir
	if backupDir == "" {
		backupDir = filepath.Join(os.Getenv("HOME"), "git-backups")
	}

	os.MkdirAll(backupDir, os.ModePerm)

	for _, repo := range repos {
		if ShouldSync(repo.GetName(), config.Repos) {
			log.Default().Println("Syncing: ", *repo.FullName)
			CloneOrUpdateRepo(repo, backupDir)
		}
	}
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
