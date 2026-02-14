package telemetry

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/AkashRajpurohit/git-sync/pkg/config"
	"github.com/AkashRajpurohit/git-sync/pkg/logger"
	"github.com/google/uuid"
	"github.com/posthog/posthog-go"
)

const posthogAPIKey = "phc_IUI4rNyfkWwpSNPvWyrO0bOQkI0byAGyDQRKIQMGc7w"

var (
	client   posthog.Client
	deviceID string
	disabled bool
	once     sync.Once
)

func Init(cfg config.TelemetryConfig) {
	once.Do(func() {
		if isOptedOut(cfg) {
			disabled = true
			logger.Debug("Telemetry is disabled")
			return
		}

		id, err := getOrCreateDeviceID()
		if err != nil {
			logger.Debug("Failed to get device ID, disabling telemetry: ", err)
			disabled = true
			return
		}
		deviceID = id

		c, err := posthog.NewWithConfig(posthogAPIKey, posthog.Config{
			Endpoint: "https://us.i.posthog.com",
		})
		if err != nil {
			logger.Debug("Failed to create PostHog client, disabling telemetry: ", err)
			disabled = true
			return
		}

		client = c
		logger.Debug("Telemetry initialized")
	})
}

func CaptureEvent(event string, properties map[string]interface{}) {
	if disabled || client == nil {
		return
	}

	props := posthog.NewProperties()
	for k, v := range properties {
		props.Set(k, v)
	}

	client.Enqueue(posthog.Capture{
		DistinctId: deviceID,
		Event:      event,
		Properties: props,
	})
}

func Close() {
	if client != nil {
		client.Close()
	}
}

func isOptedOut(cfg config.TelemetryConfig) bool {
	if !cfg.Enabled {
		return true
	}

	if os.Getenv("GIT_SYNC_NO_TELEMETRY") == "1" {
		return true
	}

	return false
}

func getOrCreateDeviceID() (string, error) {
	configDir := config.GetDefaultConfigDir()
	idFile := filepath.Join(configDir, ".device-id")

	data, err := os.ReadFile(idFile)
	if err == nil && len(data) > 0 {
		return string(data), nil
	}

	id := uuid.New().String()

	if err := os.MkdirAll(configDir, os.ModePerm); err != nil {
		return "", err
	}

	if err := os.WriteFile(idFile, []byte(id), 0600); err != nil {
		return "", err
	}

	return id, nil
}

func Reset() {
	if client != nil {
		client.Close()
		client = nil
	}
	deviceID = ""
	disabled = false
	once = sync.Once{}
}
