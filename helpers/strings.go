package helpers

import (
	"strings"
)

func ToTitle(s string) string {
	if len(s) == 0 {
		return s
	}
	if len(s) == 1 {
		return strings.ToUpper(s)
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func TrimToLastSlash(s string) string {
	lastSlashIndex := strings.LastIndex(s, "/")
	if lastSlashIndex == -1 {
		return s // No slash found, return the original string
	}
	return s[lastSlashIndex+1:]
}

func Contains(list []string, s string) bool {
	for _, item := range list {
		if item == s {
			return true
		}
	}
	return false
}
