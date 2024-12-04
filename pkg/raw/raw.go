package raw

import (
	"path/filepath"
	"strings"
	"sync"

	"github.com/AkashRajpurohit/git-sync/pkg/config"
	"github.com/AkashRajpurohit/git-sync/pkg/logger"
	gitSync "github.com/AkashRajpurohit/git-sync/pkg/sync"
)

type RawClient struct{}

func NewRawClient() *RawClient {
	return &RawClient{}
}

// extractRepoInfo extracts the owner and repo name from a git URL
func (c RawClient) extractRepoInfo(url string) (string, string) {
	// Remove .git suffix if present
	url = strings.TrimSuffix(url, ".git")

	// Split the URL into parts
	parts := strings.Split(url, "/")
	if len(parts) < 2 {
		return "raw", filepath.Base(url)
	}

	return parts[len(parts)-2], parts[len(parts)-1]
}

func (c RawClient) Sync(cfg config.Config) error {
	if len(cfg.RawGitURLs) == 0 {
		return nil
	}

	logger.Info("Total raw repositories: ", len(cfg.RawGitURLs))

	var wg sync.WaitGroup
	sem := make(chan struct{}, 10) // Concurrency of 10

	for _, url := range cfg.RawGitURLs {
		wg.Add(1)
		go func(repoURL string) {
			defer wg.Done()
			sem <- struct{}{}
			owner, name := c.extractRepoInfo(repoURL)
			gitSync.CloneOrUpdateRawRepo(owner, name, repoURL, cfg)
			<-sem
		}(url)
	}

	wg.Wait()
	logger.Info("All raw repositories synced ✅")

	return nil
}
