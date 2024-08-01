package helpers

import "path/filepath"

func IsExcludedInList(excludedOrgs []string, orgName string) bool {
	for _, excludedOrg := range excludedOrgs {
		match, err := filepath.Match(excludedOrg, orgName)
		if err == nil && match {
			return true
		}
	}
	return false
}

func IsIncludedInList(includedOrgs []string, orgName string) bool {
	for _, includedOrg := range includedOrgs {
		match, err := filepath.Match(includedOrg, orgName)
		if err == nil && match {
			return true
		}
	}
	return false
}
