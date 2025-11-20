package helper

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"google.golang.org/api/youtube/v3"
)

func ToInt64(v uint64) int64 {
	return int64(v)
}

func MakeVideoURL(videoID string) string {
	if videoID == "" {
		return ""
	}
	return "https://www.youtube.com/watch?v=" + videoID
}

func ParseISODuration(iso string) int64 {
	if iso == "" || len(iso) < 3 {
		return 0
	}

	var hours, minutes, seconds int64

	// Remove starting "PT"
	iso = strings.TrimPrefix(iso, "PT")

	num := ""
	for _, ch := range iso {
		if ch >= '0' && ch <= '9' {
			num += string(ch)
			continue
		}

		// ch is one of H, M, S — convert num accordingly
		switch ch {
		case 'H':
			hours, _ = strconv.ParseInt(num, 10, 64)
		case 'M':
			minutes, _ = strconv.ParseInt(num, 10, 64)
		case 'S':
			seconds, _ = strconv.ParseInt(num, 10, 64)
		}
		num = ""
	}

	total := hours*3600 + minutes*60 + seconds
	return total
}

// defaultString returns `fallback` if s is empty.
func DefaultString(s, fallback string) string {
	if s == "" {
		return fallback
	}
	return s
}

// Extracts highest-quality available thumbnail
func ExtractBestThumbnailFromYT(t *youtube.ThumbnailDetails) string {
	if t == nil {
		return ""
	}
	if t.Maxres != nil {
		return t.Maxres.Url
	}
	if t.High != nil {
		return t.High.Url
	}
	if t.Medium != nil {
		return t.Medium.Url
	}
	if t.Default != nil {
		return t.Default.Url
	}
	return ""
}

// Builds channel URL
func MakeChannelURL(channelID string) string {
	if channelID == "" {
		return ""
	}
	return "https://www.youtube.com/channel/" + channelID
}

func FormatViewCount(v int64) string {
	switch {
	case v >= 1_000_000_000:
		return fmt.Sprintf("%.1fB", float64(v)/1_000_000_000)
	case v >= 1_000_000:
		return fmt.Sprintf("%.1fM", float64(v)/1_000_000)
	case v >= 1_000:
		return fmt.Sprintf("%.1fK", float64(v)/1_000)
	default:
		return fmt.Sprintf("%d", v)
	}
}

func RelativeTime(now, t time.Time) string {
	d := now.Sub(t)

	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		return fmt.Sprintf("%d minutes ago", int(d.Minutes()))
	case d < 24*time.Hour:
		return fmt.Sprintf("%d hours ago", int(d.Hours()))
	case d < 7*24*time.Hour:
		return fmt.Sprintf("%d days ago", int(d.Hours()/24))
	default:
		return t.Format("02 Jan 2006")
	}
}
