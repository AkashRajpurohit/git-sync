package notification

import (
	"fmt"
	"strings"

	"github.com/AkashRajpurohit/git-sync/pkg/config"
	"github.com/AkashRajpurohit/git-sync/pkg/logger"
)

type NotificationProvider interface {
	Send(title, message string) error
}

type SyncSummary struct {
	ReposSuccess int
	ReposFailed  []string
	WikisSuccess int
	WikisFailed  []string
}

func (s *SyncSummary) HasFailures() bool {
	return len(s.ReposFailed) > 0 || len(s.WikisFailed) > 0
}

func (s *SyncSummary) FormatMessage() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("✅ Repositories: %d successfully synced\n", s.ReposSuccess))
	if len(s.ReposFailed) > 0 {
		sb.WriteString(fmt.Sprintf("❌ Failed repositories: %d\n", len(s.ReposFailed)))
		for _, repo := range s.ReposFailed {
			sb.WriteString(fmt.Sprintf("- %s\n", repo))
		}
	}

	sb.WriteString(fmt.Sprintf("✅ Wikis: %d successfully synced\n", s.WikisSuccess))
	if len(s.WikisFailed) > 0 {
		sb.WriteString(fmt.Sprintf("❌ Failed wikis: %d\n", len(s.WikisFailed)))
		for _, wiki := range s.WikisFailed {
			sb.WriteString(fmt.Sprintf("- %s\n", wiki))
		}
	}

	return sb.String()
}

func NotifyAll(cfg *config.NotificationConfig, summary *SyncSummary) error {
	if !cfg.Enabled {
		return nil
	}

	if cfg.OnlyFailures && !summary.HasFailures() {
		return nil
	}

	title := "Git-Sync Operation Summary"
	message := summary.FormatMessage()

	var errors []string

	if cfg.Ntfy != nil {
		ntfy := NewNtfyProvider(cfg.Ntfy)
		if err := ntfy.Send(title, message); err != nil {
			errors = append(errors, fmt.Sprintf("ntfy: %v", err))
		}
	}

	if cfg.Gotify != nil {
		gotify := NewGotifyProvider(cfg.Gotify)
		if err := gotify.Send(title, message); err != nil {
			errors = append(errors, fmt.Sprintf("gotify: %v", err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("notification errors: %s", strings.Join(errors, "; "))
	}

	logger.Infof("Notifications sent successfully")

	return nil
}
