package utils

import (
	"regexp"
	"strings"
)

// Slugify converts a string into a URL-friendly/file-friendly slug.
// E.g. "Persib Bandung" -> "persib_bandung"
// E.g. "Persib Jakarta" -> "persija_jakarta"
func Slugify(s string) string {
	s = strings.ToLower(s)
	// Replace non-alphanumeric with underscore
	re := regexp.MustCompile(`[^a-z0-9]+`)
	s = re.ReplaceAllString(s, "_")
	// Trim underscores
	s = strings.Trim(s, "_")
	return s
}
