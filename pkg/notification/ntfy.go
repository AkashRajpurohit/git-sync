package notification

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/AkashRajpurohit/git-sync/pkg/config"
)

type NtfyProvider struct {
	config *config.NtfyConfig
}

func NewNtfyProvider(cfg *config.NtfyConfig) *NtfyProvider {
	if cfg.Server == "" {
		cfg.Server = "https://ntfy.sh"
	}

	if cfg.Priority == 0 || cfg.Priority < 1 || cfg.Priority > 5 {
		cfg.Priority = 3
	}

	if len(cfg.Tags) == 0 {
		cfg.Tags = []string{"git-sync"}
	} else if !contains(cfg.Tags, "git-sync") {
		cfg.Tags = append(cfg.Tags, "git-sync")
	}
	return &NtfyProvider{
		config: cfg,
	}
}

func (n *NtfyProvider) Send(title, message string) error {
	if n.config.Topic == "" {
		return fmt.Errorf("ntfy topic is required")
	}

	url := fmt.Sprintf("%s/%s", strings.TrimRight(n.config.Server, "/"), n.config.Topic)

	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(message))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Title", title)
	req.Header.Set("Priority", fmt.Sprintf("%d", n.config.Priority))
	req.Header.Set("Tags", strings.Join(n.config.Tags, ","))

	if n.config.Username != "" && n.config.Password != "" {
		req.SetBasicAuth(n.config.Username, n.config.Password)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send notification: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("notification failed with status %d", resp.StatusCode)
	}

	return nil
}

func contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}
