package sync

import (
	"fmt"
	"runtime"
	"sync"

	"github.com/AkashRajpurohit/git-sync/pkg/config"
	"github.com/AkashRajpurohit/git-sync/pkg/logger"
	"github.com/AkashRajpurohit/git-sync/pkg/notification"
	"github.com/AkashRajpurohit/git-sync/pkg/telemetry"
	"github.com/AkashRajpurohit/git-sync/pkg/version"
)

type SyncStats struct {
	mu            sync.Mutex
	ReposSuccess  int
	ReposFailed   []string
	WikisSuccess  int
	WikisFailed   []string
	IssuesSuccess int
	IssuesFailed  []string
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

func recordIssuesSuccess() {
	stats.mu.Lock()
	defer stats.mu.Unlock()
	stats.IssuesSuccess++
}

func recordIssuesFailure(repoName string, err error) {
	stats.mu.Lock()
	defer stats.mu.Unlock()
	stats.IssuesFailed = append(stats.IssuesFailed, fmt.Sprintf("%s (Error: %v)", repoName, err))
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

	logger.Infof("✅ Issues: %d repositories' issues synced", stats.IssuesSuccess)
	if len(stats.IssuesFailed) > 0 {
		failedIssues := []string{}
		for _, issue := range stats.IssuesFailed {
			failedIssues = append(failedIssues, issue)
		}

		logger.Errorf("❌ Failed issues: %d", len(failedIssues))
		logger.Errorf("%s", failedIssues)
	}

	summary := &notification.SyncSummary{
		ReposSuccess:  stats.ReposSuccess,
		ReposFailed:   stats.ReposFailed,
		WikisSuccess:  stats.WikisSuccess,
		WikisFailed:   stats.WikisFailed,
		IssuesSuccess: stats.IssuesSuccess,
		IssuesFailed:  stats.IssuesFailed,
	}

	if err := notification.NotifyAll(&cfg.Notification, summary); err != nil {
		logger.Errorf("Failed to send notifications: %v", err)
	}

	telemetry.CaptureEvent("sync_completed", map[string]interface{}{
		"platform":       cfg.Platform,
		"clone_type":     cfg.CloneType,
		"concurrency":    cfg.Concurrency,
		"include_wiki":   cfg.IncludeWiki,
		"include_issues": cfg.IncludeIssues,
		"include_forks":  cfg.IncludeForks,
		"repos_success":  stats.ReposSuccess,
		"repos_failed":   len(stats.ReposFailed),
		"wikis_success":  stats.WikisSuccess,
		"wikis_failed":   len(stats.WikisFailed),
		"issues_success": stats.IssuesSuccess,
		"issues_failed":  len(stats.IssuesFailed),
		"app_version":    version.Version,
		"os":             runtime.GOOS,
		"arch":           runtime.GOARCH,
	})

	// Reset stats for next sync
	stats = &SyncStats{}
}
