package utils

import "strings"

func GetIdFromUri(uri string) string {
	parts := strings.Split(uri, ":")
	return parts[len(parts)-1]
}

func GetIdFromUrl(url string) string {
	parts := strings.Split(strings.TrimSuffix(url, "/"), "/")
	return parts[len(parts)-1]
}

func ParseIdFromString(raw string) string {
	trimmed := strings.TrimSpace(raw)
	if strings.Contains(trimmed, "spotify:") {
		return GetIdFromUri(trimmed)
	}

	if strings.Contains(trimmed, "://") {
		return GetIdFromUrl(trimmed)
	}

	return trimmed
}
