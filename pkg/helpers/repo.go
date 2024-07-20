package helpers

import "path/filepath"

func IsRepoExcluded(excludedRepos []string, repoName string) bool {
	for _, excludedRepo := range excludedRepos {
		match, err := filepath.Match(excludedRepo, repoName)
		if err == nil && match {
			return true
		}
	}
	return false
}

func IsRepoIncluded(includedRepos []string, repoName string) bool {
	for _, includedRepo := range includedRepos {
		match, err := filepath.Match(includedRepo, repoName)
		if err == nil && match {
			return true
		}
	}
	return false
}
