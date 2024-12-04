package sync

import (
	"sync"

	"github.com/AkashRajpurohit/git-sync/pkg/config"
	"github.com/AkashRajpurohit/git-sync/pkg/logger"
)

func SyncWithConcurrency(cfg config.Config, repos []string, syncFn func(string)) {
	var wg sync.WaitGroup
	sem := make(chan struct{}, cfg.Concurrency)

	for _, repo := range repos {
		wg.Add(1)
		go func(r string) {
			defer wg.Done()
			sem <- struct{}{}
			syncFn(r)
			<-sem
		}(repo)
	}

	wg.Wait()
}

func SyncReposWithConcurrency[T any](cfg config.Config, repos []T, syncFn func(T)) {
	var wg sync.WaitGroup
	sem := make(chan struct{}, cfg.Concurrency)

	for _, repo := range repos {
		wg.Add(1)
		go func(r T) {
			defer wg.Done()
			sem <- struct{}{}
			syncFn(r)
			<-sem
		}(repo)
	}

	wg.Wait()
}

func LogRepoCount(count int, repoType string) {
	logger.Info("Total ", repoType, " repositories: ", count)
}

func LogSyncComplete(repoType string) {
	logger.Info("All ", repoType, " repositories synced ✅")
}
