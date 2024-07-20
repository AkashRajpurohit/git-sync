package helpers

func IsRepoExcluded(excludedRepos []string, repoName string) bool {
	for _, excludedRepo := range excludedRepos {
		if excludedRepo == repoName {
			return true
		}
	}
	return false
}

func IsRepoIncluded(includedRepos []string, repoName string) bool {
	for _, includedRepo := range includedRepos {
		if includedRepo == repoName {
			return true
		}
	}
	return false
}
