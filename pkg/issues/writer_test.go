package issues

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func sampleIssues() []Issue {
	now := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
	closed := time.Date(2024, 1, 16, 12, 0, 0, 0, time.UTC)

	return []Issue{
		{
			Number: 1,
			Title:  "Bug report",
			Body:   "Something is broken",
			State:  "closed",
			Author: User{Login: "alice", URL: "https://github.com/alice"},
			Labels: []string{"bug", "critical"},
			Assignees: []User{
				{Login: "bob", URL: "https://github.com/bob"},
			},
			Milestone: "v1.0",
			URL:       "https://github.com/owner/repo/issues/1",
			CreatedAt: now,
			UpdatedAt: now.Add(time.Hour),
			ClosedAt:  &closed,
			Comments: []Comment{
				{
					ID:        101,
					Body:      "I can reproduce this",
					Author:    User{Login: "bob", URL: "https://github.com/bob"},
					URL:       "https://github.com/owner/repo/issues/1#issuecomment-101",
					CreatedAt: now.Add(30 * time.Minute),
					UpdatedAt: now.Add(30 * time.Minute),
				},
			},
		},
		{
			Number:    2,
			Title:     "Feature request",
			Body:      "Please add dark mode",
			State:     "open",
			Author:    User{Login: "charlie", URL: "https://github.com/charlie"},
			Labels:    []string{"enhancement"},
			Assignees: []User{},
			URL:       "https://github.com/owner/repo/issues/2",
			CreatedAt: now,
			UpdatedAt: now,
			Comments:  []Comment{},
		},
	}
}

func TestWriteIssues(t *testing.T) {
	tmpDir := t.TempDir()
	owner := "testowner"
	repo := "testrepo"
	testIssues := sampleIssues()

	err := WriteIssues(tmpDir, owner, repo, testIssues)
	if err != nil {
		t.Fatalf("WriteIssues failed: %v", err)
	}

	baseDir := filepath.Join(tmpDir, owner, repo, "issues")

	// Verify directory structure
	for _, dir := range []string{"json", "md"} {
		info, err := os.Stat(filepath.Join(baseDir, dir))
		if err != nil {
			t.Fatalf("Expected directory %s to exist: %v", dir, err)
		}
		if !info.IsDir() {
			t.Fatalf("Expected %s to be a directory", dir)
		}
	}

	// Verify files exist
	expectedFiles := []string{
		"json/1.json", "json/2.json",
		"md/1.md", "md/2.md",
		"index.json",
	}
	for _, f := range expectedFiles {
		if _, err := os.Stat(filepath.Join(baseDir, f)); err != nil {
			t.Errorf("Expected file %s to exist: %v", f, err)
		}
	}
}

func TestWriteJSON(t *testing.T) {
	tmpDir := t.TempDir()
	issue := sampleIssues()[0]

	err := writeJSON(tmpDir, issue)
	if err != nil {
		t.Fatalf("writeJSON failed: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(tmpDir, "1.json"))
	if err != nil {
		t.Fatalf("Failed to read JSON file: %v", err)
	}

	var readBack Issue
	if err := json.Unmarshal(data, &readBack); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if readBack.Number != issue.Number {
		t.Errorf("Expected number %d, got %d", issue.Number, readBack.Number)
	}
	if readBack.Title != issue.Title {
		t.Errorf("Expected title %q, got %q", issue.Title, readBack.Title)
	}
	if readBack.State != issue.State {
		t.Errorf("Expected state %q, got %q", issue.State, readBack.State)
	}
	if readBack.Author.Login != issue.Author.Login {
		t.Errorf("Expected author %q, got %q", issue.Author.Login, readBack.Author.Login)
	}
	if len(readBack.Labels) != len(issue.Labels) {
		t.Errorf("Expected %d labels, got %d", len(issue.Labels), len(readBack.Labels))
	}
	if len(readBack.Comments) != len(issue.Comments) {
		t.Errorf("Expected %d comments, got %d", len(issue.Comments), len(readBack.Comments))
	}
	if readBack.ClosedAt == nil {
		t.Error("Expected ClosedAt to be non-nil")
	}
}

func TestWriteMarkdown(t *testing.T) {
	tmpDir := t.TempDir()
	issue := sampleIssues()[0]

	err := writeMarkdown(tmpDir, issue)
	if err != nil {
		t.Fatalf("writeMarkdown failed: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(tmpDir, "1.md"))
	if err != nil {
		t.Fatalf("Failed to read markdown file: %v", err)
	}

	content := string(data)

	expectedStrings := []string{
		"# #1: Bug report",
		"**State:** closed",
		"**Author:** alice",
		"**Labels:** bug, critical",
		"**Assignees:** bob",
		"**Milestone:** v1.0",
		"Something is broken",
		"## Comments",
		"bob commented on",
		"I can reproduce this",
	}

	for _, s := range expectedStrings {
		if !strings.Contains(content, s) {
			t.Errorf("Expected markdown to contain %q", s)
		}
	}
}

func TestWriteIssuesIncrementalMerge(t *testing.T) {
	tmpDir := t.TempDir()
	owner := "testowner"
	repo := "testrepo"

	// First sync: write 2 issues
	initial := sampleIssues()
	if err := WriteIssues(tmpDir, owner, repo, initial); err != nil {
		t.Fatalf("Initial WriteIssues failed: %v", err)
	}

	// Verify ReadLastSyncTime works
	since, ok := ReadLastSyncTime(tmpDir, owner, repo)
	if !ok {
		t.Fatal("Expected ReadLastSyncTime to return true")
	}
	// Max UpdatedAt from sample issues is issue #1 (now + 1 hour)
	expected := initial[0].UpdatedAt
	if !since.Equal(expected) {
		t.Errorf("Expected last sync time %v, got %v", expected, since)
	}

	// Incremental sync: only update issue #2 (simulating it changed)
	later := time.Date(2024, 2, 1, 10, 0, 0, 0, time.UTC)
	updatedIssue := Issue{
		Number:    2,
		Title:     "Feature request (updated)",
		Body:      "Dark mode with auto-detection",
		State:     "closed",
		Author:    User{Login: "charlie", URL: "https://github.com/charlie"},
		Labels:    []string{"enhancement", "done"},
		Assignees: []User{},
		URL:       "https://github.com/owner/repo/issues/2",
		CreatedAt: initial[1].CreatedAt,
		UpdatedAt: later,
		Comments:  []Comment{},
	}
	if err := WriteIssues(tmpDir, owner, repo, []Issue{updatedIssue}); err != nil {
		t.Fatalf("Incremental WriteIssues failed: %v", err)
	}

	// Verify index has both issues (old #1 preserved + updated #2)
	indexPath := filepath.Join(tmpDir, owner, repo, "issues", "index.json")
	data, err := os.ReadFile(indexPath)
	if err != nil {
		t.Fatalf("Failed to read index.json: %v", err)
	}

	var index []IndexEntry
	if err := json.Unmarshal(data, &index); err != nil {
		t.Fatalf("Failed to unmarshal index: %v", err)
	}

	if len(index) != 2 {
		t.Fatalf("Expected 2 index entries, got %d", len(index))
	}

	// Index should be sorted by number
	if index[0].Number != 1 || index[1].Number != 2 {
		t.Errorf("Expected issues 1 and 2 in order, got %d and %d", index[0].Number, index[1].Number)
	}

	// Issue #1 should retain original title, issue #2 should be updated
	if index[0].Title != "Bug report" {
		t.Errorf("Issue #1 title should be preserved, got %q", index[0].Title)
	}
	if index[1].Title != "Feature request (updated)" {
		t.Errorf("Issue #2 title should be updated, got %q", index[1].Title)
	}
	if index[1].State != "closed" {
		t.Errorf("Issue #2 state should be updated to closed, got %q", index[1].State)
	}

	// Verify the updated MD file has new content
	mdData, err := os.ReadFile(filepath.Join(tmpDir, owner, repo, "issues", "md", "2.md"))
	if err != nil {
		t.Fatalf("Failed to read 2.md: %v", err)
	}
	if !strings.Contains(string(mdData), "Dark mode with auto-detection") {
		t.Error("Expected updated markdown content")
	}

	// ReadLastSyncTime should now return the later time
	newSince, ok := ReadLastSyncTime(tmpDir, owner, repo)
	if !ok {
		t.Fatal("Expected ReadLastSyncTime to return true after incremental sync")
	}
	if !newSince.Equal(later) {
		t.Errorf("Expected last sync time %v, got %v", later, newSince)
	}
}

func TestReadLastSyncTimeNoIndex(t *testing.T) {
	tmpDir := t.TempDir()
	_, ok := ReadLastSyncTime(tmpDir, "noowner", "norepo")
	if ok {
		t.Error("Expected ReadLastSyncTime to return false when no index exists")
	}
}

func TestWriteIssuesEmpty(t *testing.T) {
	tmpDir := t.TempDir()
	owner := "testowner"
	repo := "testrepo"

	err := WriteIssues(tmpDir, owner, repo, []Issue{})
	if err != nil {
		t.Fatalf("WriteIssues with empty input failed: %v", err)
	}

	// index.json should still be written
	indexPath := filepath.Join(tmpDir, owner, repo, "issues", "index.json")
	data, err := os.ReadFile(indexPath)
	if err != nil {
		t.Fatalf("Failed to read index.json: %v", err)
	}

	var index []IndexEntry
	if err := json.Unmarshal(data, &index); err != nil {
		t.Fatalf("Failed to unmarshal index: %v", err)
	}
	if len(index) != 0 {
		t.Errorf("Expected empty index, got %d entries", len(index))
	}
}
