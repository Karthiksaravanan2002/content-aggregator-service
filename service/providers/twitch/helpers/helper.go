package helpers

import (
	"strconv"
	"strings"
	"time"
)

func GetString(m map[string]any, key string) string {
	if v, ok := m[key].(string); ok {
		return v
	}
	return ""
}

func GetInt64(m map[string]any, key string) int64 {
	switch v := m[key].(type) {
	case int64:
		return v
	case float64:
		return int64(v)
	case int:
		return int64(v)
	case string:
		i, _ := strconv.ParseInt(v, 10, 64)
		return i
	}
	return 0
}

func MustParseTime(s string) time.Time {
	t, _ := time.Parse(time.RFC3339, s)
	return t
}

func ReplaceThumbnail(url string) string {
	return strings.Replace(url, "{width}x{height}", "1280x720", 1)
}

func MakeStreamURL(login string) string {
	if login == "" {
		return ""
	}
	return "https://www.twitch.tv/" + login
}

func MakeChannelURL(id string) string {
	if id == "" {
		return ""
	}
	return "https://www.twitch.tv/" + id
}
