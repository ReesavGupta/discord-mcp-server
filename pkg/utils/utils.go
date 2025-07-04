package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"
	"time"
)

// SanitizeFilename removes invalid characters from filename
func SanitizeFilename(filename string) string {
	// Remove invalid characters
	reg := regexp.MustCompile(`[<>:"/\\|?*]`)
	sanitized := reg.ReplaceAllString(filename, "_")

	// Limit length
	if len(sanitized) > 100 {
		sanitized = sanitized[:100]
	}

	return strings.TrimSpace(sanitized)
}

// GenerateRandomString generates a random string of specified length
func GenerateRandomString(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return ""
	}
	return hex.EncodeToString(bytes)[:length]
}

// FormatTimestamp formats time in a readable format
func FormatTimestamp(t time.Time) string {
	return t.Format("2006-01-02 15:04:05 MST")
}

// TruncateString truncates string to specified length with ellipsis
func TruncateString(s string, length int) string {
	if len(s) <= length {
		return s
	}
	return s[:length-3] + "..."
}

// ValidateDiscordID validates Discord snowflake ID format
func ValidateDiscordID(id string) bool {
	if len(id) < 17 || len(id) > 20 {
		return false
	}

	for _, char := range id {
		if char < '0' || char > '9' {
			return false
		}
	}

	return true
}

// ParseTimeFilter parses time filter string to time.Time
func ParseTimeFilter(timeStr string) (*time.Time, error) {
	formats := []string{
		time.RFC3339,
		"2006-01-02T15:04:05Z",
		"2006-01-02 15:04:05",
		"2006-01-02",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, timeStr); err == nil {
			return &t, nil
		}
	}

	return nil, fmt.Errorf("invalid time format: %s", timeStr)
}
