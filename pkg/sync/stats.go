package sync

import (
	"fmt"
	"sync"

	"github.com/AkashRajpurohit/git-sync/pkg/config"
	"github.com/AkashRajpurohit/git-sync/pkg/logger"
	"github.com/AkashRajpurohit/git-sync/pkg/notification"
)

type SyncStats struct {
	mu           sync.Mutex
	ReposSuccess int
	ReposFailed  []string
	WikisSuccess int
	WikisFailed  []string
}

var stats = &SyncStats{}

func recordRepoSuccess() {
	stats.mu.Lock()
	defer stats.mu.Unlock()
	stats.ReposSuccess++
}

func recordRepoFailure(repoName string, err error) {
	stats.mu.Lock()
	defer stats.mu.Unlock()
	stats.ReposFailed = append(stats.ReposFailed, fmt.Sprintf("%s (Error: %v)", repoName, err))
}

func recordWikiSuccess() {
	stats.mu.Lock()
	defer stats.mu.Unlock()
	stats.WikisSuccess++
}

func recordWikiFailure(wikiName string, err error) {
	stats.mu.Lock()
	defer stats.mu.Unlock()
	stats.WikisFailed = append(stats.WikisFailed, fmt.Sprintf("%s (Error: %v)", wikiName, err))
}

func LogRepoCount(count int, repoType string) {
	logger.Info("Total ", repoType, " repositories: ", count)
}

func LogSyncSummary(cfg *config.Config) {
	logger.Infof("✅ Repositories: %d successfully synced", stats.ReposSuccess)
	if len(stats.ReposFailed) > 0 {
		failedRepos := []string{}
		for _, repo := range stats.ReposFailed {
			failedRepos = append(failedRepos, repo)
		}

		logger.Errorf("❌ Failed repositories: %d", len(failedRepos))
		logger.Errorf("%s", failedRepos)
	}

	logger.Infof("✅ Wikis: %d successfully synced", stats.WikisSuccess)
	if len(stats.WikisFailed) > 0 {
		failedWikis := []string{}
		for _, wiki := range stats.WikisFailed {
			failedWikis = append(failedWikis, wiki)
		}

		logger.Errorf("❌ Failed wikis: %d", len(failedWikis))
		logger.Errorf("%s", failedWikis)
	}

	summary := &notification.SyncSummary{
		ReposSuccess: stats.ReposSuccess,
		ReposFailed:  stats.ReposFailed,
		WikisSuccess: stats.WikisSuccess,
		WikisFailed:  stats.WikisFailed,
	}

	if err := notification.NotifyAll(&cfg.Notification, summary); err != nil {
		logger.Errorf("Failed to send notifications: %v", err)
	}

	// Reset stats for next sync
	stats = &SyncStats{}
}
