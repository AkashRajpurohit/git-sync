package sync

import (
	"fmt"
	"time"

	"github.com/AkashRajpurohit/git-sync/pkg/config"
	"github.com/AkashRajpurohit/git-sync/pkg/logger"
)

func retryOperation(cfg config.Config, operation func() error, operationName string) error {
	var lastErr error

	// If retry count is 0 or negative, just execute once without retries
	if cfg.Retry.Count <= 0 {
		return operation()
	}

	for attempt := 1; attempt <= cfg.Retry.Count; attempt++ {
		err := operation()
		if err == nil {
			if attempt > 1 {
				logger.Warnf("Operation %s succeeded after %d attempts", operationName, attempt)
			}
			return nil
		}

		lastErr = err
		if attempt < cfg.Retry.Count {
			logger.Warnf("Attempt %d/%d failed for %s: %v. Retrying in %d seconds...",
				attempt, cfg.Retry.Count, operationName, err, cfg.Retry.Delay)
			time.Sleep(time.Duration(cfg.Retry.Delay) * time.Second)
		}
	}

	return fmt.Errorf("operation failed after %d attempts: %v", cfg.Retry.Count, lastErr)
}
