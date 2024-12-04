package raw

import (
	"path/filepath"
	"strings"

	"github.com/AkashRajpurohit/git-sync/pkg/config"
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

	gitSync.LogRepoCount(len(cfg.RawGitURLs), "raw")

	gitSync.SyncWithConcurrency(cfg, cfg.RawGitURLs, func(url string) {
		owner, name := c.extractRepoInfo(url)
		gitSync.CloneOrUpdateRawRepo(owner, name, url, cfg)
	})

	gitSync.LogSyncComplete("raw")
	return nil
}
