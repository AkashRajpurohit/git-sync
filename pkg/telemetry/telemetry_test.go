package telemetry

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/AkashRajpurohit/git-sync/pkg/config"
	"github.com/AkashRajpurohit/git-sync/pkg/logger"
)

func TestMain(m *testing.M) {
	logger.InitLogger("fatal")
	os.Exit(m.Run())
}

func TestIsOptedOut_ConfigDisabled(t *testing.T) {
	cfg := config.TelemetryConfig{Enabled: false}
	if !isOptedOut(cfg) {
		t.Error("expected opt-out when config.Enabled is false")
	}
}

func TestIsOptedOut_ConfigEnabled(t *testing.T) {
	cfg := config.TelemetryConfig{Enabled: true}
	if isOptedOut(cfg) {
		t.Error("expected opt-in when config.Enabled is true and no env vars set")
	}
}

func TestIsOptedOut_GitSyncNoTelemetryEnv(t *testing.T) {
	t.Setenv("GIT_SYNC_NO_TELEMETRY", "1")
	cfg := config.TelemetryConfig{Enabled: true}
	if !isOptedOut(cfg) {
		t.Error("expected opt-out when GIT_SYNC_NO_TELEMETRY=1")
	}
}

func TestGetOrCreateDeviceID(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	id1, err := getOrCreateDeviceID()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id1 == "" {
		t.Fatal("expected non-empty device ID")
	}

	// Second call should return the same ID
	id2, err := getOrCreateDeviceID()
	if err != nil {
		t.Fatalf("unexpected error on second call: %v", err)
	}
	if id1 != id2 {
		t.Errorf("expected same device ID, got %q and %q", id1, id2)
	}

	// Verify the file exists
	idFile := filepath.Join(tmpDir, ".config", "git-sync", ".device-id")
	data, err := os.ReadFile(idFile)
	if err != nil {
		t.Fatalf("expected device ID file to exist: %v", err)
	}
	if string(data) != id1 {
		t.Errorf("file content %q doesn't match returned ID %q", string(data), id1)
	}
}

func TestCaptureEventWhenDisabled(t *testing.T) {
	Reset()
	defer Reset()

	Init(config.TelemetryConfig{Enabled: false})

	// Should not panic or cause any side effects
	CaptureEvent("test_event", map[string]interface{}{
		"key": "value",
	})
}
