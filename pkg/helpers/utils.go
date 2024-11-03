package helpers

import "path/filepath"

func IsIncludedInList(list []string, value string) bool {
	for _, includeValue := range list {
		match, err := filepath.Match(includeValue, value)
		if err == nil && match {
			return true
		}
	}
	return false
}
