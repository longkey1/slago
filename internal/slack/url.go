package slack

import (
	"fmt"
	"regexp"
	"strings"
)

// URLInfo contains parsed information from a Slack URL
type URLInfo struct {
	ChannelID string
	MessageTS string
	ThreadTS  string
}

// ParseURL parses a Slack message URL and extracts channel ID and timestamps
func ParseURL(url string) (*URLInfo, error) {
	// Pattern: https://xxx.slack.com/archives/C123/p456[?thread_ts=789]
	channelRe := regexp.MustCompile(`/archives/([^/]+)/p(\d{13,16})`)
	matches := channelRe.FindStringSubmatch(url)
	if len(matches) < 3 {
		return nil, fmt.Errorf("invalid Slack URL format: %s", url)
	}

	info := &URLInfo{
		ChannelID: matches[1],
		MessageTS: normalizeTimestamp(matches[2]),
	}

	// Check for thread_ts parameter
	if strings.Contains(url, "thread_ts=") {
		threadRe := regexp.MustCompile(`thread_ts=([0-9.]+)`)
		threadMatches := threadRe.FindStringSubmatch(url)
		if len(threadMatches) > 1 {
			info.ThreadTS = normalizeTimestamp(threadMatches[1])
		}
	}

	return info, nil
}

// normalizeTimestamp converts a timestamp to the format "seconds.microseconds"
func normalizeTimestamp(raw string) string {
	if strings.Contains(raw, ".") {
		return raw
	}

	if len(raw) <= 6 {
		return raw
	}

	prefix := raw[:len(raw)-6]
	suffix := raw[len(raw)-6:]
	return prefix + "." + suffix
}
