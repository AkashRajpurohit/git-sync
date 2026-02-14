package issues

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

func ReadLastSyncTime(backupDir, owner, repo string) (time.Time, bool) {
	indexPath := filepath.Join(backupDir, owner, repo, "issues", "index.json")

	data, err := os.ReadFile(indexPath)
	if err != nil {
		return time.Time{}, false
	}

	var entries []IndexEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return time.Time{}, false
	}

	if len(entries) == 0 {
		return time.Time{}, false
	}

	var maxTime time.Time
	for _, e := range entries {
		if e.UpdatedAt.After(maxTime) {
			maxTime = e.UpdatedAt
		}
	}

	return maxTime, true
}

func WriteIssues(backupDir, owner, repo string, newIssues []Issue) error {
	baseDir := filepath.Join(backupDir, owner, repo, "issues")
	jsonDir := filepath.Join(baseDir, "json")
	mdDir := filepath.Join(baseDir, "md")

	if err := os.MkdirAll(jsonDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create json directory: %w", err)
	}
	if err := os.MkdirAll(mdDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create md directory: %w", err)
	}

	updatedEntries := make(map[int]IndexEntry, len(newIssues))
	for _, issue := range newIssues {
		if err := writeJSON(jsonDir, issue); err != nil {
			return fmt.Errorf("failed to write JSON for issue #%d: %w", issue.Number, err)
		}
		if err := writeMarkdown(mdDir, issue); err != nil {
			return fmt.Errorf("failed to write markdown for issue #%d: %w", issue.Number, err)
		}
		updatedEntries[issue.Number] = IndexEntry{
			Number:    issue.Number,
			Title:     issue.Title,
			State:     issue.State,
			UpdatedAt: issue.UpdatedAt,
		}
	}

	existingIndex := readExistingIndex(baseDir)
	for _, entry := range existingIndex {
		if _, updated := updatedEntries[entry.Number]; !updated {
			updatedEntries[entry.Number] = entry
		}
	}

	finalIndex := make([]IndexEntry, 0, len(updatedEntries))
	for _, entry := range updatedEntries {
		finalIndex = append(finalIndex, entry)
	}
	sort.Slice(finalIndex, func(i, j int) bool {
		return finalIndex[i].Number < finalIndex[j].Number
	})

	if err := writeIndex(baseDir, finalIndex); err != nil {
		return fmt.Errorf("failed to write index: %w", err)
	}

	return nil
}

func readExistingIndex(baseDir string) []IndexEntry {
	data, err := os.ReadFile(filepath.Join(baseDir, "index.json"))
	if err != nil {
		return nil
	}

	var entries []IndexEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil
	}
	return entries
}

func writeJSON(dir string, issue Issue) error {
	data, err := json.MarshalIndent(issue, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, fmt.Sprintf("%d.json", issue.Number)), data, 0644)
}

func writeMarkdown(dir string, issue Issue) error {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("# #%d: %s\n\n", issue.Number, issue.Title))
	sb.WriteString(fmt.Sprintf("- **State:** %s\n", issue.State))
	sb.WriteString(fmt.Sprintf("- **Author:** %s\n", issue.Author.Login))
	if len(issue.Labels) > 0 {
		sb.WriteString(fmt.Sprintf("- **Labels:** %s\n", strings.Join(issue.Labels, ", ")))
	}
	if len(issue.Assignees) > 0 {
		logins := make([]string, len(issue.Assignees))
		for i, a := range issue.Assignees {
			logins[i] = a.Login
		}
		sb.WriteString(fmt.Sprintf("- **Assignees:** %s\n", strings.Join(logins, ", ")))
	}
	if issue.Milestone != "" {
		sb.WriteString(fmt.Sprintf("- **Milestone:** %s\n", issue.Milestone))
	}
	sb.WriteString(fmt.Sprintf("- **Created:** %s\n", issue.CreatedAt.Format("2006-01-02 15:04:05 UTC")))
	sb.WriteString(fmt.Sprintf("- **Updated:** %s\n", issue.UpdatedAt.Format("2006-01-02 15:04:05 UTC")))
	if issue.ClosedAt != nil {
		sb.WriteString(fmt.Sprintf("- **Closed:** %s\n", issue.ClosedAt.Format("2006-01-02 15:04:05 UTC")))
	}
	if issue.URL != "" {
		sb.WriteString(fmt.Sprintf("- **URL:** %s\n", issue.URL))
	}

	sb.WriteString(fmt.Sprintf("\n---\n\n%s\n", issue.Body))

	if len(issue.Comments) > 0 {
		sb.WriteString("\n---\n\n## Comments\n")
		for _, comment := range issue.Comments {
			sb.WriteString(fmt.Sprintf("\n### %s commented on %s\n\n", comment.Author.Login, comment.CreatedAt.Format("2006-01-02 15:04:05 UTC")))
			sb.WriteString(comment.Body)
			sb.WriteString("\n")
		}
	}

	return os.WriteFile(filepath.Join(dir, fmt.Sprintf("%d.md", issue.Number)), []byte(sb.String()), 0644)
}

func writeIndex(dir string, index []IndexEntry) error {
	data, err := json.MarshalIndent(index, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, "index.json"), data, 0644)
}
