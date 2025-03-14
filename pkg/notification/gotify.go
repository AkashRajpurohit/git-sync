package notification

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/AkashRajpurohit/git-sync/pkg/config"
)

type GotifyProvider struct {
	url      string
	appToken string
	priority int
}

type gotifyMessage struct {
	Title    string `json:"title"`
	Message  string `json:"message"`
	Priority int    `json:"priority"`
}

func NewGotifyProvider(cfg *config.GotifyConfig) *GotifyProvider {
	priority := cfg.Priority
	if priority == 0 || priority < 1 || priority > 5 {
		priority = 5
	}
	return &GotifyProvider{
		url:      strings.TrimRight(cfg.URL, "/"),
		appToken: cfg.AppToken,
		priority: priority,
	}
}

func (g *GotifyProvider) Send(title, message string) error {
	if g.url == "" {
		return fmt.Errorf("gotify URL is required")
	}
	if g.appToken == "" {
		return fmt.Errorf("gotify app token is required")
	}

	msg := gotifyMessage{
		Title:    title,
		Message:  message,
		Priority: g.priority,
	}

	payload, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %v", err)
	}

	url := fmt.Sprintf("%s/message?token=%s", g.url, g.appToken)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

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
