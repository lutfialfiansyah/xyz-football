package pagination

import (
	"encoding/base64"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// ParseParams extracts the cursor, limit, and search query parameters.
func ParseParams(c *gin.Context) (cursor string, limit int, q string) {
	cursor = c.Query("cursor")
	q = c.Query("q")

	limitStr := c.Query("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 15 // Default limit
	}

	// Max limit protection
	if limit > 50 {
		limit = 50
	}

	return cursor, limit, q
}

// EncodeCursor takes a timestamp and an ID to create a composite base64 cursor.
func EncodeCursor(t time.Time, id string) string {
	raw := t.Format(time.RFC3339Nano) + "|" + id
	return base64.StdEncoding.EncodeToString([]byte(raw))
}

// DecodeCursor parses the base64 composite cursor back into a timestamp and ID.
func DecodeCursor(encoded string) (time.Time, string, error) {
	if encoded == "" {
		return time.Time{}, "", errors.New("empty cursor")
	}

	decodedBytes, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return time.Time{}, "", errors.New("invalid cursor format")
	}

	parts := strings.SplitN(string(decodedBytes), "|", 2)
	if len(parts) != 2 {
		return time.Time{}, "", errors.New("invalid cursor parts")
	}

	t, err := time.Parse(time.RFC3339Nano, parts[0])
	if err != nil {
		return time.Time{}, "", errors.New("invalid cursor time format")
	}

	return t, parts[1], nil
}
