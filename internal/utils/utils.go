package utils

import (
	"crypto/rand"
	"fmt"
	"time"
)

// GenerateUUID generates a UUID v4 using the standard library
func GenerateUUID() (string, error) {
	uuid := make([]byte, 16)
	_, err := rand.Read(uuid)
	if err != nil {
		return "", err
	}

	// Set version (4) and variant (RFC 4122)
	uuid[6] = (uuid[6] & 0x0f) | 0x40 // Version 4
	uuid[8] = (uuid[8] & 0x3f) | 0x80 // Variant RFC 4122

	return fmt.Sprintf("%x-%x-%x-%x-%x",
		uuid[0:4],
		uuid[4:6],
		uuid[6:8],
		uuid[8:10],
		uuid[10:16]), nil
}

// FormatTimeISO8601 formats a time.Time as ISO 8601 string
func FormatTimeISO8601(t time.Time) string {
	return t.UTC().Format(time.RFC3339)
}

// ParseTimeISO8601 parses an ISO 8601 string to time.Time
func ParseTimeISO8601(s string) (time.Time, error) {
	return time.Parse(time.RFC3339, s)
}

// NewCustomTime creates a new CustomTime from a time.Time
func NewCustomTime(t time.Time) interface{} {
	// This function is a placeholder for now
	// It will be implemented when we need to convert between time.Time and model.CustomTime
	// We can't directly import model here due to import cycle
	return t
}
